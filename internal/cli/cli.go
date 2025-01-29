// Package cli contains resources necessary for the execution of the command
// line interface.
package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"

	"github.com/opdev/productctl/internal/catalogapi"
	"github.com/opdev/productctl/internal/logger"
)

var (
	ErrEnvVarMissing       = errors.New("required environment variable is missing")
	ErrEnvVarInvalidFormat = errors.New("required environment variable is malformed")
)

var (
	EnvAPIToken = "CONNECT_API_TOKEN"
	EnvOrgID    = "CONNECT_ORG_ID"
)

// EnsureEnv looks for the minimum required environment variables for the CLI to function
func EnsureEnv() (int, string, error) {
	token := os.Getenv(EnvAPIToken)
	orgIDstr := os.Getenv(EnvOrgID)

	if token == "" {
		return 0, "", fmt.Errorf("%w: CONNECT_API_TOKEN must be set", ErrEnvVarMissing)
	}

	if orgIDstr == "" {
		return 0, "", fmt.Errorf("%w: CONNECT_ORG_ID must be set", ErrEnvVarMissing)
	}

	var orgID int
	var err error
	if orgID, err = strconv.Atoi(orgIDstr); err != nil {
		return 0, "", fmt.Errorf("%w: OrgID did not convert nicely to an integer which is unexpected", ErrEnvVarInvalidFormat)
	}

	return orgID, token, nil
}

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

var ErrAPIEndpointUnknown = errors.New("unknown api endpoint")

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
