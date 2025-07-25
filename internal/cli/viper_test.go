package cli

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Viper tests", func() {
	Context("Lazy Loading Viper", func() {
		When("the viper instance hasn't been initialized", func() {
			v = nil
			It("should be initialized when calling for it", func() {
				_ = viper() // we don't care about this return value for this test.
				Expect(v).ToNot(BeNil())
			})
		})
	})

	Context("Getting the project-specific Viper instance", func() {
		When("Requesting the viper instance for the project", func() {
			It("Should return a non-empty viper instance", func() {
				packageV := viper()
				packageV.Set("foo", "bar")
				Expect(viper().Get("foo")).To(Equal("bar"))
			})
		})
	})

	When("Resetting the project-specific Viper instance", func() {
		It("should properly clear the instance", func() {
			packageV := viper()
			packageV.Set("foo", "bar")
			Expect(viper().Get("foo")).To(Equal("bar"))
			reset()
			Expect(viper().Get("foo")).To(BeNil())
		})
	})
})
