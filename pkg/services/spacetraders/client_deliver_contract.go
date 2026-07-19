package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

type DeliverContractRequest struct {
	ShipSymbol  string `json:"shipSymbol"`
	TradeSymbol string `json:"tradeSymbol"`
	Units       int    `json:"units"`
}

// DeliverContract delivers cargo from the ship to a contract. The ship must
// be docked at the contract's destination waypoint and hold the required good.
//
// POST /my/contracts/{contractId}/deliver.
func (c *Client) DeliverContract(
	ctx context.Context,
	contractID string,
	req DeliverContractRequest,
) (ss.ShipCargo, ss.Contract, error) {
	var res struct {
		Cargo    ss.ShipCargo `json:"cargo"`
		Contract ss.Contract  `json:"contract"`
	}

	err := c.Request(ctx, Request{
		Method:    http.MethodPost,
		Path:      "/my/contracts/" + contractID + "/deliver",
		Body:      req,
		ResultRef: &res,
	}, RequestOptions{})
	if err != nil {
		return ss.ShipCargo{}, ss.Contract{}, err
	}

	return res.Cargo, res.Contract, nil
}
