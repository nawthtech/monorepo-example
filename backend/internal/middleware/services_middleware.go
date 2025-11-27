package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
)

// ServicesMiddleware وسيط الخدمات
type ServicesMiddleware struct {
	servicesService services.ServicesService
	cacheService    services.CacheService
}

// NewServicesMiddleware إنشاء وسيط خدمات جديد
func NewServicesMiddleware(servicesService services.ServicesService, cacheService services.CacheService) *ServicesMiddleware {
	return &ServicesMiddleware{
		servicesService: servicesService,
		cacheService:    cacheService,
	}
}

// ValidateServiceOwner وسيط التحقق من مالك الخدمة
func (m *ServicesMiddleware) ValidateServiceOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceID := c.Param("serviceId")
		if serviceID == "" {
			serviceID = c.Param("id")
		}

		if serviceID == "" {
			utils.ErrorResponse(c, http.StatusBadRequest, "معرف الخدمة مطلوب", "SERVICE_ID_REQUIRED")
			return
		}

		userID, exists := c.Get("userID")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
			return
		}

		// جلب الخدمة والتحقق من المالك
		service, err := m.servicesService.GetServiceDetails(c.Request.Context(), serviceID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusNotFound, "الخدمة غير موجودة", "SERVICE_NOT_FOUND")
			return
		}

		if service.SellerID != userID.(string) {
			// التحقق إذا كان المستخدم مسؤولاً
			userRole, _ := c.Get("userRole")
			if userRole != "admin" {
				utils.ErrorResponse(c, http.StatusForbidden, "غير مصرح بتعديل هذه الخدمة", "NOT_SERVICE_OWNER")
				return
			}
		}

		// تخزين الخدمة في السياق للاستخدام لاحقاً
		c.Set("service", service)
		c.Next()
	}
}

// ValidateServiceStatus وسيط التحقق من حالة الخدمة
func (m *ServicesMiddleware) ValidateServiceStatus(allowedStatuses ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		service, exists := c.Get("service")
		if !exists {
			utils.ErrorResponse(c, http.StatusInternalServerError, "خطأ في تحميل الخدمة", "SERVICE_LOAD_ERROR")
			return
		}

		serviceObj := service.(*models.Service)
		isAllowed := false

		for _, status := range allowedStatuses {
			if serviceObj.Status == status {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			utils.ErrorResponse(c, http.StatusBadRequest, 
				"حالة الخدمة لا تسمح بهذا الإجراء", 
				"INVALID_SERVICE_STATUS")
			return
		}

		c.Next()
	}
}

// RateLimitServices وسيط تحديد معدل طلبات الخدمات
func (m *ServicesMiddleware) RateLimitServices(requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.Next()
			return
		}

		endpoint := c.Request.URL.Path
		minute := time.Now().Format("2006-01-02T15:04")
		rateLimitKey := fmt.Sprintf("rate:services:%s:%s:%s", userID.(string), endpoint, minute)

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
			utils.ErrorResponse(c, http.StatusTooManyRequests, 
				"تم تجاوز الحد المسموح لطلبات الخدمات", 
				"SERVICES_RATE_LIMIT_EXCEEDED")
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(requestsPerMinute))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(requestsPerMinute-int(count)))
		
		c.Next()
	}
}

// ValidateServiceCreation وسيط التحقق من إنشاء الخدمة
func (m *ServicesMiddleware) ValidateServiceCreation() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
			return
		}

		// التحقق من عدد الخدمات النشطة (يمكن إضافة حدود حسب احتياجات العمل)
		services, _, err := m.servicesService.GetSellerServices(c.Request.Context(), services.GetSellerServicesParams{
			SellerID: userID.(string),
			Status:   "active",
			Limit:    1000, // رقم كبير لعد جميع الخدمات
		})

		if err == nil && len(services) >= 50 { // حد 50 خدمة نشطة
			utils.ErrorResponse(c, http.StatusBadRequest, 
				"لقد وصلت إلى الحد الأقصى للخدمات النشطة", 
				"MAX_ACTIVE_SERVICES_REACHED")
			return
		}

		c.Next()
	}
}

