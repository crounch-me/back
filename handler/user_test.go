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
	expectedError         *model.Error
}

var validUser = &model.User{
	ID: "Hello",
}

func TestSignup(t *testing.T) {
	var validBody = `
		{
			"email": "test@test.com",
			"password": "test"
		}
	`
	testCases := []userSignupTestCases{
		{
			description:           "OK - Should create user",
			createUserStorageMock: createUserStorageMock{isCalled: true, err: nil},
			getUserStorageMock:    &getUserStorageMock{result: nil, err: model.NewDatabaseError(model.ErrNotFound, nil)},
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
			description:        "KO - missing email",
			expectedStatusCode: http.StatusBadRequest,
			requestBody: `
				{
					"password": "a"
				}
			`,
			expectedError: &model.Error{
				Code:        errorcode.InvalidCode,
				Description: fmt.Sprintf(errorcode.InvalidDescription, "Email", "required"),
			},
		},
		{
			description:        "KO - missing password",
			expectedStatusCode: http.StatusBadRequest,
			requestBody: `
				{
					"email": "test@test.com"
				}
			`,
			expectedError: &model.Error{
				Code:        errorcode.InvalidCode,
				Description: fmt.Sprintf(errorcode.InvalidDescription, "Password", "required"),
			},
		},
		{
			description:        "KO - password length too short",
			expectedStatusCode: http.StatusBadRequest,
			requestBody: `
				{
					"email": "test@test.com",
					"password": "a"
				}
			`,
			expectedError: &model.Error{
				Code:        errorcode.InvalidCode,
				Description: fmt.Sprintf(errorcode.InvalidDescription, "Password", "gt"),
			},
		},
		{
			description:        "KO - email is not an email",
			expectedStatusCode: http.StatusBadRequest,
			requestBody: `
				{
					"email": "a",
					"password": "a"
				}
			`,
			expectedError: &model.Error{
				Code:        errorcode.InvalidCode,
				Description: fmt.Sprintf(errorcode.InvalidDescription, "Email", "email"),
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
			description:           "KO - unknown database error when creating user",
			createUserStorageMock: createUserStorageMock{isCalled: true, err: errors.New("unknown database error")},
			getUserStorageMock:    &getUserStorageMock{result: nil, err: model.NewDatabaseError(model.ErrNotFound, nil)},
			requestBody:           validBody,
			expectedStatusCode:    http.StatusInternalServerError,
			expectedError: &model.Error{
				Code:        errorcode.DatabaseCode,
				Description: errorcode.DatabaseDescription,
			},
		},
		{
			description:        "KO - unknown database error when getting user",
			getUserStorageMock: &getUserStorageMock{result: nil, err: errors.New("unknown database error")},
			requestBody:        validBody,
			expectedStatusCode: http.StatusInternalServerError,
			expectedError: &model.Error{
				Code:        errorcode.DatabaseCode,
				Description: errorcode.DatabaseDescription,
			},
		},
		{
			description:        "KO - unknown database error when getting user",
			getUserStorageMock: &getUserStorageMock{result: nil, err: model.NewDatabaseError(model.ErrWrongPassword, nil)},
			requestBody:        validBody,
			expectedStatusCode: http.StatusInternalServerError,
			expectedError: &model.Error{
				Code:        errorcode.DatabaseCode,
				Description: errorcode.DatabaseDescription,
			},
		},
		{
			description:        "KO - duplicated users",
			getUserStorageMock: &getUserStorageMock{result: validUser, err: nil},
			requestBody:        validBody,
			expectedStatusCode: http.StatusConflict,
			expectedError: &model.Error{
				Code:        errorcode.Duplicate,
				Description: errorcode.DuplicateDescription,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer([]byte(tc.requestBody)))

			config := &configuration.Config{
				Mock: true,
			}
			hc := NewContext(config)
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
	expectedError                  *model.Error
}

func TestLogin(t *testing.T) {
	var validBody = `
		{
			"email": "test@test.com",
			"password": "test"
		}
	`
	testCases := []userLoginTestCases{
		{
			description: "OK - Should create authorization",
			createAuthorizationStorageMock: createAuthorizationStorageMock{isCalled: true, err: nil, result: &model.Authorization{
				AccessToken: "Hello",
			}},
			expectedStatusCode: http.StatusCreated,
			requestBody:        validBody,
			expectedBody: []Body{
				{
					Path: "$.accessToken",
					Data: "Hello",
				},
			},
		},
		{
			description:        "KO - missing email",
			expectedStatusCode: http.StatusBadRequest,
			requestBody: `
				{
					"password": "a"
				}
			`,
			expectedError: &model.Error{
				Code:        errorcode.InvalidCode,
				Description: fmt.Sprintf(errorcode.InvalidDescription, "Email", "required"),
			},
		},
		{
			description:        "KO - missing password",
			expectedStatusCode: http.StatusBadRequest,
			requestBody: `
				{
					"email": "test@test.com"
				}
			`,
			expectedError: &model.Error{
				Code:        errorcode.InvalidCode,
				Description: fmt.Sprintf(errorcode.InvalidDescription, "Password", "required"),
			},
		},
		{
			description:        "KO - password length too short",
			expectedStatusCode: http.StatusBadRequest,
			requestBody: `
				{
					"email": "test@test.com",
					"password": "a"
				}
			`,
			expectedError: &model.Error{
				Code:        errorcode.InvalidCode,
				Description: fmt.Sprintf(errorcode.InvalidDescription, "Password", "gt"),
			},
		},
		{
			description:        "KO - email is not an email",
			expectedStatusCode: http.StatusBadRequest,
			requestBody: `
				{
					"email": "a",
					"password": "a"
				}
			`,
			expectedError: &model.Error{
				Code:        errorcode.InvalidCode,
				Description: fmt.Sprintf(errorcode.InvalidDescription, "Email", "email"),
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
			description:                    "KO - wrong password",
			createAuthorizationStorageMock: createAuthorizationStorageMock{isCalled: true, err: model.NewDatabaseError(model.ErrWrongPassword, nil), result: nil},
			expectedStatusCode:             http.StatusBadRequest,
			requestBody:                    validBody,
			expectedError: &model.Error{
				Code:        errorcode.WrongPasswordCode,
				Description: errorcode.WrongPasswordDescription,
			},
		},
		{
			description:                    "KO - unknown database error",
			createAuthorizationStorageMock: createAuthorizationStorageMock{isCalled: true, err: errors.New("unknown database error"), result: nil},
			expectedStatusCode:             http.StatusInternalServerError,
			requestBody:                    validBody,
			expectedError: &model.Error{
				Code:        errorcode.DatabaseCode,
				Description: errorcode.DatabaseDescription,
			},
		},
		{
			description:                    "KO - unknown database error of type DatabaseError",
			createAuthorizationStorageMock: createAuthorizationStorageMock{isCalled: true, err: model.NewDatabaseError(model.ErrCreation, nil), result: nil},
			expectedStatusCode:             http.StatusInternalServerError,
			requestBody:                    validBody,
			expectedError: &model.Error{
				Code:        errorcode.DatabaseCode,
				Description: errorcode.DatabaseDescription,
			},
		},
		{
			description:                    "KO - user not found",
			createAuthorizationStorageMock: createAuthorizationStorageMock{isCalled: true, err: model.NewDatabaseError(model.ErrNotFound, nil), result: nil},
			expectedStatusCode:             http.StatusForbidden,
			requestBody:                    validBody,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(tc.requestBody)))

			config := &configuration.Config{
				Mock: true,
			}
			hc := NewContext(config)
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
