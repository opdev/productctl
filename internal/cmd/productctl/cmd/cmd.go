package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/opdev/productctl/internal/cli"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/alpha"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/apply"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/lsp"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/newproduct"
	fetch "github.com/opdev/productctl/internal/cmd/productctl/cmd/populate"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/sanitize"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/version"
	libversion "github.com/opdev/productctl/internal/version"
)

// Execute runs the top-most command structure of the CLI.
func Execute() error {
	return rootCmd().ExecuteContext(context.Background())
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "productctl",
		Short:             "A CLI for managing Red Hat Partner Product Listings.",
		Long:              "A basic CLI useful for helping Red Hat Certification Partners define their Product Listings, and create and manage certification projects associated with those product listings.",
		Args:              cobra.MinimumNArgs(1),
		Version:           libversion.Version.String(),
		PersistentPreRunE: configureCLIPreRunE,
	}

	cmd.PersistentFlags().String(cli.FlagIDLogLevel, "warn", "The verbosity of the tool itself. Ex. error, warn, info, debug")
	cmd.PersistentFlags().String(cli.FlagIDEndpoint, "prod", "The catalog API environment to use. Choose from stage, prod")
	cmd.PersistentFlags().String(cli.FlagIDCustomEndpoint, "", "Define a custom API endpoint. Supersedes predefined environment values like \"prod\" if set")

	cmd.AddCommand(newproduct.Command())
	cmd.AddCommand(apply.Command())
	cmd.AddCommand(fetch.Command())
	cmd.AddCommand(sanitize.Command())
	cmd.AddCommand(version.Command())

	// build alpha commands.
	alpha := alpha.Command()
	alpha.AddCommand(lsp.Command())
	cmd.AddCommand(alpha)

	return cmd
}

var ErrConfiguringCLI = errors.New("failed to configure CLI")

func configureCLIPreRunE(cmd *cobra.Command, args []string) error {
	err := cmd.ParseFlags(args)
	if err != nil {
		return errors.Join(ErrConfiguringCLI, err)
	}

	loglevel, err := cmd.Flags().GetString("log-level")
	if err != nil {
		return errors.Join(ErrConfiguringCLI, err)
	}
	ctx, _, err := cli.ConfigureLogger(loglevel, os.Stderr)
	if err != nil {
		return errors.Join(ErrConfiguringCLI, err)
	}

	cmd.SetContext(ctx)

	return nil
}
