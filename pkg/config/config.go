package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	// AWS Configuration
	AWSRegion    string
	AWSAccessKey string
	AWSSecretKey string

	// LLM Configuration
	LLMAPIKey    string
	LLMModelName string

	// Web Configuration
	WebPort    int
	WebBaseURL string
}

// NewConfig creates a new configuration instance
func NewConfig() *Config {
	return &Config{
		AWSRegion:    getEnvOrDefault("AWS_REGION", "us-west-2"),
		AWSAccessKey: getEnvOrDefault("AWS_ACCESS_KEY", ""),
		AWSSecretKey: getEnvOrDefault("AWS_SECRET_KEY", ""),
		LLMAPIKey:    getEnvOrDefault("LLM_API_KEY", ""),
		LLMModelName: getEnvOrDefault("LLM_MODEL_NAME", "gpt-3.5-turbo"),
		WebPort:      8080,
		WebBaseURL:   getEnvOrDefault("WEB_BASE_URL", "http://localhost:8080"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
