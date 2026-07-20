package spacetraders

import (
	"context"
	"net/url"
	"strconv"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
	"github.com/opoccomaxao/spacetraders/pkg/utils/sequtils"
)

type ListWaypointsRequest struct {
	SystemSymbol string
	Page         int
	Limit        int
	// Traits optionally filters waypoints that expose the given trait symbols
	// (e.g. "MARKETPLACE", "SHIPYARD"). Repeatable in the query string.
	Traits []string
}

// ListSystemWaypoints returns a single page of waypoints in a system.
//
// GET /systems/{systemSymbol}/waypoints.
func (c *Client) ListSystemWaypoints(
	ctx context.Context,
	req ListWaypointsRequest,
) ([]ss.Waypoint, ss.Meta, error) {
	var res struct {
		Waypoints []ss.Waypoint `json:"waypoints"`
	}

	query := c.buildPageQuery(req.Page, req.Limit)
	for _, t := range req.Traits {
		query.Add("traits", t)
	}

	meta, err := c.requestListQuery(
		ctx,
		"/systems/"+req.SystemSymbol+"/waypoints",
		query,
		&res,
	)
	if err != nil {
		return nil, ss.Meta{}, err
	}

	return res.Waypoints, meta, nil
}

// ListSystemWaypointsSeq yields every waypoint in a system across all pages.
func (c *Client) ListSystemWaypointsSeq(
	ctx context.Context,
	req ListWaypointsRequest,
) func(yield func(*ss.Waypoint, error) bool) {
	return sequtils.Paginate(
		ctx,
		req.Page,
		req.Limit,
		defaultPageSize,
		func(page, limit int) ([]*ss.Waypoint, error) {
			items, _, err := c.ListSystemWaypoints(ctx, ListWaypointsRequest{
				SystemSymbol: req.SystemSymbol,
				Page:         page,
				Limit:        limit,
				Traits:       req.Traits,
			})

			return sequtils.ToPtrSlice(items), err
		},
	)
}

// buildPageQuery returns a url.Values pre-filled with page/limit when set.
func (*Client) buildPageQuery(page, limit int) url.Values {
	res := make(url.Values)
	if page > 0 {
		res.Set("page", strconv.Itoa(page))
	}

	if limit > 0 {
		res.Set("limit", strconv.Itoa(limit))
	}

	return res
}
