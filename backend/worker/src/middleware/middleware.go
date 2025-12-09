package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

// CORSMiddleware لتفعيل CORS
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
		origin := r.Header.Get("Origin")
		for _, o := range allowedOrigins {
			if o == origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware للتحقق من JWT (مثال بسيط)
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "UNAUTHORIZED",
			})
			return
		}
		// هنا يمكن إضافة التحقق من JWT لاحقًا
		next.ServeHTTP(w, r)
	})
}