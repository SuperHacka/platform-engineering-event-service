package app

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Set test environment variables
	os.Setenv("PORT", "9090")
	os.Setenv("ENV", "test")
	os.Setenv("PROCESSING_DELAY_MS", "500")

	config := LoadConfig()

	if config.Port != "9090" {
		t.Errorf("Expected Port 9090, got %s", config.Port)
	}
	if config.Env != "test" {
		t.Errorf("Expected Env test, got %s", config.Env)
	}
	if config.ProcessingDelayMs != 500 {
		t.Errorf("Expected ProcessingDelayMs 500, got %d", config.ProcessingDelayMs)
	}
}

func TestNew(t *testing.T) {
	config := Config{
		Port:              "8080",
		Env:               "dev",
		ProcessingDelayMs: 1000,
	}

	application := New(config)

	if application == nil {
		t.Fatal("Expected application instance, got nil")
	}

	if application.config.Port != "8080" {
		t.Errorf("Expected application to have port 8080, got %s", application.config.Port)
	}
}
