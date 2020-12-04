package ports

import "time"

type CreateListRequest struct {
	Name string `json:"name" validate:"required"`
}

type List struct {
	UUID         string         `json:"id"`
	Name         string         `json:"name"`
	CreationDate time.Time      `json:"creationDate"`
	Contributors []*Contributor `json:"contributors"`
	Products     []*Product     `json:"products"`
}

type Contributor struct {
	UUID string `json:"id"`
}

type Product struct {
	UUID string `json:"id"`
}
