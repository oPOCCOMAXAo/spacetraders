package spacetraders

import (
	"context"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
	"github.com/opoccomaxao/spacetraders/pkg/utils/sequtils"
)

type ListShipsRequest struct {
	Page  int `json:"-"`
	Limit int `json:"-"`
}

type ListMyShipsResponse struct {
	Ships []ss.Ship
	Meta  ss.Meta
}

// ListMyShips returns a single page of the agent's ships.
//
// GET /my/ships.
func (c *Client) ListMyShips(
	ctx context.Context,
	req ListShipsRequest,
) (ListMyShipsResponse, error) {
	var res struct {
		Ships []ss.Ship `json:"ships"`
	}

	meta, err := c.requestList(ctx, "/my/ships", req.Page, req.Limit, &res)
	if err != nil {
		return ListMyShipsResponse{}, err
	}

	return ListMyShipsResponse{Ships: res.Ships, Meta: meta}, nil
}

// ListMyShipsSeq yields every ship across all pages.
//
// Page-based pagination: walks page = 1..N with Limit items per page until a
// page returns fewer than Limit items (or an empty page). Errors are yielded
// in-stream and stop the iteration.
func (c *Client) ListMyShipsSeq(
	ctx context.Context,
	req ListShipsRequest,
) func(yield func(*ss.Ship, error) bool) {
	return sequtils.Paginate(
		ctx,
		req.Page,
		req.Limit,
		defaultPageSize,
		func(page, limit int) ([]*ss.Ship, error) {
			resp, err := c.ListMyShips(ctx, ListShipsRequest{Page: page, Limit: limit})

			return sequtils.ToPtrSlice(resp.Ships), err
		},
	)
}
