package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"

	"nawthtech/worker/src/config"
	"nawthtech/worker/src/handlers"
	"nawthtech/worker/src/middleware"
)

func main() {
	cfg := config.Load()

	router := chi.NewRouter()

	// ✅ Middleware شامل
	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.Cors.AllowedOrigins,
		AllowedMethods:   cfg.Cors.AllowedMethods,
		AllowedHeaders:   cfg.Cors.AllowedHeaders,
		ExposedHeaders:   cfg.Cors.ExposedHeaders,
		AllowCredentials: cfg.Cors.AllowCredentials,
	})
	router.Use(c.Handler)
	router.Use(middleware.AuthMiddleware)

	// ✅ Health
	router.Get("/health", handlers.HealthCheck)
	router.Get("/health/live", handlers.HealthLive)
	router.Get("/health/ready", handlers.HealthReady)

	// ✅ Auth
	router.Post("/auth/register", handlers.Register)
	router.Post("/auth/login", handlers.Login)
	router.Post("/auth/refresh", handlers.Refresh)
	router.Post("/auth/forgot-password", handlers.ForgotPassword)

	// ✅ User
	router.Get("/user/profile", handlers.GetProfile)
	router.Put("/user/profile", handlers.UpdateProfile)

	// ✅ Services
	router.Get("/services", handlers.GetServices)
	router.Get("/services/{id}", handlers.GetServiceByID)

	port := cfg.Port
	log.Printf("🚀 Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}