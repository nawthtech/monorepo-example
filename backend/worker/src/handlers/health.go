package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"nawthtech/utils"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	dbHealth := utils.CheckDBHealth()

	data := map[string]interface{}{
		"status":     dbHealth.Status,
		"database":   dbHealth.Type,
		"timestamp":  time.Now().UTC(),
		"environment": os.Getenv("ENVIRONMENT"),
		"version":     os.Getenv("API_VERSION"),
		"service":     "nawthtech-worker",
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Service is " + dbHealth.Status,
		"data":    data,
	})
}

func HealthReady(w http.ResponseWriter, r *http.Request) {
	dbHealth := utils.CheckDBHealth()

	if dbHealth.Status != "healthy" {
		respondJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
			"success": false,
			"error":   "SERVICE_NOT_READY",
			"message": "Database is not ready",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Service is ready",
		"data": map[string]interface{}{
			"status":    "ready",
			"database":  dbHealth.Type,
			"timestamp": time.Now().UTC(),
		},
	})
}

// ================= Helper
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}