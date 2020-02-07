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

type createListStorageMock struct {
	isCalled bool
	err      error
}

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
				Description: fmt.Sprintf(errorcode.InvalidDescription, "Name", "required"),
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
				Description: fmt.Sprintf(errorcode.InvalidDescription, "Name", "lt"),
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
				contextTest.Set(ContextUserID, userID)
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

type getOwnerListsStorageMock struct {
	isCalled bool
	result   []*model.List
	err      error
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
				contextTest.Set(ContextUserID, userID)
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

type addOFFProductListMock struct {
	isCalled bool
	err      error
}

type getListMock struct {
	isCalled bool
	result   *model.List
	err      error
}

type addOFFProductListTestCases struct {
	addOFFProductListMock addOFFProductListMock
	getListMock           getListMock
	description           string
	expectedStatusCode    int
	expectedBody          []Body
	requestBody           string
	expectedError         *model.Error
	noContext             bool
}

func TestAddOFFProduct(t *testing.T) {
	id := "product-id"
	code := "1234567890123"
	validBody := fmt.Sprintf(`
    {
      "code": "%s"
    }
  `, code)
	testCases := []addOFFProductListTestCases{
		{
			description:        "KO - error when retrieving user id from context",
			noContext:          true,
			expectedStatusCode: http.StatusInternalServerError,
			expectedError: &model.Error{
				Code:        errorcode.UserDataCode,
				Description: errorcode.UserDataDescription,
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
			description:        "KO - missing code",
			expectedStatusCode: http.StatusBadRequest,
			requestBody: `
				{}
			`,
			expectedError: &model.Error{
				Code:        errorcode.InvalidCode,
				Description: fmt.Sprintf(errorcode.InvalidDescription, "Code", "required"),
			},
		},
		{
			description:        "KO - code length too long",
			expectedStatusCode: http.StatusBadRequest,
			requestBody: `
				{
					"code": "12345678901234"
				}
			`,
			expectedError: &model.Error{
				Code:        errorcode.InvalidCode,
				Description: fmt.Sprintf(errorcode.InvalidDescription, "Code", "len"),
			},
		},
		{
			description:        "KO - code length too short",
			expectedStatusCode: http.StatusBadRequest,
			requestBody: `
				{
					"code": "123456789012"
				}
			`,
			expectedError: &model.Error{
				Code:        errorcode.InvalidCode,
				Description: fmt.Sprintf(errorcode.InvalidDescription, "Code", "len"),
			},
		},
		{
			description:        "KO - error when searching list",
			expectedStatusCode: http.StatusInternalServerError,
			requestBody:        validBody,
			getListMock:        getListMock{err: errors.New("unknown"), isCalled: true},
			expectedError: &model.Error{
				Code:        errorcode.DatabaseCode,
				Description: errorcode.DatabaseDescription,
			},
		},
		{
			description:        "KO - list not found",
			expectedStatusCode: http.StatusNotFound,
			requestBody:        validBody,
			getListMock:        getListMock{err: model.NewDatabaseError(model.ErrNotFound, nil), isCalled: true},
			expectedError: &model.Error{
				Code:        errorcode.NotFoundCode,
				Description: ListNotFoundDescription,
			},
		},
		{
			description:        "KO - user is not owner",
			expectedStatusCode: http.StatusForbidden,
			requestBody:        validBody,
			getListMock: getListMock{isCalled: true, result: &model.List{
				Owner: &model.User{
					ID: otherUserID,
				},
			}},
		},
		{
			description:        "KO - unknown error while adding product to the list",
			expectedStatusCode: http.StatusInternalServerError,
			getListMock: getListMock{isCalled: true, result: &model.List{
				Owner: &model.User{
					ID: userID,
				},
			}},
			requestBody:           validBody,
			addOFFProductListMock: addOFFProductListMock{isCalled: true, err: errors.New("unknown")},
			expectedError: &model.Error{
				Code:        errorcode.DatabaseCode,
				Description: errorcode.DatabaseDescription,
			},
		},
		{
			description:        "OK - product added to list",
			expectedStatusCode: http.StatusCreated,
			getListMock: getListMock{isCalled: true, result: &model.List{
				Owner: &model.User{
					ID: userID,
				},
			}},
			requestBody:           validBody,
			addOFFProductListMock: addOFFProductListMock{isCalled: true},
			expectedBody: []Body{
				{
					Path: "$.code",
					Data: code,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/lists/%s/offproducts", id), bytes.NewBuffer([]byte(tc.requestBody)))

			config := &configuration.Config{
				Mock: true,
			}
			hc := NewContext(config)
			gin.SetMode(gin.TestMode)

			contextTest, _ := gin.CreateTestContext(w)
			contextTest.Request = req

			if !tc.noContext {
				contextTest.Set(ContextUserID, userID)
			}

			storageMock := &storagemock.StorageMock{}

			storageMock.On("GetList", mock.Anything).Return(tc.getListMock.result, tc.getListMock.err)
			storageMock.On("AddOFFProductToList", mock.Anything, mock.Anything).Return(tc.addOFFProductListMock.err)

			hc.Storage = storageMock

			hc.AddOFFProduct(contextTest)

			if tc.getListMock.isCalled {
				storageMock.AssertCalled(t, "GetList", mock.Anything)
			} else {
				storageMock.AssertNotCalled(t, "GetList", mock.Anything)
			}

			if tc.addOFFProductListMock.isCalled {
				storageMock.AssertCalled(t, "AddOFFProductToList", mock.Anything, mock.Anything)
			} else {
				storageMock.AssertNotCalled(t, "AddOFFProductToList", mock.Anything, mock.Anything)
			}

			assert.Equal(t, tc.expectedStatusCode, w.Code)

			verify(t, tc.expectedBody, tc.expectedError, string(w.Body.Bytes()))
		})
	}
}
