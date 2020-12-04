package ports

import "time"

type CreateListRequest struct {
	Name string `json:"name" validate:"required"`
}

type ListResponse struct {
	UUID         string                 `json:"id"`
	Name         string                 `json:"name"`
	CreationDate time.Time              `json:"creationDate"`
	Contributors []*ContributorResponse `json:"contributors"`
	Products     []*ProductResponse     `json:"products"`
}

type ContributorResponse struct {
	UUID string `json:"id"`
}

type ProductResponse struct {
	UUID string `json:"id"`
}
