package fetch_test

import (
	"os"
	"syscall"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/cmd/productctl/cmd"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/testutils"
)

var _ = Describe("Fetch", func() {
	When("using the fetch command", func() {
		var tempDirPath string

		// tempDirPath scaffolded but unused until proper mocks are in place
		BeforeEach(func() {
			var err error
			tempDirPath, err = os.MkdirTemp("", "productctl-unit-test-*")
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			err := os.RemoveAll(tempDirPath)
			Expect(err).ToNot(HaveOccurred())
		})

		When("an existing product listing ID is provided", func() {
			var listingID string

			BeforeEach(func() {
				listingID = "123"
			})

			It("should fail if the minimum environment variables are not set", func() {
				output, err := testutils.ExecuteCommand(cmd.RootCmd(), "product", "fetch", listingID, "--custom-endpoint", "http://localhost:9630")
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

				It("should reach the apply phase, then fail", func() {
					// Endpoint is spoofed to avoid spamming actual endpoints with requests
					output, err := testutils.ExecuteCommand(cmd.RootCmd(), "product", "fetch", listingID, "--custom-endpoint", "http://localhost:9630")
					// We still expect an error here until business logic mocks have been implemented.
					Expect(err).To(HaveOccurred())
					Expect(output).To(ContainSubstring(syscall.ECONNREFUSED.Error()))
				})
			})
		})
	})
})
