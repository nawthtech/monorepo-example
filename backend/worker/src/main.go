package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"worker/src/handlers"
	"worker/src/utils"
)

// Middleware عام: CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ",")
		origin := r.Header.Get("Origin")
		allowed := false
		for _, o := range allowedOrigins {
			if strings.TrimSpace(o) == origin {
				allowed = true
				break
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-USER-ID")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Middleware المصادقة (يجب تمرير X-USER-ID في Header)
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-USER-ID")
		if userID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"success":false,"error":"UNAUTHORIZED"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// دالة مساعدة للحصول على متغيرات البيئة
func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func main() {
	// تهيئة D1
	if err := utils.InitD1(); err != nil {
		log.Fatalf("Failed to initialize D1: %v", err)
	}

	mux := http.NewServeMux()

	// ==== Health ====
	mux.HandleFunc("/health", handlers.HealthCheck)
	mux.HandleFunc("/health/ready", handlers.HealthReady)

	// ==== Users ====
	mux.Handle("/user/profile", authMiddleware(http.HandlerFunc(handlers.GetUserProfile)))
	mux.Handle("/users", authMiddleware(http.HandlerFunc(handlers.GetUsers)))

	// ==== Services ====
	mux.HandleFunc("/services", handlers.GetServices)
	mux.HandleFunc("/services/id", handlers.GetServiceByID)

	// ==== Test ====
	mux.HandleFunc("/test", handlers.TestHandler)

	// ==== 404 لكل المسارات الأخرى ====
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"success":false,"error":"Not Found"}`))
	})

	// تشغيل السيرفر
	port := getEnv("PORT", "3000")
	log.Printf("Server started on port %s\n", port)
	if err := http.ListenAndServe(":"+port, corsMiddleware(mux)); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}