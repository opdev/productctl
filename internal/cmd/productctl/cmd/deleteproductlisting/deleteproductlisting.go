package deleteproductlisting

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/opdev/productctl/internal/catalogapi"
	"github.com/opdev/productctl/internal/cli"
	"github.com/opdev/productctl/internal/genpyxis"
	"github.com/opdev/productctl/internal/logger"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-productlisting <productID>",
		Short: "Deletes the product listingwith the specified product ID.",
		Long: `Deletes the product listing with the specified product ID

This should be considered a destructive operation. Note that there are various reasons why the API may reject this operation. Those reasons may need to be handled directly via the Partner Connect UI.
`,
		Args: cobra.MinimumNArgs(1), // The product ID
		RunE: runE,
	}

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
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

	return run(cmd.Context(), args[0], token, endpoint)
}

func run(ctx context.Context, listingID string, token string, endpoint catalogapi.APIEndpoint) error {
	L := logger.FromContextOrDiscard(ctx)
	L.Info("deleting product listing", "_id", listingID)

	L.Debug("building graphql client")
	httpClient := catalogapi.TokenAuthenticatedHTTPClient(token, L.With("name", "httpclient"))
	client := graphql.NewClient(endpoint, httpClient)

	resp, err := genpyxis.DeleteProduct(ctx, client, listingID)
	if err != nil {
		return err
	}

	if gqlErr := resp.Update_product_listing.GetError(); gqlErr != nil {
		return catalogapi.ParseGraphQLResponseError(gqlErr)
	}

	L.Info("done")
	return nil
}
