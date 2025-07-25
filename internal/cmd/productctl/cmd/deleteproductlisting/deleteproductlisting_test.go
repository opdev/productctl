package deleteproductlisting_test

import (
	"os"
	"syscall"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/cmd/productctl/cmd"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/testutils"
)

var _ = Describe("DeleteProductlisting", func() {
	When("using the delete-productlisting command", func() {
		It("should fail if the minimum environment variables are not set", func() {
			output, err := testutils.ExecuteCommand(cmd.RootCmd(), "util", "delete-productlisting", "foo", "--custom-endpoint", "http://localhost:9630")
			Expect(err).To(HaveOccurred())
			Expect(output).To(ContainSubstring(cmd.ErrMinOneAPITokenConfig.Error()))
		})

		When("the appropriate environment variables are in place", func() {
			BeforeEach(func() {
				os.Setenv("PRODUCTCTL_API_TOKEN", "foo")
			})

			AfterEach(func() {
				os.Setenv("PRODUCTCTL_API_TOKEN", "")
			})
			It("should reach the deletion phase, then fail", func() {
				// Endpoint is spoofed to avoid spamming actual endpoints with requests
				output, err := testutils.ExecuteCommand(cmd.RootCmd(), "util", "delete-productlisting", "foo", "--custom-endpoint", "http://localhost:9630")
				// We still expect an error here until business logic mocks have been implemented.
				Expect(err).To(HaveOccurred())
				Expect(output).To(ContainSubstring(syscall.ECONNREFUSED.Error()))
			})
		})
	})
})
