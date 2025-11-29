package services

import (
	"testing"
)

func TestNewServiceContainer(t *testing.T) {
	// Test that service container can be created without panicking
	container := NewServiceContainer(nil)
	
	if container == nil {
		t.Error("Expected service container to be created")
	}
	
	// Test that all services are initialized
	if container.Auth == nil {
		t.Error("Expected Auth service to be initialized")
	}
	
	if container.User == nil {
		t.Error("Expected User service to be initialized")
	}
	
	if container.Service == nil {
		t.Error("Expected Service service to be initialized")
	}
}