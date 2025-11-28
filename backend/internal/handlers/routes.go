package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/config"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"gorm.io/gorm"
)

// ServiceContainer حاوية الخدمات
type ServiceContainer struct {
	// الخدمات الأساسية
	AuthService         services.AuthService
	UserService         services.UserService
	AdminService        services.AdminService
	ServicesService     services.ServicesService
	CacheService        services.CacheService
	HealthService       services.HealthService
	
	// خدمات المتجر والطلبات
	StoreService        services.StoreService
	CartService         services.CartService
	OrdersService       services.OrdersService
	PaymentService      services.PaymentService
	CategoryService     services.CategoryService
	
	// خدمات المحتوى والإشعارات
	ContentService      services.ContentService
	NotificationService services.NotificationService
	UploadService       services.UploadService
	
	// خدمات التقارير والتحليلات
	AnalyticsService    services.AnalyticsService
	ReportsService      services.ReportsService
	StrategiesService   services.StrategiesService
	
	// خدمات الذكاء الاصطناعي والموقع
	AIService           services.AIService
	WebsiteService      services.WebsiteService
}

// MiddlewareContainer حاوية الوسائط
type MiddlewareContainer struct {
	AuthMiddleware      gin.HandlerFunc
	AdminMiddleware     gin.HandlerFunc
	SellerMiddleware    gin.HandlerFunc
	CORSMiddleware      gin.HandlerFunc
	SecurityMiddleware  gin.HandlerFunc
	CacheMiddleware     *middleware.CacheMiddleware
	RateLimitMiddleware *middleware.RateLimitMiddleware
}

// RegisterAllRoutes تسجيل جميع المسارات
func RegisterAllRoutes(router *gin.Engine, db *gorm.DB, config *config.Config) {
	// تطبيق middleware العام على مستوى التطبيق
	applyGlobalMiddleware(router, config)
	
	// إنشاء حاوية الخدمات
	services := initializeServices(db, config)
	
	// إنشاء حاوية الوسائط
	middlewares := initializeMiddlewares(services, config)
	
	// مجموعة API الرئيسية
	api := router.Group("/api/v1")
	
	// ========== المسارات العامة (لا تتطلب مصادقة) ==========
	registerPublicRoutes(api, services, middlewares)
	
	// ========== المسارات المحمية (تتطلب مصادقة) ==========
	registerProtectedRoutes(api, services, middlewares)
	
	// ========== مسارات المسؤولين ==========
	registerAdminRoutes(api, services, middlewares)
	
	// ========== مسارات البائعين ==========
	registerSellerRoutes(api, services, middlewares)
	
	// ========== مسارات SSE والوقت الحقيقي ==========
	registerRealTimeRoutes(api, services, middlewares)
	
	// ========== مسارات الويب هووك والتحليلات ==========
	registerWebhookRoutes(api, services, middlewares)
}

// applyGlobalMiddleware تطبيق الوسائط العامة على مستوى التطبيق
func applyGlobalMiddleware(router *gin.Engine, config *config.Config) {
	// CORS middleware - يتم تطبيقه على مستوى التطبيق بالكامل
	router.Use(middleware.CORS())
	
	// Security headers middleware
	router.Use(middleware.SecurityHeaders())
	
	// Logging middleware
	router.Use(middleware.Logging())
	
	// Rate limiting middleware
	router.Use(middleware.RateLimit())
}

