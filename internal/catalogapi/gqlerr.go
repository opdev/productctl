package catalogapi

import (
	"fmt"
)

// ParseGraphQLResponseError returns a Go error for the input backendErr.
func ParseGraphQLResponseError(backendErr GraphQLResponseError) error {
	return fmt.Errorf("error sending request with status \"%d\" and detail \"%s\"", backendErr.GetStatus(), backendErr.GetDetail())
}

// GraphQLResponseError contains the methods the CatalogAPI implements for
// request errors. Generated code containing backend error data is expected to
// implement this interface
type GraphQLResponseError interface {
	GetStatus() int
	GetDetail() string
}
