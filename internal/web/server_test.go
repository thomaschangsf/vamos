package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNewServer(t *testing.T) {
	server := NewServer(8080)
	assert.NotNil(t, server)
	assert.Equal(t, 8080, server.port)
	assert.NotNil(t, server.router)
}

func TestHealthCheck(t *testing.T) {
	server := NewServer(8080)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)

	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"ok"}`, w.Body.String())
}

func TestGetStatus(t *testing.T) {
	server := NewServer(8080)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/status", nil)

	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "time")
}

func TestStart(t *testing.T) {
	server := NewServer(8080)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := server.Start(ctx)
	assert.NoError(t, err)
}

func TestStartWithError(t *testing.T) {
	server := &Server{
		port:   -1, // Invalid port
		router: gin.Default(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := server.Start(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid port number")
}
