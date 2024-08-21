package models

type Size struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	MinPages *int64 `json:"minPages,omitempty"`
	MaxPages *int64 `json:"maxPages,omitempty"`
}
