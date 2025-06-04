package catalogapi

import (
	"context"

	"github.com/Khan/genqlient/graphql"

	"github.com/opdev/productctl/internal/genpyxis"
	"github.com/opdev/productctl/internal/logger"
	"github.com/opdev/productctl/internal/resource"
)

// CleanupProduct will, if able, detach and archive all components on a product
// listing. Then, it will archive the product listing, and sanitize the listing.
func CleanupProduct(
	ctx context.Context,
	client graphql.Client,
	declaration *resource.ProductListingDeclaration,
) (*resource.ProductListingDeclaration, error) {
	L := logger.FromContextOrDiscard(ctx)
	listingExists := declaration.Spec.ID != ""

	if listingExists {
		L.Info("detaching any and all components from product listing", "productListingID", declaration.Spec.ID, "productListingName", declaration.Spec.Name)
		resp, err := genpyxis.SetComponentsForProduct(ctx, client, declaration.Spec.ID, []string{})
		if err != nil {
			return nil, err
		}

		if gqlErr := resp.Update_product_listing.GetError(); gqlErr != nil {
			return nil, ParseGraphQLResponseError(gqlErr)
		}
	}

	for _, component := range declaration.With.Components {
		if component.ID == "" {
			continue
		}

		L.Info("archiving component", "id", component.ID, "name", component.Name, "type", component.Type)
		resp, err := genpyxis.ArchiveComponent(ctx, client, component.ID)
		if err != nil {
			return nil, err
		}

		if gqlErr := resp.Update_certification_project.GetError(); gqlErr != nil {
			return nil, ParseGraphQLResponseError(gqlErr)
		}
	}

	if listingExists {
		L.Info("deleting product listing")
		resp, err := genpyxis.DeleteProduct(ctx, client, declaration.Spec.ID)
		if err != nil {
			return nil, err
		}

		if gqlErr := resp.Update_product_listing.GetError(); gqlErr != nil {
			return nil, ParseGraphQLResponseError(gqlErr)
		}
	}

	L.Info("cleanup API calls completed")
	declaration.Sanitize()

	return declaration, nil
}
