// Package fetch implements the fetch-product subcommand.
package fetch

import (
	"fmt"
	"os"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/opdev/productctl/internal/catalogapi"
	"github.com/opdev/productctl/internal/cli"
	"github.com/opdev/productctl/internal/logger"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fetch <productID>",
		Short: "Get a pre-existing product listing",
		Long: `Get data about a pre-existing product listing by its ID and generate its declaration for storage on disk.

Only components attached to this product listing with an "active" status will be returned.

This command does not overwrite an existing file, and relies in output redirection to store the contents to disk at any location you would prefer.
`,
		Args: cobra.MinimumNArgs(1),
		RunE: getProductListingRunE,
	}

	return cmd
}

func getProductListingRunE(cmd *cobra.Command, args []string) error {
	L := logger.FromContextOrDiscard(cmd.Context())
	_, token, err := cli.EnsureEnv()
	if err != nil {
		return err
	}

	productID := args[0]

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

	httpClient := catalogapi.TokenAuthenticatedHTTPClient(token, L.With("name", "httpclient"))
	client := graphql.NewClient(endpoint, httpClient)

	newListing, err := catalogapi.PopulateProduct(cmd.Context(), client, productID)
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(newListing)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(os.Stdout, string(b))
	if err != nil {
		return err
	}

	return nil
}
