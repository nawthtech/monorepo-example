package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"nawthtech-worker/handlers/health"
	"nawthtech-worker/handlers/users"
	"nawthtech-worker/utils"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8787" // منفذ افتراضي لـ Cloudflare Worker Go
	}

	http.HandleFunc("/api/v1/health/check", health.CheckHandler)
	http.HandleFunc("/api/v1/health/ready", health.ReadyHandler)

	http.HandleFunc("/api/v1/users/profile", users.GetProfileHandler)
	http.HandleFunc("/api/v1/users", users.GetUsersHandler)

	srv := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	fmt.Printf("🚀 Worker running on port %s\n", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}