// initializeServices تهيئة جميع الخدمات
func initializeServices(db *gorm.DB, config *config.Config) *ServiceContainer {
	// تهيئة خدمة التخزين المؤقت أولاً
	cacheService := services.NewCacheService(config.GetCacheConfig())
	
	// الخدمات الأساسية
	authService := services.NewAuthService(db, config.JWTSecret)
	userService := services.NewUserService(db)
	adminService := services.NewAdminService(db)
	healthService := services.NewHealthService(db)
	
	// خدمات المستودعات
	servicesRepo := services.NewServicesRepository(db)
	servicesService := services.NewServicesService(servicesRepo)
	
	// خدمات المتجر والطلبات
	storeService := services.NewStoreService(db)
	cartService := services.NewCartService(db)
	ordersService := services.NewOrdersService(db)
	paymentService := services.NewPaymentService(db)
	categoryService := services.NewCategoryService(db)
	
	// خدمات المحتوى والإشعارات
	contentService := services.NewContentService(db)
	notificationService := services.NewNotificationService(db)
	uploadService := services.NewUploadService(config.GetUploadConfig())
	
	// خدمات التقارير والتحليلات
	analyticsService := services.NewAnalyticsService(db)
	reportsService := services.NewReportsService(db)
	strategiesService := services.NewStrategiesService(db)
	
	// خدمات الذكاء الاصطناعي والموقع
	aiService := services.NewAIService(db)
	websiteService := services.NewWebsiteService(db)
	
	return &ServiceContainer{
		AuthService:         authService,
		UserService:         userService,
		AdminService:        adminService,
		ServicesService:     servicesService,
		CacheService:        cacheService,
		HealthService:       healthService,
		StoreService:        storeService,
		CartService:         cartService,
		OrdersService:       ordersService,
		PaymentService:      paymentService,
		CategoryService:     categoryService,
		ContentService:      contentService,
		NotificationService: notificationService,
		UploadService:       uploadService,
		AnalyticsService:    analyticsService,
		ReportsService:      reportsService,
		StrategiesService:   strategiesService,
		AIService:           aiService,
		WebsiteService:      websiteService,
	}
}

// initializeMiddlewares تهيئة جميع الوسائط
func initializeMiddlewares(services *ServiceContainer, config *config.Config) *MiddlewareContainer {
	// الوسائط الأساسية
	authMiddleware := middleware.AuthMiddleware(services.AuthService)
	adminMiddleware := middleware.AdminMiddleware()
	sellerMiddleware := middleware.SellerMiddleware()
	
	// وسائط التخزين المؤقت
	cacheMiddleware := middleware.NewCacheMiddleware(services.CacheService, config.Cache.Prefix)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(services.CacheService)
	
	return &MiddlewareContainer{
		AuthMiddleware:      authMiddleware,
		AdminMiddleware:     adminMiddleware,
		SellerMiddleware:    sellerMiddleware,
		CacheMiddleware:     cacheMiddleware,
		RateLimitMiddleware: rateLimitMiddleware,
	}
}

// registerPublicRoutes تسجيل المسارات العامة
func registerPublicRoutes(api *gin.RouterGroup, services *ServiceContainer, middlewares *MiddlewareContainer) {
	// معالج الصحة
	healthHandler := NewHealthHandler(services.HealthService, config.Load())
	api.GET("/health", healthHandler.Check)
	api.GET("/health/live", healthHandler.Live)
	api.GET("/health/ready", healthHandler.Ready)
	api.GET("/health/info", healthHandler.Info)
	
	// معالج المصادقة
	authHandler := NewAuthHandler(services.AuthService, services.UserService)
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/refresh", authHandler.RefreshToken)
	api.POST("/auth/forgot-password", authHandler.ForgotPassword)
	api.POST("/auth/reset-password", authHandler.ResetPassword)
	
	// معالج الموقع
	websiteHandler := NewWebsiteHandler(services.WebsiteService)
	api.GET("/website/content", websiteHandler.GetContent)
	api.GET("/website/features", websiteHandler.GetFeatures)
	api.GET("/website/testimonials", websiteHandler.GetTestimonials)
	
	// معالج الخدمات (العامة)
	servicesHandler := NewServicesHandler(services.ServicesService, services.AuthService)
	api.GET("/services", servicesHandler.GetServices)
	api.GET("/services/search", servicesHandler.SearchServices)
	api.GET("/services/featured", servicesHandler.GetFeaturedServices)
	api.GET("/services/categories", servicesHandler.GetAllCategories)
	api.GET("/services/tags/popular", servicesHandler.GetPopularTags)
	api.GET("/services/popular", servicesHandler.GetPopularServices)
	api.GET("/services/category/:category", servicesHandler.GetServicesByCategory)
	api.GET("/services/tag/:tag", servicesHandler.GetServicesByTag)
	api.GET("/services/:serviceId", servicesHandler.GetServiceDetails)
	api.GET("/services/:serviceId/recommended", servicesHandler.GetRecommendedServices)
	api.GET("/services/:serviceId/similar", servicesHandler.GetSimilarServices)
	api.GET("/services/seller/:sellerId", servicesHandler.GetSellerServices)
	api.GET("/services/:serviceId/ratings", servicesHandler.GetServiceRatings)
}

