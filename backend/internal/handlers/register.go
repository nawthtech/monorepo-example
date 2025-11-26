package handlers

import (
	"github.com/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/backend/internal/services"

	"github.com/gorilla/mux"
)

// Services يحتوي على جميع الخدمات
type Services struct {
	Admin   *services.AdminService
	User    *services.UserService
	Auth    *services.AuthService
	Store   *services.StoreService
	Cart    *services.CartService
	Payment *services.PaymentService
	AI      *services.AIService
	Email   *services.EmailService
	Upload  *services.UploadService
}

// Register تسجيل جميع المسارات
func Register(router *mux.Router, services *Services) {
	// إنشاء المعالجات
	adminHandler := NewAdminHandler(services.Admin)
	authHandler := NewAuthHandler(services.Auth)
	userHandler := NewUserHandler(services.User)
	storeHandler := NewStoreHandler(services.Store)
	cartHandler := NewCartHandler(services.Cart)
	paymentHandler := NewPaymentHandler(services.Payment)
	aiHandler := NewAIHandler(services.AI)
	uploadHandler := NewUploadHandler(services.Upload)
	healthHandler := NewHealthHandler()

	// وسيط المصادقة الأساسي
	authMiddleware := middleware.AuthMiddleware

	// المسارات العامة
	public := router.PathPrefix("/api").Subrouter()
	
	// الصحة
	public.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")
	
	// المصادقة
	public.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	public.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	public.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST")
	public.HandleFunc("/auth/forgot-password", authHandler.ForgotPassword).Methods("POST")
	public.HandleFunc("/auth/reset-password", authHandler.ResetPassword).Methods("POST")

	// المتجر (عام)
	public.HandleFunc("/store/services", storeHandler.GetServices).Methods("GET")
	public.HandleFunc("/store/services/{id}", storeHandler.GetService).Methods("GET")
	public.HandleFunc("/store/categories", storeHandler.GetCategories).Methods("GET")

	// الرفع (يحتاج مصادقة)
	uploadRouter := router.PathPrefix("/api/upload").Subrouter()
	uploadRouter.Use(authMiddleware)
	uploadRouter.HandleFunc("/image", uploadHandler.UploadImage).Methods("POST")
	uploadRouter.HandleFunc("/file", uploadHandler.UploadFile).Methods("POST")

	// المسارات المحمية
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(authMiddleware)

	// المستخدم
	protected.HandleFunc("/users/profile", userHandler.GetProfile).Methods("GET")
	protected.HandleFunc("/users/profile", userHandler.UpdateProfile).Methods("PUT")
	protected.HandleFunc("/users/change-password", userHandler.ChangePassword).Methods("POST")

	// المتجر المحمي
	protected.HandleFunc("/store/orders", storeHandler.CreateOrder).Methods("POST")
	protected.HandleFunc("/store/orders", storeHandler.GetUserOrders).Methods("GET")
	protected.HandleFunc("/store/orders/{id}", storeHandler.GetOrder).Methods("GET")

	// السلة
	protected.HandleFunc("/cart", cartHandler.GetCart).Methods("GET")
	protected.HandleFunc("/cart/items", cartHandler.AddToCart).Methods("POST")
	protected.HandleFunc("/cart/items/{id}", cartHandler.UpdateCartItem).Methods("PUT")
	protected.HandleFunc("/cart/items/{id}", cartHandler.RemoveFromCart).Methods("DELETE")
	protected.HandleFunc("/cart/clear", cartHandler.ClearCart).Methods("DELETE")

	// المدفوعات
	protected.HandleFunc("/payments/create", paymentHandler.CreatePayment).Methods("POST")
	protected.HandleFunc("/payments/{id}", paymentHandler.GetPayment).Methods("GET")
	protected.HandleFunc("/payments/{id}/verify", paymentHandler.VerifyPayment).Methods("POST")

	// الذكاء الاصطناعي
	protected.HandleFunc("/ai/analyze", aiHandler.AnalyzeContent).Methods("POST")
	protected.HandleFunc("/ai/recommend", aiHandler.GetRecommendations).Methods("GET")
	protected.HandleFunc("/ai/generate", aiHandler.GenerateContent).Methods("POST")

	// مسارات الإدارة (تحتاج صلاحيات إدارة)
	adminRouter := router.PathPrefix("/api/v1/admin").Subrouter()
	adminRouter.Use(authMiddleware)
	adminRouter.Use(middleware.AdminAuth)

	// تسجيل مسارات الإدارة
	RegisterAdminRoutes(adminRouter, adminHandler, authMiddleware)
}

