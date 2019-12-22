package model

// OFFProduct represents a product coming from open food facts
type OFFProduct struct {
	Code string `json:"code,omitempty" validate:"required,len=13"`
}
