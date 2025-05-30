package certifycontainers

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"dario.cat/mergo"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/opdev/productctl/internal/ansible"
	"github.com/opdev/productctl/internal/cli"
	"github.com/opdev/productctl/internal/libcerttoolrunner"
	"github.com/opdev/productctl/internal/libcerttoolrunner/execpodman"
	"github.com/opdev/productctl/internal/logger"
	"github.com/opdev/productctl/internal/resource"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "containers <product-declaration.yaml> <generated-mappings.yaml>",
		Short: "Run Container Certification for a given Product Listing",
		Args:  cobra.ExactArgs(2),
		RunE:  runE,
	}

	flags := cmd.Flags()
	flags.String(cli.FlagIDUserfilesDir, "", "A full path to a user files")
	flags.String(cli.FlagIDLogsDir, "", "A full path to the logs directory")
	flags.String(cli.FlagIDCatalogAPIToken, "", "The catalog API token to use for submit requests")
	flags.String(cli.FlagIDRuntimeImage, libcerttoolrunner.DefaultImageCertifyContainers, "The container image to use to certify containers")

	// Debug Flags
	flags.Bool(cli.FlagIDKeepTempDir, false, "keep the temporary directory where generated assets are stored")

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	L := logger.FromContextOrDiscard(cmd.Context())

	L.Debug("Reading product from declaration", "file", args[0])
	productFile, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer productFile.Close()

	L.Debug("Reading component certification mappings from provided mapping", "file", args[1])
	mappingFile, err := os.Open(args[1])
	if err != nil {
		return err
	}

	defer mappingFile.Close()

	declaration, err := resource.ReadProductListing(productFile)
	if err != nil {
		return err
	}

	mappingB, err := io.ReadAll(mappingFile)
	if err != nil {
		return err
	}

	var mapping ansible.MappingDeclaration
	err = yaml.Unmarshal(mappingB, &mapping)
	if err != nil {
		return err
	}

	L.Debug("Generating inventory from product and mappings")
	inventory, err := ansible.GenerateInventory(
		cmd.Context(),
		declaration,
		mapping,
	)
	if err != nil {
		return err
	}

	runBaseDir, err := os.MkdirTemp(os.TempDir(), "cert-automation-")
	if err != nil {
		return err
	}

	keep, _ := cmd.Flags().GetBool(cli.FlagIDKeepTempDir)
	if !keep {
		defer func() {
			err = os.RemoveAll(runBaseDir)
			if err != nil {
				L.Error("unable to clean up temporary directory", "errorMsg", err, "tempdir", runBaseDir)
			}
		}()
	} else {
		defer L.Debug("tempdir kept per flag", "path", runBaseDir)
	}

	inventoryDir, err := os.MkdirTemp(runBaseDir, "inventory-")
	if err != nil {
		return err
	}

	inventoryData, err := yaml.Marshal(inventory)
	if err != nil {
		return err
	}

	L.Debug("Writing generated inventory to temporary directory", "tmpdir", runBaseDir)
	err = os.WriteFile(filepath.Join(inventoryDir, "generated.product.inventory.yaml"), inventoryData, 0o600)
	if err != nil {
		return err
	}

	envVarFileName := ""

	if token, err := cmd.Flags().GetString(cli.FlagIDCatalogAPIToken); err == nil && token != "" {
		L.Debug("Writing generated envvar file to temporary directory", "tmpdir", runBaseDir)
		envVarFile, err := os.CreateTemp(runBaseDir, "envvar-")
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(envVarFile, "---\nPFLT_PYXIS_API_TOKEN: %s", token)
		if err != nil {
			return err
		}

		envVarFileName = envVarFile.Name()
		envVarFile.Close()
	}

	userfilesDir, _ := cmd.Flags().GetString(cli.FlagIDUserfilesDir)
	logsDir, _ := cmd.Flags().GetString(cli.FlagIDLogsDir)
	runtimeImage, _ := cmd.Flags().GetString(cli.FlagIDRuntimeImage)

	execContainerConfig := execpodman.DefaultConfig()
	err = mergo.Merge(execContainerConfig, &execpodman.Config{
		UserfilesDir:     userfilesDir,
		UserHostLogDir:   logsDir,
		UserInventoryDir: inventoryDir,
		EnvVarsFile:      envVarFileName,
	}, mergo.WithOverride)
	if err != nil {
		L.Error("error merging user configuration over default", "errMsg", err)
		return err
	}

	err = execpodman.Execute(context.TODO(), runtimeImage, os.Stdout, os.Stderr, L, execContainerConfig)
	if err != nil {
		L.Error("error running certification workload", "errMsg", err)
		return err
	}

	return nil
}
