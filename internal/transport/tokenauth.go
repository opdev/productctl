package transport

import (
	"log/slog"
	"net/http"

	"github.com/opdev/productctl/internal/logger"
)

// Ensure the tranport implements http.RoundTripper.
var _ http.RoundTripper = &APITokenAuthenticated{}

// APITokenAuthenticated injects the APIKey and User Agent to the HTTP
// requests made by the consuming http.Client
type APITokenAuthenticated struct {
	Wrapped        http.RoundTripper
	Token          string
	UserAgent      string
	Logger         *slog.Logger
	resolvedLogger *slog.Logger
}

func (t *APITokenAuthenticated) RoundTrip(req *http.Request) (*http.Response, error) {
	t.logger().Debug("adding api key and user agent headers to request")
	req.Header.Set("X-API-KEY", t.Token)
	req.Header.Set("User-Agent", t.UserAgent)
	return t.Wrapped.RoundTrip(req)
}

// logger returns the provided logger, or a discarding logger if one isn't
// provided. Gives the transport logging mechanism a stable interface, given the
// lack of a context.
func (t *APITokenAuthenticated) logger() *slog.Logger {
	if t.Logger != nil {
		return t.Logger
	}

	if t.resolvedLogger == nil {
		t.resolvedLogger = logger.DiscardingLogger()
	}

	return t.resolvedLogger
}
