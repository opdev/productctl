package transport_test

import (
	"bytes"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/transport"
)

var _ = Describe("Useragent", func() {
	When("using the UserAgent transport with a specified user agent", func() {
		var userAgent string
		BeforeEach(func() {
			userAgent = "foo-useragent"
		})

		It("should include the User-Agent header in the request with the right value", func() {
			t := transport.AddUserAgent{
				Wrapped:   http.DefaultTransport,
				UserAgent: userAgent,
			}

			req, err := http.NewRequest(http.MethodGet, fakeHTTPEndpoint, bytes.NewBuffer([]byte{}))
			Expect(err).ToNot(HaveOccurred())
			// We don't care about the results of the roundtrip. Just that the
			// header had the right information injected.
			t.RoundTrip(req)
			Expect(req.Header.Get("User-Agent")).To(Equal(userAgent))
		})
	})
})
