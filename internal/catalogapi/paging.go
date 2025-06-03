package catalogapi

import (
	"context"
	"errors"

	"github.com/opdev/productctl/internal/logger"
)

var ErrQueryPageFailed = errors.New("failed to query page")

// QueryAll returns all items of type T in a paginated response from
// startingPage with the set pageSize.
func QueryAll[T any](
	ctx context.Context,
	startingPage, pageSize int,
	queryPageFn func(page, pageSize int) (returnedItems []T, totalItems int, queryError error),
) ([]T, error) {
	L := logger.FromContextOrDiscard(ctx)
	L.Debug("querying all pages", "startingPage", startingPage, "pageSize", pageSize)
	allItems := []T{}
	page := startingPage
	// -1 value for remaining indicates first run.
	remaining := -1
	for remaining == -1 || remaining > 0 {
		L := L.With("page", page)
		returned, total, err := queryPageFn(page, pageSize)
		if err != nil {
			return nil, errors.Join(ErrQueryPageFailed, err)
		}

		allItems = append(allItems, returned...)
		remaining = total - len(allItems)
		L.Debug("completed page query", "returnedItems", len(returned), "totalItems", total, "remainingItems", remaining)
		page++
	}

	return allItems, nil
}

const (
	// DefaultPageSize represents the default paging used for paginated queries.
	DefaultPageSize = 100
)
