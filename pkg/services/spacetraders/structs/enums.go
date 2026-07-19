package ss

// ShipNavStatus is the navigation status of a ship.
type ShipNavStatus string

const (
	ShipNavStatusInTransit ShipNavStatus = "IN_TRANSIT"
	ShipNavStatusInOrbit   ShipNavStatus = "IN_ORBIT"
	ShipNavStatusDocked    ShipNavStatus = "DOCKED"
)

// ShipFlightMode controls the speed/fuel trade-off of travel.
type ShipFlightMode string

const (
	FlightModeDrift   ShipFlightMode = "DRIFT"
	FlightModeStealth ShipFlightMode = "STEALTH"
	FlightModeCruise  ShipFlightMode = "CRUISE"
	FlightModeBurn    ShipFlightMode = "BURN"
)

// ContractType enumerates the kinds of contract a faction can offer.
type ContractType string

const (
	ContractTypeProcurement ContractType = "PROCUREMENT"
	ContractTypeTransport   ContractType = "TRANSPORT"
	ContractTypeShuttle     ContractType = "SHUTTLE"
)

// WaypointType enumerates the kinds of waypoint in a system.
type WaypointType string

const (
	WaypointTypePlanet         WaypointType = "PLANET"
	WaypointTypeMoon           WaypointType = "MOON"
	WaypointTypeOrbitalStation WaypointType = "ORBITAL_STATION"
	WaypointTypeAsteroidField  WaypointType = "ASTEROID_FIELD"
	WaypointTypeGasGiant       WaypointType = "GAS_GIANT"
	WaypointTypeJumpGate       WaypointType = "JUMP_GATE"
	WaypointTypeNebula         WaypointType = "NEBULA"
	WaypointTypeStar           WaypointType = "STAR"
	WaypointTypeBlackHole      WaypointType = "BLACK_HOLE"
	WaypointTypeGravityWell    WaypointType = "GRAVITY_WELL"
	WaypointTypeDebrisField    WaypointType = "DEBRIS_FIELD"
	WaypointTypeGasStation     WaypointType = "GAS_STATION"
)
