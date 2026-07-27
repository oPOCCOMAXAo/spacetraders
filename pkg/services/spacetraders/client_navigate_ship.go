package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

type NavigateShipRequest struct {
	ShipSymbol     string `json:"-"`
	WaypointSymbol string `json:"waypointSymbol"`
}

type NavigateShipResponse struct {
	Nav  ss.ShipNav  `json:"nav"`
	Fuel ss.ShipFuel `json:"fuel"`
}

// NavigateShip commands a ship to travel to a waypoint in the same system.
//
// POST /my/ships/{shipSymbol}/navigate.
func (c *Client) NavigateShip(
	ctx context.Context,
	req NavigateShipRequest,
) (NavigateShipResponse, error) {
	var res NavigateShipResponse

	err := c.Request(ctx, Request{
		Method:    http.MethodPost,
		Path:      "/my/ships/" + req.ShipSymbol + "/navigate",
		Body:      req,
		ResultRef: &res,
	}, RequestOptions{})
	if err != nil {
		return NavigateShipResponse{}, err
	}

	return res, nil
}
