package resource

type ContainerComponent struct {
	DistributionMethod    ContainerDistributionMethod   `json:"distribution_method,omitempty" jsonschema:"enum=rhcc,enum=external,enum=non_registry,enum=marketplace_only"`
	PID                   string                        `json:"isv_pid,omitempty"`
	OSContentType         ContainerComponentContentType `json:"os_content_type,omitempty" jsonschema:"enum=Red Hat Enterprise Linux,enum=Red Hat Universal Base Image (UBI),enum=Operator Bundle Image,enum=Scratch Image"`
	Privileged            *bool                         `json:"privileged,omitempty"`
	Registry              string                        `json:"registry,omitempty"`
	Repository            string                        `json:"repository,omitempty"`
	RepositoryDescription string                        `json:"repository_description,omitempty"`
	RepositoryName        string                        `json:"repository_name,omitempty"`
	ShortDescription      string                        `json:"short_description,omitempty"`
	SupportPlatforms      []string                      `json:"support_platforms,omitempty"`
	Type                  ContainerComponentType        `json:"type,omitempty" jsonschema:"enum=container,enum=operator bundle image"`
	GithubUsernames       []string                      `json:"github_usernames,omitempty"`
	HostedRegistry        *bool                         `json:"hosted_registry,omitempty"`
}

type ContainerDistributionMethod = string

const (
	ContainerDistributionRHCC            ContainerDistributionMethod = "rhcc"
	ContainerDistributionExternal        ContainerDistributionMethod = "external"
	ContainerDistributionNonRegistry     ContainerDistributionMethod = "non_registry"
	ContainerDistributionMarketplaceOnly ContainerDistributionMethod = "marketplace_only"
)

type ContainerComponentType = string

const (
	ContainerTypeContainer      ContainerComponentType = "container"
	ContainerTypeOperatorBundle ContainerComponentType = "operator bundle image"
)

type ContainerComponentContentType = string

const (
	ContentTypeRHEL           ContainerComponentContentType = "Red Hat Enterprise Linux"
	ContentTypeUBI            ContainerComponentContentType = "Red Hat Universal Base Image (UBI)"
	ContentTypeOperatorBundle ContainerComponentContentType = "Operator Bundle Image"
	ContentTypeScratch        ContainerComponentContentType = "Scratch Image"
)
