package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

type FulfillContractRequest struct {
	ContractID string `json:"-"`
}

type FulfillContractResponse struct {
	Agent    ss.Agent    `json:"agent"`
	Contract ss.Contract `json:"contract"`
}

// FulfillContract fulfills a contract whose delivery terms are all met,
// paying the on-fulfillment credits.
//
// POST /my/contracts/{contractId}/fulfill.
func (c *Client) FulfillContract(
	ctx context.Context,
	req FulfillContractRequest,
) (FulfillContractResponse, error) {
	var res FulfillContractResponse

	err := c.Request(ctx, Request{
		Method:    http.MethodPost,
		Path:      "/my/contracts/" + req.ContractID + "/fulfill",
		ResultRef: &res,
	}, RequestOptions{})
	if err != nil {
		return FulfillContractResponse{}, err
	}

	return res, nil
}
