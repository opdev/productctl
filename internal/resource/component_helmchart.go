package resource

type HelmChartComponent struct {
	ApplicationCategories []ApplicationCategory       `json:"application_categories,omitempty" jsonschema:"maxItems=3,enum=Accounting,enum=AI / Machine learning,enum=API Management,enum=Application Delivery,enum=Application Server,enum=Automation,enum=Backup & Recovery,enum=Business Intelligence,enum=Business Process Management,enum=Capacity Management,enum=Cloud Management,enum=Collaboration/Groupware/Messaging,enum=Configuration Management,enum=Console,enum=Container Platform / Management,enum=Content Management/Authoring,enum=Customer Relationship Management,enum=Dashboard,enum=Database & Data Management,enum=Data Store,enum=Developer Tools,enum=Enterprise Resource Planning,enum=Identity Management,enum=Integration,enum=Logging,enum=Logging & Metrics,enum=Management,enum=Messaging,enum=Metrics,enum=Migration,enum=Middleware,enum=Mobile Application Development Platform (MADP),enum=Monitoring,enum=Network Management,enum=Networking,enum=Observability,enum=Other,enum=Operating System,enum=Performance Management,enum=Plugin,enum=Policy Enforcement,enum=Programming Languages & Runtimes,enum=Scheduling,enum=Search,enum=Security,enum=Storage,enum=Tracing,enum=Virtualization Platform,enum=Web Services"`
	ChartName             string                      `json:"chart_name,omitempty"`
	Repository            string                      `json:"repository,omitempty"`
	ShortDescription      string                      `json:"short_description,omitempty"`
	LongDescription       string                      `json:"long_description,omitempty"`
	GitHubUsernames       []string                    `json:"github_usernames,omitempty"`
	DistributionMethod    HelmChartDistributionMethod `json:"distribution_method,omitempty" jsonschema:"enum=redhat,enum=external,enum=undistributed"`
}

// HelmChartDistributionMethod is a string alias for the distribution method of a Helm chart.
type HelmChartDistributionMethod = string

const (
	HelmChartDistributionMethodRedHat        HelmChartDistributionMethod = "redhat"
	HelmChartDistributionMethodExternal      HelmChartDistributionMethod = "external"
	HelmChartDistributionMethodUndistributed HelmChartDistributionMethod = "undistributed"
)
