package apply_test

import (
	"io"
	"os"
	"path/filepath"
	"syscall"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/cli"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/testutils"
)

const testfixturesdir = "../testutils/testdata"

var fixtureMinimalProduct = filepath.Join(testfixturesdir, "fixture.minimal.product.yaml")

var _ = Describe("Apply", func() {
	When("using the apply command", func() {
		var (
			tempDirPath string
			tempProduct string
		)

		BeforeEach(func() {
			var err error
			tempDirPath, err = os.MkdirTemp("", "productctl-unit-test-*")
			Expect(err).ToNot(HaveOccurred())

			tempfile, err := os.CreateTemp(tempDirPath, "*.apply-test.product.yaml")
			Expect(err).ToNot(HaveOccurred())
			defer tempfile.Close()
			tempProduct = tempfile.Name()
			fixture, err := os.Open(fixtureMinimalProduct)
			Expect(err).ToNot(HaveOccurred())

			written, err := io.Copy(tempfile, fixture)
			Expect(err).ToNot(HaveOccurred())
			Expect(written).ToNot(BeZero())
		})

		AfterEach(func() {
			err := os.RemoveAll(tempDirPath)
			Expect(err).ToNot(HaveOccurred())
		})

		When("a declaration file is provided", func() {
			var file string
			BeforeEach(func() {
				file = tempProduct
			})

			It("should fail if the minimum environment variables are not set", func() {
				output, err := testutils.ExecuteCommand(cmd.RootCmd(), "product", "apply", file, "--custom-endpoint", "http://localhost:9630")
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

				It("should reach the apply phase, then fail", func() {
					// Endpoint is spoofed to avoid spamming actual endpoints with requests
					output, err := testutils.ExecuteCommand(cmd.RootCmd(), "product", "apply", file, "--custom-endpoint", "http://localhost:9630")
					// We still expect an error here until business logic mocks have been implemented.
					Expect(err).To(HaveOccurred())
					Expect(output).To(ContainSubstring(syscall.ECONNREFUSED.Error()))
				})
			})
		})
	})
})
