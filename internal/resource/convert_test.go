package resource_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/resource"
)

var _ = Describe("Convert", func() {
	type FooOriginal struct {
		IntKey    int    `json:"int_key"`
		StringKey string `json:"string_key"`
	}

	type FooLike struct {
		MyInt    int    `json:"int_key"`
		MyString string `json:"string_key"`
	}

	When("converting between types through the JSONConvert roundtrip", func() {
		When("like struct markers are used", func() {
			var (
				original FooOriginal
				similar  FooLike
			)

			BeforeEach(func() {
				original = FooOriginal{
					IntKey:    10,
					StringKey: "ten",
				}

				similar = FooLike{}
			})

			It("should round trip successfully", func() {
				var err error
				similar, err = resource.JSONConvert[FooLike](original)
				Expect(err).ToNot(HaveOccurred())
				Expect(similar.MyInt).ToNot(BeZero())
				Expect(original.IntKey).To(Equal(similar.MyInt))
				Expect(similar.MyString).ToNot(BeEmpty())
				Expect(original.StringKey).To(Equal(similar.MyString))
			})
		})

		When("unlike types are used", func() {
			var in string
			BeforeEach(func() {
				in = "foo"
			})
			It("should throw an error", func() {
				var err error
				_, err = resource.JSONConvert[struct{}](in)
				Expect(err).To(HaveOccurred())
			})
		})
		When("the input type cannot be marshaled", func() {
			var in func()
			BeforeEach(func() {
				in = func() {}
			})
			It("should throw an error", func() {
				_, err := resource.JSONConvert[string](in)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
