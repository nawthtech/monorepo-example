package handlers

import (
	"github.com/gin-gonic/gin"
)

func RegisterContentRoutes(router *gin.RouterGroup, contentHandler *ContentHandler, authMiddleware gin.HandlerFunc) {
	contentRoutes := router.Group("/content")
	contentRoutes.Use(authMiddleware)
	
	// ==================== إنشاء المحتوى ====================
	contentRoutes.POST("/generate", contentHandler.GenerateContent)
	contentRoutes.POST("/batch-generate", contentHandler.BatchGenerateContent)
	
	// ==================== إدارة المحتوى ====================
	contentRoutes.GET("", contentHandler.GetContent)
	contentRoutes.GET("/:id", contentHandler.GetContentByID)
	contentRoutes.PUT("/:id", contentHandler.UpdateContent)
	contentRoutes.DELETE("/:id", contentHandler.DeleteContent)
	
	// ==================== تحليل المحتوى ====================
	contentRoutes.POST("/:id/analyze", contentHandler.AnalyzeContent)
	contentRoutes.POST("/optimize", contentHandler.OptimizeContent)
	
	// ==================== تحليل الأداء ====================
	contentRoutes.GET("/:id/performance", contentHandler.GetContentPerformance)
}