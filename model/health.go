package model

type Health struct {
	Alive   bool   `json:"alive"`
	Version string `json:"version"`
}
