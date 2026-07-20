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
	return c.requestListQuery(ctx, path, c.buildPageQuery(page, limit), result)
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
