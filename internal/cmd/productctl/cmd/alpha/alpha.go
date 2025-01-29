package alpha

import (
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alpha",
		Short: "Alpha commands that may be removed or modified at any point",
		Args:  cobra.MinimumNArgs(1),
	}

	return cmd
}
