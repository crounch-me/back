package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	storagemock "github.com/Sehsyha/crounch-back/storage/mock"

	"github.com/gin-gonic/gin"
	"github.com/oliveagle/jsonpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type userSignupTestCases struct {
	createUserStorageMock createUserStorageMock
	description           string
	expectedStatusCode    int
	requestBody           string
	expectedBody          string
	expectedError         string
}

type createUserStorageMock struct {
	isCalled bool
	err      error
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
			description:           "OK - Should create user",
			expectedStatusCode:    http.StatusCreated,
			requestBody:           validBody,
			expectedBody: `
				{
					"email": "test@test.com"
				}
			`,
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

			hc.Storage = storageMock

			hc.Signup(contextTest)

			if tc.createUserStorageMock.isCalled {
				storageMock.AssertCalled(t, "CreateUser", mock.Anything)
			} else {
				storageMock.AssertNotCalled(t, "CreateUser", mock.Anything)
			}

			assert.Equal(t, tc.expectedStatusCode, w.Code)

			if tc.expectedError != "" {

				fmt.Println(string(w.Body.Bytes()))

				pattern, err := jsonpath.Compile("$.error")
				assert.NoError(t, err)

				var actualData interface{}
				json.Unmarshal([]byte(w.Body.Bytes()), &actualData)
				actualError, _ := pattern.Lookup(actualData)

				assert.Equal(t, tc.expectedError, actualError)
			} else {
				pattern, err := jsonpath.Compile("$.email")
				assert.NoError(t, err)

				var expectedData interface{}
				json.Unmarshal([]byte(tc.expectedBody), &expectedData)
				expectedEmail, _ := pattern.Lookup(expectedData)

				var actualData interface{}
				json.Unmarshal([]byte(w.Body.Bytes()), &actualData)
				actualEmail, _ := pattern.Lookup(actualData)

				assert.Equal(t, expectedEmail, actualEmail)
			}
		})
	}
}
