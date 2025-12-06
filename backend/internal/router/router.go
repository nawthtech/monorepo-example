package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/handlers"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
)

// NewRouter creates and configures the main application router
func NewRouter() *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)
	
	// Create router
	router := gin.New()
	
	// Global middleware
	router.Use(gin.Recovery()) // Recover from panics
	router.Use(middleware.Logger()) // Custom logging
	router.Use(middleware.CORS()) // CORS middleware
	
	// Static files
	router.Static("/static", "./static")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")
	
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"time":    time.Now().UTC(),
			"service": "nawthtech-backend",
		})
	})
	
	// API routes
	api := router.Group("/api")
	{
		// API v1
		v1 := api.Group("/v1")
		{
			// AI routes
			ai := v1.Group("/ai")
			ai.Use(middleware.AuthMiddleware()) // Require auth for AI endpoints
			{
				ai.POST("/generate", handlers.GenerateAIHandler)
				ai.POST("/generate/text", handlers.GenerateTextHandler)
				ai.POST("/generate/image", handlers.GenerateImageHandler)
				ai.POST("/generate/video", handlers.GenerateVideoHandler)
				ai.POST("/translate", handlers.TranslateHandler)
				ai.POST("/analyze", handlers.AnalyzeHandler)
				ai.GET("/models", handlers.GetModelsHandler)
				ai.GET("/usage", handlers.GetUsageHandler)
				ai.GET("/history", handlers.GetHistoryHandler)
			}
			
			// Content routes
			content := v1.Group("/content")
			{
				content.POST("/blog", handlers.GenerateBlogHandler)
				content.POST("/social", handlers.GenerateSocialHandler)
				content.POST("/product", handlers.GenerateProductHandler)
			}
			
			// Media routes
			media := v1.Group("/media")
			{
				media.POST("/upload", handlers.UploadMediaHandler)
				media.GET("/:id", handlers.GetMediaHandler)
			}
			
			// User routes
			users := v1.Group("/users")
			{
				users.POST("/register", handlers.RegisterHandler)
				users.POST("/login", handlers.LoginHandler)
				users.GET("/profile", middleware.AuthMiddleware(), handlers.GetProfileHandler)
				users.PUT("/profile", middleware.AuthMiddleware(), handlers.UpdateProfileHandler)
			}
			
			// Admin routes
			admin := v1.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.GET("/users", handlers.AdminGetUsersHandler)
				admin.GET("/stats", handlers.AdminGetStatsHandler)
				admin.POST("/config", handlers.AdminUpdateConfigHandler)
			}
		}
		
		// API v2 (future version)
		v2 := api.Group("/v2")
		{
			v2.GET("/status", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"version": "2.0.0",
					"status":  "coming_soon",
				})
			})
		}
	}
	
	// SSE (Server-Sent Events) for real-time updates
	router.GET("/sse/ai-progress", handlers.SSEAIProgressHandler)
	
	// WebSocket for real-time communication
	router.GET("/ws/ai", handlers.WebSocketAIHandler)
	
	// Documentation
	router.GET("/docs", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"api": map[string]string{
				"v1": "/api/v1",
				"v2": "/api/v2",
			},
			"endpoints": map[string][]string{
				"ai": {
					"POST /api/v1/ai/generate",
					"POST /api/v1/ai/generate/text",
					"POST /api/v1/ai/generate/image",
					"GET /api/v1/ai/models",
				},
				"content": {
					"POST /api/v1/content/blog",
					"POST /api/v1/content/social",
				},
				"health": {
					"GET /health",
				},
			},
		})
	})
	
	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "endpoint_not_found",
			"message": "The requested endpoint does not exist",
			"path":    c.Request.URL.Path,
			"docs":    "/docs",
		})
	})
	
	return router
}

// CORSMiddleware provides CORS configuration
func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
			"https://nawthtech.com",
			"https://*.nawthtech.com",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"X-API-Key",
			"X-Requested-With",
			"Accept",
			"Cache-Control",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// DevelopmentRouter creates a router with development settings
func DevelopmentRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	router := NewRouter()
	
	// Add development-only middleware
	router.Use(gin.Logger()) // Detailed logging in dev
	
	return router
}

// ProductionRouter creates a router with production settings
func ProductionRouter() *gin.Engine {
	router := NewRouter()
	
	// Add production-only middleware
	router.Use(middleware.RateLimitMiddleware()) // Rate limiting
	router.Use(middleware.SecurityMiddleware())  // Security headers
	
	return router
}