package cli

import (
	"io/fs"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	spfviper "github.com/spf13/viper"
)

// tempDirFS implements fs.FS for a temporary directory
type tempDirFS struct {
	baseDir string
}

func (t tempDirFS) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(t.baseDir, name))
}

var _ = Describe("Config", func() {
	BeforeEach(func() {
		// Reset viper instance before each test
		reset()
	})

	When("loading configuration", func() {
		When("calling Config() function", func() {
			It("should return a valid UserConfig with defaults when no config file exists", func() {
				cfg, err := Config()
				Expect(err).ToNot(HaveOccurred())
				Expect(cfg).ToNot(BeNil())
				Expect(cfg.LogLevel).To(Equal(DefaultLogLevel))
				Expect(cfg.Env).To(Equal(DefaultEnv))
				Expect(cfg.SourceFile()).To(BeEmpty())
			})
		})

		When("calling renderedConfig() function", func() {
			It("should render valid configuration from viper instance", func() {
				v := spfviper.New()
				v.Set("log-level", "debug")
				v.Set("env", "stage")
				v.Set("api-token", "test-token")
				v.Set("api-token-file", "/path/to/token/file")
				v.SetConfigFile("/test/config.yaml")

				cfg, err := renderedConfig(v)
				Expect(err).ToNot(HaveOccurred())
				Expect(cfg).ToNot(BeNil())
				Expect(cfg.LogLevel).To(Equal("debug"))
				Expect(cfg.Env).To(Equal("stage"))
				Expect(cfg.APIToken).To(Equal("test-token"))
				Expect(cfg.APITokenFile).To(Equal("/path/to/token/file"))
				Expect(cfg.SourceFile()).To(Equal("/test/config.yaml"))
			})

			It("should handle empty viper configuration", func() {
				v := spfviper.New()
				// Don't set any values, just test with empty config

				cfg, err := renderedConfig(v)
				Expect(err).ToNot(HaveOccurred())
				Expect(cfg).ToNot(BeNil())
				Expect(cfg.LogLevel).To(BeEmpty())
				Expect(cfg.Env).To(BeEmpty())
				Expect(cfg.APIToken).To(BeEmpty())
				Expect(cfg.APITokenFile).To(BeEmpty())
				Expect(cfg.SourceFile()).To(BeEmpty())
			})

			It("should handle viper configuration with partial values", func() {
				v := spfviper.New()
				v.Set("log-level", "warn")
				v.Set("api-token", "partial-token")
				// Don't set env or api-token-file

				cfg, err := renderedConfig(v)
				Expect(err).ToNot(HaveOccurred())
				Expect(cfg).ToNot(BeNil())
				Expect(cfg.LogLevel).To(Equal("warn"))
				Expect(cfg.Env).To(BeEmpty())
				Expect(cfg.APIToken).To(Equal("partial-token"))
				Expect(cfg.APITokenFile).To(BeEmpty())
			})
		})

		When("calling UserConfig.Token() method", func() {
			When("a token file is configured", func() {
				var (
					tempDir      string
					tokenPath    string
					tokenContent string
				)

				BeforeEach(func() {
					// Create a temporary directory and token file
					tempDir = GinkgoT().TempDir()
					tokenContent = "file-token-content"
					tokenPath = filepath.Join(tempDir, "token.txt")
					err := os.WriteFile(tokenPath, []byte(tokenContent), 0o644)
					Expect(err).ToNot(HaveOccurred())
				})

				It("should read token from file when api-token-file is configured", func() {
					cfg := &UserConfig{
						APITokenFile: "token.txt", // Relative path for the tempDirFS
					}

					token, err := cfg.readTokenFile(tempDirFS{baseDir: tempDir}, cfg.APITokenFile)
					Expect(err).ToNot(HaveOccurred())
					Expect(token).To(Equal("file-token-content"))
				})

				It("should prioritize direct API token over token file", func() {
					cfg := &UserConfig{
						APIToken:     "direct-token",
						APITokenFile: "token.txt", // Relative path for the tempDirFS
					}

					token, err := cfg.readTokenFile(tempDirFS{baseDir: tempDir}, cfg.APITokenFile)
					Expect(err).ToNot(HaveOccurred())
					Expect(token).To(Equal("file-token-content"))

					// But the Token() method should prioritize the direct token
					token, err = cfg.Token()
					Expect(err).ToNot(HaveOccurred())
					Expect(token).To(Equal("direct-token"))
				})
			})

			It("should return API token when directly configured", func() {
				cfg := &UserConfig{
					APIToken: "direct-token",
				}
				token, err := cfg.Token()
				Expect(err).ToNot(HaveOccurred())
				Expect(token).To(Equal("direct-token"))
			})

			It("should return error when no token configuration is found", func() {
				cfg := &UserConfig{}
				token, err := cfg.Token()
				Expect(err).To(HaveOccurred())
				Expect(token).To(BeEmpty())
				Expect(err).To(MatchError("no API token configuration found in config file"))
			})

			It("should return error when token file does not exist", func() {
				cfg := &UserConfig{
					APITokenFile: "/nonexistent/path/to/token",
				}
				token, err := cfg.Token()
				Expect(err).To(HaveOccurred())
				Expect(token).To(BeEmpty())
			})
		})

		When("calling UserConfig.SourceFile() method", func() {
			It("should return the config file source path", func() {
				cfg := &UserConfig{
					configFileSource: "/path/to/config.yaml",
				}
				source := cfg.SourceFile()
				Expect(source).To(Equal("/path/to/config.yaml"))
			})

			It("should return empty string when no config file was used", func() {
				cfg := &UserConfig{}
				source := cfg.SourceFile()
				Expect(source).To(BeEmpty())
			})
		})
	})
})
