package config

import (
	"testing"
)

func TestConfigDefaults(t *testing.T) {
	t.Skip("Skipping test - requires environment setup")
	
	cfg := Load()
	if cfg == nil {
		t.Fatal("Expected config to be loaded")
	}
	
	// Test basic defaults
	if cfg.Port == "" {
		t.Error("Expected default port to be set")
	}
	
	if cfg.Environment == "" {
		t.Error("Expected default environment to be set")
	}
	
	if cfg.Version == "" {
		t.Error("Expected version to be set")
	}
}

func TestDatabaseConfig(t *testing.T) {
	t.Skip("Skipping test - requires environment setup")
	
	cfg := Load()
	if cfg == nil {
		t.Fatal("Config should be loaded")
	}
	
	if cfg.Database.URL == "" {
		t.Error("Database URL should be set")
	}
	
	if cfg.Database.Name == "" {
		t.Error("Database name should be set")
	}
}

func TestIsProduction(t *testing.T) {
	t.Skip("Skipping test - requires environment setup")
	
	cfg := Load()
	if cfg == nil {
		t.Fatal("Config should be loaded")
	}
	
	// Test environment detection
	if cfg.Environment == "production" {
		if !cfg.IsProduction() {
			t.Error("IsProduction should return true for production environment")
		}
	} else {
		if cfg.IsProduction() {
			t.Error("IsProduction should return false for non-production environment")
		}
	}
}

func TestIsDevelopment(t *testing.T) {
	t.Skip("Skipping test - requires environment setup")
	
	cfg := Load()
	if cfg == nil {
		t.Fatal("Config should be loaded")
	}
	
	// Test environment detection
	if cfg.Environment == "development" {
		if !cfg.IsDevelopment() {
			t.Error("IsDevelopment should return true for development environment")
		}
	} else {
		if cfg.IsDevelopment() && cfg.Environment != "development" {
			t.Error("IsDevelopment should return false for non-development environment")
		}
	}
}