package sanitize_test

import (
	"io"
	"io/fs"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/cmd/productctl/cmd/sanitize"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/testutils"
)

const testProductFixture = "testdata/fixture.test.product.yaml"

var _ = Describe("Sanitize", func() {
	var (
		tempDirPath string
		tempProduct string
	)

	BeforeEach(func() {
		var err error
		tempDirPath, err = os.MkdirTemp("", "productctl-unit-test-*")
		Expect(err).ToNot(HaveOccurred())

		tempfile, err := os.CreateTemp(tempDirPath, "*.sanitize-test.product.yaml")
		Expect(err).ToNot(HaveOccurred())
		defer tempfile.Close()
		tempProduct = tempfile.Name()
		fixture, err := os.Open(testProductFixture)
		Expect(err).ToNot(HaveOccurred())

		written, err := io.Copy(tempfile, fixture)
		Expect(err).ToNot(HaveOccurred())
		Expect(written).ToNot(BeZero())
	})

	AfterEach(func() {
		err := os.RemoveAll(tempDirPath)
		Expect(err).ToNot(HaveOccurred())
	})

	When("sanitizing a resource declaration", func() {
		When("the file is not found", func() {
			It("should throw an error", func() {
				_, err := testutils.ExecuteCommand(sanitize.Command(), "foofile")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(fs.ErrNotExist))
			})
		})
		When("the file is not a product listing", func() {
			It("should throw an error", func() {
				// a file that exists but isn't a listing
				_, err := testutils.ExecuteCommand(sanitize.Command(), "sanitize.go")
				Expect(err).To(HaveOccurred())
			})
		})

		It("should not longer contain data pertinent to the original resource", func() {
			output, err := testutils.ExecuteCommand(sanitize.Command(), tempProduct)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(ContainSubstring("kind: ProductListing"))
			Expect(output).ToNot(ContainSubstring("org_id"))
			Expect(output).ToNot(ContainSubstring("_id"))
		})
	})
})
