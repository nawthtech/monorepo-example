package handlers

import (
	"net/http"

	"nawthtech/worker/src/utils"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	if userID == nil {
		utils.JSONResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "UNAUTHORIZED",
		})
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]string{
			"id": userID.(string),
		},
	})
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Profile updated successfully",
	})
}