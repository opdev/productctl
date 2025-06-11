package jsonschema_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/cmd/productctl/cmd/jsonschema"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/testutils"
	"github.com/opdev/productctl/internal/resource"
)

var _ = Describe("jsonschema", func() {
	When("generating resource schemas", func() {
		var (
			cmdOut string
			cmdErr error
			args   []string
		)

		BeforeEach(func() {
			cmdOut, cmdErr = testutils.ExecuteCommand(jsonschema.Command(), args...)
		})

		It("should not be empty", func() {
			Expect(cmdErr).ToNot(HaveOccurred())
			Expect(cmdOut).ToNot(BeEmpty())
		})

		It("should be proper JSON", func() {
			Expect(cmdErr).ToNot(HaveOccurred())
			x := map[string]any{}
			err := json.Unmarshal([]byte(cmdOut), &x)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should contain relevant information", func() {
			Expect(cmdErr).ToNot(HaveOccurred())
			// A random smattering of values from defined resources.
			Expect(cmdOut).To(ContainSubstring("os_content_type"))
			Expect(cmdOut).To(ContainSubstring("distribution_method"))
			Expect(cmdOut).To(ContainSubstring("github_usernames"))
			Expect(cmdOut).To(ContainSubstring(resource.ContentTypeUBI))
		})
	})
})
