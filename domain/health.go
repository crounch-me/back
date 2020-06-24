package domain

type Health struct {
	Alive   bool   `json:"alive"`
	Version string `json:"version"`
}
