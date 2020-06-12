package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/crounch-me/back/configuration"
)

type healthTestCases struct {
	description        string
	expectedStatusCode int
}

func TestHealth(t *testing.T) {
	testCases := []healthTestCases{
		{
			description:        "Test GET should succeed",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/health", nil)

			config := &configuration.Config{
				Mock: true,
			}
			hc := NewContext(config)
			gin.SetMode(gin.TestMode)

			contextTest, _ := gin.CreateTestContext(w)
			contextTest.Request = req
			hc.Health(contextTest)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
		})
	}
}
