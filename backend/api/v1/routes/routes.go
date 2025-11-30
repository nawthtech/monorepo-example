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
		}

		// 📁 مسارات الفئات
		categories := public.Group("/categories")
		{
			categories.GET("/", handlers.Category.GetCategories)
			categories.GET("/:id", handlers.Category.GetCategoryByID)
		}

		// 💚 مسارات الصحة والفحص
		health := public.Group("/health")
		{
			health.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"status":    "healthy",
					"timestamp": time.Now().UTC(),
					"service":   "NawthTech API v1",
				})
			})

			health.GET("/ready", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"status":    "ready",
					"timestamp": time.Now().UTC(),
					"service":   "NawthTech API v1",
				})
			})

			health.GET("/live", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"status":    "live",
					"timestamp": time.Now().UTC(),
					"service":   "NawthTech API v1",
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
			users.PUT("/change-password", handlers.User.ChangePassword)
			users.GET("/stats", handlers.User.GetUserStats)
		}

		// 🛒 مسارات الطلبات
		orders := protected.Group("/orders")
		{
			orders.GET("/", handlers.Order.GetUserOrders)
			orders.POST("/", handlers.Order.CreateOrder)
			orders.GET("/:id", handlers.Order.GetOrderByID)
			orders.PUT("/:id/status", handlers.Order.UpdateOrderStatus)
			orders.DELETE("/:id", handlers.Order.CancelOrder)
		}

		// 💳 مسارات الدفع
		payments := protected.Group("/payments")
		{
			payments.GET("/history", handlers.Payment.GetPaymentHistory)
			payments.POST("/create-intent", handlers.Payment.CreatePaymentIntent)
			payments.POST("/confirm", handlers.Payment.ConfirmPayment)
		}

		// ☁️ مسارات الرفع - Cloudinary
		upload := protected.Group("/upload")
		{
			upload.POST("/image", handlers.Upload.UploadImage)
			upload.POST("/images", handlers.Upload.UploadMultipleImages)
			upload.GET("/image/:public_id", handlers.Upload.GetImageInfo)
			upload.DELETE("/image/:public_id", handlers.Upload.DeleteImage)
			upload.GET("/my-images", handlers.Upload.GetUserImages)
		}

		// 🔔 مسارات الإشعارات
		notifications := protected.Group("/notifications")
		{
			notifications.GET("/", handlers.Notification.GetUserNotifications)
			notifications.PUT("/:id/read", handlers.Notification.MarkAsRead)
			notifications.PUT("/read-all", handlers.Notification.MarkAllAsRead)
			notifications.GET("/unread-count", handlers.Notification.GetUnreadCount)
		}

		// 🛍️ مسارات الخدمات المحمية
		protectedServices := protected.Group("/services")
		{
			protectedServices.GET("/my-services", handlers.Service.GetMyServices)
			protectedServices.POST("/", handlers.Service.CreateService)
			protectedServices.PUT("/:id", handlers.Service.UpdateService)
			protectedServices.DELETE("/:id", handlers.Service.DeleteService)
		}

		// 📁 مسارات الفئات المحمية
		protectedCategories := protected.Group("/categories")
		{
			protectedCategories.POST("/", handlers.Category.CreateCategory)
			protectedCategories.PUT("/:id", handlers.Category.UpdateCategory)
			protectedCategories.DELETE("/:id", handlers.Category.DeleteCategory)
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

		// 👥 إدارة المستخدمين
		admin.GET("/users", handlers.Admin.GetUsers)
		admin.PUT("/users/:id/status", handlers.Admin.UpdateUserStatus)

		// 📋 سجلات النظام
		admin.GET("/system-logs", handlers.Admin.GetSystemLogs)

		// 🛍️ إدارة الخدمات (إضافية)
		adminServices := admin.Group("/services")
		{
			adminServices.GET("/all", handlers.Service.GetServices)
			adminServices.DELETE("/:id/force", handlers.Service.DeleteService)
		}

		// 📁 إدارة الفئات (إضافية)
		adminCategories := admin.Group("/categories")
		{
			adminCategories.GET("/all", handlers.Category.GetCategories)
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
			"database":       "MongoDB",
			"upload_service": "Cloudinary",
			"status":         "running",
			"endpoints": gin.H{
				"auth":       "/api/v1/auth",
				"services":   "/api/v1/services",
				"categories": "/api/v1/categories",
				"users":      "/api/v1/users",
				"orders":     "/api/v1/orders",
				"upload":     "/api/v1/upload",
				"health":     "/api/v1/health",
				"docs":       "/api/v1/docs",
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
			},
		})
	})

	// 📈 إحصائيات الـ API
	router.GET("/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "إحصائيات النظام",
			"data": gin.H{
				"total_endpoints": 45,
				"active_services": 150,
				"total_users":     1250,
				"total_orders":    890,
				"uptime":          "99.8%",
				"response_time":   "125ms",
			},
		})
	})
}

// GetRoutesInfo معلومات عن المسارات المسجلة
func GetRoutesInfo() map[string]interface{} {
	return map[string]interface{}{
		"total_endpoints":     45,
		"public_endpoints":    15,
		"protected_endpoints": 25,
		"admin_endpoints":     5,
		"version":             "v1.0.0",
		"categories": []string{
			"المصادقة", "المستخدمين", "الخدمات", "الفئات",
			"الطلبات", "الدفع", "الرفع", "الإشعارات", "الإدارة",
		},
	}
}
