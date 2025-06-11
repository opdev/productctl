package resource_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/resource"
)

var _ = Describe("Resource", func() {
	When("working with product listing declarations", func() {
		var declaration *resource.ProductListingDeclaration

		When("sanitizing the declaration", func() {
			BeforeEach(func() {
				new := resource.NewProductListing()
				new.Spec = resource.ProductListing{
					ID:             "abc123",
					Name:           "test-fixture",
					OrgID:          1234,
					LastUpdateDate: &time.Time{},
					Type:           resource.ProductListingTypeContainerStack,
					CreationDate:   &time.Time{},
				}
				new.With = resource.Inclusions{
					Components: []*resource.Component{
						{
							ID:                "c123",
							CertificationDate: &time.Time{},
							Name:              "cname",
							OrgID:             1234,
							ProjectStatus:     "unset",
							Type:              resource.ComponentTypeContainer,
							CreationDate:      &time.Time{},
							LastUpdateDate:    &time.Time{},
						},
					},
				}

				declaration = &new
			})

			JustBeforeEach(func() {
				declaration.Sanitize()
			})
			It("should unset created, update, etc. dates", func() {
				Expect(declaration.Spec.LastUpdateDate).To(BeNil())
				Expect(declaration.Spec.CreationDate).To(BeNil())
				for _, c := range declaration.With.Components {
					Expect(c.LastUpdateDate).To(BeNil())
					Expect(c.CertificationDate).To(BeNil())
				}
			})
			It("should unset the top-level ID and all component IDs", func() {
				Expect(declaration.Spec.ID).To(BeEmpty())
				Expect(declaration.Spec.OrgID).To(BeZero())
				for _, c := range declaration.With.Components {
					Expect(c.ID).To(BeEmpty())
					Expect(c.OrgID).To(BeZero())
				}
			})
			It("should reset the component statuses", func() {
				for _, c := range declaration.With.Components {
					Expect(c.ProjectStatus).To(Equal(resource.ProjectStatusActive))
				}
			})
			It("should leave component information intact", func() {
				Expect(declaration.Spec.Name).To(Equal("test-fixture"))
				Expect(declaration.Spec.Type).To(Equal(resource.ProductListingTypeContainerStack))
				for _, c := range declaration.With.Components {
					Expect(c.Name).To(Equal("cname"))
				}
			})
		})

		When("the declaration has components", func() {
			BeforeEach(func() {
				declaration.With = resource.Inclusions{Components: []*resource.Component{{ID: "id"}}}
			})
			It("should return a corresponding boolean from the helper function", func() {
				Expect(declaration.HasComponents()).To(BeTrue())
			})
		})
	})
})
