package models

type Era struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	MinYear *int64 `json:"minYear,omitempty"`
	MaxYear *int64 `json:"maxYear,omitempty"`
}
