package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
	pkgerrors "github.com/pkg/errors"
)

type GetShipCooldownRequest struct {
	ShipSymbol string `json:"-"`
}

type GetShipCooldownResponse struct {
	Cooldown ss.Cooldown
}

// GetShipCooldown returns the ship's current reactor cooldown.
//
// GET /my/ships/{shipSymbol}/cooldown
//
// When the ship has no active cooldown the API responds with 204 No Content;
// in that case a zero-value Cooldown and nil error are returned.
func (c *Client) GetShipCooldown(
	ctx context.Context,
	req GetShipCooldownRequest,
) (GetShipCooldownResponse, error) {
	var cooldown ss.Cooldown

	err := c.Request(ctx, Request{
		Method:    http.MethodGet,
		Path:      "/my/ships/" + req.ShipSymbol + "/cooldown",
		ResultRef: &cooldown,
	}, RequestOptions{})
	if err != nil {
		return GetShipCooldownResponse{}, pkgerrors.WithStack(err)
	}

	return GetShipCooldownResponse{Cooldown: cooldown}, nil
}
