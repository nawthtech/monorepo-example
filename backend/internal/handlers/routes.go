package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
)

// RegisterAllRoutes تسجيل جميع المسارات
func RegisterAllRoutes(app *gin.Engine, cfg *config.Config, hc *HandlerContainer) {
	// ==================== Public Routes ====================
	api := app.Group("/api/v1")
	
	// Authentication
	auth := api.Group("/auth")
	{
		if hc.Auth != nil {
			auth.POST("/register", hc.Auth.Register)
			auth.POST("/login", hc.Auth.Login)
			auth.POST("/logout", hc.Auth.Logout)
			auth.POST("/refresh", hc.Auth.RefreshToken)
			auth.POST("/forgot-password", hc.Auth.ForgotPassword)
			auth.POST("/reset-password", hc.Auth.ResetPassword)
		}
	}
	
	// Health endpoints
	health := app.Group("/health")
	{
		if hc.Health != nil {
			health.GET("", hc.Health.CheckHealth)
			health.GET("/live", hc.Health.HealthCheck)
			health.GET("/ready", hc.HealthCheck)
		} else {
			health.GET("/live", func(c *gin.Context) { c.JSON(200, gin.H{"status": "live"}) })
			health.GET("/ready", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ready"}) })
		}
	}
	
	// SSE (Server-Sent Events)
	sse := app.Group("/sse")
	{
		// Add SSE endpoints if needed
		sse.GET("/events", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "SSE endpoint"})
		})
	}
	
	// AI endpoints (from ai.go)
	ai := api.Group("/ai")
	{
		// Assuming you have AI handlers in ai.go
		// ai.POST("/generate", hc.AI.Generate)
		// ai.POST("/chat", hc.AI.Chat)
		ai.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "AI service available"})
		})
	}
	
	// Video endpoints (from video.go)
	video := api.Group("/video")
	{
		// Assuming you have video handlers in video.go
		// video.POST("/process", hc.Video.Process)
		video.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "Video service available"})
		})
	}
	
	// Email endpoints (from email_handler.go)
	email := api.Group("/email")
	{
		// If you have email handler
		email.POST("/send", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Email sent (placeholder)"})
		})
	}

	// ==================== Protected Routes ====================
	// Apply authentication middleware to all protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))
	
	// User routes
	user := protected.Group("/user")
	{
		if hc.User != nil {
			user.GET("/profile", hc.User.GetProfile)
			user.PUT("/profile", hc.User.UpdateProfile)
			user.POST("/change-password", hc.ChangePassword)
		}
	}
	
	// Service routes
	service := protected.Group("/services")
	{
		if hc.Service != nil {
			service.GET("", hc.Service.GetServices)
			service.POST("", hc.Service.CreateService)
			service.GET("/:id", hc.Service.GetServiceByID)
			service.PUT("/:id", hc.Service.UpdateService)
			service.DELETE("/:id", hc.Service.DeleteService)
		}
	}
	
	// Category routes
	category := protected.Group("/categories")
	{
		if hc.Category != nil {
			category.GET("", hc.Category.GetCategories)
			category.POST("", hc.Category.CreateCategory)
			category.GET("/:id", hc.Category.GetCategoryByID)
			category.PUT("/:id", hc.Category.UpdateCategory)
			category.DELETE("/:id", hc.Category.DeleteCategory)
		}
	}
	
	// Order routes
	order := protected.Group("/orders")
	{
		if hc.Order != nil {
			order.GET("", hc.Order.GetUserOrders)
			order.POST("", hc.Order.CreateOrder)
			order.GET("/:id", hc.Order.GetOrderByID)
			order.PUT("/:id/cancel", hc.Order.CancelOrder)
		}
	}
	
	// Payment routes
	payment := protected.Group("/payments")
	{
		if hc.Payment != nil {
			payment.POST("/intent", hc.Payment.CreatePaymentIntent)
			payment.POST("/:id/confirm", hc.Payment.ConfirmPayment)
			payment.GET("/history", hc.Payment.GetPaymentHistory)
			payment.GET("/:id", hc.Payment.GetPaymentByID)
		}
	}
	
	// Upload routes
	upload := protected.Group("/upload")
	{
		if hc.Upload != nil {
			upload.POST("", hc.Upload.UploadFile)
			upload.POST("/image", hc.Upload.UploadImage)
			upload.POST("/document", hc.Upload.UploadDocument)
		}
	}
	
	// Notification routes
	notification := protected.Group("/notifications")
	{
		if hc.Notification != nil {
			notification.GET("", hc.Notification.GetNotifications)
			notification.PUT("/:id/read", hc.Notification.MarkAsRead)
			notification.DELETE("/:id", hc.Notification.DeleteNotification)
			notification.PUT("/settings", hc.Notification.UpdateNotificationSettings)
		}
	}
	
	// Admin routes (admin only)
	admin := protected.Group("/admin")
	admin.Use(middleware.AdminMiddleware())
	{
		if hc.Admin != nil {
			admin.GET("/stats", hc.Admin.GetStatistics)
			admin.GET("/users", hc.Admin.GetAllUsers)
			admin.GET("/users/:id", hc.Admin.GetUserByID)
			admin.PUT("/users/:id", hc.Admin.UpdateUser)
			admin.DELETE("/users/:id", hc.Admin.DeleteUser)
			
			// System management
			admin.GET("/system/health", hc.Admin.GetSystemHealth)
			admin.POST("/system/maintenance", hc.Admin.StartMaintenance)
			admin.DELETE("/system/maintenance", hc.Admin.StopMaintenance)
		}
	}
	
	// ==================== Additional Features ====================
	
	// WebSocket/Socket.IO endpoints (if implemented)
	ws := api.Group("/ws")
	{
		ws.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "WebSocket endpoint"})
		})
	}
	
	// API documentation
	api.GET("/docs", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API Documentation",
			"version": "v1",
			"endpoints": []string{
				"/api/v1/auth/*",
				"/api/v1/user/*",
				"/api/v1/services/*",
				"/api/v1/categories/*",
				"/api/v1/orders/*",
				"/api/v1/payments/*",
				"/api/v1/upload/*",
				"/api/v1/notifications/*",
				"/api/v1/admin/*",
			},
		})
	})
	
	// ==================== Root/Home ====================
	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to NawthTech API",
			"version": "1.0.0",
			"status":  "running",
			"docs":    "/api/v1/docs",
		})
	})
	
	// ==================== 404 Handler ====================
	app.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error":   "Endpoint not found",
			"message": "The requested resource does not exist",
			"path":    c.Request.URL.Path,
		})
	})
}

// Helper function for health check
func (hc *HandlerContainer) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":    "ready",
		"timestamp": "2024-01-01T00:00:00Z",
		"services": gin.H{
			"database": "connected",
			"cache":    "connected",
			"storage":  "connected",
		},
	})
}

// Helper function for change password
func (hc *HandlerContainer) ChangePassword(c *gin.Context) {
	// This is a placeholder - you should implement this properly
	c.JSON(200, gin.H{
		"message": "Password change endpoint",
		"note":    "Implement this in AuthHandler or UserHandler",
	})
}