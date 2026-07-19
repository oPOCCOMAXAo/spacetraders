package spacetraders

import "encoding/json"

// responseEnvelope mirrors the standard SpaceTraders success envelope:
//
//	{"data": <T>, "meta": {...}}
//
// Both fields are kept as RawMessage so each method can decode them into its
// own per-endpoint response/meta structs (the shape of "data" varies per
// endpoint, and "meta" is only present on list endpoints).
type responseEnvelope struct {
	Data json.RawMessage `json:"data"`
	Meta json.RawMessage `json:"meta,omitempty"`
}
