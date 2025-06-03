package catalogapi_test

import (
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/catalogapi"
	"github.com/opdev/productctl/internal/genpyxis"
)

var _ = Describe("Gqlerr", func() {
	var gqlErr catalogapi.GraphQLResponseError

	When("parsing GraphQL errors", func() {
		BeforeEach(func() {
			gqlErr = &genpyxis.MutateProductListingCommonResponseError{
				Status: 404,
				Detail: "some error detail",
			}
		})
		It("should return an error message containing the Status and Detail values", func() {
			parsed := catalogapi.ParseGraphQLResponseError(gqlErr)
			statusStr := strconv.Itoa(gqlErr.GetStatus())
			Expect(parsed.Error()).To(ContainSubstring(gqlErr.GetDetail()))
			Expect(parsed.Error()).To(ContainSubstring(statusStr))
		})
	})
})
