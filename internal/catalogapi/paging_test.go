package catalogapi_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opdev/productctl/internal/catalogapi"
)

func mockData(size uint) []int {
	d := make([]int, size)
	for i := range d {
		// value is 1-based, to align with the size input.
		d[i] = i + 1
	}

	return d
}

func mockSuccessfulPaginatedQueryWithInput(inputData []int) func(page, pageSize int) ([]int, int, error) {
	// maximum number of "queries" allowed.
	max := 8
	counter := 0
	return func(page, pageSize int) ([]int, int, error) {
		if counter > max {
			return nil, -1, errors.New("TESTING CONSTRAINT REACHED: Max query count reached in mock function - Adjust unit tests to use less queries")
		}
		defer func() { counter++ }()
		start := (page - 1) * pageSize
		// account for start and end points that are outside of inputData bounds.
		if start > len(inputData) {
			return []int{}, 0, nil
		}
		end := start + pageSize
		if end > len(inputData) {
			end = len(inputData)
		}
		return inputData[start:end], len(inputData), nil
	}
}

var _ = Describe("Paging", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.TODO()
	})

	When("querying all records", func() {
		var (
			inputData    []int
			startingPage int
			pageSize     int
		)

		BeforeEach(func() {
			inputData = mockData(5)
			startingPage = 1
			pageSize = catalogapi.DefaultPageSize
		})

		When("the query function returns an error", func() {
			It("should return library errors and the provided error", func() {
				returnedErr := errors.New("query error")
				_, err := catalogapi.QueryAll(
					ctx,
					startingPage,
					pageSize,
					func(page, pageSize int) (returnedItems []int, totalItems int, queryError error) {
						return nil, 0, returnedErr
					},
				)
				Expect(err).To(MatchError(returnedErr))
				Expect(err).To(MatchError(catalogapi.ErrQueryPageFailed))
			})
		})

		When("you start at page 0", func() {
			BeforeEach(func() {
				startingPage = 1
			})

			It("should query all records", func() {
				records, err := catalogapi.QueryAll(
					ctx,
					startingPage,
					pageSize,
					mockSuccessfulPaginatedQueryWithInput(inputData),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(records).To(BeEquivalentTo(inputData))
			})
		})

		When("you start at a specific page", func() {
			BeforeEach(func() {
				inputData = mockData(3)
				startingPage = 2
				// pageSize out of scope, but adjusted to make this assertion easier.
				pageSize = 1
			})

			It("should only contain the records from that point forward", func() {
				records, err := catalogapi.QueryAll(
					ctx,
					startingPage,
					pageSize,
					mockSuccessfulPaginatedQueryWithInput(inputData),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(records).To(BeEquivalentTo(inputData[(startingPage-1)*pageSize:]))
			})
		})
		When("you define a page size", func() {
			BeforeEach(func() {
				inputData = mockData(5)
				startingPage = 1
				pageSize = 1
			})
			It("should use the expected number of queries to query all records", func() {
				totalQueryCount := 0
				_, err := catalogapi.QueryAll(
					ctx,
					startingPage,
					pageSize,
					func(page, pageSize int) (returnedItems []int, totalItems int, queryError error) {
						totalQueryCount++
						return mockSuccessfulPaginatedQueryWithInput(inputData)(page, pageSize)
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(totalQueryCount).To(BeEquivalentTo(len(inputData)))
			})
		})
	})
})
