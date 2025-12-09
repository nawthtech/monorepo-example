package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"worker/src/handlers"
	"worker/src/middleware"
	"worker/src/utils"
)

// EnvVariables تُخزن إعدادات البيئة
var EnvVariables map[string]string

func init() {
	EnvVariables = map[string]string{
		"ENVIRONMENT": getEnv("ENVIRONMENT", "development"),
		"API_VERSION": getEnv("API_VERSION", "v1"),
	}

	// تهيئة اتصال D1
	if err := utils.InitDatabase(); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()

	// ✅ الصحة
	mux.HandleFunc("/health", middleware.CORSMiddleware(handlers.CheckHealthHandler))
	mux.HandleFunc("/health/live", middleware.CORSMiddleware(handlers.LiveHandler))
	mux.HandleFunc("/health/ready", middleware.CORSMiddleware(handlers.ReadyHandler))

	// ✅ المصادقة
	mux.HandleFunc("/auth/register", middleware.CORSMiddleware(handlers.RegisterHandler))
	mux.HandleFunc("/auth/login", middleware.CORSMiddleware(handlers.LoginHandler))
	mux.HandleFunc("/auth/refresh", middleware.CORSMiddleware(handlers.RefreshHandler))
	mux.HandleFunc("/auth/forgot-password", middleware.CORSMiddleware(handlers.ForgotPasswordHandler))

	// ✅ المستخدمين (مسارات محمية)
	mux.Handle("/user/profile", middleware.CORSMiddleware(middleware.AuthMiddleware(http.HandlerFunc(handlers.GetProfileHandler))))
	mux.Handle("/user/profile/update", middleware.CORSMiddleware(middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateProfileHandler))))

	// ✅ الخدمات
	mux.HandleFunc("/services", middleware.CORSMiddleware(handlers.GetServicesHandler))
	mux.HandleFunc("/services/", middleware.CORSMiddleware(handlers.GetServiceByIDHandler))

	// ✅ اختبار
	mux.HandleFunc("/test", middleware.CORSMiddleware(handlers.TestHandler))

	// ✅ أي مسار غير معروف
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Not Found",
		})
	})

	port := getEnv("PORT", "8787")
	log.Printf("🚀 Worker running on port %s in %s mode", port, EnvVariables["ENVIRONMENT"])
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}

// getEnv يقرأ متغيرات البيئة مع قيمة افتراضية
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return strings.TrimSpace(value)
	}
	return defaultValue
}