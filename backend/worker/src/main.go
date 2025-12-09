package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"nawthtech/worker/src/handlers"
	"nawthtech/worker/src/middleware"
	"nawthtech/worker/src/utils"
)

func main() {
	// الاتصال بقاعدة D1
	d1 := utils.GetD1()
	if err := d1.Connect(); err != nil {
		log.Fatalf("❌ Failed to connect to D1: %v", err)
	}
	defer d1.Disconnect(context.Background())

	mux := http.NewServeMux()

	// ✅ مسارات الصحة
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/health/ready", handlers.HealthReadyHandler)

	// ✅ مسارات المستخدم
	mux.Handle("/user/profile", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetProfileHandler)))

	// ✅ مسارات الخدمات
	mux.Handle("/services", http.HandlerFunc(handlers.GetServicesHandler))
	mux.Handle("/services/", http.HandlerFunc(handlers.GetServiceByIDHandler))

	// ✅ مسار الاختبار
	mux.Handle("/test", http.HandlerFunc(handlers.TestHandler))

	// ✅ مسار 404
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "NOT_FOUND",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Printf("🚀 Server running on port %s\n", port)
	if err := http.ListenAndServe(":"+port, middleware.CORSMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}