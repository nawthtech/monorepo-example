package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"nawthtech/worker/src/utils"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	d1 := utils.GetD1()
	status, _ := d1.HealthCheck(ctx)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Service is " + status,
		"data": map[string]interface{}{
			"status":      status,
			"database":    "d1",
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
			"environment": os.Getenv("ENVIRONMENT"),
			"version":     os.Getenv("API_VERSION"),
			"service":     "nawthtech-worker",
		},
	})
}

func HealthReadyHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	d1 := utils.GetD1()
	status, err := d1.HealthCheck(ctx)
	if err != nil || status != "healthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "SERVICE_NOT_READY",
			"message": "Database is not ready",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Service is ready",
		"data": map[string]interface{}{
			"status":    "ready",
			"database":  "d1",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	})
}