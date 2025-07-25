// Package cli contains resources necessary for the execution of the command
// line interface.
package cli

import (
	"context"
	"errors"
	"io"
	"log/slog"

	"github.com/opdev/productctl/internal/catalogapi"
	"github.com/opdev/productctl/internal/logger"
)

var ErrAPIEndpointUnknown = errors.New("unknown api endpoint")

// ConfigureLogger serves as a convenience function for configuring the CLI logger,
// populating a context with it, and returning it to the user.
func ConfigureLogger(logLevel string, logTarget io.Writer) (context.Context, *slog.Logger, error) {
	l, err := logger.New(logLevel, logTarget)
	if err != nil {
		return nil, nil, err
	}

	appContext := logger.NewContextWithLogger(context.Background(), l)
	return appContext, l, nil
}

// ResolveAPIEndpoint parses a referenceString for known endpoint abbreviations
// possible catalog environments, and returns the APIEndpoint that corresponds.
// Common use is parsing a CLI environment flag value.
func ResolveAPIEndpoint(referenceString string) (catalogapi.APIEndpoint, error) {
	switch referenceString {
	case "prod":
		return catalogapi.EndpointProduction, nil
	case "stage":
		return catalogapi.EndpointStage, nil
	case "qa":
		return catalogapi.EndpointQA, nil
	case "uat":
		return catalogapi.EndpointUAT, nil
	default:
		return "", ErrAPIEndpointUnknown
	}
}
