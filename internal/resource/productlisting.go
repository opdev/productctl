package resource

import "time"

type ProductListingType = string

const (
	ProductListingTypeContainerStack         ProductListingType = "container stack"
	ProductListingTypeTraditionalApplication ProductListingType = "traditional application"
	ProductListingTypeOpenStackInfra         ProductListingType = "openstack infra"
)

type ProductListing struct {
	ID             string                      `json:"_id,omitempty"`
	Name           string                      `json:"name,omitempty"`
	OrgID          int                         `json:"org_id,omitempty"`
	LastUpdateDate *time.Time                  `json:"last_update_date,omitempty"`
	Type           ProductListingType          `json:"type,omitempty" jsonschema:"enum=container stack,enum=traditional application,enum=openstack infra"`
	Descriptions   *ProductListingDescriptions `json:"descriptions,omitempty"`
	Contacts       []ProductListingContact     `json:"contacts,omitempty"`
	CreationDate   *time.Time                  `json:"creation_date,omitempty"`
	CertProjects   []string                    `json:"cert_projects,omitempty"`
	Support        *ProductListingSupport      `json:"support,omitempty"`
	Legal          *ProductListingLegal        `json:"legal,omitempty"`
}

func (p *ProductListing) HasName() bool {
	return p.Name != ""
}

func (p *ProductListing) HasID() bool {
	return p.ID != ""
}

type ProductListingSupport struct {
	Description  string `json:"description,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	URL          string `json:"url,omitempty"`
}

type ProductListingLegal struct {
	Description         string `json:"description,omitempty"`
	LicenseAgreementURL string `json:"license_agreement_url,omitempty"`
	PrivacyPolicyURL    string `json:"privacy_policy_url,omitempty"`
}

type ProductListingDescriptions struct {
	Long  string `json:"long,omitempty"`
	Short string `json:"short,omitempty"`
}

type ProductListingContact struct {
	EmailAddress string `json:"email_address,omitempty"`
	Type         string `json:"type,omitempty"`
}

/*
Below PoC leverages genqlient's generated interface methods and allows us to
convert them to the on-disk types contained in this package. A

"generator"-like type exists so that the function signature for ContactsFrom
doesn't contain []interface, which doesn't work nicely IIRC
*/

// ContactInfoProvider is an interface definition matching the various generated
// GraphQL types that might contain contact information.
type ContactInfoProvider interface {
	GetEmail_address() string
	GetType() string
}

// ContactsFrom produces a ProductListingContact to support the varied GraphQL
// types we might generate.
func ContactsFrom[T ContactInfoProvider](in []T) []ProductListingContact {
	out := make([]ProductListingContact, 0, len(in))
	for i := range in {
		out = append(out, ProductListingContact{
			EmailAddress: in[i].GetEmail_address(),
			Type:         in[i].GetType(),
		})
	}

	return out
}
