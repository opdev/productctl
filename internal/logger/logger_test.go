package logger_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/cli"
	"github.com/opdev/productctl/internal/logger"
)

var _ = Describe("Logger", func() {
	When("creating a new logger instance", func() {
		var L *slog.Logger

		AfterEach(func() {
			L = nil
		})

		When("using the discarding logger helper", func() {
			BeforeEach(func() {
				L = logger.DiscardingLogger()
			})
			It("should have a discarding handler", func() {
				Expect(L.Handler()).To(Equal(slog.DiscardHandler))
			})
		})

		When("the caller creates a logger with a custom configuration", func() {
			var (
				loglevel  string
				logTarget *bytes.Buffer
			)

			BeforeEach(func() {
				loglevel = "error"
				logTarget = bytes.NewBuffer([]byte{})
			})

			It("should write to the provided io.Writer", func() {
				_, L, err := cli.ConfigureLogger(loglevel, logTarget)
				Expect(err).ToNot(HaveOccurred())
				msg := "hello from test case"
				L.Error(msg)
				b, err := io.ReadAll(logTarget)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(b)).To(ContainSubstring(msg))
			})

			When("an invalid loglevel is provided", func() {
				BeforeEach(func() {
					loglevel = "invalid"
				})

				It("should fall back to info", func() {
					var err error
					L, err = logger.New(loglevel, logTarget)
					Expect(err).ToNot(HaveOccurred())
					Expect(L.Handler().Enabled(context.TODO(), slog.LevelInfo)).To(BeTrue())
				})
			})
		})

		When("working with a logger context embedding", func() {
			var ctx context.Context

			BeforeEach(func() {
				ctx = context.TODO()
			})

			When("a logger is added to the context", func() {
				BeforeEach(func() {
					L = logger.DiscardingLogger().With("embedding-test", "true")
					ctx = logger.NewContextWithLogger(ctx, L)
				})
				It("should be extractable from the context", func() {
					extracted, err := logger.FromContext(ctx)
					Expect(err).ToNot(HaveOccurred())
					Expect(extracted).To(Equal(L))
				})
			})

			When("a logger is not in the context when expected", func() {
				It("should throw an appropriate error", func() {
					_, err := logger.FromContext(ctx)
					Expect(err).To(MatchError(logger.ErrLoggerNotFoundInContext))
				})

				When("the caller wants no logs if one is not found", func() {
					It("should produce a discarding logger", func() {
						L := logger.FromContextOrDiscard(ctx)
						Expect(L.Handler()).To(Equal(slog.DiscardHandler))
					})
				})
			})
		})
	})

	When("Marshaling JSON", func() {
		When("the input data cannot be converted to JSON", func() {
			It("should produce a known failure message in the output", func() {
				b := logger.MarshalJSON(func() {})
				Expect(b).To(Equal([]byte("failedToMarshal")))
			})
		})

		When("the input data can be converted to JSON successfully", func() {
			It("should produce the known-good output", func() {
				b := logger.MarshalJSON("foo")
				Expect(b).To(Equal([]byte("\"foo\"")))
			})
		})
	})
})
