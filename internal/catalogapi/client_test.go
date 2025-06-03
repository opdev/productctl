package catalogapi_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/catalogapi"
)

var _ = Describe("Client", func() {
	var (
		testLogger *slog.Logger
		testServer *httptest.Server
	)

	BeforeEach(func() {
		testLogger = slog.New(slog.DiscardHandler)
		testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
	})
	AfterEach(func() {
		testServer.Close()
	})
	When("when building an HTTP client", func() {
		var testToken string

		BeforeEach(func() {
			testToken = "test-token"
		})
		When("a token is provided", func() {
			It("should be included in the client", func() {
				client := catalogapi.TokenAuthenticatedHTTPClient(testToken, testLogger)
				req, err := http.NewRequest(http.MethodGet, testServer.URL, bytes.NewBuffer([]byte("testRequest")))
				Expect(err).ToNot(HaveOccurred())
				_, err = client.Do(req)
				Expect(err).ToNot(HaveOccurred())
				Expect(req.Header.Get("X-API-KEY")).To(Equal(testToken))
			})
		})

		It("should have the appropriate user agent configured", func() {
			client := catalogapi.TokenAuthenticatedHTTPClient(testToken, testLogger)
			req, err := http.NewRequest(http.MethodGet, testServer.URL, bytes.NewBuffer([]byte("testRequest")))
			Expect(err).ToNot(HaveOccurred())
			_, err = client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(req.Header.Get("User-Agent")).To(Equal(catalogapi.UserAgent))
		})
	})
})
