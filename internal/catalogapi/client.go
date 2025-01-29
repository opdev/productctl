package catalogapi

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Khan/genqlient/graphql"

	"github.com/opdev/productctl/internal/transport"
	"github.com/opdev/productctl/internal/version"
)

// TODO: UserAgent lives in package catalogapi but this UserAgent is CLI
// specific. It does not exist in CLI because of cyclic dependencies on this
// package. If it makes sense, find a new home for this.
var UserAgent = fmt.Sprintf("%s/%s (%s)", version.Version.BaseName, version.Version.Version, version.Version.Name)

// Ensure the client implements the graphql.Doer interface.
var _ graphql.Doer = TokenAuthenticatedHTTPClient("", nil)

// TokenAuthenticatedHTTPClient returns an HTTP client with the token and user
// agent injected at the appropriate headers.
func TokenAuthenticatedHTTPClient(
	token string,
	logger *slog.Logger,
) *http.Client {
	httpClient := http.DefaultClient

	httpClient.Transport = buildTransport(
		http.DefaultTransport,
		func(rt http.RoundTripper) http.RoundTripper {
			return &transport.RequestAndResponseLogger{
				Wrapped: rt,
				Logger:  logger,
			}
		},
		func(rt http.RoundTripper) http.RoundTripper {
			return &transport.AddUserAgent{
				Wrapped:   rt,
				UserAgent: UserAgent,
			}
		},
		func(rt http.RoundTripper) http.RoundTripper {
			return &transport.APITokenAuthenticated{
				Wrapped: rt,
				Token:   token,
				Logger:  logger,
			}
		},
	)

	// 30s timeout aligns with dialer timeouts on http.DefaultTransport.
	httpClient.Timeout = 30 * time.Second

	return httpClient
}

// buildTransport produces RoundTripper wrapped by all of those defined in
// wrapFns. Provided wrapFns are executed in order. In effect, final will be the
// innermost RoundTripper, and the last item in wrapFns will be the outermost.
func buildTransport(final http.RoundTripper, wrapFns ...func(http.RoundTripper) http.RoundTripper) http.RoundTripper {
	rt := final
	for _, fn := range wrapFns {
		rt = fn(rt)
	}

	return rt
}
