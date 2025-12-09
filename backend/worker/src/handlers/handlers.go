package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"worker/src/middleware"
	"worker/src/utils"
)

// ========== Helper ==========

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// ========== Health Handlers ==========

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	dbHealth := utils.DB.HealthCheck() // D1 check

	data := map[string]interface{}{
		"status":      dbHealth.Status,
		"database":    dbHealth.Type,
		"timestamp":   time.Now().UTC(),
		"environment": utils.Env.Environment,
		"version":     utils.Env.API_VERSION,
		"service":     "nawthtech-worker",
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Service is " + dbHealth.Status,
		"data":    data,
	})
}

func HealthReady(w http.ResponseWriter, r *http.Request) {
	dbHealth := utils.DB.HealthCheck()

	if dbHealth.Status != "healthy" {
		writeJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
			"success": false,
			"error":   "SERVICE_NOT_READY",
			"message": "Database is not ready",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Service is ready",
		"data": map[string]interface{}{
			"status":    "ready",
			"database":  dbHealth.Type,
			"timestamp": time.Now().UTC(),
		},
	})
}

// ========== Auth Handlers ==========

func AuthRegister(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Register endpoint (stub)",
	})
}

func AuthLogin(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Login endpoint (stub)",
	})
}

func AuthRefresh(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Refresh token endpoint (stub)",
	})
}

func AuthForgotPassword(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Forgot password endpoint (stub)",
	})
}

// ========== Users Handlers ==========

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "UNAUTHORIZED",
		})
		return
	}

	user, err := utils.DB.GetUserByID(userID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "USER_NOT_FOUND",
		})
		return
	}

	// إخفاء كلمة السر
	user.Password = ""

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

// ========== Services Handlers ==========

func GetServices(w http.ResponseWriter, r *http.Request) {
	services, _ := utils.DB.GetServices(50)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    services,
	})
}

func GetServiceByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	service, err := utils.DB.GetServiceByID(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "SERVICE_NOT_FOUND",
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    service,
	})
}

// ========== Test Handler ==========

func TestEndpoint(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Test endpoint is working!",
	})
}