package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

type DeliverContractRequest struct {
	ContractID  string `json:"-"`
	ShipSymbol  string `json:"shipSymbol"`
	TradeSymbol string `json:"tradeSymbol"`
	Units       int    `json:"units"`
}

type DeliverContractResponse struct {
	Cargo    ss.ShipCargo `json:"cargo"`
	Contract ss.Contract  `json:"contract"`
}

// DeliverContract delivers cargo from the ship to a contract. The ship must
// be docked at the contract's destination waypoint and hold the required good.
//
// POST /my/contracts/{contractId}/deliver.
func (c *Client) DeliverContract(
	ctx context.Context,
	req DeliverContractRequest,
) (DeliverContractResponse, error) {
	var res DeliverContractResponse

	err := c.Request(ctx, Request{
		Method:    http.MethodPost,
		Path:      "/my/contracts/" + req.ContractID + "/deliver",
		Body:      req,
		ResultRef: &res,
	}, RequestOptions{})
	if err != nil {
		return DeliverContractResponse{}, err
	}

	return res, nil
}
