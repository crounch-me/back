package builders

import (
	"testing"
	"time"

	"github.com/crounch-me/back/internal/categories"
	"github.com/crounch-me/back/internal/list"
	"github.com/crounch-me/back/internal/products"
	"github.com/crounch-me/back/internal/users"
	"gotest.tools/assert"
)

func TestGetListOK(t *testing.T) {
	builder := &ListBuilder{}

	productID1 := "productID1"
	productName1 := "productName1"
	productID2 := "productID2"
	productName2 := "productName2"

	categoryID := "categoryID"
	categoryName := "categoryName"

	listID := "listID"
	listName := "listName"
	creationDate := time.Now()
	archivationDate := time.Now()

	userID := "userID"
	email := "email"

	user := &users.User{
		ID:    userID,
		Email: email,
	}

	list := &list.List{
		ID:              listID,
		Name:            listName,
		CreationDate:    creationDate,
		ArchivationDate: &archivationDate,
		Contributors: []*users.User{
			user,
		},
		Products: []*list.ProductInList{
			{
				Product: &products.Product{
					ID:   productID1,
					Name: productName1,
					Category: &categories.Category{
						ID:   categoryID,
						Name: categoryName,
					},
				},
				Bought: false,
			},
			{
				Product: &products.Product{
					ID:   productID2,
					Name: productName2,
					Category: &categories.Category{
						ID:   categoryID,
						Name: categoryName,
					},
				},
				Bought: true,
			},
		},
	}

	expectedListResponse := &GetListResponse{
		ID:              listID,
		Name:            listName,
		CreationDate:    creationDate,
		ArchivationDate: &archivationDate,
		Contributors: []*users.User{
			user,
		},
		Categories: []*CategoryInGetListResponse{
			{
				ID:   categoryID,
				Name: categoryName,
				Products: []*ProductInGetListResponse{
					{
						ID:     productID1,
						Name:   productName1,
						Bought: false,
						Category: &categories.Category{
							ID:   categoryID,
							Name: categoryName,
						},
					},
					{
						ID:     productID2,
						Name:   productName2,
						Bought: true,
						Category: &categories.Category{
							ID:   categoryID,
							Name: categoryName,
						},
					},
				},
			},
		},
	}

	listResponse := builder.GetList(list)

	assert.DeepEqual(t, expectedListResponse, listResponse)
}
