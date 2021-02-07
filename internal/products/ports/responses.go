package ports

type CategoryResponse struct {
	UUID string `json:"id"`
	Name string `json:"name"`
}

type ProductResponse struct {
	UUID     string            `json:"id"`
	Name     string            `json:"name"`
	Category *CategoryResponse `json:"category,omitempty"`
}
