package bridge_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/cmd/productctl/cmd/bridge"
)

var _ = Describe("Bridge", func() {
	When("using the bridge command", func() {
		When("the caller provides a name and description", func() {
			var name, desc string

			BeforeEach(func() {
				name = "name"
				desc = "desc"
			})
			It("should produce a command with the given name and description", func() {
				cmd := bridge.Command(name, desc)
				Expect(cmd.Name()).To((Equal(name)))
				Expect(cmd.Short).To(Equal(desc))
			})
		})
	})
})
