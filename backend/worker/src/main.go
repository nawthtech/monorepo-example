package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"nawthtech/handlers" // تأكد من تعديل المسار حسب مشروعك
)

// ===================== Middleware شامل =====================

// CORS middleware
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
		origin := r.Header.Get("Origin")
		for _, o := range allowedOrigins {
			if o == origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-API-Key,X-User-ID")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Auth middleware بسيط للتحقق من JWT (استبدال بـ مكتبة JWT لاحقًا)
func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// TODO: تحقق من JWT باستخدام المفتاح من env: JWT_SECRET
		// الآن نعتبره صالح للنسخ
		next.ServeHTTP(w, r)
	})
}

// ===================== Helpers =====================
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// ===================== Main =====================
func main() {
	mux := http.NewServeMux()

	// ===== مسارات الصحة =====
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		handlers.HealthCheck(w, r)
	})
	mux.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		handlers.HealthReady(w, r)
	})

	// ===== مسارات Auth =====
	mux.Handle("/auth/register", http.HandlerFunc(handlers.Register))
	mux.Handle("/auth/login", http.HandlerFunc(handlers.Login))
	mux.Handle("/auth/refresh", http.HandlerFunc(handlers.Refresh))
	mux.Handle("/auth/forgot-password", http.HandlerFunc(handlers.ForgotPassword))

	// ===== مسارات المستخدم =====
	mux.Handle("/user/profile", auth(http.HandlerFunc(handlers.GetProfile)))
	// يمكن إضافة updateProfile بنفس الشكل:
	// mux.Handle("/user/profile", auth(http.HandlerFunc(handlers.UpdateProfile)))

	// ===== مسارات الخدمات =====
	mux.Handle("/services", http.HandlerFunc(handlers.GetServices))
	mux.Handle("/services/", http.HandlerFunc(handlers.GetServiceByID)) // /services/:id

	// ===== مسارات اختبار =====
	mux.Handle("/test", http.HandlerFunc(handlers.TestHandler))

	// ===== جميع المسارات الأخرى =====
	handler := cors(mux)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Nawthtech Worker running on port %s\n", port)
	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}