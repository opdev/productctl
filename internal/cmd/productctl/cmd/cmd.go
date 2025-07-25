package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/opdev/productctl/internal/cli"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/apply"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/archivecomponent"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/bridge"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/cleanup"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/create"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/deleteproductlisting"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/fetch"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/jsonschema"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/sanitize"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/version"
	libversion "github.com/opdev/productctl/internal/version"
)

// Execute runs the top-most command structure of the CLI.
func Execute() error {
	return RootCmd().ExecuteContext(context.Background())
}

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "productctl",
		Short:   "A CLI for managing Red Hat Partner Product Listings.",
		Long:    "A basic CLI useful for helping Red Hat Certification Partners define their Product Listings, and create and manage certification projects associated with those product listings.",
		Version: libversion.Version.String(),
	}

	// NOTE(komish): commonFlags is a hack to get flags that are reused across
	// various subcommand trees stored into the viper configuration.
	//
	// pflag allows you to define pflag.Flag, but they take a pflag.Value
	// implementation. We're just using strings (and similar base types). The
	// pflag lib has pflag.Value implementations for these but doesn't expose
	// them.
	//
	// https://github.com/spf13/pflag/issues/334
	//
	// In effect, we work around that by defining it once here in the
	// commonFlags flagset, then extracting it and re-using it.
	commonFlags := pflag.NewFlagSet("common", pflag.ContinueOnError)
	commonFlags.String(cli.FlagIDEnv, cli.DefaultEnv, "The catalog API environment to use. Choose from stage, prod")
	commonFlags.String(cli.FlagIDCustomEndpoint, "", "Define a custom API endpoint. Supersedes predefined environment values like \"prod\" if set")
	envFlag := commonFlags.Lookup(cli.FlagIDEnv)
	customEndpointFlag := commonFlags.Lookup(cli.FlagIDCustomEndpoint)

	cmd.AddCommand(version.Command())
	cmd.PersistentFlags().String(cli.FlagIDLogLevel, cli.DefaultLogLevel, "The verbosity of the tool itself. Ex. error, warn, info, debug")
	util := bridge.Command("util", "Utilities for the management of your Partner Connect account")
	util.PersistentFlags().AddFlag(envFlag)
	util.PersistentFlags().AddFlag(customEndpointFlag)
	util.AddCommand(archivecomponent.Command())
	util.AddCommand(deleteproductlisting.Command())
	cmd.AddCommand(util)

	// Build the product management command tree.
	product := bridge.Command("product", "Manage your Product Listing")
	product.PersistentFlags().AddFlag(envFlag)
	product.PersistentFlags().AddFlag(customEndpointFlag)
	product.AddCommand(create.Command())
	product.AddCommand(apply.Command())
	product.AddCommand(fetch.Command())
	product.AddCommand(sanitize.Command())
	product.AddCommand(cleanup.Command())
	product.AddCommand(jsonschema.Command())
	cmd.AddCommand(product)

	// These commands and their subcommands require an API token to be
	// configured. cobra.MatchAll is a bit of a misnomer, intended for use with
	// cobra.PositionalArgs. In effect, we use it to chain together multiple
	// pre-run functions.
	product.PersistentPreRunE = cobra.MatchAll(
		configureCLIPreRunE,
		ensureAtLeastOneTokenConfigured,
	)
	util.PersistentPreRunE = cobra.MatchAll(
		configureCLIPreRunE,
		ensureAtLeastOneTokenConfigured,
	)

	// Bind flags to configuration
	rawC := cli.RawConfig()
	_ = rawC.BindPFlag(cli.FlagIDLogLevel, cmd.PersistentFlags().Lookup(cli.FlagIDLogLevel))
	_ = rawC.BindPFlag(cli.FlagIDEnv, commonFlags.Lookup(cli.FlagIDEnv))

	return cmd
}

var ErrConfiguringCLI = errors.New("failed to configure CLI")

func configureCLIPreRunE(cmd *cobra.Command, args []string) error {
	cfg, err := cli.Config()
	if err != nil {
		return errors.Join(ErrConfiguringCLI, err)
	}

	ctx, L, err := cli.ConfigureLogger(cfg.LogLevel, os.Stderr)
	if err != nil {
		return errors.Join(ErrConfiguringCLI, err)
	}

	if cfg.SourceFile() != "" {
		L.Info("using config file", "file", cfg.SourceFile())
	}

	cmd.SetContext(ctx)
	return nil
}

var ErrMinOneAPITokenConfig = errors.New("either api-token or api-token-file must be configured in your config file or environment")

func ensureAtLeastOneTokenConfigured(_ *cobra.Command, _ []string) error {
	cfg, err := cli.Config()
	if err != nil {
		return errors.Join(ErrConfiguringCLI, err)
	}
	if cfg.APIToken == "" && cfg.APITokenFile == "" {
		return errors.Join(ErrMinOneAPITokenConfig)
	}

	return nil
}
