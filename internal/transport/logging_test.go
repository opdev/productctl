package transport_test

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/transport"
)

var _ = Describe("Logging", func() {
	When("using the RequestLogger transport with a configured logger", func() {
		var (
			t          transport.RequestLogger
			logger     *slog.Logger
			logBuffer  *bytes.Buffer
			testServer *httptest.Server
			traceID    string
		)

		BeforeEach(func() {
			logBuffer = bytes.NewBuffer([]byte{})
			logger = slog.New(slog.NewTextHandler(logBuffer, &slog.HandlerOptions{Level: slog.LevelDebug}))
			// The trace just needs to be something relatively unique per run.
			traceID = fmt.Sprintf("test-trace-%s", time.Now().Format(time.RFC3339Nano))
			testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Trace_id", traceID)
				w.WriteHeader(http.StatusOK)
			}))

			t = transport.RequestLogger{
				Wrapped: http.DefaultTransport,
				Logger:  logger,
			}
		})

		AfterEach(func() {
			testServer.Close()
		})

		It("should include the request body in the log output", func() {
			body := bytes.NewBufferString("foo-body")
			req, err := http.NewRequest(http.MethodGet, testServer.URL, body)
			Expect(err).ToNot(HaveOccurred())
			// We don't care about the results of the roundtrip. Just that the
			// header had the right information injected.
			t.RoundTrip(req)
			Expect(bytes.Contains(logBuffer.Bytes(), body.Bytes())).To(BeTrue())
		})

		It("should include the resposne Trace ID, if set", func() {
			req, err := http.NewRequest(http.MethodGet, testServer.URL, bytes.NewBuffer([]byte{}))
			Expect(err).ToNot(HaveOccurred())
			// We don't care about the results of the roundtrip. Just that the
			// header had the right information injected.
			resp, err := t.RoundTrip(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Header.Get("Trace_id")).To(Equal(traceID))
			Expect(logBuffer.String()).To(ContainSubstring(traceID))
		})
	})
})