// registerProtectedRoutes تسجيل المسارات المحمية
func registerProtectedRoutes(api *gin.RouterGroup, services *ServiceContainer, middlewares *MiddlewareContainer) {
	protected := api.Group("")
	protected.Use(middlewares.AuthMiddleware)
	
	// معالج المستخدم
	userHandler := NewUserHandler(services.UserService)
	protected.GET("/user/profile", userHandler.GetProfile)
	protected.PUT("/user/profile", userHandler.UpdateProfile)
	protected.PUT("/user/password", userHandler.ChangePassword)
	protected.GET("/user/notifications", userHandler.GetNotifications)
	protected.PUT("/user/notifications/:id/read", userHandler.MarkNotificationRead)
	
	// معالج السلة
	cartHandler := NewCartHandler(services.CartService)
	protected.GET("/cart", cartHandler.GetCart)
	protected.POST("/cart/items", cartHandler.AddToCart)
	protected.PUT("/cart/items/:itemId", cartHandler.UpdateCartItem)
	protected.DELETE("/cart/items/:itemId", cartHandler.RemoveFromCart)
	protected.DELETE("/cart", cartHandler.ClearCart)
	
	// معالج الطلبات
	ordersHandler := NewOrdersHandler(services.OrdersService)
	protected.GET("/orders", ordersHandler.GetUserOrders)
	protected.GET("/orders/:orderId", ordersHandler.GetOrderDetails)
	protected.POST("/orders", ordersHandler.CreateOrder)
	protected.PUT("/orders/:orderId/cancel", ordersHandler.CancelOrder)
	
	// معالج التقييمات
	protected.POST("/services/:serviceId/ratings", servicesHandler.AddRating)
	protected.PUT("/ratings/:ratingId", servicesHandler.UpdateRating)
	protected.DELETE("/ratings/:ratingId", servicesHandler.DeleteRating)
	protected.GET("/my/ratings", servicesHandler.GetMyRatings)
	
	// معالج المفضلة
	protected.POST("/services/:serviceId/favorite", servicesHandler.AddToFavorites)
	protected.DELETE("/services/:serviceId/favorite", servicesHandler.RemoveFromFavorites)
	protected.GET("/user/favorites", servicesHandler.GetFavoriteServices)
	
	// معالج التوفر
	protected.POST("/services/:serviceId/check-availability", servicesHandler.CheckAvailability)
}

