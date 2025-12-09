package middleware

import (
	"context"
	"net/http"
	"strings"

	"nawthtech/worker/src/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// تحقق فقط للمسارات المحمية
		if strings.HasPrefix(r.URL.Path, "/user") || strings.HasPrefix(r.URL.Path, "/services") {
			token := r.Header.Get("Authorization")
			if token == "" {
				utils.JSONResponse(w, http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "UNAUTHORIZED",
				})
				return
			}

			userID, err := utils.ValidateJWT(token)
			if err != nil {
				utils.JSONResponse(w, http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "INVALID_TOKEN",
				})
				return
			}

			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		next.ServeHTTP(w, r)
	})
}