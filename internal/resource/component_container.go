package resource

type BuildCategory string

const (
	BuildCategoryStandaloneImage BuildCategory = "Standalone image"
	BuildCategoryComponentImage  BuildCategory = "Component image"
	BuildCategoryOperatorImage   BuildCategory = "Operator image"
	BuildCategoryOperatorBundle  BuildCategory = "Operator bundle"
)

type ContainerComponent struct {
	ApplicationCategories []ApplicationCategory `json:"application_categories,omitempty" jsonschema:"maxItems=3,enum=Accounting,enum=AI / Machine learning,enum=API Management,enum=Application Delivery,enum=Application Server,enum=Automation,enum=Backup & Recovery,enum=Business Intelligence,enum=Business Process Management,enum=Capacity Management,enum=Cloud Management,enum=Collaboration/Groupware/Messaging,enum=Configuration Management,enum=Console,enum=Container Platform / Management,enum=Content Management/Authoring,enum=Customer Relationship Management,enum=Dashboard,enum=Database & Data Management,enum=Data Store,enum=Developer Tools,enum=Enterprise Resource Planning,enum=Identity Management,enum=Integration,enum=Logging,enum=Logging & Metrics,enum=Management,enum=Messaging,enum=Metrics,enum=Migration,enum=Middleware,enum=Mobile Application Development Platform (MADP),enum=Monitoring,enum=Network Management,enum=Networking,enum=Observability,enum=Other,enum=Operating System,enum=Performance Management,enum=Plugin,enum=Policy Enforcement,enum=Programming Languages & Runtimes,enum=Scheduling,enum=Search,enum=Security,enum=Storage,enum=Tracing,enum=Virtualization Platform,enum=Web Services"`
	// NOTE: this is plural in the json and in the API spec, but doesn't
	// actually contain multiple values.
	BuildCategory         BuildCategory                 `json:"build_categories,omitempty" jsonschema:"enum=Standalone image,enum=Component image,enum=Operator image,enum=Operator bundle"`
	DistributionMethod    ContainerDistributionMethod   `json:"distribution_method,omitempty" jsonschema:"enum=rhcc,enum=external,enum=non_registry,enum=marketplace_only"`
	PID                   string                        `json:"isv_pid,omitempty"`
	OSContentType         ContainerComponentContentType `json:"os_content_type,omitempty" jsonschema:"enum=Red Hat Enterprise Linux,enum=Red Hat Universal Base Image (UBI),enum=Operator Bundle Image,enum=Scratch Image"`
	Privileged            *bool                         `json:"privileged,omitempty"`
	Registry              string                        `json:"registry,omitempty"`
	Repository            string                        `json:"repository,omitempty"`
	RepositoryDescription string                        `json:"repository_description,omitempty"`
	RepositoryName        string                        `json:"repository_name,omitempty"`
	ReleaseCategories     []ReleaseCategory             `json:"release_categories,omitempty" jsonschema:"enum=Generally Available,enum=Beta"`
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

type ReleaseCategory string

const (
	ReleaseCategoryGA   ReleaseCategory = "Generally Available"
	ReleaseCategoryBeta ReleaseCategory = "Beta"
)

type ApplicationCategory = string

const (
	ApplicationCategoryAccounting                      ApplicationCategory = "Accounting"
	ApplicationCategoryAIML                            ApplicationCategory = "AI / Machine learning"
	ApplicationCategoryAPIManagement                   ApplicationCategory = "API Management"
	ApplicationCategoryApplicationDelivery             ApplicationCategory = "Application Delivery"
	ApplicationCategoryApplicationServer               ApplicationCategory = "Application Server"
	ApplicationCategoryAutomation                      ApplicationCategory = "Automation"
	ApplicationCategoryBackupRecovery                  ApplicationCategory = "Backup & Recovery"
	ApplicationCategoryBusinessIntelligence            ApplicationCategory = "Business Intelligence"
	ApplicationCategoryBusinessProcessManagement       ApplicationCategory = "Business Process Management"
	ApplicationCategoryCapacityManagement              ApplicationCategory = "Capacity Management"
	ApplicationCategoryCloudManagement                 ApplicationCategory = "Cloud Management"
	ApplicationCategoryCollaborationGroupwareMessaging ApplicationCategory = "Collaboration/Groupware/Messaging"
	ApplicationCategoryConfigurationManagement         ApplicationCategory = "Configuration Management"
	ApplicationCategoryConsole                         ApplicationCategory = "Console"
	ApplicationCategoryContainerPlatformManagement     ApplicationCategory = "Container Platform / Management"
	ApplicationCategoryContentManagementAuthoring      ApplicationCategory = "Content Management/Authoring"
	ApplicationCategoryCRM                             ApplicationCategory = "Customer Relationship Management"
	ApplicationCategoryDashboard                       ApplicationCategory = "Dashboard"
	ApplicationCategoryDatabaseDataManagement          ApplicationCategory = "Database & Data Management"
	ApplicationCategoryDataStore                       ApplicationCategory = "Data Store"
	ApplicationCategoryDeveloperTools                  ApplicationCategory = "Developer Tools"
	ApplicationCategoryERP                             ApplicationCategory = "Enterprise Resource Planning"
	ApplicationCategoryIdentityManagement              ApplicationCategory = "Identity Management"
	ApplicationCategoryIntegration                     ApplicationCategory = "Integration"
	ApplicationCategoryLogging                         ApplicationCategory = "Logging"
	ApplicationCategoryLoggingMetrics                  ApplicationCategory = "Logging & Metrics"
	ApplicationCategoryManagement                      ApplicationCategory = "Management"
	ApplicationCategoryMessaging                       ApplicationCategory = "Messaging"
	ApplicationCategoryMetrics                         ApplicationCategory = "Metrics"
	ApplicationCategoryMigration                       ApplicationCategory = "Migration"
	ApplicationCategoryMiddleware                      ApplicationCategory = "Middleware"
	ApplicationCategoryMADP                            ApplicationCategory = "Mobile Application Development Platform (MADP)"
	ApplicationCategoryMonitoring                      ApplicationCategory = "Monitoring"
	ApplicationCategoryNetworkManagement               ApplicationCategory = "Network Management"
	ApplicationCategoryNetworking                      ApplicationCategory = "Networking"
	ApplicationCategoryObservability                   ApplicationCategory = "Observability"
	ApplicationCategoryOther                           ApplicationCategory = "Other"
	ApplicationCategoryOS                              ApplicationCategory = "Operating System"
	ApplicationCategoryPerformanceManagement           ApplicationCategory = "Performance Management"
	ApplicationCategoryPlugin                          ApplicationCategory = "Plugin"
	ApplicationCategoryPolicyEnforcement               ApplicationCategory = "Policy Enforcement"
	ApplicationCategoryProgLangsRuntimes               ApplicationCategory = "Programming Languages & Runtimes"
	ApplicationCategoryScheduling                      ApplicationCategory = "Scheduling"
	ApplicationCategorySearch                          ApplicationCategory = "Search"
	ApplicationCategorySecurity                        ApplicationCategory = "Security"
	ApplicationCategoryStorage                         ApplicationCategory = "Storage"
	ApplicationCategoryTracing                         ApplicationCategory = "Tracing"
	ApplicationCategoryVirtualizationPlatform          ApplicationCategory = "Virtualization Platform"
	ApplicationCategoryWebServices                     ApplicationCategory = "Web Services"
)
