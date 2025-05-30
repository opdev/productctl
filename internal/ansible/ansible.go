package ansible

import (
	"context"
	"errors"
	"strings"

	"github.com/opdev/productctl/internal/logger"
	"github.com/opdev/productctl/internal/resource"
)

var ErrNoComponentsDeclared = errors.New("no components declared")

const (
	AnsibleConnectionLocal = "local"
	InventoryKeyContainer  = "container_components"
	InventoryKeyOperator   = "operator_components"
	InventoryKeyHelmChart  = "helm_chart_components"
)

// HelmChartComponentHostVars contains the variables for a container component
type HelmChartComponentHostVars struct {
	AnsibleConnection string              `json:"ansible_connection"`
	ChartURI          string              `json:"chart_uri"`
	Component         *resource.Component `json:"component"`
	Product           *ProductMeta        `json:"product"`
	// BUG: sigs.k8s.io/yaml does not seem to support inlining, or else we'd be
	// inlining this key. For now, we'll set a key.
	ToolFlags map[string]any `json:"tool_flags,inline,omitempty"`
}

// ContainerComponentHostVars contains the variables for a container component
// cert.
type ContainerComponentHostVars struct {
	AnsibleConnection string              `json:"ansible_connection"`
	Image             string              `json:"image"`
	Component         *resource.Component `json:"component"`
	Product           *ProductMeta        `json:"product"`
	// BUG: sigs.k8s.io/yaml does not seem to support inlining, or else we'd be
	// inlining this key. For now, we'll set a key.
	ToolFlags map[string]any `json:"tool_flags,inline,omitempty"`
}

type ProductMeta struct {
	ProductName string `json:"product_name"`
	ProductID   string `json:"product_id"`
}

// GenerateInventory produces the Ansible Inventory content representing the
// given product's components and the associated mapping containing
// certification targets for each component.
func GenerateInventory(
	ctx context.Context,
	product *resource.ProductListingDeclaration,
	mapping MappingDeclaration,
) (map[string]any, error) {
	L := logger.FromContextOrDiscard(ctx)
	if !product.HasComponents() {
		return nil, ErrNoComponentsDeclared
	}

	L.Debug("Checking product for container components")
	containerHosts, err := generateContainerComponentInventory(ctx, product, mapping.ContainerComponents)
	if err != nil {
		return nil, err
	}

	L.Debug("Checking product for helm chart components")
	helmHosts, err := generateHelmComponentInventory(ctx, product, mapping.HelmChartComponents)
	if err != nil {
		return nil, err
	}

	L.Debug("Checking product for operator components")
	operatorHosts, err := generateOperatorComponentInventory(ctx, product, mapping.OperatorComponents)
	if err != nil {
		return nil, err
	}

	inventory := map[string]any{
		InventoryKeyContainer: map[string]any{
			"hosts": containerHosts,
		},
		InventoryKeyOperator: map[string]any{
			"hosts": operatorHosts,
		},
		InventoryKeyHelmChart: map[string]any{
			"hosts": helmHosts,
		},
	}

	return inventory, nil
}

// generateContainerComponentInventory produces an Ansible inventory merging
// product and mapping.
func generateContainerComponentInventory(
	ctx context.Context,
	product *resource.ProductListingDeclaration,
	mapping ComponentCertificationConfig[ContainerCertTarget],
) (HostMap[ContainerComponentHostVars], error) { // nolint:unparam // ignoring lint flagging error return value that's always nil
	L := logger.FromContextOrDiscard(ctx)

	hosts := map[string]*ContainerComponentHostVars{}

	for _, cmp := range product.With.Components {
		if cmp.Container == nil {
			L.Debug(
				"Skipping component that does not have container metadata",
				"component_id", cmp.ID,
				"component_name", cmp.Name,
			)
			continue
		}

		if cmp.Container.OSContentType == resource.ContentTypeOperatorBundle {
			L.Debug(
				"Skipping component that isn't an application container type, inferred from content type metadata",
				"component_id", cmp.ID,
				"component_name", cmp.Name,
				"os_content_type", cmp.Container.OSContentType,
			)
			continue
		}

		var certificationCfg ContainerCertTarget
		var entryExists bool
		if certificationCfg, entryExists = mapping[cmp.ID]; !entryExists {
			L.Debug(
				"component did not have a mapping entry. refusing to generate an inventory entry for component",
				"component_id", cmp.ID,
				"component_name", cmp.Name,
			)
			continue
		}

		for _, tagConfig := range certificationCfg.Tags {
			imageTarget := strings.Join([]string{certificationCfg.ImageRef, tagConfig.Tag}, ":")
			vars := ContainerComponentHostVars{
				ToolFlags:         certificationCfg.ToolFlags,
				AnsibleConnection: AnsibleConnectionLocal,
				Image:             imageTarget,
				Product: &ProductMeta{
					ProductName: product.Spec.Name,
					ProductID:   product.Spec.ID,
				},
				Component: cmp,
			}

			// if tagConfig.ToolFlags != nil {
			if len(tagConfig.ToolFlags) > 0 {
				L.Debug("container tag has custom tool flags.")
				vars.ToolFlags = tagConfig.ToolFlags
			}

			// TODO: if it doesn't have a component ID, we should omit it or
			// throw an error.
			name := strings.Join([]string{cmp.ID, vars.Image}, "-")
			name = normalizeContainerHostname(name)
			hosts[name] = &vars
		}
	}

	return hosts, nil
}

