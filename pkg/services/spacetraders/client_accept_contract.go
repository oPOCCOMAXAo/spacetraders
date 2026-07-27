package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

type AcceptContractRequest struct {
	ContractID string `json:"-"`
}

type AcceptContractResponse struct {
	Agent    ss.Agent    `json:"agent"`
	Contract ss.Contract `json:"contract"`
}

// AcceptContract accepts an offered contract, paying the upfront credits.
//
// POST /my/contracts/{contractId}/accept.
func (c *Client) AcceptContract(
	ctx context.Context,
	req AcceptContractRequest,
) (AcceptContractResponse, error) {
	var res AcceptContractResponse

	err := c.Request(ctx, Request{
		Method:    http.MethodPost,
		Path:      "/my/contracts/" + req.ContractID + "/accept",
		ResultRef: &res,
	}, RequestOptions{})
	if err != nil {
		return AcceptContractResponse{}, err
	}

	return res, nil
}
