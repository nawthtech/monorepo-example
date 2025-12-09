package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"nawthtech/worker/src/utils"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "UNAUTHORIZED",
		})
		return
	}

	ctx := context.Background()
	db, err := utils.GetD1().GetDB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "DB_CONNECTION_FAILED",
		})
		return
	}

	row := db.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id = ?", userID)
	user := User{}
	err = row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "USER_NOT_FOUND",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    user,
	})
}