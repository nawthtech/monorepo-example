package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"nawthtech-worker/src/utils"
)

// ==========================
// Response موحد
// ==========================
type JSONResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func respondJSON(w http.ResponseWriter, status int, resp JSONResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

// ==========================
// Middleware CORS موحد
// ==========================
func CorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS_ALLOWED_ORIGINS"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// ==========================
// Health Handlers
// ==========================
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	dbHealth := utils.HealthCheck()

	healthData := map[string]interface{}{
		"status":      dbHealth.Status,
		"database":    dbHealth.Type,
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
		"environment": os.Getenv("ENVIRONMENT"),
		"version":     os.Getenv("API_VERSION"),
		"service":     "nawthtech-worker",
	}

	respondJSON(w, http.StatusOK, JSONResponse{
		Success: true,
		Message: "Service is " + dbHealth.Status,
		Data:    healthData,
	})
}

func HealthReadyHandler(w http.ResponseWriter, r *http.Request) {
	dbHealth := utils.HealthCheck()
	if dbHealth.Status != "healthy" {
		respondJSON(w, http.StatusServiceUnavailable, JSONResponse{
			Success: false,
			Error:   "SERVICE_NOT_READY",
			Message: "Database is not ready",
		})
		return
	}

	data := map[string]interface{}{
		"status":    "ready",
		"database":  dbHealth.Type,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	respondJSON(w, http.StatusOK, JSONResponse{
		Success: true,
		Message: "Service is ready",
		Data:    data,
	})
}

// ==========================
// Users Handlers
// ==========================
func GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.Header.Get("X-User-ID") // استبدل حسب مصادقة JWT

	if userID == "" {
		respondJSON(w, http.StatusUnauthorized, JSONResponse{
			Success: false,
			Error:   "UNAUTHORIZED",
		})
		return
	}

	user, err := utils.GetUserByID(ctx, userID)
	if err != nil {
		respondJSON(w, http.StatusNotFound, JSONResponse{
			Success: false,
			Error:   "USER_NOT_FOUND",
		})
		return
	}

	// اخفاء كلمة السر
	user.Password = ""

	respondJSON(w, http.StatusOK, JSONResponse{
		Success: true,
		Data:    user,
	})
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := utils.GetUsers(ctx, 50)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, JSONResponse{
			Success: false,
			Error:   "FAILED_FETCH_USERS",
		})
		return
	}

	respondJSON(w, http.StatusOK, JSONResponse{
		Success: true,
		Data:    users,
	})
}

// ==========================
// Services Handlers
// ==========================
func GetServiceByIDHandler(w http.ResponseWriter, r *http.Request, serviceID string) {
	ctx := r.Context()
	service, err := utils.GetServiceByID(ctx, serviceID)
	if err != nil {
		respondJSON(w, http.StatusNotFound, JSONResponse{
			Success: false,
			Error:   "SERVICE_NOT_FOUND",
		})
		return
	}

	respondJSON(w, http.StatusOK, JSONResponse{
		Success: true,
		Data:    service,
	})
}

func GetServicesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	services, err := utils.GetServices(ctx, 50)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, JSONResponse{
			Success: false,
			Error:   "FAILED_FETCH_SERVICES",
		})
		return
	}

	respondJSON(w, http.StatusOK, JSONResponse{
		Success: true,
		Data:    services,
	})
}

// ==========================
// Test Handler
// ==========================
func TestHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, JSONResponse{
		Success: true,
		Message: "Test endpoint working correctly",
	})
}