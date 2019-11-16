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
	userID := "user-id"
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
			description:           "KO - missing name",
			createListStorageMock: createListStorageMock{isCalled: false, err: nil},
			expectedStatusCode:    http.StatusBadRequest,
			requestBody: `
				{}
			`,
			expectedError: &model.Error{
				Code:        errorcode.InvalidCode,
				Description: fmt.Sprintf(errorcode.InvalidDescription, "Name", "required"),
			},
		},
		{
			description:           "KO - name length too long",
			createListStorageMock: createListStorageMock{isCalled: false, err: nil},
			expectedStatusCode:    http.StatusBadRequest,
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
			description:           "KO - unmarshall error",
			createListStorageMock: createListStorageMock{isCalled: false, err: nil},
			expectedStatusCode:    http.StatusBadRequest,
			requestBody:           "",
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
