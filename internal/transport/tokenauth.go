package transport

import (
	"net/http"
)

// Ensure the tranport implements http.RoundTripper.
var _ http.RoundTripper = &APITokenAuthenticated{}

// APITokenAuthenticated injects the APIKey and User Agent to the HTTP
// requests made by the consuming http.Client
type APITokenAuthenticated struct {
	Wrapped http.RoundTripper
	Token   string
}

func (t *APITokenAuthenticated) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-API-KEY", t.Token)
	return t.Wrapped.RoundTrip(req)
}
