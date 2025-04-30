package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Server represents a web server
type Server struct {
	port   int
	router *gin.Engine
}

// NewServer creates a new web server
func NewServer(port int) *Server {
	router := gin.Default()
	server := &Server{
		port:   port,
		router: router,
	}

	// Register routes
	router.GET("/health", server.healthCheck)
	router.GET("/api/status", server.getStatus)

	return server
}

// healthCheck handles the health check endpoint
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// getStatus handles the status endpoint
func (s *Server) getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// Start starts the web server
func (s *Server) Start(ctx context.Context) error {
	if s.port < 0 || s.port > 65535 {
		return fmt.Errorf("invalid port number: %d", s.port)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.router,
	}

	errChan := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}
