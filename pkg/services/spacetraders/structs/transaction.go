package ss

import "time"

// Transaction records a market sell/purchase event.
type Transaction struct {
	WaypointSymbol string    `json:"waypointSymbol"`
	ShipSymbol     string    `json:"shipSymbol"`
	TradeSymbol    string    `json:"tradeSymbol"`
	Type           string    `json:"type"`
	Units          int       `json:"units"`
	PricePerUnit   int64     `json:"pricePerUnit"`
	TotalPrice     int64     `json:"totalPrice"`
	Timestamp      time.Time `json:"timestamp"`
}
