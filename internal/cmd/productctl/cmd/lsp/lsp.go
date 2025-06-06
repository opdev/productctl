package lsp

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/spf13/cobra"

	"github.com/opdev/productctl/internal/resource"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lsp-completion",
		Short: "Generate resource schema for LSPs that support it.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			// TODO: If this functionality will remain and be useful, we should
			// add comments, enums, and warnings when things are immutable after
			// being applied to the resource declaration itself.
			schema := jsonschema.Reflect(&resource.ProductListingDeclaration{})
			b, err := json.MarshalIndent(schema, "", "  ")
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), string(b))
			return nil
		},
	}

	return cmd
}
