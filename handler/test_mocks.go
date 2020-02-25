package handler

import "github.com/Sehsyha/crounch-back/model"

const (
	userIDMock    = "user-id"
	listIDMock    = "a33eb4fb-4df4-4392-a9fa-0a477663143b"
	productIDMock = "5ded8977-f4f8-4204-a300-19abf44613f2"
	otherUserID   = "other-user-id"
)

// LIST MOCKS

type createListStorageMock struct {
	isCalled bool
	err      error
}

type getListStorageMock struct {
	isCalled bool
	result   *model.List
	err      error
}

type getOwnerListsStorageMock struct {
	isCalled bool
	result   []*model.List
	err      error
}

type getProductInListMock struct {
	isCalled bool
	result   *model.ProductInList
	err      error
}

type addProductToListMock struct {
	isCalled bool
	err      error
}

var notOwnerListMock = &model.List{
	ID: productIDMock,
	Owner: &model.User{
		ID: otherUserID,
	},
}

var ownerListMock = &model.List{
	ID: productIDMock,
	Owner: &model.User{
		ID: userIDMock,
	},
}

var productInListMock = &model.ProductInList{
	ProductID: productIDMock,
	ListID:    listIDMock,
}

// PRODUCT MOCKS

type createProductStorageMock struct {
	isCalled bool
	err      error
}

type getProductMock struct {
	isCalled bool
	result   *model.Product
	err      error
}

var notOwnerProductMock = &model.Product{
	ID: productIDMock,
	Owner: &model.User{
		ID: otherUserID,
	},
}

var ownerProductMock = &model.Product{
	ID: productIDMock,
	Owner: &model.User{
		ID: userIDMock,
	},
}

// USER MOCKS

type createAuthorizationStorageMock struct {
	isCalled bool
	result   *model.Authorization
	err      error
}

type createUserStorageMock struct {
	isCalled bool
	err      error
}

type getUserStorageMock struct {
	result *model.User
	err    error
}
