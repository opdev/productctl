package cleanup

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/opdev/productctl/internal/catalogapi"
	"github.com/opdev/productctl/internal/cli"
	"github.com/opdev/productctl/internal/file"
	"github.com/opdev/productctl/internal/logger"
	"github.com/opdev/productctl/internal/resource"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cleanup [my.product.yaml]",
		Short: "Detaches and archives components. Deletes the product listing. This is destructive. Use with caution.",
		Args:  cobra.MinimumNArgs(1), // The product declaration
		RunE:  runE,
	}

	cmd.Flags().Bool(cli.FlagIDCreateBackupOnOverwrite, false, "Create a backup of the declaration on overwrite. Note that this backups the on-disk declaration that was created/applied before overwriting it with new content.")

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	L := logger.FromContextOrDiscard(cmd.Context())
	_, token, err := cli.EnsureEnv()
	if err != nil {
		return err
	}

	var endpoint string
	if cmd.Flags().Changed(cli.FlagIDCustomEndpoint) {
		endpoint, _ = cmd.Flags().GetString(cli.FlagIDCustomEndpoint)
		L.Debug("custom endpoint set, using it over env value", "endpoint", endpoint)
	} else {
		env, _ := cmd.Flags().GetString(cli.FlagIDEndpoint)
		endpoint, err = cli.ResolveAPIEndpoint(env)
		if err != nil {
			return err
		}
		L.Debug("endpoint resolved", "endpoint", endpoint)
	}

	if args[0] == "-" {
		return runCleanup(cmd.Context(), os.Stdin, os.Stdout, token, endpoint)
	}

	// This is a read-only open.
	f, err := os.Open(args[0])
	if err != nil {
		return err
	}

	backupOnOverwrite, _ := cmd.Flags().GetBool(cli.FlagIDCreateBackupOnOverwrite)

	updateFileOnSuccess := file.LazyOverwriter{
		Filename:       args[0],
		DoBackup:       backupOnOverwrite,
		OptionalLogger: L.With("name", "fileIO"),
	}

	defer f.Close()
	return runCleanup(cmd.Context(), f, &updateFileOnSuccess, token, endpoint)
}

func runCleanup(ctx context.Context, in io.Reader, outOnCompletion io.Writer, token string, endpoint catalogapi.APIEndpoint) error {
	L := logger.FromContextOrDiscard(ctx)

	L.Info("reading in product listing")
	declaration, err := resource.ReadProductListing(in)
	if err != nil {
		return err
	}

	L.Debug("building graphql client")
	httpClient := catalogapi.TokenAuthenticatedHTTPClient(token, L.With("name", "httpclient"))
	client := graphql.NewClient(endpoint, httpClient)

	L.Debug("starting cleanup")
	cleaned, err := catalogapi.CleanupProduct(ctx, client, declaration)
	if err != nil {
		return err
	}

	L.Info("Updating provided resource declaration.")
	b, err := yaml.Marshal(cleaned)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(outOnCompletion, string(b))
	if err != nil {
		return err
	}

	return nil
}
