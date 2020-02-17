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
