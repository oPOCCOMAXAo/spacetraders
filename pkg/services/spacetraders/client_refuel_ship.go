package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

// RefuelShip refuels the ship from the current marketplace. The ship must be
// DOCKED at a waypoint with the Marketplace trait.
//
// POST /my/ships/{shipSymbol}/refuel.
func (c *Client) RefuelShip(
	ctx context.Context,
	shipSymbol string,
) (ss.Agent, ss.Transaction, ss.ShipFuel, error) {
	var res struct {
		Agent       ss.Agent       `json:"agent"`
		Transaction ss.Transaction `json:"transaction"`
		Fuel        ss.ShipFuel    `json:"fuel"`
	}

	err := c.Request(ctx, Request{
		Method:    http.MethodPost,
		Path:      "/my/ships/" + shipSymbol + "/refuel",
		ResultRef: &res,
	}, RequestOptions{})
	if err != nil {
		return ss.Agent{}, ss.Transaction{}, ss.ShipFuel{}, err
	}

	return res.Agent, res.Transaction, res.Fuel, nil
}
