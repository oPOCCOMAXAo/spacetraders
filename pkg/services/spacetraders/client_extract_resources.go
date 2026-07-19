package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

// Extraction is the result of an extract action: what was extracted and how much.
type Extraction struct {
	ShipSymbol string          `json:"shipSymbol"`
	Yield      ExtractionYield `json:"yield"`
}

// ExtractionYield describes the good and units produced by a single extract.
type ExtractionYield struct {
	Symbol string `json:"symbol"`
	Units  int    `json:"units"`
}

// ExtractResources mines resources from the ship's current waypoint. Requires
// the ship to be IN_ORBIT at an extractable waypoint with an extractor mount.
//
// POST /my/ships/{shipSymbol}/extract.
func (c *Client) ExtractResources(
	ctx context.Context,
	shipSymbol string,
) (Extraction, ss.Cooldown, ss.ShipCargo, error) {
	var res struct {
		Extract  Extraction   `json:"extract"`
		Cooldown ss.Cooldown  `json:"cooldown"`
		Cargo    ss.ShipCargo `json:"cargo"`
	}

	err := c.Request(ctx, Request{
		Method:    http.MethodPost,
		Path:      "/my/ships/" + shipSymbol + "/extract",
		ResultRef: &res,
	}, RequestOptions{})
	if err != nil {
		return Extraction{}, ss.Cooldown{}, ss.ShipCargo{}, err
	}

	return res.Extract, res.Cooldown, res.Cargo, nil
}
