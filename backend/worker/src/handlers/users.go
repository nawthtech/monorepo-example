package handlers

import (
	"encoding/json"
	"net/http"
	"worker/src/utils"
)

func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	db, err := utils.GetD1Database()
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	userID := r.Header.Get("user_id")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// استعلام D1
	row := db.DB.QueryRow("SELECT id, username, email FROM users WHERE id = ?", userID)
	var id, username, email string
	err = row.Scan(&id, &username, &email)
	if err != nil {
		http.Error(w, "User not found", 404)
		return
	}

	user := map[string]string{
		"id":       id,
		"username": username,
		"email":    email,
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

func UserListHandler(w http.ResponseWriter, r *http.Request) {
	db, err := utils.GetD1Database()
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	rows, err := db.DB.Query("SELECT id, username, email FROM users LIMIT 50")
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	defer rows.Close()

	var users []map[string]string
	for rows.Next() {
		var id, username, email string
		rows.Scan(&id, &username, &email)
		users = append(users, map[string]string{
			"id":       id,
			"username": username,
			"email":    email,
		})
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    users,
	})
}