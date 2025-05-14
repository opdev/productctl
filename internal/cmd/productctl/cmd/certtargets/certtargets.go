package certtargets

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/opdev/productctl/internal/logger"
	"github.com/opdev/productctl/internal/resource"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-component-mappings [new-product.yaml]",
		Short: "Generates a certification mapping file for a given product.",
		RunE:  runE,
		Args:  cobra.MinimumNArgs(1),
	}

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	L := logger.FromContextOrDiscard(cmd.Context())
	productFile, err := os.Open(args[0])
	if err != nil {
		return err
	}

	defer productFile.Close()

	declaration, err := resource.ReadProductListing(productFile)
	if err != nil {
		return err
	}

	components := declaration.With.Components
	if len(components) == 0 {
		return errors.New("no components in this product")
	}

	// Sort the components by their types.
	helmComponents := []*resource.Component{}
	operatorComponents := []*resource.Component{}
	containerComponents := []*resource.Component{}

	L.Debug("sorting containers by type")
	for _, cmp := range components {
		// Components must have IDs to be mapped and processed correctly. If a
		// component does not have an ID, we assume it does not exist in the
		// upstream.
		if cmp.ID == "" {
			L.Warn("Skipping component that does not have an ID in provided product declaration", "component_name", cmp.Name)
			continue
		}

		switch cmp.Type {
		case resource.ComponentTypeHelmChart:
			L.Debug("component is a helm chart", "component_id", cmp.ID)
			helmComponents = append(helmComponents, cmp)
		case resource.ComponentTypeContainer:
			L.Debug("component is a container", "component_id", cmp.ID)
			if cmp.Container == nil {
				L.Debug(
					"Skipping container component that doesn't have container details",
					"component_id", cmp.ID,
					"component_name", cmp.Name,
				)
				continue
			}

			if cmp.Container.OSContentType == resource.ContentTypeOperatorBundle && cmp.Container.Type == resource.ContainerTypeOperatorBundle {
				L.Debug("more specifically, the component is an operator bundle", "component_id", cmp.ID)
				operatorComponents = append(operatorComponents, cmp)
				continue
			}

			if cmp.Container.Type == resource.ContainerTypeContainer {
				L.Debug("more specifically, the component is a standard container")
				containerComponents = append(containerComponents, cmp)
				continue
			}

			L.Warn("container component could not be sorted", "_id", cmp.ID)
		default:
			L.Warn("component was skipped because it was not of an expected type", "componentType", cmp.Type)
			continue
		}
	}
	L.Debug("done sorting containers by type")

	// buffer out contains all of the data to be written. It's stored here so
	// that it can be written to the runtime output (e.g. stdout) all at once.
	out := &bytes.Buffer{}

	containerComponentsWritten := 0
	for _, cmp := range containerComponents {
		image := "registry.example.com/placeholder/placeholder"
		if cmp.Container.Registry != "" || cmp.Container.Repository != "" || cmp.Container.RepositoryName != "" {
			image = strings.Join([]string{cmp.Container.Registry, cmp.Container.Repository, cmp.Container.RepositoryName}, "/")
		}
		L.Debug("executing container template")
		b := &bytes.Buffer{}

		// Only write the section title just before the first component is
		// written.
		if containerComponentsWritten == 0 {
			fmt.Fprintln(b, "container_components:")
		}
		err := containerTemplate(image, cmp.ID, cmp.Name, b)
		if err != nil {
			L.Error("failed to generate template for component", "errMsg", err)
			continue
		}

		// Template executions were already successful, so we can just copy all
		// of the content through to the final buffer.
		if _, err := io.Copy(out, b); err != nil {
			return err
		}
		containerComponentsWritten++
	}

	helmComponentsWritten := 0
	for _, cmp := range helmComponents {
		chartURI := "https://example.com/path/to/your/chart-0.1.1.tgz"

		b := &bytes.Buffer{}
		if helmComponentsWritten == 0 {
			fmt.Fprintln(b, "helm_chart_components:")
		}
		err := helmChartTemplate(chartURI, cmp.ID, cmp.Name, b)
		if err != nil {
			L.Error("failed to generate template for component", "errMsg", err)
			continue
		}

		if _, err := io.Copy(out, b); err != nil {
			return err
		}
		helmComponentsWritten++
	}

	operatorComponentsWritten := 0
	for _, cmp := range operatorComponents {
		image := "registry.example.com/operatorbundle/placeholder"
		if cmp.Container.Registry != "" || cmp.Container.Repository != "" || cmp.Container.RepositoryName != "" {
			image = strings.Join([]string{cmp.Container.Registry, cmp.Container.Repository, cmp.Container.RepositoryName}, "/")
		}
		L.Debug("executing operator template")
		b := &bytes.Buffer{}
		// Make sure the parent key is written before the first component.
		if operatorComponentsWritten == 0 {
			fmt.Fprintln(b, "operator_components:")
		}
		err := operatorTemplate(image, cmp.ID, cmp.Name, b)
		if err != nil {
			L.Error("failed to generate template for component", "errMsg", err)
			continue
		}

		if _, err := io.Copy(out, b); err != nil {
			return err
		}
		operatorComponentsWritten++
	}

	fmt.Fprint(os.Stdout, out)
	L.Info(
		"Total components written successfully",
		"container_components", containerComponentsWritten,
		"operator_components", operatorComponentsWritten,
		"helm_chart_components", helmComponentsWritten,
	)
	return nil
}

