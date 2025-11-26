package handlers

import (
	"github.com/gin-gonic/gin"
)

func RegisterStrategiesRoutes(router *gin.RouterGroup, strategiesHandler *StrategiesHandler, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {
	strategiesRoutes := router.Group("/strategies")
	strategiesRoutes.Use(authMiddleware)
	
	// ==================== إدارة الاستراتيجيات ====================
	strategiesRoutes.POST("", adminMiddleware, strategiesHandler.CreateStrategy)
	strategiesRoutes.GET("", adminMiddleware, strategiesHandler.GetStrategies)
	strategiesRoutes.GET("/:id", adminMiddleware, strategiesHandler.GetStrategyByID)
	strategiesRoutes.PUT("/:id", adminMiddleware, strategiesHandler.UpdateStrategy)
	strategiesRoutes.DELETE("/:id", adminMiddleware, strategiesHandler.DeleteStrategy)
	
	// ==================== تحليل الاستراتيجيات ====================
	strategiesRoutes.POST("/:id/analyze", adminMiddleware, strategiesHandler.AnalyzeStrategy)
	strategiesRoutes.GET("/:id/performance", adminMiddleware, strategiesHandler.GetStrategyPerformance)
	
	// ==================== التوصيات الذكية ====================
	strategiesRoutes.POST("/recommend", adminMiddleware, strategiesHandler.GetStrategyRecommendations)
}