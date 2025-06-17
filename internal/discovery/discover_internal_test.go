package discovery

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	discoverworkload "github.com/opdev/discover-workload/discovery"
)

var _ = Describe("Discovery (internal)", func() {
	When("determining the most frequent container name in a container list", func() {
		var containers []discoverworkload.DiscoveredContainer

		When("there one name present in the list is used more than onee", func() {
			BeforeEach(func() {
				containers = []discoverworkload.DiscoveredContainer{
					{Name: "container-1"},
					{Name: "container-1"},
					{Name: "container-2"},
				}
			})

			It("should return the container name with the larger count", func() {
				actual, count, err := mostFrequentName(containers)
				Expect(err).ToNot(HaveOccurred())
				Expect(count).To(Equal(2))
				Expect(actual).To(Equal("container-1"))
			})
		})
		When("there are no containers", func() {
			BeforeEach(func() {
				containers = []discoverworkload.DiscoveredContainer{}
			})
			It("should return the expected error", func() {
				_, _, err := mostFrequentName(containers)
				Expect(err).To(MatchError(ErrNoContainers))
			})
		})

		When("the containers input has malformed data", func() {
			BeforeEach(func() {
				containers = []discoverworkload.DiscoveredContainer{
					{Name: ""},
				}
			})
			It("should return the expected error", func() {
				_, _, err := mostFrequentName(containers)
				Expect(err).To(MatchError(ErrCountingContainerNameOccurrences))
			})
		})
	})
})
