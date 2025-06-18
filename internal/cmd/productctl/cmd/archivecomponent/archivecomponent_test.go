package archivecomponent_test

import (
	"os"
	"syscall"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/cli"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/testutils"
)

var _ = Describe("ArchiveComponent", func() {
	When("using the archive-component command", func() {
		It("should fail if the minimum environment variables are not set", func() {
			output, err := testutils.ExecuteCommand(cmd.RootCmd(), "util", "archive-component", "foo", "--custom-endpoint", "http://localhost:9630")
			Expect(err).To(HaveOccurred())
			Expect(output).To(ContainSubstring(cli.ErrEnvVarMissing.Error()))
		})

		When("the appropriate environment variables are in place", func() {
			BeforeEach(func() {
				os.Setenv(cli.EnvAPIToken, "foo")
				os.Setenv(cli.EnvOrgID, "123")
			})

			AfterEach(func() {
				os.Setenv(cli.EnvAPIToken, "")
				os.Setenv(cli.EnvOrgID, "")
			})

			It("should reach the archive phase, then fail", func() {
				// Endpoint is spoofed to avoid spamming actual endpoints with requests
				output, err := testutils.ExecuteCommand(cmd.RootCmd(), "util", "archive-component", "foo", "--custom-endpoint", "http://localhost:9630")
				// We still expect an error here until business logic mocks have been implemented.
				Expect(err).To(HaveOccurred())
				Expect(output).To(ContainSubstring(syscall.ECONNREFUSED.Error()))
			})
		})
	})
})
