package config

import (
 "fmt"
	"os"
	"strconv"
 "strings"
 "time"

 "github.com/caarlos0/env/v11"
)

func TestLoad(t *testing.T) {
	// حفظ الإعدادات الحالية
	originalEnv := os.Getenv("ENVIRONMENT")
	originalDB := os.Getenv("DATABASE_URL")
	defer func() {
		os.Setenv("ENVIRONMENT", originalEnv)
		os.Setenv("DATABASE_URL", originalDB)
	}()
	
	// تعيين بيئة اختبار
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test")
	
	// إعادة تعيين appConfig لفرض إعادة التحميل
	appConfig = nil
	
	cfg := Load()
	
	if cfg == nil {
		t.Fatal("Expected config to be loaded, got nil")
	}
	
	if cfg.Environment != "test" {
		t.Errorf("Expected environment 'test', got '%s'", cfg.Environment)
	}
	
	if cfg.Port == "" {
		t.Error("Expected port to be set")
	}
	
	if cfg.Version == "" {
		t.Error("Expected version to be set")
	}
}

func TestIsDevelopment(t *testing.T) {
	cfg := &Config{Environment: "development"}
	
	if !cfg.IsDevelopment() {
		t.Error("Expected IsDevelopment to return true for 'development' environment")
	}
	
	cfg.Environment = "production"
	if cfg.IsDevelopment() {
		t.Error("Expected IsDevelopment to return false for 'production' environment")
	}
}

func TestIsProduction(t *testing.T) {
	cfg := &Config{Environment: "production"}
	
	if !cfg.IsProduction() {
		t.Error("Expected IsProduction to return true for 'production' environment")
	}
	
	cfg.Environment = "development"
	if cfg.IsProduction() {
		t.Error("Expected IsProduction to return false for 'development' environment")
	}
}

func TestGetDSN(t *testing.T) {
	cfg := &Config{
		Database: DatabaseConfig{
			URL: "postgres://user:pass@localhost:5432/db",
		},
	}
	
	dsn := cfg.GetDSN()
	if dsn != "postgres://user:pass@localhost:5432/db" {
		t.Errorf("Expected DSN to match, got '%s'", dsn)
	}
}