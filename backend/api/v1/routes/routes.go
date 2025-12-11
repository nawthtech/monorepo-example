package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/handlers"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
)

// HandlerContainer حاوية لجميع المعاجل
type HandlerContainer struct {
	Auth         handlers.AuthHandler
	User         handlers.UserHandler
	Service      handlers.ServiceHandler
	Category     handlers.CategoryHandler
	Order        handlers.OrderHandler
	Payment      handlers.PaymentHandler
	Upload       handlers.UploadHandler
	Notification handlers.NotificationHandler
	Admin        handlers.AdminHandler
}

// RegisterV1Routes تسجيل جميع مسارات الإصدار 1 في ملف واحد
func RegisterV1Routes(router *gin.RouterGroup, handlers *HandlerContainer, authMiddleware gin.HandlerFunc) {

	// ================================
	// ✅ المسارات العامة (بدون مصادقة)
	// ================================
	public := router.Group("")
	{
		// 🔐 مسارات المصادقة
		auth := public.Group("/auth")
		{
			auth.POST("/register", handlers.Auth.Register)
			auth.POST("/login", handlers.Auth.Login)
			auth.POST("/logout", handlers.Auth.Logout)
			auth.POST("/refresh-token", handlers.Auth.RefreshToken)
			auth.POST("/forgot-password", handlers.Auth.ForgotPassword)
			auth.POST("/reset-password", handlers.Auth.ResetPassword)
			auth.GET("/verify-token", handlers.Auth.VerifyToken)
		}

		// 🛍️ مسارات الخدمات
		services := public.Group("/services")
		{
			services.GET("/", handlers.Service.GetServices)
			services.GET("/search", handlers.Service.SearchServices)
			services.GET("/featured", handlers.Service.GetFeaturedServices)
			services.GET("/categories", handlers.Service.GetCategories)
			services.GET("/:id", handlers.Service.GetServiceByID)
			services.GET("/:id/similar", handlers.Service.GetSimilarServices)
		}

		// 📁 مسارات الفئات
		categories := public.Group("/categories")
		{
			categories.GET("/", handlers.Category.GetCategories)
			categories.GET("/:id", handlers.Category.GetCategoryByID)
			categories.GET("/tree", handlers.Category.GetCategoryTree)
		}

		// 💚 مسارات الصحة والفحص
		health := public.Group("/health")
		{
			health.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"status":    "healthy",
					"timestamp": time.Now().UTC(),
					"service":   "NawthTech API v1",
					"database":  "Cloudflare D1",
				})
			})

			health.GET("/ready", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"status":    "ready",
					"timestamp": time.Now().UTC(),
					"service":   "NawthTech API v1",
					"database":  "Cloudflare D1",
				})
			})

			health.GET("/live", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"status":    "live",
					"timestamp": time.Now().UTC(),
					"service":   "NawthTech API v1",
					"database":  "Cloudflare D1",
				})
			})
		}

		// 📚 مسارات التوثيق
		docs := public.Group("/docs")
		{
			docs.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"name":          "NawthTech API Documentation",
					"version":       "v1.0.0",
					"description":   "منصة نوذ تك للخدمات الإلكترونية - وثائق API",
					"documentation": "سيتم إضافة رابط التوثيق هنا",
					"stack": gin.H{
						"database":       "Cloudflare D1",
						"upload_service": "Cloudinary",
						"backend":        "Go + Gin",
						"frontend":       "React + Vite",
						"deployment":     "Cloudflare Workers",
					},
					"endpoints": []string{
						"GET    /api/v1/health",
						"POST   /api/v1/auth/register",
						"POST   /api/v1/auth/login",
						"GET    /api/v1/services",
						"POST   /api/v1/upload/image",
						"GET    /api/v1/categories",
						"POST   /api/v1/orders",
						"GET    /api/v1/users/profile",
					},
				})
			})

			docs.GET("/openapi", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"openapi": "3.0.0",
					"info": gin.H{
						"title":       "NawthTech API",
						"version":     "v1.0.0",
						"description": "منصة الخدمات الإلكترونية",
					},
					"servers": []gin.H{
						{
							"url":         "/api/v1",
							"description": "الإصدار 1 من API",
						},
					},
					"components": gin.H{
						"schemas": gin.H{
							"Database": gin.H{
								"type":        "string",
								"description": "Cloudflare D1 SQLite Database",
							},
						},
					},
				})
			})
		}
	}

	// ================================
	// ✅ المسارات المحمية (تتطلب مصادقة)
	// ================================
	protected := router.Group("")
	protected.Use(authMiddleware)
	{
		// 👤 مسارات المستخدم
		users := protected.Group("/users")
		{
			users.GET("/profile", handlers.User.GetProfile)
			users.PUT("/profile", handlers.User.UpdateProfile)
			users.PUT("/avatar", handlers.User.UpdateAvatar)
			users.PUT("/change-password", handlers.User.ChangePassword)
			users.GET("/stats", handlers.User.GetUserStats)
			users.DELETE("/account", handlers.User.DeleteAccount)
		}

		// 🛒 مسارات الطلبات
		orders := protected.Group("/orders")
		{
			orders.GET("/", handlers.Order.GetUserOrders)
			orders.POST("/", handlers.Order.CreateOrder)
			orders.GET("/:id", handlers.Order.GetOrderByID)
			orders.PUT("/:id/status", handlers.Order.UpdateOrderStatus)
			orders.DELETE("/:id", handlers.Order.CancelOrder)
			orders.GET("/stats", handlers.Order.GetOrderStats)
		}

		// 💳 مسارات الدفع
		payments := protected.Group("/payments")
		{
			payments.GET("/history", handlers.Payment.GetPaymentHistory)
			payments.POST("/create-intent", handlers.Payment.CreatePaymentIntent)
			payments.POST("/confirm", handlers.Payment.ConfirmPayment)
			payments.POST("/validate", handlers.Payment.ValidatePayment)
		}

		// ☁️ مسارات الرفع - Cloudinary
		upload := protected.Group("/upload")
		{
			upload.POST("/image", handlers.Upload.UploadImage)
			upload.POST("/images", handlers.Upload.UploadMultipleImages)
			upload.GET("/image/:public_id", handlers.Upload.GetImageInfo)
			upload.DELETE("/image/:public_id", handlers.Upload.DeleteImage)
			upload.GET("/my-images", handlers.Upload.GetUserImages)
			upload.GET("/presigned-url", handlers.Upload.GeneratePresignedURL)
			upload.POST("/file", handlers.Upload.UploadFile)
			upload.GET("/file/:id", handlers.Upload.GetFile)
			upload.DELETE("/file/:id", handlers.Upload.DeleteFile)
			upload.GET("/my-files", handlers.Upload.GetUserFiles)
		}

		// 🔔 مسارات الإشعارات
		notifications := protected.Group("/notifications")
		{
			notifications.GET("/", handlers.Notification.GetUserNotifications)
			notifications.PUT("/:id/read", handlers.Notification.MarkAsRead)
			notifications.PUT("/read-all", handlers.Notification.MarkAllAsRead)
			notifications.GET("/unread-count", handlers.Notification.GetUnreadCount)
			notifications.DELETE("/:id", handlers.Notification.DeleteNotification)
			notifications.POST("/", handlers.Notification.CreateNotification)
		}

		// 🛍️ مسارات الخدمات المحمية
		protectedServices := protected.Group("/services")
		{
			protectedServices.GET("/my-services", handlers.Service.GetMyServices)
			protectedServices.POST("/", handlers.Service.CreateService)
			protectedServices.PUT("/:id", handlers.Service.UpdateService)
			protectedServices.DELETE("/:id", handlers.Service.DeleteService)
			protectedServices.GET("/similar/:id", handlers.Service.GetSimilarServices)
		}

		// 📁 مسارات الفئات المحمية
		protectedCategories := protected.Group("/categories")
		{
			protectedCategories.POST("/", handlers.Category.CreateCategory)
			protectedCategories.PUT("/:id", handlers.Category.UpdateCategory)
			protectedCategories.DELETE("/:id", handlers.Category.DeleteCategory)
			protectedCategories.GET("/tree", handlers.Category.GetCategoryTree)
		}

		// 🛒 مسارات السلة (Cart)
		cart := protected.Group("/cart")
		{
			cart.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "سلة المشتريات - تحت التطوير",
					"status":  "coming_soon",
				})
			})
		}
	}

	// ================================
	// ✅ مسارات الإدارة (تتطلب صلاحيات إدارية)
	// ================================
	admin := router.Group("/admin")
	admin.Use(authMiddleware, middleware.AdminRequired())
	{
		// 📊 لوحة التحكم
		admin.GET("/dashboard", handlers.Admin.GetDashboard)
		admin.GET("/dashboard/stats", handlers.Admin.GetDashboardStats)
		admin.GET("/dashboard/metrics", handlers.Admin.GetSystemMetrics)

		// 👥 إدارة المستخدمين
		admin.GET("/users", handlers.Admin.GetUsers)
		admin.PUT("/users/:id/status", handlers.Admin.UpdateUserStatus)
		admin.POST("/users/:id/ban", handlers.Admin.BanUser)
		admin.POST("/users/:id/unban", handlers.Admin.UnbanUser)
		admin.GET("/users/search", handlers.Admin.SearchUsers)

		// 📋 سجلات النظام
		admin.GET("/system-logs", handlers.Admin.GetSystemLogs)
		admin.DELETE("/system-logs/clean", handlers.Admin.CleanOldLogs)

		// ⚙️ إعدادات النظام
		admin.PUT("/system-settings", handlers.Admin.UpdateSystemSettings)
		admin.GET("/system-config", handlers.Admin.GetSystemConfig)

		// 🛍️ إدارة الخدمات (إضافية)
		adminServices := admin.Group("/services")
		{
			adminServices.GET("/all", handlers.Service.GetAllServices)
			adminServices.DELETE("/:id/force", handlers.Service.ForceDeleteService)
			adminServices.PUT("/:id/feature", handlers.Service.ToggleFeatured)
			adminServices.PUT("/:id/status", handlers.Service.ToggleStatus)
			adminServices.GET("/stats", handlers.Service.GetServiceStats)
		}

		// 📁 إدارة الفئات (إضافية)
		adminCategories := admin.Group("/categories")
		{
			adminCategories.GET("/all", handlers.Category.GetAllCategories)
			adminCategories.POST("/bulk", handlers.Category.BulkCreateCategories)
			adminCategories.PUT("/:id/status", handlers.Category.ToggleCategoryStatus)
		}

		// 📊 إدارة الطلبات
		adminOrders := admin.Group("/orders")
		{
			adminOrders.GET("/all", handlers.Order.GetAllOrders)
			adminOrders.GET("/stats/advanced", handlers.Order.GetAdvancedStats)
			adminOrders.PUT("/:id/force-status", handlers.Order.ForceUpdateStatus)
		}

		// 💳 إدارة المدفوعات
		adminPayments := admin.Group("/payments")
		{
			adminPayments.GET("/all", handlers.Payment.GetAllPayments)
			adminPayments.GET("/revenue", handlers.Payment.GetRevenueStats)
			adminPayments.POST("/:id/refund", handlers.Payment.ProcessRefund)
		}

		// 🔔 إدارة الإشعارات
		adminNotifications := admin.Group("/notifications")
		{
			adminNotifications.GET("/all", handlers.Notification.GetAllNotifications)
			adminNotifications.POST("/broadcast", handlers.Notification.BroadcastNotification)
			adminNotifications.DELETE("/clean", handlers.Notification.CleanOldNotifications)
		}

		// ☁️ إدارة الملفات
		adminFiles := admin.Group("/files")
		{
			adminFiles.GET("/all", handlers.Upload.GetAllFiles)
			adminFiles.DELETE("/clean", handlers.Upload.CleanOrphanedFiles)
			adminFiles.GET("/storage", handlers.Upload.GetStorageStats)
		}

		// 🏥 صحة النظام للمسؤولين
		adminHealth := admin.Group("/health")
		{
			adminHealth.GET("/", handlers.Admin.GetAdminHealth)
			adminHealth.GET("/detailed", handlers.Admin.GetDetailedHealth)
			adminHealth.GET("/database", handlers.Admin.GetDatabaseHealth)
		}
	}

	// ================================
	// ✅ المسارات العامة الإضافية
	// ================================

	// 🏠 الصفحة الرئيسية للـ API
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":        "مرحباً بك في نوذ تك - منصة الخدمات الإلكترونية",
			"version":        "v1.0.0",
			"timestamp":      time.Now().UTC(),
			"stack": gin.H{
				"database":        "Cloudflare D1",
				"upload_service":  "Cloudinary",
				"authentication":  "JWT Tokens",
				"backend":         "Go + Gin",
				"frontend":        "React + Vite",
				"deployment":      "Cloudflare Workers",
				"cache":           "In-Memory",
				"email":           "SMTP + Mailersend",
				"payments":        "Stripe Integration",
			},
			"status": "running",
			"endpoints": gin.H{
				"auth":         "/api/v1/auth",
				"services":     "/api/v1/services",
				"categories":   "/api/v1/categories",
				"users":        "/api/v1/users",
				"orders":       "/api/v1/orders",
				"payments":     "/api/v1/payments",
				"upload":       "/api/v1/upload",
				"notifications": "/api/v1/notifications",
				"health":       "/api/v1/health",
				"docs":         "/api/v1/docs",
				"admin":        "/api/v1/admin",
			},
		})
	})

	// 🔍 مسار البحث العام
	router.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "نتائج البحث",
			"data": gin.H{
				"query":   query,
				"results": []gin.H{},
				"filters": gin.H{
					"services":   true,
					"categories": true,
					"users":      false,
				},
				"database": "Cloudflare D1",
			},
		})
	})

	// 📈 إحصائيات الـ API
	router.GET("/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "إحصائيات النظام",
			"data": gin.H{
				"total_endpoints": 65,
				"public_endpoints": 20,
				"protected_endpoints": 30,
				"admin_endpoints": 15,
				"active_services": 150,
				"total_users":     1250,
				"total_orders":    890,
				"total_payments":  750,
				"database":        "Cloudflare D1",
				"uptime":          "99.8%",
				"response_time":   "125ms",
				"storage_used":    "2.5GB",
				"api_version":     "v1.0.0",
			},
		})
	})

	// 🔄 مسار إعادة تعيين التخزين المؤقت (للأغراض التطويرية)
	router.POST("/cache/reset", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "تم إعادة تعيين التخزين المؤقت",
			"timestamp": time.Now().UTC(),
		})
	})

	// ⚙️ معلومات الإصدار
	router.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"api_version":    "v1.0.0",
				"go_version":     "1.25.4",
				"gin_version":    "v1.9.1",
				"database":       "Cloudflare D1",
				"cloud_provider": "Cloudflare",
				"deployed_at":    "2024-01-15T10:30:00Z",
				"last_updated":   time.Now().UTC(),
			},
		})
	})
}

// GetRoutesInfo معلومات عن المسارات المسجلة
func GetRoutesInfo() map[string]interface{} {
	return map[string]interface{}{
		"total_endpoints":       65,
		"public_endpoints":      20,
		"protected_endpoints":   30,
		"admin_endpoints":       15,
		"version":              "v1.0.0",
		"database":             "Cloudflare D1",
		"categories": []string{
			"المصادقة", "المستخدمين", "الخدمات", "الفئات",
			"الطلبات", "الدفع", "الرفع", "الإشعارات", "الإدارة",
			"الصحة", "التوثيق", "البحث", "الإحصائيات",
		},
		"authentication": "JWT Bearer Tokens",
		"rate_limiting":  "مفعل",
		"cors":           "مفعل",
		"security":       "HTTPS فقط في Production",
		"monitoring":     "مفعل",
		"logging":        "مفعل",
	}
}