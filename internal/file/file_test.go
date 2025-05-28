package file

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("File", func() {
	When("using the LazyOverwriter", func() {
		var (
			workingDirectory string
			targetFileName   string
		)

		BeforeEach(func() {
			var err error
			workingDirectory, err = os.MkdirTemp("", "workingdir-file_test.go-*")
			Expect(err).ToNot(HaveOccurred())

			targetFile, err := os.CreateTemp(workingDirectory, "target-file-file_test.go-*")
			Expect(err).ToNot(HaveOccurred())
			targetFileName = targetFile.Name()
			targetFile.Close()
		})

		AfterEach(func() {
			os.RemoveAll(workingDirectory)
		})

		When("the file contains data", func() {
			var originalFileContent string

			BeforeEach(func() {
				originalFileContent = "olddata"
				f, err := os.Create(targetFileName)
				Expect(err).ToNot(HaveOccurred())
				defer f.Close()
				written, err := f.Write([]byte(originalFileContent))
				Expect(err).ToNot(HaveOccurred())
				Expect(written).ToNot(BeZero())
			})

			It("should not truncate the data at instantiation of the LazyOverwriter", func() {
				_ = &LazyOverwriter{Filename: targetFileName}

				content, err := os.ReadFile(targetFileName)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(content)).ToNot(BeZero())
				Expect(string(content)).To(Equal(originalFileContent))
			})

			It("should truncate the data when Write is called", func() {
				lw := &LazyOverwriter{Filename: targetFileName}

				newFileContent := "newdata"
				written, err := lw.Write([]byte(newFileContent))
				Expect(err).ToNot(HaveOccurred())
				Expect(written).ToNot(BeZero())

				content, err := os.ReadFile(targetFileName)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(content)).ToNot(BeZero())
				Expect(string(content)).To(Equal(newFileContent))
			})

			When("the backup flag is enabled", func() {
				When("writing new content", func() {
					It("should produce a backup before overwriting", func() {
						newNameFn := func(s string) string {
							return fmt.Sprintf("GENBACKUP-%s", s)
						}

						lw := &LazyOverwriter{
							Filename:            targetFileName,
							DoBackup:            true,
							BackupFilenameGenFn: newNameFn,
						}

						written, err := lw.Write([]byte("newdata"))
						Expect(err).ToNot(HaveOccurred())
						Expect(written).ToNot(BeZero())

						expectedBackupFilePath := filepath.Join(workingDirectory, newNameFn(filepath.Base(targetFileName)))
						Expect(expectedBackupFilePath).To(BeAnExistingFile())

						content, err := os.ReadFile(expectedBackupFilePath)
						Expect(err).ToNot(HaveOccurred())
						Expect(string(content)).To(Equal(originalFileContent))
					})
				})
			})
		})
	})
})
