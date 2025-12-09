package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"worker/src/utils"
)

// ResponseHelper هيكل للرد الموحد
type ResponseHelper struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ====== Health Handlers ======

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	status, err := utils.HealthCheck()
	data := map[string]interface{}{
		"status":      status,
		"timestamp":   time.Now().UTC(),
		"service":     "nawthtech-worker",
		"environment": getEnv("ENVIRONMENT", "development"),
		"version":     getEnv("API_VERSION", "v1"),
	}

	if err != nil {
		writeJSON(w, http.Status503ServiceUnavailable, ResponseHelper{
			Success: false,
			Error:   "SERVICE_UNHEALTHY",
			Message: err.Error(),
			Data:    data,
		})
		return
	}

	writeJSON(w, http.StatusOK, ResponseHelper{
		Success: true,
		Message: fmt.Sprintf("Service is %s", status),
		Data:    data,
	})
}

func HealthReady(w http.ResponseWriter, r *http.Request) {
	status, err := utils.HealthCheck()
	if err != nil || status != "healthy" {
		writeJSON(w, http.StatusServiceUnavailable, ResponseHelper{
			Success: false,
			Error:   "SERVICE_NOT_READY",
			Message: "Database is not ready",
		})
		return
	}

	writeJSON(w, http.StatusOK, ResponseHelper{
		Success: true,
		Message: "Service is ready",
		Data: map[string]interface{}{
			"status":    "ready",
			"database":  "D1",
			"timestamp": time.Now().UTC(),
		},
	})
}

// ====== User Handlers ======

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-USER-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, ResponseHelper{
			Success: false,
			Error:   "UNAUTHORIZED",
		})
		return
	}

	row := utils.QueryRow("SELECT id, name, email FROM users WHERE id = ?", userID)
	var id, name, email string
	if err := row.Scan(&id, &name, &email); err != nil {
		writeJSON(w, http.StatusNotFound, ResponseHelper{
			Success: false,
			Error:   "USER_NOT_FOUND",
		})
		return
	}

	writeJSON(w, http.StatusOK, ResponseHelper{
		Success: true,
		Data: map[string]string{
			"id":    id,
			"name":  name,
			"email": email,
		},
	})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := utils.QueryRows("SELECT id, name, email FROM users LIMIT 50")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ResponseHelper{
			Success: false,
			Error:   "DB_ERROR",
			Message: err.Error(),
		})
		return
	}
	defer rows.Close()

	users := []map[string]string{}
	for rows.Next() {
		var id, name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			log.Println("Row scan error:", err)
			continue
		}
		users = append(users, map[string]string{
			"id":    id,
			"name":  name,
			"email": email,
		})
	}

	writeJSON(w, http.StatusOK, ResponseHelper{
		Success: true,
		Data:    users,
	})
}

// ====== Services Handlers ======

func GetServices(w http.ResponseWriter, r *http.Request) {
	rows, err := utils.QueryRows("SELECT id, title, price FROM services LIMIT 50")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ResponseHelper{
			Success: false,
			Error:   "DB_ERROR",
			Message: err.Error(),
		})
		return
	}
	defer rows.Close()

	services := []map[string]interface{}{}
	for rows.Next() {
		var id, title string
		var price float64
		if err := rows.Scan(&id, &title, &price); err != nil {
			log.Println("Row scan error:", err)
			continue
		}
		services = append(services, map[string]interface{}{
			"id":    id,
			"title": title,
			"price": price,
		})
	}

	writeJSON(w, http.StatusOK, ResponseHelper{
		Success: true,
		Data:    services,
	})
}

func GetServiceByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ResponseHelper{
			Success: false,
			Error:   "INVALID_ID",
		})
		return
	}

	row := utils.QueryRow("SELECT id, title, price FROM services WHERE id = ?", id)
	var title string
	var price float64
	if err := row.Scan(&id, &title, &price); err != nil {
		writeJSON(w, http.StatusNotFound, ResponseHelper{
			Success: false,
			Error:   "SERVICE_NOT_FOUND",
		})
		return
	}

	writeJSON(w, http.StatusOK, ResponseHelper{
		Success: true,
		Data: map[string]interface{}{
			"id":    id,
			"title": title,
			"price": price,
		},
	})
}

// ====== Test Handler ======

func TestHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, ResponseHelper{
		Success: true,
		Message: "Test endpoint working!",
	})
}

// ====== Helpers ======

func writeJSON(w http.ResponseWriter, status int, payload ResponseHelper) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}