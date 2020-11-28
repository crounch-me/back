package ports

type CreateListRequest struct {
	Name string `json:"name" validate="required"`
}
