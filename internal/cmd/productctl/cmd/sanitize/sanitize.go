package sanitize

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/opdev/productctl/internal/resource"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sanitize-product",
		Short: "Cleans declaration for re-use and emits to stdout",
		Long:  "Strips data from the product declaration on disk that associates a product with an entry in the backend. Does not impact the backend, or overwrite the input file.",
		Args:  cobra.ExactArgs(1),
		RunE:  sanitizeProductCmdRunE,
	}

	return cmd
}

func sanitizeProductCmdRunE(_ *cobra.Command, args []string) error {
	f, err := os.Open(args[0])
	if err != nil {
		return err
	}

	r, err := resource.ReadProductListing(f)
	if err != nil {
		return err
	}

	r.Sanitize()

	b, err := yaml.Marshal(r)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(os.Stdout, string(b))
	if err != nil {
		return err
	}

	return nil
}