// registerAdminRoutes تسجيل مسارات المسؤولين
func registerAdminRoutes(api *gin.RouterGroup, services *ServiceContainer, middlewares *MiddlewareContainer) {
	admin := api.Group("/admin")
	admin.Use(middlewares.AuthMiddleware, middlewares.AdminMiddleware)
	
	// معالج الإدارة
	adminHandler := NewAdminHandler(services.AdminService)
	admin.GET("/dashboard", adminHandler.GetDashboard)
	admin.GET("/dashboard/stats", adminHandler.GetDashboardStats)
	admin.GET("/system/health", adminHandler.GetSystemHealth)
	admin.GET("/system/metrics", adminHandler.GetSystemMetrics)
	admin.POST("/system/maintenance", adminHandler.SetMaintenanceMode)
	admin.GET("/users", adminHandler.GetUsers)
	admin.PUT("/users/:userId/status", adminHandler.UpdateUserStatus)
	admin.PUT("/users/:userId/role", adminHandler.UpdateUserRole)
	
	// معالج الخدمات (الإدارة)
	servicesHandler := NewServicesHandler(services.ServicesService, services.AuthService)
	admin.GET("/services", servicesHandler.GetAllServices)
	admin.PUT("/services/:serviceId/status", servicesHandler.AdminUpdateServiceStatus)
	admin.PUT("/services/:serviceId/featured", servicesHandler.AdminUpdateFeaturedStatus)
	admin.DELETE("/services/:serviceId", servicesHandler.AdminDeleteService)
	admin.GET("/services/stats/overview", servicesHandler.GetAdminServicesStats)
	admin.GET("/services/stats/categories", servicesHandler.GetCategoriesStats)
	admin.GET("/services/stats/growth", servicesHandler.GetAdminServicesGrowth)
	admin.GET("/services/reports/popular", servicesHandler.GetPopularServicesReport)
	admin.DELETE("/ratings/:ratingId", servicesHandler.AdminDeleteRating)
	admin.GET("/ratings/reported", servicesHandler.GetReportedRatings)
	
	// معالج التقارير
	reportsHandler := NewReportsHandler(services.ReportsService)
	admin.GET("/reports/services", reportsHandler.GetServicesReport)
	admin.GET("/reports/users", reportsHandler.GetUsersReport)
	admin.GET("/reports/orders", reportsHandler.GetOrdersReport)
	admin.GET("/reports/financial", reportsHandler.GetFinancialReport)
	admin.POST("/reports/generate", reportsHandler.GenerateReport)
	
	// معالج التحليلات
	analyticsHandler := NewAnalyticsHandler(services.AnalyticsService)
	admin.GET("/analytics/overview", analyticsHandler.GetOverview)
	admin.GET("/analytics/services", analyticsHandler.GetServicesAnalytics)
	admin.GET("/analytics/users", analyticsHandler.GetUsersAnalytics)
	admin.GET("/analytics/sales", analyticsHandler.GetSalesAnalytics)
	
	// معالج التخزين المؤقت (الإدارة)
	cacheHandler := NewCacheHandler(services.CacheService)
	admin.DELETE("/cache/flush", cacheHandler.Flush)
	admin.DELETE("/cache/flush-pattern", cacheHandler.FlushPattern)
	admin.GET("/cache/stats", cacheHandler.Stats)
	admin.GET("/cache/health", cacheHandler.Health)
	
	// معالج الصحة (الإدارة)
	healthHandler := NewHealthHandler(services.HealthService, config.Load())
	admin.GET("/health/admin", healthHandler.AdminHealth)
	admin.GET("/health/detailed", healthHandler.Detailed)
	admin.GET("/health/metrics", healthHandler.Metrics)
}

// registerSellerRoutes تسجيل مسارات البائعين
func registerSellerRoutes(api *gin.RouterGroup, services *ServiceContainer, middlewares *MiddlewareContainer) {
	seller := api.Group("/seller")
	seller.Use(middlewares.AuthMiddleware, middlewares.SellerMiddleware)
	
	// معالج الخدمات (البائعين)
	servicesHandler := NewServicesHandler(services.ServicesService, services.AuthService)
	seller.POST("/services", servicesHandler.CreateService)
	seller.PUT("/services/:serviceId", servicesHandler.UpdateService)
	seller.PATCH("/services/:serviceId/status", servicesHandler.UpdateServiceStatus)
	seller.DELETE("/services/:serviceId", servicesHandler.DeleteService)
	seller.GET("/services/my", servicesHandler.GetMyServices)
	seller.GET("/services/my/stats", servicesHandler.GetServicesStats)
	seller.GET("/services/my/stats/status", servicesHandler.GetServicesStatusCount)
	seller.GET("/services/my/growth", servicesHandler.GetServicesGrowth)
	seller.POST("/services/search/advanced", servicesHandler.AdvancedSearch)
	seller.POST("/services/:serviceId/time-slots", servicesHandler.CreateTimeSlot)
	seller.GET("/services/:serviceId/time-slots", servicesHandler.GetTimeSlots)
	
	// معالج الطلبات (البائعين)
	ordersHandler := NewOrdersHandler(services.OrdersService)
	seller.GET("/orders", ordersHandler.GetSellerOrders)
	seller.PUT("/orders/:orderId/status", ordersHandler.UpdateOrderStatus)
	seller.GET("/orders/stats", ordersHandler.GetSellerOrdersStats)
	
	// معالج الرفع (البائعين)
	uploadsHandler := NewUploadsHandler(services.UploadService)
	seller.POST("/uploads/images", uploadsHandler.UploadImage)
	seller.POST("/uploads/documents", uploadsHandler.UploadDocument)
	seller.DELETE("/uploads/:fileId", uploadsHandler.DeleteFile)
}

