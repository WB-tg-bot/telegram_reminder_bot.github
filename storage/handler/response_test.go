// response_test.go
package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"telegram_reminder_bot/handler"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/error", func(c *gin.Context) {
		handler.NewErrorResponse(c, http.StatusBadRequest, "this is an error")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/error", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"message":"this is an error"}`, w.Body.String())
}