// CacheServices وسيط التخزين المؤقت للخدمات
func (m *ServicesMiddleware) CacheServices(ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// تجاهل طلبات غير GET
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// إنشاء مفتاح التخزين المؤقت
		cacheKey := m.generateServicesCacheKey(c)

		// محاولة جلب البيانات من التخزين المؤقت
		cachedData, err := m.cacheService.Get(c.Request.Context(), cacheKey)
		if err == nil && cachedData != nil {
			c.JSON(http.StatusOK, cachedData)
			c.Abort()
			return
		}

		// استبدال الكاتب لاعتراض الاستجابة
		writer := &servicesResponseWriter{
			ResponseWriter: c.Writer,
			body:           make([]byte, 0),
			statusCode:     http.StatusOK,
		}
		c.Writer = writer

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

// ValidateServiceSearchParams وسيط التحقق من معاملات البحث
func (m *ServicesMiddleware) ValidateServiceSearchParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		// التحقق من الصفحة والحد
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

		if page < 1 {
			page = 1
		}

		if limit < 1 || limit > 100 {
			limit = 20
		}

		// التحقق من السعر
		minPrice, err := strconv.ParseFloat(c.Query("minPrice"), 64)
		if err != nil && c.Query("minPrice") != "" {
			utils.ErrorResponse(c, http.StatusBadRequest, "سعر البداية غير صالح", "INVALID_MIN_PRICE")
			return
		}

		maxPrice, err := strconv.ParseFloat(c.Query("maxPrice"), 64)
		if err != nil && c.Query("maxPrice") != "" {
			utils.ErrorResponse(c, http.StatusBadRequest, "سعر النهاية غير صالح", "INVALID_MAX_PRICE")
			return
		}

		if minPrice > 0 && maxPrice > 0 && minPrice > maxPrice {
			utils.ErrorResponse(c, http.StatusBadRequest, "سعر البداية يجب أن يكون أقل من سعر النهاية", "INVALID_PRICE_RANGE")
			return
		}

		// تخزين المعاملات المصححة في السياق
		c.Set("validatedPage", page)
		c.Set("validatedLimit", limit)
		c.Set("validatedMinPrice", minPrice)
		c.Set("validatedMaxPrice", maxPrice)

		c.Next()
	}
}

// LogServicesActivity وسيط تسجيل نشاط الخدمات
func (m *ServicesMiddleware) LogServicesActivity() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		c.Next()

		duration := time.Since(start)
		userID, _ := c.Get("userID")
		method := c.Request.Method
		path := c.Request.URL.Path
		status := c.Writer.Status()

		// تسجيل النشاط (يمكن إرساله إلى نظام التحليلات)
		activity := map[string]interface{}{
			"user_id":    userID,
			"method":     method,
			"path":       path,
			"status":     status,
			"duration":   duration.Milliseconds(),
			"timestamp":  time.Now(),
			"user_agent": c.Request.UserAgent(),
			"ip":         c.ClientIP(),
		}

		// تخزين في التخزين المؤقت للنشاطات الحديثة
		activityKey := fmt.Sprintf("activity:services:%s:%d", userID, time.Now().UnixNano())
		m.cacheService.Set(c.Request.Context(), activityKey, activity, 24*time.Hour)
	}
}

// ========== الدوال المساعدة ==========

// generateServicesCacheKey إنشاء مفتاح التخزين المؤقت للخدمات
func (m *ServicesMiddleware) generateServicesCacheKey(c *gin.Context) string {
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	userID, _ := c.Get("userID")

	key := "services:" + path
	if query != "" {
		key += "?" + query
	}
	if userID != nil {
		key += ":user:" + userID.(string)
	}

	return utils.GenerateMD5Hash(key)
}

// servicesResponseWriter كاتب استجابة مخصص للخدمات
type servicesResponseWriter struct {
	gin.ResponseWriter
	body       []byte
	statusCode int
}

// Write اعتراض كتابة البيانات
func (w *servicesResponseWriter) Write(data []byte) (int, error) {
	w.body = append(w.body, data...)
	return w.ResponseWriter.Write(data)
}

// WriteHeader اعتراض كتابة رأس الاستجابة
func (w *servicesResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}