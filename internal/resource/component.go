package resource

import "time"

type Component struct {
	ID                   string              `json:"_id,omitempty"`
	CertificationDate    *time.Time          `json:"certification_date,omitempty"`
	CertificationLevel   string              `json:"certification_level,omitempty"`
	Contacts             []ComponentContacts `json:"contacts,omitempty"`
	Container            *ContainerComponent `json:"container,omitempty"`
	Name                 string              `json:"name,omitempty"`
	OperatorDistribution string              `json:"operator_distribution,omitempty"`
	OrgID                int                 `json:"org_id,omitempty"`
	PublishedBy          string              `json:"published_by,omitempty"`
	Badges               []string            `json:"badges,omitempty"`
	Type                 ComponentType       `json:"type,omitempty" jsonschema:"enum=Containers,enum=Helm Chart,enum=OpenShift-cnf"`
	CreationDate         *time.Time          `json:"creation_date,omitempty"`
	HelmChart            *HelmChartComponent `json:"helm_chart,omitempty"`
	LastUpdateDate       *time.Time          `json:"last_update_date,omitempty"`
}

type ComponentContacts struct {
	EmailAddress string `json:"email_address,omitempty"`
	Type         string `json:"type,omitempty" jsonschema:"enum=Technical contact"`
}

type ComponentType = string

const (
	ComponentTypeContainer ComponentType = "Containers"
	ComponentTypeHelmChart ComponentType = "Helm Chart"
	ComponentTypeCNF       ComponentType = "OpenShift-cnf"
	// ComponentTypeOCPVirt                    ComponentType = "OpenShift-virtualization"
	// ComponentTypeOpenStackInfraContainer    ComponentType = "OpenStack-infra-container"
	// ComponentTypeOpenStackInfraNonContainer ComponentType = "OpenStack-infra-noncontainer"
	// ComponentTypeOpenStackVNF               ComponentType = "OpenStack-vnf"
	// ComponentTypeOpenStackAppContainer      ComponentType = "OpenStack-app-container"
	// ComponentTypeRHEL                       ComponentType = "RHEL"
)
