package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/nawthtech/nawthtech/backend/worker/src/utils"
	"strings"
)

func ServiceListHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := utils.GetD1Database()

	rows, _ := db.DB.Query("SELECT id, name, description FROM services")
	defer rows.Close()

	var services []map[string]string
	for rows.Next() {
		var id, name, desc string
		rows.Scan(&id, &name, &desc)
		services = append(services, map[string]string{
			"id":          id,
			"name":        name,
			"description": desc,
		})
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    services,
	})
}

func ServiceDetailHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := utils.GetD1Database()
	id := strings.TrimPrefix(r.URL.Path, "/services/")

	row := db.DB.QueryRow("SELECT id, name, description FROM services WHERE id=?", id)
	var sid, name, desc string
	err := row.Scan(&sid, &name, &desc)
	if err != nil {
		http.Error(w, "Service not found", 404)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]string{
			"id":          sid,
			"name":        name,
			"description": desc,
		},
	})
}