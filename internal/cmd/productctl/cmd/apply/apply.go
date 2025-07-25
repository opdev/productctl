package apply

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
		Use:   "apply <your-declaration.yaml>",
		Short: "Apply changes to Partner product listings from the input file.",
		Long:  "Apply changes to partner product listings based on the provided configuration file",
		Args:  cobra.ExactArgs(1),
		RunE:  applyProductRunE,
	}

	cmd.Flags().Bool(cli.FlagIDCreateBackupOnOverwrite, false, "Create a backup of the declaration on overwrite. Note that this backups the on-disk declaration that was created/applied before overwriting it with new content.")

	return cmd
}

func applyProductRunE(cmd *cobra.Command, args []string) error {
	L := logger.FromContextOrDiscard(cmd.Context())
	cfg, err := cli.Config()
	if err != nil {
		return err
	}

	token, err := cfg.Token()
	if err != nil {
		return err
	}

	var endpoint string
	if cmd.Flags().Changed(cli.FlagIDCustomEndpoint) {
		endpoint, _ = cmd.Flags().GetString(cli.FlagIDCustomEndpoint)
		L.Debug("custom endpoint set, using it over env value", "endpoint", endpoint)
	} else {
		endpoint, err = cli.ResolveAPIEndpoint(cfg.Env)
		if err != nil {
			return err
		}
		L.Debug("endpoint resolved", "endpoint", endpoint)
	}

	if args[0] == "-" {
		return runApply(cmd.Context(), os.Stdin, os.Stdout, token, endpoint)
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
	return runApply(cmd.Context(), f, &updateFileOnSuccess, token, endpoint)
}

func runApply(ctx context.Context, in io.Reader, outOnCompletion io.Writer, token string, endpoint catalogapi.APIEndpoint) error {
	L := logger.FromContextOrDiscard(ctx)

	L.Info("reading in desired product listing")
	declaration, err := resource.ReadProductListing(in)
	if err != nil {
		return err
	}

	L.Debug("building graphql client")
	httpClient := catalogapi.TokenAuthenticatedHTTPClient(token, L.With("name", "httpclient"))
	client := graphql.NewClient(endpoint, httpClient)

	applied, err := catalogapi.ApplyProduct(ctx, client, declaration)
	if err != nil {
		return err
	}

	L.Info("Updating provided resource declaration.")
	b, err := yaml.Marshal(applied)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(outOnCompletion, string(b))
	if err != nil {
		return err
	}

	return nil
}
