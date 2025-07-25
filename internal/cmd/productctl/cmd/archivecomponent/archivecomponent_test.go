package archivecomponent_test

import (
	"fmt"
	"os"
	"syscall"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/cmd/productctl/cmd"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/testutils"
)

var _ = Describe("ArchiveComponent", func() {
	When("using the archive-component command", func() {
		When("the appropriate environment variables are in place", func() {
			BeforeEach(func() {
				os.Setenv("PRODUCTCTL_API_TOKEN", "foo")
			})

			AfterEach(func() {
				os.Setenv("PRODUCTCTL_API_TOKEN", "")
			})

			It("should reach the archive phase, then fail", func() {
				fmt.Println(os.Getenv("PRODUCTCTL_API_TOKEN"))
				// Endpoint is spoofed to avoid spamming actual endpoints with requests
				output, err := testutils.ExecuteCommand(cmd.RootCmd(), "util", "archive-component", "foo", "--custom-endpoint", "http://localhost:9630")
				// We still expect an error here until business logic mocks have been implemented.
				Expect(err).To(HaveOccurred())
				Expect(output).To(ContainSubstring(syscall.ECONNREFUSED.Error()))
			})
		})
	})
})
