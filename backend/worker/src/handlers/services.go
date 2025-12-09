package handlers

import (
	"net/http"

	"nawthtech/worker/src/utils"
)

func GetServices(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    []string{},
	})
}

func GetServiceByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]string{
			"id": id,
		},
	})
}