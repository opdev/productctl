// Package resource contains bits of code necessary for the version of resources
// stored on a user's disk. These will look similar to the API types, but will
// differ in subtle ways that enable the on-disk representation.
package resource

type ProductListingDeclaration struct {
	Kind string         `json:"kind"`
	Spec ProductListing `json:"spec"`
	With Inclusions     `json:"with,omitempty"`
}

// NewProductListing returns a net-new product listing declaration.
func NewProductListing() ProductListingDeclaration {
	return ProductListingDeclaration{
		Kind: "ProductListing",
		Spec: ProductListing{},
	}
}

// Inclusions enumerates any additional resources that would otherwise be tied
// directly to the Product Listing.
type Inclusions struct {
	// Components represent the certification projects associated with the given
	// product listing.
	Components []*Component `json:"components,omitempty"`
}

// Sanitize removes identifiers that tie this declaration to a specific entity
// in the Catalog. Common for cases where a given declaration is going to be
// stored for re-use.
func (d *ProductListingDeclaration) Sanitize() {
	d.Spec.CertProjects = nil
	d.Spec.CreationDate = nil
	d.Spec.LastUpdateDate = nil
	d.Spec.ID = ""
	d.Spec.OrgID = 0

	for i := range d.With.Components {
		d.With.Components[i].ID = ""
		d.With.Components[i].OrgID = 0
		if c := d.With.Components[i].Container; c != nil {
			c.PID = ""
		}
		d.With.Components[i].CreationDate = nil
		d.With.Components[i].LastUpdateDate = nil
		d.With.Components[i].CertificationDate = nil
	}
}

func (d *ProductListingDeclaration) HasComponents() bool {
	return len(d.With.Components) > 0
}
