package resource

import "time"

type ProductListingType = string

const (
	ProductListingTypeContainerStack         ProductListingType = "container stack"
	ProductListingTypeTraditionalApplication ProductListingType = "traditional application"
	ProductListingTypeOpenStackInfra         ProductListingType = "openstack infra"
)

type ProductListing struct {
	ID                      string                                `json:"_id,omitempty"`
	Name                    string                                `json:"name,omitempty"`
	OrgID                   int                                   `json:"org_id,omitempty"`
	LastUpdateDate          *time.Time                            `json:"last_update_date,omitempty"`
	Type                    ProductListingType                    `json:"type,omitempty" jsonschema:"enum=container stack,enum=traditional application,enum=openstack infra"`
	Descriptions            *ProductListingDescriptions           `json:"descriptions,omitempty"`
	Contacts                []ProductListingContact               `json:"contacts,omitempty" jsonschema:"minItems=1,maxItems=10"`
	CreationDate            *time.Time                            `json:"creation_date,omitempty"`
	CertProjects            []string                              `json:"cert_projects,omitempty"`
	Support                 *ProductListingSupport                `json:"support,omitempty"`
	Legal                   *ProductListingLegal                  `json:"legal,omitempty"`
	LinkedResources         []ProductListingLinkedResource        `json:"linked_resources,omitempty" jsonschema:"minItems=3,maxItems=8"`
	FAQs                    []FAQ                                 `json:"faqs,omitempty"`
	SearchAliases           []SearchAlias                         `json:"search_aliases,omitempty" jsonschema:"maxItems=5" jsonschema_description:"A collection of key value pairs used assist users searching for your product listing"`
	FunctionalCategory      []FunctionalCategory                  `json:"functional_categories,omitempty" jsonschema:"minItems=1,maxItems=3,enum=AI/ML,enum=Analytics,enum=App dev,enum=App modernization,enum=Automation,enum=Backup & Recovery,enum=Cloud,enum=Compute,enum=Content management,enum=Data management,enum=Developer tools,enum=DevOps,enum=Edge,enum=Infrastructure,enum=IT & management tools,enum=Migration,enum=Networking,enum=Observability,enum=Orchestration,enum=OS & platforms,enum=Security,enum=Storage,enum=Virtualization"`
	QuickStartConfiguration ProductListingQuickStartConfiguration `json:"quick_start_configuration,omitempty"`
	Features                []ProductListingFeature               `json:"features,omitempty"`
}

func (p *ProductListing) HasName() bool {
	return p.Name != ""
}

func (p *ProductListing) HasID() bool {
	return p.ID != ""
}

type ProductListingQuickStartConfiguration struct {
	Instructions string `json:"instructions,omitempty" jsonschema:"minLength=1,maxLength=10000" jsonschema_description:"Quick start instructions for your users. Supports HTML formatting."`
}

type ProductListingFeature struct {
	Title       string `json:"title,omitempty" jsonschema:"maxLength=60" jsonschema_description:"The title of a supported feature."`
	Description string `json:"description,omitempty" jsonschema:"maxLength=1000" jsonschema_description:"A description of the titled feature. Supports HTML Formatting."`
}

type ProductListingSupport struct {
	URL          string `json:"url,omitempty"`
	Description  string `json:"description,omitempty" jsonschema:"minLength=1,maxLength-=500"`
	EmailAddress string `json:"email_address,omitempty"`
	PhoneNumber  string `json:"phone_number,omitempty" jsonschema:"minLength=1,maxLength=50"`
}

type ProductListingLegal struct {
	LicenseAgreementURL string `json:"license_agreement_url,omitempty"`
	PrivacyPolicyURL    string `json:"privacy_policy_url,omitempty"`
}

type ProductListingDescriptions struct {
	Long  string `json:"long,omitempty"`
	Short string `json:"short,omitempty" jsonschema:"minLength=50"`
}

type ProductListingContact struct {
	EmailAddress string `json:"email_address,omitempty"`
	Type         string `json:"type,omitempty" jsonschema:"enum=Marketing contact,enum=Technical contact"`
}

type ProductListingLinkedResource struct {
	Title       string             `json:"title,omitempty"`
	Description string             `json:"description,omitempty"`
	Type        LinkedResourceType `json:"type,omitempty" jsonschema:"enum=Video,enum=Article,enum=Documentation,enum=Website,enum=Podcasts,enum=On-demand Events"`
	URL         string             `json:"url,omitempty"`
}

type LinkedResourceType = string

const (
	LinkedResourceTypeVideo         LinkedResourceType = "Video"
	LinkedResourceTypeArticle       LinkedResourceType = "Article"
	LinkedResourceTypeDocumentation LinkedResourceType = "Documentation"
	LinkedResourceTypeWebsite       LinkedResourceType = "Website"
	LinkedResourceTypePodcasts      LinkedResourceType = "Podcasts"
	LinkedResourceTypeOnDemandEvent LinkedResourceType = "On-demand Events"
)

type FAQ struct {
	Question string `json:"question,omitempty" jsonschema:"maxLength=500" jsonschema_description:"Common questions"`
	Answer   string `json:"answer,omitempty" jsonschema:"maxLength=10000" jsonschema_description:"Answers to your questions. May contain HTML."`
}

type SearchAlias struct {
	Key   string `json:"key,omitempty" jsonschema:"minLength=1,maxLength=50" jsonschema_description:"Acronyms, or short identifiers related to the project. E.g. \"RHEL\""`
	Value string `json:"value,omitempty" jsonschema:"minLength=1,maxLength=100" jsonschema_description:"Related long-form meaning of the related acronym or short identifier. E.g. \"Red Hat Enterprise Linux\""`
}

type FunctionalCategory = string

const (
	FunctionalCategoryAIML              FunctionalCategory = "AI/ML"
	FunctionalCategoryAnalytics         FunctionalCategory = "Analytics"
	FunctionalCategoryAppDev            FunctionalCategory = "App dev"
	FunctionalCategoryAppModernization  FunctionalCategory = "App modernization"
	FunctionalCategoryAutomation        FunctionalCategory = "Automation"
	FunctionalCategoryBackupRecovery    FunctionalCategory = "Backup & Recovery"
	FunctionalCategoryCloud             FunctionalCategory = "Cloud"
	FunctionalCategoryCompute           FunctionalCategory = "Compute"
	FunctionalCategoryContentManagement FunctionalCategory = "Content management"
	FunctionalCategoryDataManagement    FunctionalCategory = "Data management"
	FunctionalCategoryDeveloperTools    FunctionalCategory = "Developer tools"
	FunctionalCategoryDevOps            FunctionalCategory = "DevOps"
	FunctionalCategoryEdge              FunctionalCategory = "Edge"
	FunctionalCategoryInfrastructure    FunctionalCategory = "Infrastructure"
	FunctionalCategoryITManagementTools FunctionalCategory = "IT & management tools"
	FunctionalCategoryMigration         FunctionalCategory = "Migration"
	FunctionalCategoryNetworking        FunctionalCategory = "Networking"
	FunctionalCategoryObservability     FunctionalCategory = "Observability"
	FunctionalCategoryOrchestration     FunctionalCategory = "Orchestration"
	FunctionalCategoryOSPlatforms       FunctionalCategory = "OS & platforms"
	FunctionalCategorySecurity          FunctionalCategory = "Security"
	FunctionalCategoryStorage           FunctionalCategory = "Storage"
	FunctionalCategoryVirtualization    FunctionalCategory = "Virtualization"
)

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
