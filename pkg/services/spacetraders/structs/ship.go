package ss

import "time"

// Ship is a vessel owned by the agent.
type Ship struct {
	Symbol       string           `json:"symbol"`
	Registration ShipRegistration `json:"registration"`
	Nav          ShipNav          `json:"nav"`
	Crew         ShipCrew         `json:"crew"`
	Fuel         ShipFuel         `json:"fuel"`
	Frame        ShipFrame        `json:"frame"`
	Reactor      ShipReactor      `json:"reactor"`
	Engine       ShipEngine       `json:"engine"`
	Modules      []ShipModule     `json:"modules"`
	Mounts       []ShipMount      `json:"mounts"`
	Cargo        ShipCargo        `json:"cargo"`
	Cooldown     *Cooldown        `json:"cooldown,omitempty"`
}

type ShipRegistration struct {
	Name          string `json:"name"`
	FactionSymbol string `json:"factionSymbol"`
	Role          string `json:"role"`
}

type ShipNav struct {
	SystemSymbol   string         `json:"systemSymbol"`
	WaypointSymbol string         `json:"waypointSymbol"`
	Route          ShipNavRoute   `json:"route"`
	Status         ShipNavStatus  `json:"status"`
	FlightMode     ShipFlightMode `json:"flightMode"`
	Speed          int            `json:"speed"`
}

type ShipNavRoute struct {
	Destination   ShipNavWaypoint `json:"destination"`
	Departure     ShipNavWaypoint `json:"departure"`
	Arrival       time.Time       `json:"arrival"`
	DepartureTime time.Time       `json:"departureTime"`
}

type ShipNavWaypoint struct {
	Symbol string `json:"symbol"`
	Type   string `json:"type,omitempty"`
}

type ShipCrew struct {
	Current  int    `json:"current"`
	Required int    `json:"required"`
	Capacity int    `json:"capacity"`
	Rotation string `json:"rotation"`
	Morale   int    `json:"morale"`
	Wages    int64  `json:"wages"`
}

type ShipFuel struct {
	Current  int               `json:"current"`
	Capacity int               `json:"capacity"`
	Consumed *ShipFuelConsumed `json:"consumed,omitempty"`
}

type ShipFuelConsumed struct {
	Amount    int       `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

type ShipFrame struct {
	Symbol       string           `json:"symbol"`
	Name         string           `json:"name"`
	Condition    float64          `json:"condition"`
	ModuleSlots  int              `json:"moduleSlots"`
	FuelCapacity int              `json:"fuelCapacity"`
	Requirements ShipRequirements `json:"requirements"`
}

type ShipRequirements struct {
	Power int `json:"power,omitempty"`
	Crew  int `json:"crew,omitempty"`
	Slots int `json:"slots,omitempty"`
}

type ShipReactor struct {
	Symbol       string           `json:"symbol"`
	Name         string           `json:"name"`
	Condition    float64          `json:"condition"`
	PowerOutput  int              `json:"powerOutput"`
	Requirements ShipRequirements `json:"requirements"`
}

type ShipEngine struct {
	Symbol       string           `json:"symbol"`
	Name         string           `json:"name"`
	Condition    float64          `json:"condition"`
	Integrity    float64          `json:"integrity"`
	Speed        int              `json:"speed"`
	Requirements ShipRequirements `json:"requirements"`
}

type ShipModule struct {
	Symbol       string           `json:"symbol"`
	Name         string           `json:"name"`
	Capacity     int              `json:"capacity,omitempty"`
	Range        int              `json:"range,omitempty"`
	Requirements ShipRequirements `json:"requirements"`
}

type ShipMount struct {
	Symbol       string           `json:"symbol"`
	Name         string           `json:"name"`
	Strength     int              `json:"strength,omitempty"`
	Deposits     []string         `json:"deposits,omitempty"`
	Requirements ShipRequirements `json:"requirements"`
}

type ShipCargo struct {
	Capacity  int             `json:"capacity"`
	Units     int             `json:"units"`
	Inventory []ShipCargoItem `json:"inventory"`
}

type ShipCargoItem struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Units       int    `json:"units"`
}

// Cooldown is the reactor cooldown that gates some ship actions.
type Cooldown struct {
	ShipSymbol       string    `json:"shipSymbol"`
	TotalSeconds     int       `json:"totalSeconds"`
	RemainingSeconds int       `json:"remainingSeconds"`
	Expiration       time.Time `json:"expiration"`
}

// Expired reports whether the cooldown has fully elapsed.
func (c Cooldown) Expired() bool {
	return c.RemainingSeconds <= 0 || !c.Expiration.IsZero() && time.Now().After(c.Expiration)
}
