package handlers

import (
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/services"

	"github.com/go-chi/chi/v5"
)

// RegisterAdminRoutes تسجيل جميع مسارات الإدارة
func RegisterAdminRoutes(router chi.Router, adminService *services.AdminService) {
	adminHandler := NewAdminHandler(adminService)

	router.Route("/api/v1/admin", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware, middleware.AdminAuth)

		// ==================== تحديث النظام ====================
		r.Route("/system", func(system chi.Router) {
			// بدء عملية تحديث النظام مع تحليلات الذكاء الاصطناعي
			system.With(
				middleware.RateLimiter(5, 30*60*1000),
				middleware.ValidateAdminAction,
				middleware.UpdateMiddleware,
			).Post("/update", adminHandler.InitiateSystemUpdate)
			
			// الحصول على حالة النظام الشاملة مع تحليلات الذكاء الاصطناعي
			system.With(middleware.RateLimiter(30, 2*60*1000)).Get("/status", adminHandler.GetSystemStatus)
			
			// تحليل النظام باستخدام الذكاء الاصطناعي للحصول على رؤى متقدمة
			system.With(middleware.RateLimiter(15, 5*60*1000)).Get("/ai-analytics", adminHandler.GetAIAnalytics)
			
			// فحص صحة النظام بشكل مفصل مع تحليلات الذكاء الاصطناعي
			system.With(middleware.RateLimiter(10, 5*60*1000)).Get("/health", adminHandler.GetSystemHealth)
			
			// تفعيل/تعطيل وضع الصيانة مع التخطيط الذكي
			system.With(
				middleware.RateLimiter(10, 10*60*1000),
				middleware.ValidateAdminAction,
			).Post("/maintenance", adminHandler.SetMaintenanceMode)
			
			// الحصول على سجلات النظام مع تحليل الذكاء الاصطناعي
			system.With(middleware.RateLimiter(30, 2*60*1000)).Get("/logs", adminHandler.GetSystemLogs)
			
			// إنشاء نسخة احتياطية ذكية
			system.With(
				middleware.RateLimiter(5, 60*60*1000),
				middleware.ValidateAdminAction,
			).Post("/backup", adminHandler.CreateSystemBackup)
			
			// تحسين النظام تلقائياً باستخدام الذكاء الاصطناعي
			system.With(
				middleware.RateLimiter(10, 30*60*1000),
				middleware.ValidateAdminAction,
			).Post("/optimize", adminHandler.PerformOptimization)
		})

		// ==================== إدارة المستخدمين المتقدمة ====================
		r.Route("/users", func(users chi.Router) {
			// تحليل سلوك المستخدمين باستخدام الذكاء الاصطناعي
			users.With(middleware.RateLimiter(20, 10*60*1000)).Get("/analytics", adminHandler.GetUserAnalytics)
		})
	})
}