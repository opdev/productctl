// Package discovery adds relevant library functions for working with workloads
// discovered via the discover-workload library.
// https://github.com/opdev/discover-workload.
package discovery

import (
	"errors"
	"fmt"
	"maps"

	"github.com/opdev/discover-workload/discovery"

	"github.com/opdev/productctl/internal/resource"
)

var (
	ErrDuplicateComponentName           = errors.New("duplicate component name")
	ErrNoImagesDiscovered               = errors.New("no images found in discovery manifest")
	ErrCountingContainerNameOccurrences = errors.New("failed counting container names")
	ErrNoContainers                     = errors.New("no containers in entry for image")
)

// ComponentsFromDiscoveryManifest converts discovered workloads into Component
// declarations. DiscoveredImages are treated as container components. The
// specifics of the container component declaration are assumed, and left to the
// user to change before use.
func ComponentsFromDiscoveryManifest(manifest discovery.Manifest) ([]*resource.Component, error) {
	if len(manifest.DiscoveredImages) == 0 {
		return nil, ErrNoImagesDiscovered
	}

	processedNames := map[string]any{}

	components := make([]*resource.Component, 0, len(manifest.DiscoveredImages))
	for _, image := range manifest.DiscoveredImages {
		mostCommonContainerName, _, err := mostFrequentName(image.Containers)
		if err != nil {
			return nil, err
		}
		c := resource.Component{
			Container: &resource.ContainerComponent{
				DistributionMethod: resource.ContainerDistributionExternal,
				OSContentType:      resource.ContentTypeUBI,
				Type:               resource.ContainerTypeContainer,
			},
			Name: mostCommonContainerName,
			Type: resource.ComponentTypeContainer,
		}

		components = append(components, &c)
		// Components with the same name are not allowed.
		if _, exists := processedNames[c.Name]; exists {
			return nil, fmt.Errorf("%w: component with name %s is defined more than once", ErrDuplicateComponentName, c.Name)
		}

		processedNames[c.Name] = nil
	}

	return components, nil
}

// mostFrequentName returns the most frequent discovered container name for
// items in containers, with a goal of providing the user with a likely name to
// use for a certification component for the specified container. Returned is
// the actual container name and how many times it was found in the input data.
func mostFrequentName(containers []discovery.DiscoveredContainer) (string, int, error) {
	if len(containers) == 0 {
		return "", 0, ErrNoContainers
	}

	count := map[string]int{}
	for _, c := range containers {
		if _, exists := count[c.Name]; !exists {
			count[c.Name] = 0
		}
		count[c.Name] += 1
	}

	keys := maps.Keys(count)
	var mostFrequent string
	occurrences := 0
	for key := range keys {
		if count[key] > occurrences {
			mostFrequent = key
			occurrences = count[key]
		}
	}

	if mostFrequent == "" || occurrences == 0 {
		return "", 0, ErrCountingContainerNameOccurrences
	}

	return mostFrequent, occurrences, nil
}
