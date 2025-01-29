package transport

import (
	"net/http"
)

// Ensure the tranport implements http.RoundTripper.
var _ http.RoundTripper = &AddUserAgent{}

// AddUserAgent injects UserAgent specifiet to the User-Agent header.
type AddUserAgent struct {
	Wrapped   http.RoundTripper
	UserAgent string
}

func (t *AddUserAgent) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.UserAgent)
	return t.Wrapped.RoundTrip(req)
}
