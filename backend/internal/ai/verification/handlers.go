// backend/internal/ai/verification/handlers.go
package verification

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Handler معالج طلبات التحقق
type Handler struct {
	verifier *Verifier
	logger   *logrus.Logger
}

// NewHandler إنشاء معالج جديد
func NewHandler(verifier *Verifier, logger *logrus.Logger) *Handler {
	return &Handler{
		verifier: verifier,
		logger:   logger,
	}
}

// VerifyHandler معالج التحقق الفردي
func (h *Handler) VerifyHandler(c *gin.Context) {
	startTime := time.Now()
	
	var req VerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warn("Invalid verification request")
		c.JSON(http.StatusBadRequest, VerificationResponse{
			Success: false,
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}
	
	// التحقق من صحة الإدخال
	if err := ValidateInput(req.Content, 10000); err != nil {
		c.JSON(http.StatusBadRequest, VerificationResponse{
			Success: false,
			Error:   "Invalid input",
			Message: err.Error(),
		})
		return
	}
	
	// بناء خيارات التحقق
	opts := []VerificationOption{}
	if req.Model != "" {
		opts = append(opts, WithModel(req.Model))
	}
	if req.Type != "" {
		opts = append(opts, WithType(req.Type))
	}
	if req.Context != "" {
		opts = append(opts, WithContext(req.Context))
	}
	
	// تطبيق الخيارات الإضافية
	if req.Options != nil {
		if temp, ok := req.Options["temperature"].(float64); ok {
			opts = append(opts, WithTemperature(float32(temp)))
		}
		if maxTokens, ok := req.Options["maxTokens"].(float64); ok {
			opts = append(opts, WithMaxTokens(int(maxTokens)))
		}
	}
	
	// تنفيذ التحقق
	result, err := h.verifier.Verify(c.Request.Context(), req.Content, opts...)
	if err != nil {
		h.logger.WithError(err).Error("Verification failed")
		c.JSON(http.StatusInternalServerError, VerificationResponse{
			Success: false,
			Error:   "Verification failed",
			Message: err.Error(),
		})
		return
	}
	
	// إضافة مقاييس الأداء
	result.Metrics.Latency = time.Since(startTime).Milliseconds()
	
	// تسجيل النتيجة
	h.logger.WithFields(logrus.Fields{
		"isValid":    result.IsValid,
		"confidence": result.Confidence,
		"latency":    result.Metrics.Latency,
		"model":      result.Metrics.Model,
	}).Info("Verification completed")
	
	c.JSON(http.StatusOK, VerificationResponse{
		Success: true,
		Data:    result,
	})
}

// BatchVerifyHandler معالج التحقق الدفعي
func (h *Handler) BatchVerifyHandler(c *gin.Context) {
	startTime := time.Now()
	
	var req BatchVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warn("Invalid batch verification request")
		c.JSON(http.StatusBadRequest, BatchVerificationResponse{
			Success: false,
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}
	
	// التحقق من صحة الإدخال
	if len(req.Contents) == 0 {
		c.JSON(http.StatusBadRequest, BatchVerificationResponse{
			Success: false,
			Error:   "No content provided",
			Message: "At least one content item is required",
		})
		return
	}
	
	if len(req.Contents) > 100 {
		c.JSON(http.StatusBadRequest, BatchVerificationResponse{
			Success: false,
			Error:   "Too many items",
			Message: "Maximum 100 items per batch",
		})
		return
	}
	
	// بناء خيارات التحقق
	opts := []VerificationOption{}
	if req.Model != "" {
		opts = append(opts, WithModel(req.Model))
	}
	if req.Type != "" {
		opts = append(opts, WithType(req.Type))
	}
	
	// تطبيق الخيارات الإضافية
	if req.Options != nil {
		if temp, ok := req.Options["temperature"].(float64); ok {
			opts = append(opts, WithTemperature(float32(temp)))
		}
		if maxTokens, ok := req.Options["maxTokens"].(float64); ok {
			opts = append(opts, WithMaxTokens(int(maxTokens)))
		}
		if batchSize, ok := req.Options["batchSize"].(float64); ok {
			// Note: Batch size is handled by the verifier
		}
	}
	
	// تنفيذ التحقق الدفعي
	results, err := h.verifier.BatchVerify(c.Request.Context(), req.Contents, opts...)
	if err != nil {
		h.logger.WithError(err).Error("Batch verification failed")
		c.JSON(http.StatusInternalServerError, BatchVerificationResponse{
			Success: false,
			Error:   "Batch verification failed",
			Message: err.Error(),
		})
		return
	}
	
	// حساب الإحصائيات
	stats := CalculateBatchStats(results)
	stats.ProcessingTime = time.Since(startTime)
	
	// تسجيل النتيجة
	h.logger.WithFields(logrus.Fields{
		"total":            stats.Total,
		"valid":            stats.Valid,
		"invalid":          stats.Invalid,
		"average_confidence": stats.AverageConfidence,
		"processing_time":   stats.ProcessingTime.Milliseconds(),
	}).Info("Batch verification completed")
	
	c.JSON(http.StatusOK, BatchVerificationResponse{
		Success: true,
		Data:    stats,
	})
}

// HealthHandler معالج فحص الصحة
func (h *Handler) HealthHandler(c *gin.Context) {
	result := HealthCheckResult{
		Status:      "healthy",
		LLMProvider: "openai",
		Model:       h.verifier.model,
		Connected:   true,
		Timestamp:   time.Now(),
		Message:     "Service is running",
	}
	
	// اختبار الاتصال بـ LLM
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	
	startTime := time.Now()
	connected, err := h.verifier.TestConnection(ctx)
	result.Latency = time.Since(startTime).Milliseconds()
	
	if err != nil || !connected {
		result.Status = "unhealthy"
		result.Connected = false
		result.Message = fmt.Sprintf("LLM connection failed: %v", err)
		
		h.logger.WithError(err).Warn("Health check failed")
		c.JSON(http.StatusServiceUnavailable, result)
		return
	}
	
	c.JSON(http.StatusOK, result)
}

// StatsHandler معالج الإحصائيات
func (h *Handler) StatsHandler(c *gin.Context) {
	// في تطبيق حقيقي، هذه البيانات ستأتي من قاعدة البيانات
	stats := Stats{
		TotalVerifications: 1000,
		Successful:         850,
		Failed:             150,
		AverageConfidence:  0.78,
		TotalCost:          12.50,
		TotalTokens:        125000,
		MostCommonIssues: map[string]int{
			"toxicity":     45,
			"factuality":   32,
			"coherence":    28,
			"safety":       25,
			"relevance":    20,
		},
		VerificationByType: map[string]int{
			"content_safety": 600,
			"fact_check":     250,
			"quality":        150,
		},
		LastUpdated: time.Now(),
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// ConfigHandler معالج التكوين
func (h *Handler) ConfigHandler(c *gin.Context) {
	config := gin.H{
		"provider": "openai",
		"model":    h.verifier.model,
		"criteria": h.verifier.criteria,
		"limits": gin.H{
			"max_retries":  h.verifier.maxRetries,
			"timeout_ms":   h.verifier.timeout.Milliseconds(),
			"temperature":  h.verifier.temperature,
			"max_tokens":   h.verifier.maxTokens,
			"max_input_length": 10000,
			"max_batch_size":   100,
		},
		"features": []string{
			"single_verification",
			"batch_verification",
			"toxicity_check",
			"factuality_check",
			"safety_check",
			"moderation",
		},
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    config,
	})
}

// SetupRoutes إعداد مسارات API
func SetupRoutes(router *gin.Engine, handler *Handler) {
	api := router.Group("/api/v1/verification")
	{
		api.POST("/verify", handler.VerifyHandler)
		api.POST("/verify/batch", handler.BatchVerifyHandler)
		api.GET("/health", handler.HealthHandler)
		api.GET("/stats", handler.StatsHandler)
		api.GET("/config", handler.ConfigHandler)
		
		// مسارات الإدارة (تتطلب صلاحيات)
		admin := api.Group("/admin")
		admin.Use(AdminMiddleware())
		{
			admin.GET("/analytics", handler.StatsHandler)
			admin.POST("/test", handler.HealthHandler)
		}
	}
}

// AdminMiddleware middleware للتحقق من صلاحيات المسؤول
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// في تطبيق حقيقي، تحقق من التوكن أو الجلسة
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Admin access required",
			})
			c.Abort()
			return
		}
		
		// تحقق من صلاحيات المسؤول
		// (هذا مثال بسيط، في التطبيق الحقيقي استخدم نظام صلاحيات كامل)
		c.Next()
	}
}