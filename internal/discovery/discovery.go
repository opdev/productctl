// Package discovery adds relevant library functions for working with worklodas
// discovered via the discover-workload library.
// https://github.com/opdev/discover-workload.
package discovery

import (
	"errors"

	"github.com/opdev/discover-workload/discovery"

	"github.com/opdev/productctl/internal/resource"
)

var ErrDuplicateComponentName = errors.New("duplicate component name")

// ParseDiscoveryManifest converts discovered workloads into Component
// declarations. DiscoveredImages are treated as container components. The
// specifics of the container component declaration are assumed, and left to the
// user to change before use.
func ParseDiscoveryManifest(
	manifest discovery.Manifest,
) ([]*resource.Component, error) {
	if len(manifest.DiscoveredImages) == 0 {
		return nil, errors.New("could not find discovered images in discovery manifest")
	}

	processedNames := map[string]any{}

	components := make([]*resource.Component, 0, len(manifest.DiscoveredImages))
	for _, image := range manifest.DiscoveredImages {
		c := resource.Component{
			Container: &resource.ContainerComponent{
				DistributionMethod: resource.ContainerDistributionRHCC,
				OSContentType:      resource.ContentTypeUBI,
				Type:               resource.ContainerTypeContainer,
			},
			Name:          image.ContainerName,
			ProjectStatus: resource.ProjectStatusActive,
			Type:          resource.ComponentTypeContainer,
		}

		components = append(components, &c)
		// Components with the same name are not allowed.
		if _, exists := processedNames[c.Name]; exists {
			return nil, ErrDuplicateComponentName
		}

		processedNames[c.Name] = nil
	}

	return components, nil
}
