package spacetraders

import (
	"context"
	"net/http"
	"net/url"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

// requestList performs a GET against a paginated endpoint, decoding "data"
// into result (a *struct{<Key>: []T}) and "meta" into a returned Meta.
func (c *Client) requestList(
	ctx context.Context,
	path string,
	page, limit int,
	result any,
) (ss.Meta, error) {
	return c.requestListQuery(ctx, path, buildPageQuery(page, limit), result)
}

// requestListQuery is the query-bearing variant for endpoints that add
// extra filters (e.g. waypoint traits).
func (c *Client) requestListQuery(
	ctx context.Context,
	path string,
	query url.Values,
	result any,
) (ss.Meta, error) {
	var meta ss.Meta

	err := c.Request(ctx, Request{
		Method:    http.MethodGet,
		Path:      path,
		Query:     query,
		ResultRef: result,
		MetaRef:   &meta,
	}, RequestOptions{})

	return meta, err
}

// paginate wraps page/limit pagination in an iter.Seq2-style generator.
//
// It walks pages starting from startPage (default 1) with the given limit
// (default defaultPageSize) until a page returns fewer than `limit` items or
// an empty page, yielding each element in turn. A fetch error is yielded
// in-stream and stops iteration.
func paginate[T any](
	ctx context.Context,
	startPage, limit int,
	fetch func(page, limit int) ([]*T, error),
) func(yield func(*T, error) bool) {
	page := startPage
	if page <= 0 {
		page = 1
	}

	pageLimit := limit
	if pageLimit <= 0 {
		pageLimit = defaultPageSize
	}

	return func(yield func(*T, error) bool) {
		for ctx.Err() == nil {
			items, err := fetch(page, pageLimit)
			if err != nil {
				yield(nil, err)

				return
			}

			if len(items) == 0 {
				return
			}

			for _, item := range items {
				if !yield(item, nil) {
					return
				}
			}

			if len(items) < pageLimit {
				return
			}

			page++
		}
	}
}

// toPtrSlice returns a slice of pointers into the given items.
func toPtrSlice[T any](items []T) []*T {
	res := make([]*T, len(items))
	for i := range items {
		res[i] = &items[i]
	}

	return res
}
