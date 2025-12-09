package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"worker/src/utils"
)

// HealthResponse هيكل الرد الصحي
type HealthResponse struct {
	Status     string `json:"status"`
	Database   string `json:"database"`
	Timestamp  string `json:"timestamp"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
	Service     string `json:"service"`
}

// ResponseHelper للردود الموحدة
type ResponseHelper struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// CheckHealthHandler يتحقق من حالة الخدمة
func CheckHealthHandler(w http.ResponseWriter, r *http.Request, env map[string]string) {
	db := utils.GetDatabase()
	dbStatus := "healthy"

	// فحص الاتصال بقاعدة البيانات
	if err := db.Ping(); err != nil {
		dbStatus = "unhealthy"
	}

	healthData := HealthResponse{
		Status:     dbStatus,
		Database:   "D1 Cloudflare",
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Environment: env["ENVIRONMENT"],
		Version:     env["API_VERSION"],
		Service:     "nawthtech-worker",
	}

	respondJSON(w, http.StatusOK, ResponseHelper{
		Success: true,
		Message: "Service is " + dbStatus,
		Data:    healthData,
	})
}

// ReadyHandler يتحقق إذا كانت الخدمة جاهزة للإستعمال
func ReadyHandler(w http.ResponseWriter, r *http.Request, env map[string]string) {
	db := utils.GetDatabase()
	dbStatus := "healthy"

	if err := db.Ping(); err != nil {
		dbStatus = "unhealthy"
	}

	if dbStatus != "healthy" {
		respondJSON(w, http.StatusServiceUnavailable, ResponseHelper{
			Success: false,
			Error:   "SERVICE_NOT_READY",
			Message: "Database is not ready",
		})
		return
	}

	healthData := HealthResponse{
		Status:    "ready",
		Database:  "D1 Cloudflare",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	respondJSON(w, http.StatusOK, ResponseHelper{
		Success: true,
		Message: "Service is ready",
		Data:    healthData,
	})
}

// respondJSON تساعد على إرسال الردود بصيغة JSON
func respondJSON(w http.ResponseWriter, status int, payload ResponseHelper) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}