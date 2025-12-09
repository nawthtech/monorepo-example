package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"worker/src/utils"
)

// User يمثل هيكل بيانات المستخدم
type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"` // استبعاد كلمة السر من JSON
}

// ResponseHelper هيكل للردود الموحدة
type ResponseHelper struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// GetProfileHandler يسترجع بيانات المستخدم بواسطة ID
func GetProfileHandler(w http.ResponseWriter, r *http.Request, userID int64) {
	db := utils.GetDatabase()

	query := `SELECT id, name, email FROM users WHERE id = ? LIMIT 1`
	rows, err := db.Query(query, userID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ResponseHelper{
			Success: false,
			Error:   "DATABASE_ERROR",
			Message: err.Error(),
		})
		return
	}

	if len(rows) == 0 {
		respondJSON(w, http.StatusNotFound, ResponseHelper{
			Success: false,
			Error:   "USER_NOT_FOUND",
		})
		return
	}

	user := User{
		ID:    rows[0]["id"].(int64),
		Name:  fmt.Sprintf("%v", rows[0]["name"]),
		Email: fmt.Sprintf("%v", rows[0]["email"]),
	}

	respondJSON(w, http.StatusOK, ResponseHelper{
		Success: true,
		Data:    user,
	})
}

// GetUsersHandler يسترجع قائمة المستخدمين
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.GetDatabase()

	query := `SELECT id, name, email FROM users LIMIT 50`
	rows, err := db.Query(query)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ResponseHelper{
			Success: false,
			Error:   "DATABASE_ERROR",
			Message: err.Error(),
		})
		return
	}

	users := []User{}
	for _, row := range rows {
		user := User{
			ID:    row["id"].(int64),
			Name:  fmt.Sprintf("%v", row["name"]),
			Email: fmt.Sprintf("%v", row["email"]),
		}
		users = append(users, user)
	}

	respondJSON(w, http.StatusOK, ResponseHelper{
		Success: true,
		Data:    users,
	})
}

// respondJSON تساعد على إرسال الردود بصيغة JSON
func respondJSON(w http.ResponseWriter, status int, payload ResponseHelper) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// ExtractUserIDFromRequest يمكن أن تُعدل لتستخرج ID المستخدم من JWT أو سياق request
func ExtractUserIDFromRequest(r *http.Request) int64 {
	// مثال: hardcoded للعرض فقط
	return 1
}

// Middleware لاستخدام مع http.HandlerFunc إذا أردت
func WithUserID(handler func(http.ResponseWriter, *http.Request, int64)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := ExtractUserIDFromRequest(r)
		handler(w, r, userID)
	}
}