package resource_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/genpyxis"
	"github.com/opdev/productctl/internal/resource"
)

var _ = Describe("Productlisting", func() {
	DescribeTable(
		"when using boolean helper methods for product listing state",
		func(actual, expected bool) {
			Expect(actual).To(Equal(expected))
		},
		Entry("name is present", (&resource.ProductListing{Name: "foo"}).HasName(), true),
		Entry("name is missing", (&resource.ProductListing{}).HasName(), false),
		Entry("ID is present", (&resource.ProductListing{ID: "123"}).HasID(), true),
		Entry("ID is missing", (&resource.ProductListing{}).HasID(), false),
	)
	When("converting a ContactInfoProvider into a corresponding product listing contact list", func() {
		var provider *genpyxis.ProductListingSupportedFieldsContactsContactsItems
		BeforeEach(func() {
			provider = &genpyxis.ProductListingSupportedFieldsContactsContactsItems{
				Email_address: "user@example.com",
				Type:          "testcontact",
			}
		})
		It("should map the Email Address and Type to the right places", func() {
			contacts := resource.ContactsFrom([]*genpyxis.ProductListingSupportedFieldsContactsContactsItems{provider})
			Expect(contacts).To(HaveLen(1))
			Expect(contacts[0].EmailAddress).To(Equal(provider.GetEmail_address()))
			Expect(contacts[0].Type).To(Equal(provider.GetType()))
		})
	})
})
