package discovery_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	libdiscovery "github.com/opdev/productctl/internal/discovery"

	"github.com/opdev/discover-workload/discovery"
)

var _ = Describe("Discovery", func() {
	When("converting workload discovery manifest to component resources", func() {
		It("should return an error when discovery manifest contains no entries", func() {
			_, err := libdiscovery.ComponentsFromDiscoveryManifest(discovery.Manifest{})
			Expect(err).To(MatchError(libdiscovery.ErrNoImagesDiscovered))
		})

		It("should return an error for duplicate component names", func() {
			manifestWithDuplicateNames := discovery.Manifest{
				DiscoveredImages: []discovery.DiscoveredImage{
					{ContainerName: "component1"},
					{ContainerName: "component1"},
				},
			}
			_, err := libdiscovery.ComponentsFromDiscoveryManifest(manifestWithDuplicateNames)
			Expect(err).To(MatchError(libdiscovery.ErrDuplicateComponentName))
		})

		It("should return components for valid manifest", func() {
			manifest := discovery.Manifest{
				DiscoveredImages: []discovery.DiscoveredImage{
					{ContainerName: "component1"},
					{ContainerName: "component2"},
				},
			}
			components, err := libdiscovery.ComponentsFromDiscoveryManifest(manifest)
			Expect(err).NotTo(HaveOccurred())
			Expect(components).To(HaveLen(len(manifest.DiscoveredImages)))
		})
	})
})
