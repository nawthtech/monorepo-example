package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/handlers/health"
	"gorm.io/gorm"
)

// RegisterHealthRoutes تسجيل مسارات الصحة
func RegisterHealthRoutes(router *gin.RouterGroup, db *gorm.DB, healthService services.HealthService, version, environment string, adminMiddleware gin.HandlerFunc) {
	healthHandler := health.NewHealthHandler(db, healthService, version, environment)
	
	healthRoutes := router.Group("/health")
	{
		// مسارات عامة
		healthRoutes.GET("", healthHandler.Check)
		healthRoutes.GET("/live", healthHandler.Live)
		healthRoutes.GET("/ready", healthHandler.Ready)
		healthRoutes.GET("/info", healthHandler.Info)
		healthRoutes.GET("/detailed", healthHandler.Detailed)
		healthRoutes.GET("/metrics", healthHandler.Metrics)
		
		// مسارات المسؤولين
		adminHealth := healthRoutes.Group("/admin")
		adminHealth.Use(adminMiddleware)
		{
			adminHealth.GET("", healthHandler.AdminHealth)
		}
	}
}