package catalogapi

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"github.com/opdev/productctl/internal/genpyxis"
	"github.com/opdev/productctl/internal/logger"
	"github.com/opdev/productctl/internal/resource"
)

var (
	ErrMissingName         = errors.New("listing did not have a name and it is required")
	ErrDetachingComponents = errors.New("unable to detach components from product listing")
)

// APIEndpoint represents a full URL to a given Catalog API instance.
type APIEndpoint = string

const (
	EndpointProduction APIEndpoint = "https://catalog.redhat.com/api/containers/graphql/"
	EndpointStage      APIEndpoint = "https://catalog.stage.redhat.com/api/containers/graphql/"
	EndpointQA         APIEndpoint = "https://catalog.qa.redhat.com/api/containers/graphql/"
	EndpointUAT        APIEndpoint = "https://catalog.uat.redhat.com/api/containers/graphql/"
)

// ApplyProduct will update an existing Product Listing if it exists (identified
// by the presence of an ID) or create the product listing if it does not
// already exists.
func ApplyProduct(
	ctx context.Context,
	client graphql.Client,
	declaration *resource.ProductListingDeclaration,
) (*resource.ProductListingDeclaration, error) {
	L := logger.FromContextOrDiscard(ctx)

	updateListing := false

	if declaration.Spec.HasID() {
		updateListing = true
	}

	if !declaration.Spec.HasName() {
		return nil, ErrMissingName
	}

	if updateListing {
		L = L.With("operation", "update")
	} else {
		L = L.With("operation", "create")
	}

	if updateListing {
		if len(declaration.With.Components) == 0 {
			L.Info("declaration enumerated no components. detaching all components from product (if necessary)")
			resp, err := genpyxis.SetComponentsForProduct(ctx, client, declaration.Spec.ID, []string{})
			if err != nil {
				return nil, errors.Join(ErrDetachingComponents, err)
			}

			if gqlErr := resp.Update_product_listing.GetError(); gqlErr != nil {
				return nil, errors.Join(ErrDetachingComponents, ParseGraphQLResponseError(gqlErr))
			}

			declaration.Spec.LastUpdateDate = resp.Update_product_listing.GetData().GetLast_update_date()
			declaration.Spec.CertProjects = resp.Update_product_listing.Data.GetCert_projects()
		}
	}

	newComponents := []genpyxis.CertificationProjectInput{}
	existingComponents := []genpyxis.CertificationProjectInput{}
	associatedComponentIDs := []string{}

	// Treat components that have IDs on-disk as pre-existing.
	for _, c := range declaration.With.Components {
		cInput, err := resource.JSONConvert[genpyxis.CertificationProjectInput](c)
		if err != nil {
			return nil, err
		}

		if cInput.Id == "" {
			L.Debug("component lacking id, treating as new", "component", logger.MarshalJSON(cInput))
			newComponents = append(newComponents, cInput)
			continue
		}

		L.Debug("component contained id, assuming pre-existing", "component", logger.MarshalJSON(cInput))

		existingComponents = append(existingComponents, cInput)
		associatedComponentIDs = append(associatedComponentIDs, cInput.Id)
	}

	// We assume components without IDs must be created.
	for _, newC := range newComponents {
		L.Debug("creating new component in backend", "component", logger.MarshalJSON(newC))

		// The backend complains if the project_status value isn't set for new
		// components, so we'll set it if the user hasn't.
		if newC.Project_status == "" {
			newC.Project_status = "active"
		}

		resp, err := genpyxis.NewComponent(ctx, client, &newC)
		if err != nil {
			return nil, err
		}

		if gqlErr := resp.Create_certification_project.GetError(); gqlErr != nil {
			return nil, ParseGraphQLResponseError(gqlErr)
		}

		associatedComponentIDs = append(associatedComponentIDs, resp.Create_certification_project.Data.GetId())
	}

	// existing components need to be applied
	for _, existingC := range existingComponents {
		L.Debug("applying pre-existing component's configuration", "name", existingC.Name, "id", existingC.Id)
		existingComponentAsInput, err := resource.JSONConvert[genpyxis.CertificationProjectInput](existingC)
		if err != nil {
			return nil, err
		}

		resp, err := genpyxis.ApplyComponent(ctx, client, existingC.Id, &existingComponentAsInput)
		if err != nil {
			return nil, err
		}

		if gqlErr := resp.Update_certification_project.GetError(); gqlErr != nil {
			return nil, ParseGraphQLResponseError(gqlErr)
		}
	}

	declaration.Spec.CertProjects = associatedComponentIDs
	L.Debug("components associated", "components", logger.MarshalJSON(declaration.Spec.CertProjects))

	// Create the ProductListing
	var input genpyxis.ProductListingInput
	input, err := resource.JSONConvert[genpyxis.ProductListingInput](declaration.Spec)
	if err != nil {
		return nil, err
	}

	var response *genpyxis.MutateProductListingCommonResponse
	var requestError error

	// This is ugly, but avoids having duplicate code paths (i.e previous and
	// following code) for apply/create workflows
	if updateListing {
		L.Debug("applying product listing")
		resp, err := genpyxis.ApplyProductListing(ctx, client, input.GetId(), &input)
		response = resp.GetUpdate_product_listing()
		requestError = err
	} else {
		L.Debug("creating product listing")
		resp, err := genpyxis.NewProductListing(ctx, client, &input)
		response = resp.GetCreate_product_listing()
		requestError = err
	}

	if requestError != nil {
		return nil, requestError
	}

	if gqlErr := response.GetError(); gqlErr != nil {
		return nil, ParseGraphQLResponseError(gqlErr)
	}

	returnedListing := response.GetData()

	L.Debug("retrieving updated data for associated components")
	updatedComponents, err := QueryAll(
		ctx,
		0,
		DefaultPageSize,
		func(page, pageSize int) (returnedItems []*genpyxis.ComponentSupportedFields, totalItems int, queryError error) {
			resp, err := genpyxis.ComponentsForListing(ctx, client, returnedListing.Id, page, pageSize)
			if err != nil {
				return nil, -10, err
			}

			if gqlErr := resp.Find_product_listing_certification_projects.GetError(); gqlErr != nil {
				return nil, -10, ParseGraphQLResponseError(gqlErr)
			}

			return resp.GetFind_product_listing_certification_projects().GetData(), resp.GetFind_product_listing_certification_projects().GetTotal(), nil
		},
	)
	if err != nil {
		return nil, err
	}

	L.Debug("updating manifest with updated component metadata")
	newComponentResources := make([]*resource.Component, 0, len(updatedComponents))
	for _, updatedC := range updatedComponents {
		converted, err := resource.JSONConvert[resource.Component](updatedC)
		if err != nil {
			return nil, err
		}
		newComponentResources = append(newComponentResources, &converted)
	}

	declaration.With.Components = newComponentResources

	// Finally update the product listing.
	L.Debug("updating manifest with updated product listing")
	finalListing, err := resource.JSONConvert[resource.ProductListing](returnedListing)
	if err != nil {
		return nil, err
	}

	declaration.Spec = finalListing

	return declaration, nil
}

