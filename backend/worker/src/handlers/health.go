package health

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
	"nawthtech-worker/utils"
)

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDatabase()
	dbStatus := db.HealthCheck()

	resp := map[string]interface{}{
		"success": true,
		"message": "Service is " + dbStatus,
		"data": map[string]interface{}{
			"status":     dbStatus,
			"database":   "D1",
			"timestamp":  time.Now().Format(time.RFC3339),
			"environment": os.Getenv("ENVIRONMENT"),
			"version":    os.Getenv("API_VERSION"),
			"service":    "nawthtech-worker",
		},
	}

	json.NewEncoder(w).Encode(resp)
}

func ReadyHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDatabase()
	dbStatus := db.HealthCheck()

	if dbStatus != "healthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "SERVICE_NOT_READY",
			"message": "Database is not ready",
		})
		return
	}

	resp := map[string]interface{}{
		"success": true,
		"message": "Service is ready",
		"data": map[string]interface{}{
			"status":    "ready",
			"database":  "D1",
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}

	json.NewEncoder(w).Encode(resp)
}