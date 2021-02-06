package ports

type CreateProductRequest struct {
	Name string `json:"name" validate:"required"`
}
