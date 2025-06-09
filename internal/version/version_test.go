package version

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version", func() {
	When("using the version instance", func() {
		It("should include the hard-coded values for name and project", func() {
			Expect(Version.BaseName).To(Equal(baseName))
			Expect(Version.Name).To(Equal(projectName))
			Expect(Version.Version).To(Equal(version))
			Expect(Version.Commit).To(Equal(commit))
		})

		It("should stringify including the version and commit information", func() {
			Version.Version = "test-version"
			Version.Commit = "test-commit"
			actual := Version.String()
			Expect(actual).To(ContainSubstring(Version.Version))
			Expect(actual).To(ContainSubstring(Version.Commit))
			Expect(actual).ToNot(BeEmpty())
		})
	})
})
