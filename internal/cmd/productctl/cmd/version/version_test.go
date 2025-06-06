package version_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/opdev/productctl/internal/cli"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/testutils"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/version"
	libversion "github.com/opdev/productctl/internal/version"
)

var _ = Describe("Version", func() {
	var cmd *cobra.Command
	BeforeEach(func() {
		cmd = version.Command()
	})
	When("using the version command", func() {
		It("should print the version of the tooling", func() {
			output, err := testutils.ExecuteCommand(cmd)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(ContainSubstring(libversion.Version.Version))
			Expect(output).To(ContainSubstring(libversion.Version.Commit))
		})

		When("requesting the version as a JSON object", func() {
			It("should produce a JSON object", func() {
				v := libversion.Version
				output, err := testutils.ExecuteCommand(
					cmd,
					fmt.Sprintf("--%s", cli.FlagIDVersionAsJSON),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(output).To(ContainSubstring(v.Version))
				Expect(output).To(ContainSubstring(libversion.Version.Commit))
			})
		})
	})
})
