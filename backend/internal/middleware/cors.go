package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/config"
)

// CORS middleware لإعدادات CORS الديناميكية
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// الحصول على إعدادات CORS بناءً على المسار
		corsConfig := config.GetCORSConfig(c.Request.URL.Path)
		
		origin := c.Request.Header.Get("Origin")
		
		// التحقق من النطاق المسموح به
		if !config.ValidateOrigin(origin) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Not allowed by CORS",
				"message": "النطاق غير مسموح به",
			})
			return
		}
		
		// تعيين رؤوس CORS
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}
		
		c.Header("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowedMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowedHeaders, ", "))
		c.Header("Access-Control-Expose-Headers", strings.Join(corsConfig.ExposedHeaders, ", "))
		
		if corsConfig.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		
		if corsConfig.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", string(rune(corsConfig.MaxAge)))
		}
		
		// معالجة طلبات Preflight
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// SecurityHeaders middleware لأمان إضافي
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		c.Next()
	}
}