package spacetraders

import (
	"context"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
	"github.com/opoccomaxao/spacetraders/pkg/utils/sequtils"
)

type ListContractsRequest struct {
	Page  int
	Limit int
}

// ListMyContracts returns a single page of the agent's contracts.
//
// GET /my/contracts.
func (c *Client) ListMyContracts(
	ctx context.Context,
	req ListContractsRequest,
) ([]ss.Contract, ss.Meta, error) {
	var res struct {
		Contracts []ss.Contract `json:"contracts"`
	}

	meta, err := c.requestList(ctx, "/my/contracts", req.Page, req.Limit, &res)
	if err != nil {
		return nil, ss.Meta{}, err
	}

	return res.Contracts, meta, nil
}

// ListMyContractsSeq yields every contract across all pages.
func (c *Client) ListMyContractsSeq(
	ctx context.Context,
	req ListContractsRequest,
) func(yield func(*ss.Contract, error) bool) {
	return sequtils.Paginate(
		ctx,
		req.Page,
		req.Limit,
		defaultPageSize,
		func(page, limit int) ([]*ss.Contract, error) {
			items, _, err := c.ListMyContracts(ctx, ListContractsRequest{Page: page, Limit: limit})

			return sequtils.ToPtrSlice(items), err
		},
	)
}