// PopulateProduct will return a ProductListingDeclaration for the provided
// listingID.
func PopulateProduct(
	ctx context.Context,
	client graphql.Client,
	listingID string,
) (*resource.ProductListingDeclaration, error) {
	L := logger.FromContextOrDiscard(ctx)

	L.Debug("querying product by ID", "listingID", listingID)
	resp, err := genpyxis.ProductByID(ctx, client, listingID)
	if err != nil {
		return nil, err
	}

	if gqlErr := resp.Get_product_listing.GetError(); gqlErr != nil {
		return nil, ParseGraphQLResponseError(gqlErr)
	}

	newListing := resource.NewProductListing()
	newListing.Spec, err = resource.JSONConvert[resource.ProductListing](resp.GetGet_product_listing().GetData())
	if err != nil {
		return nil, err
	}

	if len(newListing.Spec.CertProjects) > 0 {
		L.Debug("populating components attached to product", "count", len(newListing.Spec.CertProjects))
		associatedComponents, err := QueryAll(
			ctx,
			0,
			DefaultPageSize,
			func(page, pageSize int) (returnedItems []*genpyxis.ComponentSupportedFields, totalItems int, queryError error) {
				resp, err := genpyxis.ComponentsForListing(ctx, client, newListing.Spec.ID, page, pageSize)
				if err != nil {
					return nil, -10, err
				}

				if gqlErr := resp.Find_product_listing_certification_projects.GetError(); gqlErr != nil {
					return nil, -10, ParseGraphQLResponseError(gqlErr)
				}

				return resp.GetFind_product_listing_certification_projects().GetData(), resp.GetFind_product_listing_certification_projects().GetTotal(), nil
			},
		)
		if err != nil {
			return nil, err
		}

		attachedIDs := make([]string, len(associatedComponents))
		newListing.With.Components = make([]*resource.Component, 0, len(associatedComponents))
		for _, v := range associatedComponents {
			attachedIDs = append(attachedIDs, v.Id)
			converted, err := resource.JSONConvert[resource.Component](v)
			if err != nil {
				return nil, err
			}

			newListing.With.Components = append(newListing.With.Components, &converted)
		}

		slices.SortStableFunc(newListing.With.Components, func(a, b *resource.Component) int { return strings.Compare(a.ID, b.ID) })

		// The cert_projects field in the product listing will contain
		// attached-but-archived components. ComponentsForListing (above)  only
		// returns active components, so we'll true this up on the client side.
		slices.Sort(attachedIDs)
		slices.Sort(newListing.Spec.CertProjects)

		L.Info("foo")
		if !slices.Equal(attachedIDs, newListing.Spec.CertProjects) {
			L.Debug("product listing has attached components that are not marked as active.")
			L.Debug(
				"replacing product_listing.spec.cert_projects with only attached IDs",
				"original", newListing.Spec.CertProjects,
				"replacement", attachedIDs,
			)
			newListing.Spec.CertProjects = attachedIDs
		}
	}

	return &newListing, nil
}
