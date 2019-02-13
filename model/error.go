package model

type Error struct {
	Code        string `json:"error"`
	Description string `json:"errorDescription"`
}
