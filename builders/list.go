package builders

import (
	"time"

	"github.com/crounch-me/back/internal/categories"
	"github.com/crounch-me/back/internal/list"
	"github.com/crounch-me/back/internal/users"
)

const (
	DefaultCategoryID   = "00000000-0000-0000-0000-000000000000"
	DefaultCategoryName = "Divers"
)

// ListBuilder is a builder for list responses
type ListBuilder struct{}

// CategoryInGetListResponse represents a category in get list responses
type CategoryInGetListResponse struct {
	ID       string                      `json:"id"`
	Name     string                      `json:"name"`
	Products []*ProductInGetListResponse `json:"products,omitempty"`
}

// ProductInGetListResponse represents a product in get list responses
type ProductInGetListResponse struct {
	ID       string               `json:"id"`
	Name     string               `json:"name" validate:"required,lt=61"`
	Owner    *users.User          `json:"owner,omitempty"`
	Bought   bool                 `json:"bought"`
	Category *categories.Category `json:"category,omitempty"`
}

// GetListResponse is the response to the get list request
type GetListResponse struct {
	ID              string                       `json:"id"`
	Name            string                       `json:"name" validate:"required,lt=61"`
	CreationDate    time.Time                    `json:"creationDate"`
	ArchivationDate *time.Time                   `json:"archivationDate,omitempty"`
	Contributors    []*users.User                `json:"contributors,omitempty"`
	Categories      []*CategoryInGetListResponse `json:"categories"`
}

// GetList builds the response to GetList request from a regular List
func (lb *ListBuilder) GetList(list *list.List) *GetListResponse {
	listResponse := &GetListResponse{
		ID:              list.ID,
		Name:            list.Name,
		CreationDate:    list.CreationDate,
		ArchivationDate: list.ArchivationDate,
		Contributors:    list.Contributors,
	}

	categoriesInGetListResponse := make(map[string]*CategoryInGetListResponse, 0)

	for _, product := range list.Products {
		productInGetListResponse := &ProductInGetListResponse{
			ID:     product.ID,
			Name:   product.Name,
			Owner:  product.Owner,
			Bought: product.Bought,
		}

		categoryKey := DefaultCategoryID
		categoryName := DefaultCategoryName

		if product.Category != nil {
			categoryKey = product.Category.ID
			categoryName = product.Category.Name
		}

		productInGetListResponse.Category = &categories.Category{
			ID:   categoryKey,
			Name: categoryName,
		}

		if category, ok := categoriesInGetListResponse[categoryKey]; ok {
			category.Products = append(category.Products, productInGetListResponse)
		} else {
			categoriesInGetListResponse[categoryKey] = &CategoryInGetListResponse{
				ID:   categoryKey,
				Name: categoryName,
				Products: []*ProductInGetListResponse{
					productInGetListResponse,
				},
			}
		}
	}

	listCategories := make([]*CategoryInGetListResponse, 0)

	for _, category := range categoriesInGetListResponse {
		listCategories = append(listCategories, category)
	}

	listResponse.Categories = listCategories

	return listResponse
}