// generateOperatorComponentInventory produces an ansible inventory items for
// operators, merging product declaration and certification target mappings.
func generateOperatorComponentInventory(
	ctx context.Context,
	product *resource.ProductListingDeclaration,
	mapping ComponentCertificationConfig[OperatorCertTarget],
) (HostMap[OperatorComponentHostVars], error) { // nolint:unparam // ignoring lint flagging error return value that's always nil
	L := logger.FromContextOrDiscard(ctx)

	hosts := HostMap[OperatorComponentHostVars]{}

	for _, cmp := range product.With.Components {
		if cmp.Container == nil {
			L.Debug(
				"Skipping component that does not have container metadata",
				"component_id", cmp.ID,
				"component_name", cmp.Name,
			)
			continue
		}

		if cmp.Container.OSContentType != resource.ContentTypeOperatorBundle {
			L.Debug(
				"Skipping component that isn't an operator container type, inferred from content type metadata",
				"component_id", cmp.ID,
				"component_name", cmp.Name,
				"os_content_type", cmp.Container.OSContentType,
			)
			continue
		}

		var certificationCfg OperatorCertTarget
		var entryExists bool
		if certificationCfg, entryExists = mapping[cmp.ID]; !entryExists {
			L.Debug(
				"component did not have a mapping entry. refusing to generate an inventory entry for component",
				"component_id", cmp.ID,
				"component_name", cmp.Name,
			)
			continue
		}

		// for i, imageTarget := range imageRefMatrix(certificationCfg.ImageRef, tagIter) {
		for _, tagConfig := range certificationCfg.Tags {
			imageTarget := strings.Join([]string{certificationCfg.ImageRef, tagConfig.Tag}, ":")
			vars := OperatorComponentHostVars{
				ToolFlags:         certificationCfg.ToolFlags,
				AnsibleConnection: AnsibleConnectionLocal,
				Image:             imageTarget,
				Product: &ProductMeta{
					ProductName: product.Spec.Name,
					ProductID:   product.Spec.ID,
				},
				Component:  cmp,
				IndexImage: certificationCfg.IndexImage,
			}

			// All "hosts" get the base index image unless they've specified one
			// along with a specific tag.
			if tagConfig.IndexImage != "" {
				L.Debug("operator tag has custom index image.")
				vars.IndexImage = tagConfig.IndexImage
			}

			if tagConfig.ToolFlags != nil {
				L.Debug("operator tag has custom tool flags.")
				vars.ToolFlags = tagConfig.ToolFlags
			}

			// TODO: if it doesn't have a component ID, we should omit it or
			// throw an error.
			name := strings.Join([]string{cmp.ID, vars.Image}, "-")
			name = normalizeContainerHostname(name)
			hosts[name] = &vars
		}
	}

	return hosts, nil
}

// generateHelmComponentInventory produces the helm components of the Ansible
// inventory.
func generateHelmComponentInventory(
	ctx context.Context,
	product *resource.ProductListingDeclaration,
	mapping ComponentCertificationConfig[HelmCertTarget],
) (HostMap[HelmChartComponentHostVars], error) { // nolint:unparam // ignoring lint flagging error return value that's always nil
	L := logger.FromContextOrDiscard(ctx)

	hosts := HostMap[HelmChartComponentHostVars]{}

	for _, cmp := range product.With.Components {
		if cmp.HelmChart == nil {
			L.Debug(
				"Skipping component that isn't a Helm chart, inferred from content type metadata",
				"component_id", cmp.ID,
				"component_name", cmp.Name,
				"component_type", cmp.Type,
			)
			continue
		}

		var certificationCfg HelmCertTarget
		var entryExists bool
		if certificationCfg, entryExists = mapping[cmp.ID]; !entryExists {
			L.Debug(
				"component did not have a mapping entry. refusing to generate an inventory entry for component",
				"component_id", cmp.ID,
				"component_name", cmp.Name,
			)
			continue
		}

		vars := HelmChartComponentHostVars{
			ToolFlags:         certificationCfg.ToolFlags,
			ChartURI:          certificationCfg.ChartURI,
			AnsibleConnection: AnsibleConnectionLocal,
			Product: &ProductMeta{
				ProductName: product.Spec.Name,
				ProductID:   product.Spec.ID,
			},
			Component: cmp,
		}

		name := strings.Join([]string{cmp.ID, vars.ChartURI}, "-")
		name = normalizeChartHostname(name)
		hosts[name] = &vars
	}

	return hosts, nil
}

