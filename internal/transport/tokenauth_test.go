package transport_test

import (
	"bytes"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/transport"
)

var _ = Describe("Tokenauth", func() {
	When("using the APITokenAuthenticated transport with a specified token", func() {
		var fakeToken string
		BeforeEach(func() {
			fakeToken = "foo-token"
		})

		It("should include the User-Agent header in the request with the right value", func() {
			t := transport.APITokenAuthenticated{
				Wrapped: http.DefaultTransport,
				Token:   fakeToken,
			}

			req, err := http.NewRequest(http.MethodGet, fakeHTTPEndpoint, bytes.NewBuffer([]byte{}))
			Expect(err).ToNot(HaveOccurred())
			// We don't care about the results of the roundtrip. Just that the
			// header had the right information injected.
			t.RoundTrip(req)
			Expect(req.Header.Get("X-API-KEY")).To(Equal(fakeToken))
		})
	})
})
