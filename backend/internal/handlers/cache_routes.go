package handlers

import (
	"github.com/gin-gonic/gin"
)

// RegisterCacheRoutes تسجيل مسارات التخزين المؤقت
func RegisterCacheRoutes(router *gin.RouterGroup, cacheHandler *CacheHandler, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {
	cacheRoutes := router.Group("/cache")
	cacheRoutes.Use(authMiddleware) // جميع مسارات التخزين المؤقت تتطلب مصادقة
	
	{
		// العمليات الأساسية
		cacheRoutes.POST("", cacheHandler.Set)
		cacheRoutes.GET("", cacheHandler.Get)
		cacheRoutes.DELETE("", cacheHandler.Delete)
		
		// عمليات التحقق
		cacheRoutes.GET("/exists", cacheHandler.Exists)
		cacheRoutes.GET("/ttl", cacheHandler.TTL)
		
		// العمليات الرقمية
		cacheRoutes.POST("/increment", cacheHandler.Increment)
		
		// عمليات القوائم
		cacheRoutes.POST("/lpush", cacheHandler.LPush)
		cacheRoutes.GET("/lrange", cacheHandler.LRange)
		
		// عمليات الهاش
		cacheRoutes.POST("/hset", cacheHandler.HSet)
		cacheRoutes.GET("/hget", cacheHandler.HGet)
		cacheRoutes.GET("/hgetall", cacheHandler.HGetAll)
		
		// البحث
		cacheRoutes.GET("/keys", cacheHandler.Keys)
	}
	
	// مسارات الإدارة (للمسؤولين فقط)
	adminCacheRoutes := cacheRoutes.Group("/admin")
	adminCacheRoutes.Use(adminMiddleware)
	{
		adminCacheRoutes.DELETE("/flush", cacheHandler.Flush)
		adminCacheRoutes.DELETE("/flush-pattern", cacheHandler.FlushPattern)
		adminCacheRoutes.GET("/stats", cacheHandler.Stats)
		adminCacheRoutes.GET("/health", cacheHandler.Health)
	}
}