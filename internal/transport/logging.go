package transport

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/opdev/productctl/internal/logger"
)

// Ensure the tranport implements http.RoundTripper.
var _ http.RoundTripper = &RequestAndResponseLogger{}

// RequestAndResponseLogger logs the body of the request before sending it, and
// logs relevant information from the response useful for debugging.
type RequestAndResponseLogger struct {
	Wrapped        http.RoundTripper
	Logger         *slog.Logger
	resolvedLogger *slog.Logger
}

func (t *RequestAndResponseLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	L := t.logger()

	b, getErr := req.GetBody()
	rb, readErr := io.ReadAll(b)

	if getErr == nil && readErr == nil {
		L.Debug("catalog api graphql operation request", "body", string(rb))
	} else {
		L.Debug("catalog api request body could not be parsed", "getbodyErr", getErr, "readbodyErr", readErr)
	}

	resp, requestErr := t.Wrapped.RoundTrip(req)
	if requestErr != nil {
		// short-circuit if the request didn't succeed or the response was empty
		L.Debug("request failed", "error", requestErr)
		return resp, requestErr
	}

	returnedTraceID := resp.Header.Get("Trace_id")
	if returnedTraceID != "" {
		L.Debug("catalog api response", "status", resp.Status, "traceID", returnedTraceID)
	}

	return resp, requestErr
}

// logger returns the provided logger, or a discarding logger if one isn't
// provided. Gives the transport logging mechanism a stable interface, given the
// lack of a context.
func (t *RequestAndResponseLogger) logger() *slog.Logger {
	if t.Logger != nil {
		return t.Logger
	}

	if t.resolvedLogger == nil {
		t.resolvedLogger = logger.DiscardingLogger()
	}

	return t.resolvedLogger
}
