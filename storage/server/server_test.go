// server_test.go
package server_test

import (
	"context"
	"net/http"
	"telegram_reminder_bot/server"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerRunAndShutdown(t *testing.T) {
	srv := server.Server{}

	go func() {
		err := srv.Run("8080", http.NewServeMux())
		assert.NoError(t, err)
	}()

	err := srv.Shutdown(context.Background())
	assert.NoError(t, err)
}
