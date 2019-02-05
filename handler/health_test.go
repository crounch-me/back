package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type healthTestCases struct {
	description        string
	requestHTTPMethod  string
	expectedStatusCode int
}

func TestHealth(t *testing.T) {
	testCases := []healthTestCases{
		{
			description:        "Test GET should succeed",
			requestHTTPMethod:  http.MethodGet,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.requestHTTPMethod, "/health", nil)

			hc := NewContext()
			gin.SetMode(gin.TestMode)

			contextTest, _ := gin.CreateTestContext(w)
			contextTest.Request = req
			hc.Health(contextTest)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
		})
	}
}
