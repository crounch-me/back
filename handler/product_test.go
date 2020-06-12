package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crounch-me/back/configuration"
	"github.com/crounch-me/back/errorcode"
	"github.com/crounch-me/back/model"
	storagemock "github.com/crounch-me/back/storage/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type productCreateTestCases struct {
	createProductStorageMock createProductStorageMock
	description              string
	expectedStatusCode       int
	requestBody              string
	expectedBody             []Body
	expectedError            *model.Error
	noContext                bool
}

func TestCreateProduct(t *testing.T) {
	validBody := `
		{
			"name": "Mon produit"
		}
	`
	testCases := []productCreateTestCases{
		{
			description:              "OK - Should create product",
			createProductStorageMock: createProductStorageMock{isCalled: true, err: nil},
			expectedStatusCode:       http.StatusCreated,
			requestBody:              validBody,
			expectedBody: []Body{
				{
					Path: "$.name",
					Data: "Mon produit",
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
			description:              "KO - unknown database error when creating product",
			createProductStorageMock: createProductStorageMock{isCalled: true, err: errors.New("unknown database error")},
			requestBody:              validBody,
			expectedStatusCode:       http.StatusInternalServerError,
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
			req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer([]byte(tc.requestBody)))

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

			storageMock.On("CreateProduct", mock.Anything).Return(tc.createProductStorageMock.err)

			hc.Storage = storageMock

			hc.CreateProduct(contextTest)

			if tc.createProductStorageMock.isCalled {
				storageMock.AssertCalled(t, "CreateProduct", mock.Anything)
			} else {
				storageMock.AssertNotCalled(t, "CreateProduct", mock.Anything)
			}

			assert.Equal(t, tc.expectedStatusCode, w.Code)

			verify(t, tc.expectedBody, tc.expectedError, string(w.Body.Bytes()))
		})
	}
}
