package ss

import "time"

// Waypoint is a location within a system that ships can travel to.
type Waypoint struct {
	Symbol       string            `json:"symbol"`
	Type         WaypointType      `json:"type"`
	SystemSymbol string            `json:"systemSymbol"`
	X            int               `json:"x"`
	Y            int               `json:"y"`
	Orbitals     []WaypointOrbital `json:"orbitals"`
	Faction      *WaypointFaction  `json:"faction,omitempty"`
	Traits       []WaypointTrait   `json:"traits"`
	Chart        *WaypointChart    `json:"chart,omitempty"`
}

type WaypointOrbital struct {
	Symbol string `json:"symbol"`
}

type WaypointFaction struct {
	Symbol string `json:"symbol"`
}

type WaypointTrait struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type WaypointChart struct {
	Waypoint    WaypointOrbital `json:"waypoint"`
	SubmittedBy string          `json:"submittedBy"`
	SubmittedOn time.Time       `json:"submittedOn"`
}
