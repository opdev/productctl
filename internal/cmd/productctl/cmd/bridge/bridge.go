package bridge

import "github.com/spf13/cobra"

// Command provides a cobra command with the provided inputs. This Command is
// not expected to have any functionality, other than wrapping other
// subcommands.
func Command(name, shortDesc string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: shortDesc,
	}

	return cmd
}
