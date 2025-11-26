package handlers

import (
	"github.com/gin-gonic/gin"
)

func RegisterWebsiteRoutes(router *gin.RouterGroup, websiteHandler *WebsiteHandler, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {
	websiteRoutes := router.Group("/website")
	
	// ==================== الإعدادات العامة ====================
	websiteRoutes.GET("/settings", websiteHandler.GetSettings)
	websiteRoutes.PUT("/settings", adminMiddleware, websiteHandler.UpdateSettings)
	
	// ==================== الذكاء الاصطناعي والإعدادات ====================
	websiteRoutes.GET("/settings/ai-optimized", adminMiddleware, websiteHandler.GetAIOptimizedSettings)
	websiteRoutes.PATCH("/settings/ai-optimize", adminMiddleware, websiteHandler.AIOptimizeSettings)
	
	// ==================== استراتيجيات المحتوى ====================
	websiteRoutes.POST("/strategy/generate", adminMiddleware, websiteHandler.GenerateContentStrategy)
	
	// ==================== التحليلات والرؤى ====================
	websiteRoutes.GET("/analytics/ai-insights", adminMiddleware, websiteHandler.GetAIAnalyticsInsights)
	
	// ==================== إنشاء المحتوى ====================
	websiteRoutes.POST("/content/generate", adminMiddleware, websiteHandler.GenerateContent)
	
	// ==================== التوقعات والتنبؤات ====================
	websiteRoutes.GET("/predictions/performance", adminMiddleware, websiteHandler.GetPerformancePredictions)
	
	// ==================== تحليل الجمهور ====================
	websiteRoutes.GET("/audience/insights", adminMiddleware, websiteHandler.GetAudienceInsights)
}