// RegisterAdminRoutes تسجيل مسارات الإدارة
func RegisterAdminRoutes(router *mux.Router, adminHandler *AdminHandler, authMiddleware mux.MiddlewareFunc) {
	// مسارات تحديث النظام
	router.HandleFunc("/system/update", 
		middleware.RateLimiter(5, 30*60*1000)(
			middleware.ValidateAdminAction(
				middleware.UpdateMiddleware(
					adminHandler.InitiateSystemUpdate,
				),
			),
		),
	).Methods("POST")
	
	// مسارات حالة النظام
	router.HandleFunc("/system/status", 
		middleware.RateLimiter(30, 2*60*1000)(
			adminHandler.GetSystemStatus,
		),
	).Methods("GET")
	
	router.HandleFunc("/system/health", 
		middleware.RateLimiter(10, 5*60*1000)(
			adminHandler.GetSystemHealth,
		),
	).Methods("GET")
	
	// مسارات تحليلات الذكاء الاصطناعي
	router.HandleFunc("/system/ai-analytics", 
		middleware.RateLimiter(15, 5*60*1000)(
			adminHandler.GetAIAnalytics,
		),
	).Methods("GET")
	
	// مسارات تحليلات المستخدمين
	router.HandleFunc("/users/analytics", 
		middleware.RateLimiter(20, 10*60*1000)(
			adminHandler.GetUserAnalytics,
		),
	).Methods("GET")
	
	// مسارات الصيانة
	router.HandleFunc("/system/maintenance", 
		middleware.RateLimiter(10, 10*60*1000)(
			middleware.ValidateAdminAction(
				adminHandler.SetMaintenanceMode,
			),
		),
	).Methods("POST")
	
	// مسارات السجلات
	router.HandleFunc("/system/logs", 
		middleware.RateLimiter(30, 2*60*1000)(
			adminHandler.GetSystemLogs,
		),
	).Methods("GET")
	
	// مسارات النسخ الاحتياطي
	router.HandleFunc("/system/backup", 
		middleware.RateLimiter(5, 60*60*1000)(
			middleware.ValidateAdminAction(
				adminHandler.CreateSystemBackup,
			),
		),
	).Methods("POST")
	
	// مسارات التحسين
	router.HandleFunc("/system/optimize", 
		middleware.RateLimiter(10, 30*60*1000)(
			middleware.ValidateAdminAction(
				adminHandler.PerformOptimization,
			),
		),
	).Methods("POST")

	// إدارة المستخدمين
	router.HandleFunc("/users", adminHandler.GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", adminHandler.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", adminHandler.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", adminHandler.DeleteUser).Methods("DELETE")

	// إدارة الطلبات
	router.HandleFunc("/orders", adminHandler.GetOrders).Methods("GET")
	router.HandleFunc("/orders/{id}", adminHandler.GetOrder).Methods("GET")
	router.HandleFunc("/orders/{id}/status", adminHandler.UpdateOrderStatus).Methods("PUT")

	// إدارة الخدمات
	router.HandleFunc("/services", adminHandler.GetServices).Methods("GET")
	router.HandleFunc("/services", adminHandler.CreateService).Methods("POST")
	router.HandleFunc("/services/{id}", adminHandler.UpdateService).Methods("PUT")
	router.HandleFunc("/services/{id}", adminHandler.DeleteService).Methods("DELETE")
}
