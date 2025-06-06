package cmd_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/cmd/productctl/cmd"
)

var _ = Describe("Cmd", func() {
	When("building the root command", func() {
		// rootCmd has no business logic by itself.
		It("should contain subcommands", func() {
			cmd := cmd.RootCmd()
			Expect(cmd.HasAvailableSubCommands()).To(BeTrue())
		})
	})
})
