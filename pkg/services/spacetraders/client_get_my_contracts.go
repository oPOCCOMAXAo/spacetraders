package spacetraders

import (
	"context"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
	"github.com/opoccomaxao/spacetraders/pkg/utils/sequtils"
)

type ListContractsRequest struct {
	Page  int `json:"-"`
	Limit int `json:"-"`
}

type ListMyContractsResponse struct {
	Contracts []ss.Contract
	Meta      ss.Meta
}

// ListMyContracts returns a single page of the agent's contracts.
//
// GET /my/contracts.
func (c *Client) ListMyContracts(
	ctx context.Context,
	req ListContractsRequest,
) (ListMyContractsResponse, error) {
	var res struct {
		Contracts []ss.Contract `json:"contracts"`
	}

	meta, err := c.requestList(ctx, "/my/contracts", req.Page, req.Limit, &res)
	if err != nil {
		return ListMyContractsResponse{}, err
	}

	return ListMyContractsResponse{Contracts: res.Contracts, Meta: meta}, nil
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
			resp, err := c.ListMyContracts(ctx, ListContractsRequest{Page: page, Limit: limit})

			return sequtils.ToPtrSlice(resp.Contracts), err
		},
	)
}
