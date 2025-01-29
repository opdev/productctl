package version

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/opdev/productctl/internal/cli"
	"github.com/opdev/productctl/internal/version"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints the version information",
		Run: func(cmd *cobra.Command, _ []string) {
			out := version.Version.String()

			// If the user requested JSON, try and produce the JSON blob,
			// setting that as the output text if successful.
			if asJSON, err := cmd.Flags().GetBool(cli.FlagIDVersionAsJSON); err == nil && asJSON {
				if b, err := json.Marshal(version.Version); err == nil {
					out = string(b)
				}
			}

			fmt.Print(out)
		},
	}

	cmd.Flags().Bool(cli.FlagIDVersionAsJSON, false, "prints the version info as a JSON blob")

	return cmd
}
