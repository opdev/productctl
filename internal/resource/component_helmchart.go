package resource

type HelmChartDistributionMethod = string

type HelmChartComponent struct {
	ChartName          string   `json:"chart_name,omitempty"`
	Repository         string   `json:"repository,omitempty"`
	ShortDescription   string   `json:"short_description,omitempty"`
	LongDescription    string   `json:"long_description,omitempty"`
	GitHubUsernames    []string `json:"github_usernames,omitempty"`
	DistributionMethod string   `json:"distribution_method,omitempty"`
}
