package ss

// Meta holds pagination metadata returned by list endpoints.
type Meta struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
