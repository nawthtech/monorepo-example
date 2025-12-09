package handlers

import (
	"net/http"
	"time"

	"nawthtech/worker/src/utils"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Service is healthy",
		"data": map[string]interface{}{
			"status":    "healthy",
			"database":  "D1",
			"timestamp": time.Now(),
			"service":   "nawthtech-worker",
		},
	})
}

func HealthLive(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"status":  "live",
	})
}

func HealthReady(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"status":  "ready",
	})
}