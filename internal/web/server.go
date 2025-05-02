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
	router.POST("/handle-request", server.handleRequest)
	router.GET("/api/weather", server.getWeather)

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

// handleRequest handles the POST request with custom headers and JSON payload
func (s *Server) handleRequest(c *gin.Context) {
	// Extract headers
	tenantID := c.GetHeader("hawking-tenant-id")
	appContext := c.GetHeader("x-sfdc-app-context")
	clientTraceID := c.GetHeader("x-client-trace-id")

	// Validate required headers
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hawking-tenant-id header is required"})
		return
	}

	// Parse request body
	var requestBody struct {
		Namespace   string `json:"namespace"`
		Tenant      string `json:"tenant"`
		ChildFolder string `json:"childFolder"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate request body
	if requestBody.Namespace == "" || requestBody.Tenant == "" || requestBody.ChildFolder == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "namespace, tenant, and childFolder are required"})
		return
	}

	// Process the request (you can add your business logic here)
	response := gin.H{
		"status":      "success",
		"tenant_id":   tenantID,
		"app_context": appContext,
		"trace_id":    clientTraceID,
		"request":     requestBody,
	}

	c.JSON(http.StatusOK, response)
}

// WeatherResponse represents the weather forecast response
type WeatherResponse struct {
	City    string         `json:"city"`
	State   string         `json:"state"`
	Forecast []DailyForecast `json:"forecast"`
}

// DailyForecast represents the weather forecast for a single day
type DailyForecast struct {
	Date        string  `json:"date"`
	Temperature struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"temperature"`
	Description string  `json:"description"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed"`
}

// getWeather handles the weather forecast endpoint
func (s *Server) getWeather(c *gin.Context) {
	city := c.Query("city")
	state := c.Query("state")

	if city == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "city and state parameters are required",
		})
		return
	}

	// TODO: Replace with actual API call to a weather service
	// This is a mock response for demonstration
	response := WeatherResponse{
		City:  city,
		State: state,
		Forecast: []DailyForecast{
			{
				Date: time.Now().Format("2006-01-02"),
				Temperature: struct {
					Min float64 `json:"min"`
					Max float64 `json:"max"`
				}{
					Min: 15.0,
					Max: 25.0,
				},
				Description: "Sunny",
				Humidity:    65,
				WindSpeed:   5.5,
			},
			// Add more days as needed
		},
	}

	c.JSON(http.StatusOK, response)
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
