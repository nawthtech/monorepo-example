package users

import (
	"encoding/json"
	"net/http"
	"nawthtech-worker/utils"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
}

// GET /users/profile
func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDatabase()

	// مثال: المستخدم الحالي (في الإنتاج يأتي من JWT)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "UNAUTHORIZED",
		})
		return
	}

	// مثال استعلام D1
db := utils.GetDatabase()
results, err := db.Query("SELECT id, email, name FROM users WHERE id = ?", userID)
	if err != nil || len(users) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "USER_NOT_FOUND",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    users[0],
	})
}

// GET /users
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDatabase()

	usersList, err := db.Query("SELECT id, email, name FROM users LIMIT 50")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "DB_ERROR",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    usersList,
	})
}