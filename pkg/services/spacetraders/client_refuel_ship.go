package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

type RefuelShipRequest struct {
	ShipSymbol string `json:"-"`
}

type RefuelShipResponse struct {
	Agent       ss.Agent       `json:"agent"`
	Transaction ss.Transaction `json:"transaction"`
	Fuel        ss.ShipFuel    `json:"fuel"`
}

// RefuelShip refuels the ship from the current marketplace. The ship must be
// DOCKED at a waypoint with the Marketplace trait.
//
// POST /my/ships/{shipSymbol}/refuel.
func (c *Client) RefuelShip(
	ctx context.Context,
	req RefuelShipRequest,
) (RefuelShipResponse, error) {
	var res RefuelShipResponse

	err := c.Request(ctx, Request{
		Method:    http.MethodPost,
		Path:      "/my/ships/" + req.ShipSymbol + "/refuel",
		ResultRef: &res,
	}, RequestOptions{})
	if err != nil {
		return RefuelShipResponse{}, err
	}

	return res, nil
}