// registerRealTimeRoutes تسجيل مسارات الوقت الحقيقي
func registerRealTimeRoutes(api *gin.RouterGroup, services *ServiceContainer, middlewares *MiddlewareContainer) {
	// مسارات SSE
	sseHandler := sse.NewSSEHandler()
	api.GET("/sse/events", middlewares.AuthMiddleware, sseHandler.Handler)
	api.GET("/sse/notifications", middlewares.AuthMiddleware, sseHandler.NotificationHandler)
	api.GET("/sse/admin/events", middlewares.AuthMiddleware, middlewares.AdminMiddleware, sseHandler.AdminHandler)
	
	// معالج الإشعارات
	notificationHandler := NewNotificationHandler(services.NotificationService)
	api.GET("/notifications", middlewares.AuthMiddleware, notificationHandler.GetNotifications)
	api.PUT("/notifications/:id/read", middlewares.AuthMiddleware, notificationHandler.MarkAsRead)
	api.PUT("/notifications/read-all", middlewares.AuthMiddleware, notificationHandler.MarkAllAsRead)
	api.POST("/notifications", middlewares.AuthMiddleware, middlewares.AdminMiddleware, notificationHandler.SendNotification)
}

// registerWebhookRoutes تسجيل مسارات الويب هووك والتحليلات
func registerWebhookRoutes(api *gin.RouterGroup, services *ServiceContainer, middlewares *MiddlewareContainer) {
	// مسارات ويب هووك للخدمات الخارجية
	webhook := api.Group("/webhook")
	{
		// Stripe webhooks
		webhook.POST("/stripe/payment-success", services.PaymentService.HandleStripeWebhook)
		webhook.POST("/stripe/payment-failed", services.PaymentService.HandleStripeWebhook)
		webhook.POST("/stripe/subscription", services.PaymentService.HandleStripeWebhook)
		
		// Cloudinary webhooks
		webhook.POST("/cloudinary/upload", services.UploadService.HandleCloudinaryWebhook)
		webhook.POST("/cloudinary/delete", services.UploadService.HandleCloudinaryWebhook)
	}
	
	// مسارات التحليلات البديلة
	analytics := api.Group("/analytics")
	{
		// Plausible analytics
		analytics.POST("/plausible/event", services.AnalyticsService.HandlePlausibleEvent)
		analytics.GET("/plausible/stats", services.AnalyticsService.GetPlausibleStats)
		
		// Matomo analytics
		analytics.POST("/matomo/track", services.AnalyticsService.HandleMatomoTrack)
		analytics.GET("/matomo/stats", services.AnalyticsService.GetMatomoStats)
	}
}

// HealthHandler معالج الصحة المبسط للمسارات الأساسية
type HealthHandler struct {
	healthService services.HealthService
	config        *config.Config
}

func NewHealthHandler(healthService services.HealthService, config *config.Config) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
		config:        config,
	}
}

func (h *HealthHandler) Check(c *gin.Context) {
	response := gin.H{
		"status":    "healthy",
		"service":   "nawthtech-backend",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   h.config.Version,
	}
	c.JSON(200, response)
}

func (h *HealthHandler) Live(c *gin.Context) {
	response := gin.H{
		"status":  "alive",
		"message": "الخدمة حية وتعمل",
	}
	c.JSON(200, response)
}

func (h *HealthHandler) Ready(c *gin.Context) {
	response := gin.H{
		"status":  "ready",
		"message": "الخدمة جاهزة لمعالجة الطلبات",
	}
	c.JSON(200, response)
}

func (h *HealthHandler) Info(c *gin.Context) {
	response := gin.H{
		"name":        "NawthTech Backend",
		"version":     h.config.Version,
		"environment": h.config.Environment,
		"timestamp":   time.Now().Format(time.RFC3339),
	}
	c.JSON(200, response)
}