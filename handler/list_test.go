package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Sehsyha/crounch-back/configuration"
	"github.com/Sehsyha/crounch-back/errorcode"
	"github.com/Sehsyha/crounch-back/model"
	storagemock "github.com/Sehsyha/crounch-back/storage/mock"
)

type listCreateTestCases struct {
	createListStorageMock createListStorageMock
	description           string
	expectedStatusCode    int
	requestBody           string
	expectedBody          []Body
	expectedError         *model.Error
	noContext             bool
}

func TestCreateList(t *testing.T) {
	validBody := `
		{
			"name": "Ma liste de course"
		}
	`
	testCases := []listCreateTestCases{
		{
			description:           "OK - Should create list",
			createListStorageMock: createListStorageMock{isCalled: true, err: nil},
			expectedStatusCode:    http.StatusCreated,
			requestBody:           validBody,
			expectedBody: []Body{
				{
					Path: "$.name",
					Data: "Ma liste de course",
				},
			},
		},
		{
			description:        "KO - missing name",
			expectedStatusCode: http.StatusBadRequest,
			requestBody: `
				{}
			`,
			expectedError: &model.Error{
				Code:        errorcode.InvalidCode,
				Description: fmt.Sprintf(errorcode.InvalidDescription, "Name", "required", ""),
			},
		},
		{
			description:        "KO - name length too long",
			expectedStatusCode: http.StatusBadRequest,
			requestBody: `
				{
					"name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
				}
			`,
			expectedError: &model.Error{
				Code:        errorcode.InvalidCode,
				Description: fmt.Sprintf(errorcode.InvalidDescription, "Name", "lt", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
			},
		},
		{
			description:        "KO - unmarshall error",
			expectedStatusCode: http.StatusBadRequest,
			requestBody:        "",
			expectedError: &model.Error{
				Code:        errorcode.UnmarshalCode,
				Description: errorcode.UnmarshalDescription,
			},
		},
		{
			description:           "KO - unknown database error when creating list",
			createListStorageMock: createListStorageMock{isCalled: true, err: errors.New("unknown database error")},
			requestBody:           validBody,
			expectedStatusCode:    http.StatusInternalServerError,
			expectedError: &model.Error{
				Code:        errorcode.DatabaseCode,
				Description: errorcode.DatabaseDescription,
			},
		},
		{
			description:        "KO - error when retrieving user id from context",
			noContext:          true,
			requestBody:        validBody,
			expectedStatusCode: http.StatusInternalServerError,
			expectedError: &model.Error{
				Code:        errorcode.UserDataCode,
				Description: errorcode.UserDataDescription,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/list", bytes.NewBuffer([]byte(tc.requestBody)))

			config := &configuration.Config{
				Mock: true,
			}
			hc := NewContext(config)
			gin.SetMode(gin.TestMode)

			contextTest, _ := gin.CreateTestContext(w)
			contextTest.Request = req

			if !tc.noContext {
				contextTest.Set(ContextUserID, userIDMock)
			}

			storageMock := &storagemock.StorageMock{}

			storageMock.On("CreateList", mock.Anything).Return(tc.createListStorageMock.err)

			hc.Storage = storageMock

			hc.CreateList(contextTest)

			if tc.createListStorageMock.isCalled {
				storageMock.AssertCalled(t, "CreateList", mock.Anything)
			} else {
				storageMock.AssertNotCalled(t, "CreateList", mock.Anything)
			}

			assert.Equal(t, tc.expectedStatusCode, w.Code)

			verify(t, tc.expectedBody, tc.expectedError, string(w.Body.Bytes()))
		})
	}
}

type listOwnerGetTestCases struct {
	getOwnerListsStorageMock getOwnerListsStorageMock
	description              string
	expectedStatusCode       int
	expectedBody             []Body
	expectedError            *model.Error
	noContext                bool
}

func TestGetOwnerLists(t *testing.T) {
	lists := []*model.List{
		&model.List{
			ID:   "list id",
			Name: "list name",
		},
	}
	testCases := []listOwnerGetTestCases{
		{
			description:              "KO - error when retrieving user id from context",
			noContext:                true,
			expectedStatusCode:       http.StatusInternalServerError,
			getOwnerListsStorageMock: getOwnerListsStorageMock{},
			expectedError: &model.Error{
				Code:        errorcode.UserDataCode,
				Description: errorcode.UserDataDescription,
			},
		},
		{
			description:              "KO - error when retrieving lists from database",
			expectedStatusCode:       http.StatusInternalServerError,
			getOwnerListsStorageMock: getOwnerListsStorageMock{isCalled: true, err: errors.New("unknown")},
			expectedError: &model.Error{
				Code:        errorcode.DatabaseCode,
				Description: errorcode.DatabaseDescription,
			},
		},
		{
			description:              "OK - Return owner lists",
			expectedStatusCode:       http.StatusOK,
			getOwnerListsStorageMock: getOwnerListsStorageMock{isCalled: true, result: lists},
			expectedBody: []Body{
				{
					Path: "$.id[0]",
					Data: "list id",
				},
				{
					Path: "$.name[0]",
					Data: "list name",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/lists", nil)

			config := &configuration.Config{
				Mock: true,
			}
			hc := NewContext(config)
			gin.SetMode(gin.TestMode)

			contextTest, _ := gin.CreateTestContext(w)
			contextTest.Request = req

			if !tc.noContext {
				contextTest.Set(ContextUserID, userIDMock)
			}

			storageMock := &storagemock.StorageMock{}

			storageMock.On("GetOwnerLists", mock.Anything).Return(tc.getOwnerListsStorageMock.result, tc.getOwnerListsStorageMock.err)

			hc.Storage = storageMock

			hc.GetOwnerLists(contextTest)

			if tc.getOwnerListsStorageMock.isCalled {
				storageMock.AssertCalled(t, "GetOwnerLists", mock.Anything)
			} else {
				storageMock.AssertNotCalled(t, "GetOwnerLists", mock.Anything)
			}

			assert.Equal(t, tc.expectedStatusCode, w.Code)

			verify(t, tc.expectedBody, tc.expectedError, string(w.Body.Bytes()))
		})
	}
}

type productToListPostTestCases struct {
	getListMock          getListStorageMock
	getProductMock       getProductMock
	getProductInListMock getProductInListMock
	addProductToListMock addProductToListMock
	description          string
	expectedStatusCode   int
	listID               string
	productID            string
	expectedBody         []Body
	expectedError        *model.Error
	noContext            bool
}

func TestAddProductToList(t *testing.T) {
	testCases := []productToListPostTestCases{
		{
			description:          "OK",
			listID:               listIDMock,
			productID:            productIDMock,
			getListMock:          getListStorageMock{isCalled: true, result: ownerListMock},
			getProductMock:       getProductMock{isCalled: true, result: ownerProductMock},
			getProductInListMock: getProductInListMock{isCalled: true, err: model.NewDatabaseError(model.ErrNotFound, nil)},
			addProductToListMock: addProductToListMock{isCalled: true},
			expectedStatusCode:   http.StatusCreated,
			expectedBody: []Body{
				{
					Path: "$.product_id",
					Data: productIDMock,
				},
				{
					Path: "$.list_id",
					Data: listIDMock,
				},
			},
		},
		{
			description:          "KO - unknown error when adding product to the list",
			listID:               listIDMock,
			productID:            productIDMock,
			getListMock:          getListStorageMock{isCalled: true, result: ownerListMock},
			getProductMock:       getProductMock{isCalled: true, result: ownerProductMock},
			getProductInListMock: getProductInListMock{isCalled: true, err: model.NewDatabaseError(model.ErrNotFound, nil)},
			addProductToListMock: addProductToListMock{isCalled: true, err: errors.New("unknown")},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedError: &model.Error{
				Code:        errorcode.DatabaseCode,
				Description: errorcode.DatabaseDescription,
			},
		},
		{
			description:          "KO - product already added in list",
			listID:               listIDMock,
			productID:            productIDMock,
			getListMock:          getListStorageMock{isCalled: true, result: ownerListMock},
			getProductMock:       getProductMock{isCalled: true, result: ownerProductMock},
			getProductInListMock: getProductInListMock{isCalled: true, result: productInListMock},
			expectedStatusCode:   http.StatusConflict,
			expectedError: &model.Error{
				Code:        errorcode.DuplicateCode,
				Description: errorcode.DuplicateDescription,
			},
		},
		{
			description:          "KO - unknwon error when getting product in list",
			listID:               listIDMock,
			productID:            productIDMock,
			getListMock:          getListStorageMock{isCalled: true, result: ownerListMock},
			getProductMock:       getProductMock{isCalled: true, result: ownerProductMock},
			getProductInListMock: getProductInListMock{isCalled: true, err: errors.New("unknown")},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedError: &model.Error{
				Code:        errorcode.DatabaseCode,
				Description: errorcode.DatabaseDescription,
			},
		},
		{
			description:        "KO - user not owner of product",
			listID:             listIDMock,
			productID:          productIDMock,
			getListMock:        getListStorageMock{isCalled: true, result: ownerListMock},
			getProductMock:     getProductMock{isCalled: true, result: notOwnerProductMock},
			expectedStatusCode: http.StatusForbidden,
		},
		{
			description:        "KO - product not found",
			listID:             listIDMock,
			productID:          productIDMock,
			getListMock:        getListStorageMock{isCalled: true, result: ownerListMock},
			getProductMock:     getProductMock{isCalled: true, err: model.NewDatabaseError(model.ErrNotFound, nil)},
			expectedStatusCode: http.StatusNotFound,
			expectedError: &model.Error{
				Code:        errorcode.NotFoundCode,
				Description: ProductNotFoundDescription,
			},
		},
		{
			description:        "KO - unknown error while getting product",
			listID:             listIDMock,
			productID:          productIDMock,
			getListMock:        getListStorageMock{isCalled: true, result: ownerListMock},
			getProductMock:     getProductMock{isCalled: true, err: errors.New("unknwon")},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError: &model.Error{
				Code:        errorcode.DatabaseCode,
				Description: errorcode.DatabaseDescription,
			},
		},
		{
			description:        "KO - user not owner of list",
			listID:             listIDMock,
			productID:          productIDMock,
			getListMock:        getListStorageMock{isCalled: true, result: notOwnerListMock},
			expectedStatusCode: http.StatusForbidden,
		},
		{
			description:        "KO - list not found",
			listID:             listIDMock,
			productID:          productIDMock,
			getListMock:        getListStorageMock{isCalled: true, err: model.NewDatabaseError(model.ErrNotFound, nil)},
			expectedStatusCode: http.StatusNotFound,
			expectedError: &model.Error{
				Code:        errorcode.NotFoundCode,
				Description: ListNotFoundDescription,
			},
		},
		{
			description:        "KO - unknown error while getting list",
			listID:             listIDMock,
			productID:          productIDMock,
			getListMock:        getListStorageMock{isCalled: true, err: errors.New("unknown")},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError: &model.Error{
				Code:        errorcode.DatabaseCode,
				Description: errorcode.DatabaseDescription,
			},
		},
		{
			description:        "KO - product id is not an uuid",
			listID:             listIDMock,
			productID:          "",
			expectedStatusCode: http.StatusBadRequest,
			expectedError: &model.Error{
				Code:        errorcode.InvalidCode,
				Description: fmt.Sprintf(errorcode.InvalidDescription, "product ID", "uuid", ""),
			},
		},
		{
			description:        "KO - list id is not an uuid",
			expectedStatusCode: http.StatusBadRequest,
			listID:             "",
			productID:          "",
			expectedError: &model.Error{
				Code:        errorcode.InvalidCode,
				Description: fmt.Sprintf(errorcode.InvalidDescription, "list ID", "uuid", ""),
			},
		},
		{
			description:        "KO - error when retrieving user id from context",
			noContext:          true,
			expectedStatusCode: http.StatusInternalServerError,
			expectedError: &model.Error{
				Code:        errorcode.UserDataCode,
				Description: errorcode.UserDataDescription,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/lists/%s/products/%s", tc.listID, tc.productID), nil)

			config := &configuration.Config{
				Mock: true,
			}

			hc := NewContext(config)
			gin.SetMode(gin.TestMode)

			contextTest, _ := gin.CreateTestContext(w)

			contextTest.Request = req

			contextTest.Params = []gin.Param{
				{Key: "listID", Value: tc.listID},
				{Key: "productID", Value: tc.productID},
			}

			if !tc.noContext {
				contextTest.Set(ContextUserID, userIDMock)
			}

			storageMock := &storagemock.StorageMock{}
			storageMock.On("GetList", mock.Anything).Return(tc.getListMock.result, tc.getListMock.err)
			storageMock.On("GetProductInList", mock.Anything, mock.Anything).Return(tc.getProductInListMock.result, tc.getProductInListMock.err)
			storageMock.On("GetProduct", mock.Anything).Return(tc.getProductMock.result, tc.getProductMock.err)
			storageMock.On("AddProductToList", mock.Anything, mock.Anything).Return(tc.addProductToListMock.err)
			hc.Storage = storageMock

			hc.AddProductToList(contextTest)

			assert.Equal(t, tc.expectedStatusCode, w.Code)

			verify(t, tc.expectedBody, tc.expectedError, string(w.Body.Bytes()))

			if tc.getListMock.isCalled {
				storageMock.AssertCalled(t, "GetList", tc.listID)
			} else {
				storageMock.AssertNotCalled(t, "GetList")
			}

			if tc.getProductMock.isCalled {
				storageMock.AssertCalled(t, "GetProduct", tc.productID)
			} else {
				storageMock.AssertNotCalled(t, "GetProduct")
			}

			if tc.getProductInListMock.isCalled {
				storageMock.AssertCalled(t, "GetProductInList", tc.productID, tc.listID)
			} else {
				storageMock.AssertNotCalled(t, "GetProductInList")
			}

			if tc.addProductToListMock.isCalled {
				storageMock.AssertCalled(t, "AddProductToList", tc.productID, tc.listID)
			} else {
				storageMock.AssertNotCalled(t, "AddProductToList")
			}
		})
	}
}
