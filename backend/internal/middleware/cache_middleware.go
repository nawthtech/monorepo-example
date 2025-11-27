package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
)

// CacheMiddleware وسيط التخزين المؤقت
type CacheMiddleware struct {
	cacheService services.CacheService
	cacheHelper  *utils.CacheHelper
}

// NewCacheMiddleware إنشاء وسيط تخزين مؤقت جديد
func NewCacheMiddleware(cacheService services.CacheService, prefix string) *CacheMiddleware {
	return &CacheMiddleware{
		cacheService: cacheService,
		cacheHelper:  utils.NewCacheHelper(prefix),
	}
}

// CacheResponse وسيط تخزين استجابات API
func (m *CacheMiddleware) CacheResponse(ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// تجاهل طلبات غير GET
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// تجاهل بعض المسارات
		if m.shouldBypassCache(c.Request.URL.Path) {
			c.Next()
			return
		}

		// إنشاء مفتاح التخزين المؤقت
		cacheKey := m.generateCacheKey(c)

		// محاولة جلب البيانات من التخزين المؤقت
		cachedData, err := m.cacheService.Get(c.Request.Context(), cacheKey)
		if err == nil && cachedData != nil {
			// البيانات موجودة في التخزين المؤقت، إرجاعها مباشرة
			c.JSON(http.StatusOK, cachedData)
			c.Abort()
			return
		}

		// استبدال الكاتب لاعتراض الاستجابة
		writer := &cacheResponseWriter{
			ResponseWriter: c.Writer,
			body:           make([]byte, 0),
			statusCode:     http.StatusOK,
		}
		c.Writer = writer

		// معالجة الطلب
		c.Next()

		// تخزين الاستجابة في التخزين المؤقت إذا كانت ناجحة
		if writer.statusCode == http.StatusOK && len(writer.body) > 0 {
			var responseData interface{}
			if err := utils.ParseJSON(writer.body, &responseData); err == nil {
				m.cacheService.Set(c.Request.Context(), cacheKey, responseData, ttl)
			}
		}
	}
}

// InvalidateCache وسيط إبطال التخزين المؤقت
func (m *CacheMiddleware) InvalidateCache(patterns ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// معالجة الطلب أولاً
		c.Next()

		// إبطال التخزين المؤقت بعد التعديلات
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			for _, pattern := range patterns {
				m.cacheService.FlushPattern(c.Request.Context(), pattern)
			}
		}
	}
}

// RateLimitMiddleware وسيط تحديد المعدل باستخدام التخزين المؤقت
func (m *CacheMiddleware) RateLimitMiddleware(requestsPerMinute int, keyPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		endpoint := c.Request.URL.Path
		minute := time.Now().Format("2006-01-02T15:04")

		rateLimitKey := m.cacheHelper.GenerateRateLimitKey(clientIP+":"+minute, endpoint)

		// زيادة العداد
		count, err := m.cacheService.Increment(c.Request.Context(), rateLimitKey, 1)
		if err != nil {
			c.Next()
			return
		}

		// تعيين وقت الصلاحية إذا كانت هذه هي الزيادة الأولى
		if count == 1 {
			m.cacheService.Expire(c.Request.Context(), rateLimitKey, time.Minute)
		}

		// التحقق من تجاوز الحد
		if count > int64(requestsPerMinute) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "تم تجاوز الحد المسموح للطلبات",
				"success": false,
				"details": "الرجاء المحاولة مرة أخرى لاحقاً",
			})
			c.Abort()
			return
		}

		// إضافة الرؤوس للاستجابة
		c.Header("X-RateLimit-Limit", string(rune(requestsPerMinute)))
		c.Header("X-RateLimit-Remaining", string(rune(requestsPerMinute-int(count))))
		c.Header("X-RateLimit-Reset", time.Now().Add(time.Minute).Format(time.RFC1123))

		c.Next()
	}
}

// CacheControlMiddleware وسيط التحكم في التخزين المؤقت
func (m *CacheMiddleware) CacheControlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// إضافة رؤوس التحكم في التخزين المؤقت
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		c.Next()
	}
}

// CacheHealthMiddleware وسيط صحة التخزين المؤقت
func (m *CacheMiddleware) CacheHealthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// فحص صحة التخزين المؤقت
		health, err := m.cacheService.HealthCheck(c.Request.Context())
		if err != nil || health.Status != "healthy" {
			c.Header("X-Cache-Status", "unhealthy")
		} else {
			c.Header("X-Cache-Status", "healthy")
		}

		c.Next()
	}
}

// ========== الدوال المساعدة ==========

// shouldBypassCache التحقق إذا كان يجب تجاوز التخزين المؤقت
func (m *CacheMiddleware) shouldBypassCache(path string) bool {
	// قائمة المسارات التي يجب تجاوز التخزين المؤقت لها
	bypassPaths := []string{
		"/api/v1/auth",
		"/api/v1/admin",
		"/api/v1/health",
		"/api/v1/cache",
	}

	for _, bypassPath := range bypassPaths {
		if strings.HasPrefix(path, bypassPath) {
			return true
		}
	}

	return false
}

// generateCacheKey إنشاء مفتاح التخزين المؤقت
func (m *CacheMiddleware) generateCacheKey(c *gin.Context) string {
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	userID, _ := c.Get("userID")

	// إنشاء مفتاح فريد بناءً على المسار والاستعلام والمستخدم
	key := path
	if query != "" {
		key += "?" + query
	}
	if userID != nil {
		key += ":user:" + userID.(string)
	}

	return m.cacheHelper.GenerateAPIKey(path, utils.GenerateMD5Hash(key))
}

// cacheResponseWriter كاتب استجابة مخصص لاعتراض البيانات
type cacheResponseWriter struct {
	gin.ResponseWriter
	body       []byte
	statusCode int
}

// Write اعتراض كتابة البيانات
func (w *cacheResponseWriter) Write(data []byte) (int, error) {
	w.body = append(w.body, data...)
	return w.ResponseWriter.Write(data)
}

// WriteHeader اعتراض كتابة رأس الاستجابة
func (w *cacheResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}