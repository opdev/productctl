package cli

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	spfviper "github.com/spf13/viper"

	"github.com/opdev/productctl/internal/version"
)

var (
	ErrRenderingConfig = errors.New("failed to render configuration")
	ErrResolvingConfig = errors.New("failed to resolve configuration")
)

// RawConfig is an alias to the underlying viper instance coordinating this
// application's core config management logic.
var RawConfig = viper

func Config() (*UserConfig, error) {
	v := viper()
	initConfig(v)
	registerConfigDefaults(v)
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(spfviper.ConfigFileNotFoundError); !ok {
			return nil, errors.Join(ErrResolvingConfig, err)
		}
	}

	return renderedConfig(v)
}

func renderedConfig(v *spfviper.Viper) (*UserConfig, error) {
	var cfg UserConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, errors.Join(ErrRenderingConfig, err)
	}

	cfg.configFileSource = v.ConfigFileUsed()
	return &cfg, nil
}

// UserConfig is a strongly typed representation of configuration file data.
type UserConfig struct {
	APIToken     string `mapstructure:"api-token"`
	APITokenFile string `mapstructure:"api-token-file"`
	LogLevel     string `mapstructure:"log-level"`
	Env          string `mapstructure:"env"`

	configFileSource string
}

func (cfg *UserConfig) SourceFile() string {
	return cfg.configFileSource
}

func (cfg *UserConfig) Token() (string, error) {
	if cfg.APIToken != "" {
		return cfg.APIToken, nil
	}

	if cfg.APITokenFile != "" {
		return cfg.readTokenFile(os.DirFS("."))
	}

	return "", errors.New("no API token configuration found in config file")
}

func (cfg *UserConfig) readTokenFile(fs fs.FS) (string, error) {
	file, err := fs.Open(cfg.APITokenFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	token, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(token), nil
}

// initConfig initializes the CLI configuration instance, environment variables,
// handles file precedence, and baseline defaults.
func initConfig(v *spfviper.Viper) {
	v.SetEnvPrefix(version.Version.BaseName)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	// These environment variables don't have flags associated with them.
	// Bind them so either the config or the environment can be used.
	_ = v.BindEnv("api-token")
	_ = v.BindEnv("api-token-file")
	v.AutomaticEnv()

	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// e.g. $PWD/.productctl/config.yaml, highest precedence
	v.AddConfigPath(filepath.Join(".", fmt.Sprintf(".%s", version.Version.BaseName)))

	// e.g. ~/.config/productctl/config.yaml, second highest precedence
	if userConfigDir, err := os.UserConfigDir(); err == nil {
		v.AddConfigPath(filepath.Join(userConfigDir, version.Version.BaseName))
	}

	// e.g. ~/.productctl/config.yaml, lowest precedence, fallback to allow home
	// directory config dir. just in case a system does not have a user config
	// dir.
	if userHomeDir, err := os.UserConfigDir(); err == nil {
		v.AddConfigPath(filepath.Join(userHomeDir, fmt.Sprintf(".%s", version.Version.BaseName)))
	}
}

func registerConfigDefaults(v *spfviper.Viper) {
	v.SetDefault(FlagIDLogLevel, DefaultLogLevel)
	v.SetDefault(FlagIDEnv, DefaultEnv)
}
