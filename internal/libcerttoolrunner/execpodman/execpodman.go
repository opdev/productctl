// Package podmanexec contains an proof-of-concept library call demonstrating
// how the certification Ansible Runner images might run using podman and Golang's
// os/exec library.
package execpodman

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"path/filepath"
)

const (
	// DefaultUserFilesDir is the default location where a user may place files.
	DefaultUserFilesDir = "./userfiles"
	// DefaultUserInventoryDir is the default location where inventory files
	// should live. Users are responsible for populating this directory with
	// inventory files.
	DefaultUserInventoryDir = "./inventory"
	// DefaultUserHostLogDir is the default location where certification
	// tooling logs will be written.
	DefaultUserHostLogDir = "./cert-logs"
)

const (
	containerPathInventoryDir = "/runner/inventory"
	containerPathUserfilesDir = "/runner/userfiles"
	containerPathCertLogsDir  = "/runner/cert-logs"
	containerPathEnvVarsFile  = "/runner/env/envvars"
)

var (
	ErrInvalidConfig              = errors.New("invalid configuration")
	ErrMissingRequiredConfigValue = errors.New("required config value is missing")
)

// Config contains configuration for running
// PodmanExecContainers.
type Config struct {
	// UserInventoryDir is the path on the host where the inventory file is
	// defined. This is mounted to the expected place in the runtime image.
	//
	// Required.
	UserInventoryDir string
	// UserfilesDir is the path on the host where the userfiles directory is
	// defined. This is mounted to the expected place in the runtime image.
	UserfilesDir string
	// UserHostLogDir is the path on the host where the image's logs should
	// be written. This is mounted to the expected place in the runtime image.
	UserHostLogDir string
	// EnvVars File is the path on the host where environment variables are
	// written. This is mounted to the expected place in the runtime image.
	EnvVarsFile string
}

// IsValid confirms that c contains minimum values to operate. Does not validate
// runtime such as "does c.RuntimeImage exist". Returns errors detailing what is
// lacking from the configurations in order to run.
func (c *Config) Validate() error {
	if c.UserInventoryDir == "" {
		return fmt.Errorf("%s: %s", ErrMissingRequiredConfigValue, "an inventory directory was not provided and is required")
	}

	if c.UserHostLogDir == "" {
		return fmt.Errorf("%s: %s", ErrMissingRequiredConfigValue, "an inventory directory was not provided and is required")
	}

	return nil
}

// DefaultConfig returns a configuration with library-specific defaults.
func DefaultConfig() *Config {
	return &Config{
		UserInventoryDir: DefaultUserInventoryDir,
		UserHostLogDir:   DefaultUserHostLogDir,
		UserfilesDir:     DefaultUserFilesDir,
	}
}

func Execute(
	_ context.Context,
	containerImage string,
	stdout, stderr io.Writer,
	logger *slog.Logger,
	cfg *Config,
) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("configuration invalid: %w", err)
	}

	args := []string{"run", "--rm", "--interactive", "--tty", "--net=host"}

	// NOTE: Podman will interpret a static string like "foo" as a volume
	// instead of a relative path. For now, we don't support volumes, so we'll
	// force an absolute path for these values to make sure Podman will
	// interpret the value as a path instead of a volume.
	userInventoryDir, err := filepath.Abs(cfg.UserInventoryDir)
	if err != nil {
		return err
	}
	args = append(args, "--volume", fmt.Sprintf("%s:%s:Z,ro", userInventoryDir, containerPathInventoryDir))

	if cfg.UserfilesDir != "" {
		logger.Debug("user provided userfiles directory, passing that through to runtime workload")
		userfilesDir, err := filepath.Abs(cfg.UserfilesDir)
		if err != nil {
			return err
		}
		args = append(args, "--volume", fmt.Sprintf("%s:%s:Z,ro", userfilesDir, containerPathUserfilesDir))
	}

	// TODO: If the user doesn't provide this, should we create one in /tmp or
	// otherwise? Should we delegate that to the caller?
	userHostLogDir, err := filepath.Abs(cfg.UserHostLogDir)
	if err != nil {
		return err
	}

	args = append(args, "--volume", fmt.Sprintf("%s:%s:Z", userHostLogDir, containerPathCertLogsDir))

	if cfg.EnvVarsFile != "" {
		logger.Debug("user provided environment variables, so passing that through to runtime workload")
		envVarsFile, err := filepath.Abs(cfg.EnvVarsFile)
		if err != nil {
			return err
		}

		args = append(args, "--volume", fmt.Sprintf("%s:%s:Z,ro", envVarsFile, containerPathEnvVarsFile))
	}

	args = append(args, containerImage)

	podman := exec.Command("podman", args...)
	podman.Stderr = stderr
	podman.Stdout = stdout

	err = podman.Run()
	if err != nil {
		logger.Error("error running container", "errMsg", err)
		return err
	}
	return nil
}
