package handlers

import (
	"net/http"

	"nawthtech/worker/src/utils"
)

func TestEndpoint(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Test endpoint working!",
	})
}