func containerTemplate(img, key, name string, out io.Writer) error {
	t := `  ## {{ .Name }}
  {{ .Key }}:
    image_ref: {{ .Image }}
    tags:
      ## A tag you want to certify.
      - tag: placeholder
      ## If this tag should use different tool_flags than what you configure at
      ## the top level for this component, uncomment this line and specify that
      ## here.
      # tool_flags: []
`

	entryTemplate, err := template.New("entry").Parse(t)
	if err != nil {
		return err
	}

	err = entryTemplate.Execute(out, map[string]string{
		"Image": img,
		"Name":  name,
		"Key":   key,
	})
	if err != nil {
		return err
	}

	return nil
}

func helmChartTemplate(chartURI, key, name string, out io.Writer) error {
	t := `  ## {{ .Name }}
  {{ .Key }}:
    chart_uri: {{ .ChartURI }}
`

	entryTemplate, err := template.New("entry").Parse(t)
	if err != nil {
		return err
	}

	err = entryTemplate.Execute(out, map[string]string{
		"ChartURI": chartURI,
		"Name":     name,
		"Key":      key,
	})
	if err != nil {
		return err
	}

	return nil
}

func operatorTemplate(img, key, name string, out io.Writer) error {
	t := `  ## {{ .Name }}
  {{ .Key }}:
    image_ref: {{ .Image }}
    ## The operator index image (or catalog) that contains the image_ref at the
    ## listed tags. Per-tag index_images can also be configured alongside the tag
    ## definition
    index_image: placeholder.example.com/namespace/index-image:latest
    ## The tool_flags directive contains any flags to set on the
    ## certification tooling.
    tool_flags:
	  # kubeconfig is a path to the kubeconfig to use for this component.
	  kubeconfig: component-kubeconfig
    ## tags you wish to certify are configured here.
    tags:
      - tag: placeholder
        ## Different tags different sets of flags to pass to certification tools.
        ## If that's the case, you can configure those on a per-tag basis by uncommenting
        ## the below configuration
        # tool_flags: {}
        ## If this tag exists in a custom index image, uncommnent and set that here.
        # index_image: placeholder.example.com/namespace/per-tag-index-image:latest
`

	entryTemplate, err := template.New("entry").Parse(t)
	if err != nil {
		return err
	}

	err = entryTemplate.Execute(out, map[string]string{
		"Image": img,
		"Name":  name,
		"Key":   key,
	})
	if err != nil {
		return err
	}

	return nil
}