// MappingDeclaration represents the extra data provided to be merged with a
// product declaration in order to produce a usable Ansible inventory.
//
// TODO: We should offer a JSONSchema corresponding to this data, given that
// this data will be hand-typed by a user.
type MappingDeclaration struct {
	ContainerComponents ComponentCertificationConfig[ContainerCertTarget] `json:"container_components"`
	HelmChartComponents ComponentCertificationConfig[HelmCertTarget]      `json:"helm_chart_components"`
	OperatorComponents  ComponentCertificationConfig[OperatorCertTarget]  `json:"operator_components"`
}

// HelmCertTarget contains required inputs necessary to run helm certification
// tooling.
type HelmCertTarget struct {
	ChartURI  string         `json:"chart_uri"`
	ToolFlags map[string]any `json:"tool_flags,omitempty"`
}

// OperatorCertTarget contains information necessary to run Operator
// certification tooling.
type OperatorCertTarget struct {
	// ImageRef is the image URI for the operator to certify. No tags should be
	// specified here.
	ImageRef string `json:"image_ref"`

	// Tags contains the list of tags to certify.
	Tags []OperatorTagandOptions `json:"tags"`

	// IndexImage is the catalog that contains the operator used for
	// certification. Users are required to build these catalogs.
	IndexImage string `json:"index_image"`
	// BUG: sigs.k8s.io/yaml does not seem to support inlining, or else we'd be
	// inlining this key. For now, we'll set a key.
	ToolFlags map[string]any `json:"tool_flags,inline,omitempty"`
}

// OperatorTagandOptions is a tag to certify, and any configurables that can
// change on a per-tag basis.
type OperatorTagandOptions struct {
	// The image tag to use for certification. E.g. v1.0.0.
	Tag string `json:"tag"`
	// IndexImage represents a custom index image to use for this particular
	// tag. If this is not set, the fallback IndexImage, set at the
	// OperatorCertTarget is used instead.
	IndexImage string `json:"index_image,omitempty"`
	// ToolFlags contains inputs necessary to run the certification tooling
	// against a specific tag.
	ToolFlags map[string]any `json:"tool_flags,inline,omitempty"`
}

// ContainertagAndOptions is a tag to certify, and any configurables that can
// change on a per-tag basis.
type ContainerTagAndOptions struct {
	// The image tag to use for certification. E.g. v1.0.0.
	Tag string `json:"tag"`
	// ToolFlags contains inputs necessary to run the certification tooling
	// against a specific tag.
	ToolFlags map[string]any `json:"tool_flags,inline,omitempty"`
}

// ContainerCertTarget represents the image's URI and all tags that the user
// wants to certify
type ContainerCertTarget struct {
	ImageRef string                   `json:"image_ref"`
	Tags     []ContainerTagAndOptions `json:"tags"`
	// BUG: sigs.k8s.io/yaml does not seem to support inlining, or else we'd be
	// inlining this key. For now, we'll set a key.
	ToolFlags map[string]any `json:"tool_flags,inline,omitempty"`
}

// ComponentCertificationConfig is a map of Component IDs to the corresponding
// certification inputs.
//
// In effect, this is a mapping of a component's ID (as it would exist in a
// product's declaration) to the configuration inputs needed to run
// certification tools against that product.
type ComponentCertificationConfig[T ContainerCertTarget | HelmCertTarget | OperatorCertTarget] map[string]T

type OperatorComponentHostVars struct {
	AnsibleConnection string              `json:"ansible_connection"`
	Image             string              `json:"image"`
	IndexImage        string              `json:"index_image"`
	Component         *resource.Component `json:"component"`
	Product           *ProductMeta        `json:"product"`
	// BUG: sigs.k8s.io/yaml does not seem to support inlining, or else we'd be
	// inlining this key. For now, we'll set a key.
	ToolFlags map[string]any `json:"tool_flags,inline,omitempty"`
}

// HostMap map is a representation of an ansible Host map, where a given host
// name corresponds to a collection of variables.
type HostMap[T ContainerComponentHostVars | HelmChartComponentHostVars | OperatorComponentHostVars] map[string]*T

// normalizeStringFn defines a string manipulation to be used in a processing strings.
type normalizeStringFn = func(s string) string

// normalizeString processes s with all normalizationFns, in order, and returns
// the result.
func normalizeString(s string, normalizationFns ...normalizeStringFn) string {
	for _, fn := range normalizationFns {
		s = fn(s)
	}

	return s
}

func normalizeContainerHostname(s string) string {
	return normalizeString(
		s,
		// Remove the colon from the tag
		func(n string) string { return strings.Replace(n, ":", "_", -1) },
		// Remove slashes from URI
		func(n string) string { return strings.Replace(n, "/", "_", -1) },
	)
}

func normalizeChartHostname(s string) string {
	return normalizeString(
		s,
		// Strip protocol
		func(n string) string { return strings.Replace(n, "https://", "remotechart-", -1) },
		func(n string) string { return strings.Replace(n, "http://", "remotechart-", -1) },
		// Remove remaining slashes from URI
		func(n string) string { return strings.Replace(n, "/", "_", -1) },
	)
}
