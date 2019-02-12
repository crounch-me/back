package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sehsyha/crounch-back/model"
	storagemock "github.com/Sehsyha/crounch-back/storage/mock"

	"github.com/gin-gonic/gin"
	"github.com/oliveagle/jsonpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type createUserStorageMock struct {
	isCalled bool
	err      error
}

type getUserStorageMock struct {
	result *model.User
	err    error
}

type userSignupTestCases struct {
	createUserStorageMock createUserStorageMock
	getUserStorageMock    *getUserStorageMock
	description           string
	expectedStatusCode    int
	requestBody           string
	expectedBody          []Body
	expectedError         string
}

func TestSignup(t *testing.T) {
	var validBody = `
		{
			"email": "test@test.com",
			"pasword": "test"
		}
	`
	testCases := []userSignupTestCases{
		{
			createUserStorageMock: createUserStorageMock{isCalled: true, err: nil},
			getUserStorageMock:    &getUserStorageMock{result: nil, err: model.NewDatabaseError(model.ErrNotFound, nil)},
			description:           "OK - Should create user",
			expectedStatusCode:    http.StatusCreated,
			requestBody:           validBody,
			expectedBody: []Body{
				{
					Path: "$.email",
					Data: "test@test.com",
				},
			},
		},
		{
			createUserStorageMock: createUserStorageMock{isCalled: false, err: nil},
			description:           "KO - unmarshall error",
			expectedStatusCode:    http.StatusBadRequest,
			requestBody:           "",
			expectedError:         "unmarshall error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer([]byte(tc.requestBody)))

			hc := NewContext()
			gin.SetMode(gin.TestMode)

			contextTest, _ := gin.CreateTestContext(w)
			contextTest.Request = req

			storageMock := &storagemock.StorageMock{}

			storageMock.On("CreateUser", mock.Anything).Return(tc.createUserStorageMock.err)
			if tc.getUserStorageMock != nil {
				storageMock.On("GetUserByEmail", mock.Anything).Return(tc.getUserStorageMock.result, tc.getUserStorageMock.err)
			}

			hc.Storage = storageMock

			hc.Signup(contextTest)

			if tc.createUserStorageMock.isCalled {
				storageMock.AssertCalled(t, "CreateUser", mock.Anything)
			} else {
				storageMock.AssertNotCalled(t, "CreateUser", mock.Anything)
			}

			assert.Equal(t, tc.expectedStatusCode, w.Code)

			verify(t, tc.expectedBody, tc.expectedError, string(w.Body.Bytes()))
		})
	}
}

type createAuthorizationStorageMock struct {
	isCalled bool
	result   *model.Authorization
	err      error
}

type userLoginTestCases struct {
	createAuthorizationStorageMock createAuthorizationStorageMock
	description                    string
	expectedStatusCode             int
	requestBody                    string
	expectedBody                   []Body
	expectedError                  string
}

type Body struct {
	Path string
	Data string
}

func TestLogin(t *testing.T) {
	var validBody = `
		{
			"email": "test@test.com",
			"pasword": "test"
		}
	`
	testCases := []userLoginTestCases{
		{
			createAuthorizationStorageMock: createAuthorizationStorageMock{isCalled: true, err: nil, result: &model.Authorization{
				AccessToken: "Hello",
			}},
			description:        "OK - Should create authorization",
			expectedStatusCode: http.StatusCreated,
			requestBody:        validBody,
			expectedBody: []Body{
				{
					Path: "$.accessToken",
					Data: "Hello",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(tc.requestBody)))

			hc := NewContext()
			gin.SetMode(gin.TestMode)

			contextTest, _ := gin.CreateTestContext(w)
			contextTest.Request = req

			storageMock := &storagemock.StorageMock{}

			storageMock.On("CreateAuthorization", mock.Anything).Return(tc.createAuthorizationStorageMock.result, tc.createAuthorizationStorageMock.err)

			hc.Storage = storageMock

			hc.Login(contextTest)

			if tc.createAuthorizationStorageMock.isCalled {
				storageMock.AssertCalled(t, "CreateAuthorization", mock.Anything)
			} else {
				storageMock.AssertNotCalled(t, "CreateAuthorization", mock.Anything)
			}

			assert.Equal(t, tc.expectedStatusCode, w.Code)

			verify(t, tc.expectedBody, tc.expectedError, string(w.Body.Bytes()))
		})
	}
}

func verify(t *testing.T, expectedBody []Body, expectedError, actualBody string) {
	fmt.Println(actualBody)

	if expectedError != "" {
		path, err := jsonpath.Compile("$.error")
		assert.NoError(t, err)

		var actualData interface{}
		json.Unmarshal([]byte(actualBody), &actualData)
		actualError, _ := path.Lookup(actualData)

		assert.Equal(t, expectedError, actualError)
	} else {
		for _, eb := range expectedBody {
			pattern, err := jsonpath.Compile(eb.Path)
			assert.NoError(t, err)

			var actualData interface{}
			json.Unmarshal([]byte(actualBody), &actualData)
			actualDataExtracted, _ := pattern.Lookup(actualData)

			assert.Equal(t, eb.Data, actualDataExtracted)
		}
	}
}
