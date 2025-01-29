package resource

import (
	"io"

	"sigs.k8s.io/yaml"
)

// ReadProductListing reads the ProductListing resource from the io.Reader. It
// assumes YAML-formatted contents. Caller is responsible for assuring that the
// returned struct contains the necessary data. Does not fail if extra values
// are found in the input data.
func ReadProductListing(in io.Reader) (*ProductListingDeclaration, error) {
	listing := NewProductListing()

	b, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(b, &listing); err != nil {
		return nil, err
	}

	return &listing, nil
}
