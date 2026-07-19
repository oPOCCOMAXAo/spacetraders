package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

type NavigateShipRequest struct {
	WaypointSymbol string `json:"waypointSymbol"`
}

// NavigateShip commands a ship to travel to a waypoint in the same system.
//
// POST /my/ships/{shipSymbol}/navigate.
func (c *Client) NavigateShip(
	ctx context.Context,
	shipSymbol string,
	req NavigateShipRequest,
) (ss.ShipNav, ss.ShipFuel, error) {
	var res struct {
		Nav  ss.ShipNav  `json:"nav"`
		Fuel ss.ShipFuel `json:"fuel"`
	}

	err := c.Request(ctx, Request{
		Method:    http.MethodPost,
		Path:      "/my/ships/" + shipSymbol + "/navigate",
		Body:      req,
		ResultRef: &res,
	}, RequestOptions{})
	if err != nil {
		return ss.ShipNav{}, ss.ShipFuel{}, err
	}

	return res.Nav, res.Fuel, nil
}
