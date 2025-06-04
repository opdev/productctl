package cli_test

import (
	"bytes"
	"fmt"
	"io"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/catalogapi"
	"github.com/opdev/productctl/internal/cli"
	"github.com/opdev/productctl/internal/logger"
)

var _ = Describe("CLI", func() {
	When("ensuring the environment", func() {
		var (
			expectedToken string
			expectedOrgID string
		)

		AfterEach(func() {
			os.Setenv(cli.EnvAPIToken, "")
			os.Setenv(cli.EnvOrgID, "")
		})

		When("the caller sets the expected environment variables", func() {
			BeforeEach(func() {
				expectedToken = "foo"
				expectedOrgID = "1234"
				os.Setenv(cli.EnvAPIToken, expectedToken)
				os.Setenv(cli.EnvOrgID, expectedOrgID)
			})
			It("should return the orgID as an integer equivalent of the string input", func() {
				resolvedOrgID, _, err := cli.EnsureEnv()
				Expect(err).ToNot(HaveOccurred())
				Expect(expectedOrgID).To(BeEquivalentTo(fmt.Sprintf("%d", resolvedOrgID)))
			})
			It("should return the token value from the environment", func() {
				_, resolvedToken, err := cli.EnsureEnv()
				Expect(err).ToNot(HaveOccurred())
				Expect(resolvedToken).To(Equal(expectedToken))
			})

			When("the orgID format is not a valid integer", func() {
				// The catalog API expects OrgID to be an integer type, but
				// environment variables are always strings, so we do the
				// conversion and throw an error if it doesn't succeed.
				//
				// This also implies that OrgIDs can't lead with 0 because
				// integer conversion would drop that. Therefore, we assume
				// OrgIDs can never have a leading 0.
				BeforeEach(func() {
					expectedOrgID = "abcd"
					os.Setenv(cli.EnvOrgID, expectedOrgID)
				})
				It("should throw an error indicating the value is malformed", func() {
					_, _, err := cli.EnsureEnv()
					Expect(err).To(MatchError(cli.ErrEnvVarInvalidFormat))
				})
			})
		})

		When("the caller is missing the token", func() {
			BeforeEach(func() {
				os.Setenv(cli.EnvAPIToken, "")
				os.Setenv(cli.EnvOrgID, "1234")
			})
			It("should throw the expected error when the token is missing", func() {
				_, _, err := cli.EnsureEnv()
				Expect(err).To(MatchError(cli.ErrEnvVarMissing))
			})
		})

		When("the caller is missing the orgID", func() {
			BeforeEach(func() {
				os.Setenv(cli.EnvAPIToken, "foo")
				os.Setenv(cli.EnvOrgID, "")
			})
			It("should throw the expected error when the org ID is missing", func() {
				_, _, err := cli.EnsureEnv()
				Expect(err).To(MatchError(cli.ErrEnvVarMissing))
			})
		})
	})

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
