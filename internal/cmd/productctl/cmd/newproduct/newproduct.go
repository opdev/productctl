package newproduct

import (
	"encoding/json"
	"os"

	"github.com/opdev/discover-workload/discovery"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	libdiscovery "github.com/opdev/productctl/internal/discovery"

	"github.com/opdev/productctl/internal/resource"
)

var flagDiscoveredWorkloadManifestFilename string

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-product filename.yaml",
		Short: "Start building a new product listing.",
		Long:  "Scaffolds a new product listing for you to the specified filename. The contents will be a base template for you to update.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return writeInitialTemplate(args[0])
		},
	}

	cmd.PersistentFlags().StringVar(
		&flagDiscoveredWorkloadManifestFilename,
		"from-discovery-json",
		"",
		"Path to the discovered workload manifest.",
	)

	return cmd
}

func writeInitialTemplate(toFile string) error {
	base := resource.NewProductListing()

	base.Spec.Name = "My New Product"
	base.Spec.Type = resource.ProductListingTypeContainerStack
	base.Spec.Descriptions = &resource.ProductListingDescriptions{
		Long:  "This can contain long form content about your product.",
		Short: "A brief synopsis",
	}

	if flagDiscoveredWorkloadManifestFilename != "" {
		components, err := parseDiscoveredWorkloads(flagDiscoveredWorkloadManifestFilename)
		if err != nil {
			return err
		}

		base.With.Components = components
	}

	b, err := yaml.Marshal(base)
	if err != nil {
		return err
	}

	err = os.WriteFile(toFile, b, 0o644)
	if err != nil {
		return err
	}

	return nil
}

// parseDiscoveredWorkloads reads a discovered workload manifest and converts it
// to container components with some sane defaults.
func parseDiscoveredWorkloads(fromFile string) ([]*resource.Component, error) {
	manifestData, err := os.ReadFile(fromFile)
	if err != nil {
		return nil, err
	}

	var manifest discovery.Manifest
	err = json.Unmarshal(manifestData, &manifest)
	if err != nil {
		return nil, err
	}

	return libdiscovery.ParseDiscoveryManifest(manifest)
}
