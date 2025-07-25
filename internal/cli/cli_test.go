package cli_test

import (
	"bytes"
	"io"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/catalogapi"
	"github.com/opdev/productctl/internal/cli"
	"github.com/opdev/productctl/internal/logger"
)

var _ = Describe("CLI", func() {
	When("configuring the logger", func() {
		var (
			loglevel  string
			logTarget *bytes.Buffer
		)

		BeforeEach(func() {
			loglevel = "info"
			logTarget = bytes.NewBuffer([]byte{})
		})

		It("should return a context containing the logger", func() {
			ctx, L, err := cli.ConfigureLogger(loglevel, logTarget)
			Expect(err).ToNot(HaveOccurred())
			Expect(ctx).ToNot(BeNil())
			loggerFromContext, err := logger.FromContext(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(loggerFromContext).To(Equal(L))
		})
		It("should write to the provided io.Writer", func() {
			_, L, err := cli.ConfigureLogger(loglevel, logTarget)
			Expect(err).ToNot(HaveOccurred())
			msg := "hello from test case"
			L.Info(msg)
			b, err := io.ReadAll(logTarget)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(b)).To(ContainSubstring(msg))
		})
	})

	When("resolving catalog API endpoints", func() {
		DescribeTable("given known shortnames",
			func(endpointShortname string, expectedEndpoint catalogapi.APIEndpoint) {
				resolved, err := cli.ResolveAPIEndpoint(endpointShortname)
				Expect(err).ToNot(HaveOccurred())
				Expect(resolved).To(Equal(expectedEndpoint))
			},
			Entry("specifying prod", "prod", catalogapi.EndpointProduction),
			Entry("specifying stage", "stage", catalogapi.EndpointStage),
			Entry("specifying uat", "uat", catalogapi.EndpointUAT),
			Entry("specifying qa", "qa", catalogapi.EndpointQA),
		)

		When("the user provides an unknown shortname", func() {
			It("should throw an error", func() {
				_, err := cli.ResolveAPIEndpoint("foo")
				Expect(err).To(MatchError(cli.ErrAPIEndpointUnknown))
			})
		})
	})
})
