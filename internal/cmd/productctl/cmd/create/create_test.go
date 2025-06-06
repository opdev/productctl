package create_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/cli"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/create"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/testutils"
)

const testfixturesdir = "../testutils/testdata"

var fixtureDiscoveryJSON = filepath.Join(testfixturesdir, "fixture.minimal.discovery.json")

var _ = Describe("Create", func() {
	var (
		tempDirPath       string
		tempDiscoveryJSON string
	)

	BeforeEach(func() {
		var err error
		tempDirPath, err = os.MkdirTemp("", "productctl-unit-test-*")
		Expect(err).ToNot(HaveOccurred())

		tempfile, err := os.CreateTemp(tempDirPath, "*.create-test.discovery.json")
		Expect(err).ToNot(HaveOccurred())
		defer tempfile.Close()
		tempDiscoveryJSON = tempfile.Name()
		fixture, err := os.Open(fixtureDiscoveryJSON)
		Expect(err).ToNot(HaveOccurred())

		written, err := io.Copy(tempfile, fixture)
		Expect(err).ToNot(HaveOccurred())
		Expect(written).ToNot(BeZero())
	})

	AfterEach(func() {
		err := os.RemoveAll(tempDirPath)
		Expect(err).ToNot(HaveOccurred())
	})

	When("creating a new product", func() {
		When("no output file is provided", func() {
			It("should return an error", func() {
				_, err := testutils.ExecuteCommand(create.Command())
				Expect(err).To(HaveOccurred())
			})
		})
		When("the output location is provided", func() {
			var outputFile string
			BeforeEach(func() {
				tempfile, err := os.CreateTemp(tempDirPath, "*.create-test.outputfile.product.yaml")
				Expect(err).ToNot(HaveOccurred())
				outputFile = tempfile.Name()
				tempfile.Close()
			})
			It("should write to the provided location", func() {
				_, err := testutils.ExecuteCommand(create.Command(), outputFile)
				Expect(err).ToNot(HaveOccurred())
				stat, err := os.Stat(outputFile)
				Expect(err).ToNot(HaveOccurred())
				Expect(stat.Size).ToNot(BeZero())
			})

			When("a discover JSON is provided", func() {
				It("should include components in the discovery", func() {
					_, err := testutils.ExecuteCommand(create.Command(), fmt.Sprintf("--%s", cli.FlagIDFromDiscoveryJSON), tempDiscoveryJSON, outputFile)
					Expect(err).ToNot(HaveOccurred())
					stat, err := os.Stat(outputFile)
					Expect(err).ToNot(HaveOccurred())
					Expect(stat.Size).ToNot(BeZero())
				})
			})
		})
	})
})
