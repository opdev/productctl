package resource_test

import (
	"bytes"
	"io"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/resource"
)

var _ = Describe("IO", func() {
	When("provided a reader", func() {
		var fixture io.Reader
		When("it contains a product listing", func() {
			BeforeEach(func() {
				fixture = bytes.NewBufferString(`---
kind: ProductListing
spec:
  descriptions:
    long: This can contain long form content about your product.
    short: A brief synopsis
  name: Test Product 1
  type: container stack
with:
  components:
  - certification_status: Started
    container:
      distribution_method: rhcc
      hosted_registry: false
      os_content_type: Red Hat Universal Base Image (UBI)
      type: container
    name: test-component
    project_status: active
    type: Containers`)
			})

			It("should marshal correctly", func() {
				listing, err := resource.ReadProductListing(fixture)
				Expect(err).ToNot(HaveOccurred())
				Expect(listing.Kind).To(Equal("ProductListing"))
				Expect(listing.HasComponents()).To(BeTrue())
				Expect(listing.Spec.Name).To(Equal("Test Product 1"))
				Expect(listing.With.Components[0].Name).To(Equal("test-component"))
			})
		})

		When("it does not contain a product listing", func() {
			BeforeEach(func() {
				fixture = bytes.NewBufferString("package resource")
			})
			It("should not marshal correctly", func() {
				_, err := resource.ReadProductListing(fixture)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
