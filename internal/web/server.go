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
	router *gin.Engine
	port   int
}

// NewServer creates a new web server
func NewServer(port int) *Server {
	router := gin.Default()
	return &Server{
		router: router,
		port:   port,
	}
}

// Start starts the web server
func (s *Server) Start(ctx context.Context) error {
	// Add routes
	s.router.GET("/health", s.healthCheck)
	s.router.GET("/api/status", s.getStatus)

	// Create HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.router,
	}

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Shutdown server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return server.Shutdown(shutdownCtx)
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

func (s *Server) getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}
