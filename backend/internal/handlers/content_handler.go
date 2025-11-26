package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/middleware"
	"github.com/nawthtech/nawthtech/backend/internal/models"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
)

type ContentHandler struct {
	contentService services.ContentService
	authService    services.AuthService
}

func NewContentHandler(contentService services.ContentService, authService services.AuthService) *ContentHandler {
	return &ContentHandler{
		contentService: contentService,
		authService:    authService,
	}
}

// GenerateContentRequest - طلب إنشاء المحتوى
type GenerateContentRequest struct {
	Topic       string   `json:"topic" binding:"required"`
	Platform    string   `json:"platform" binding:"required"`
	ContentType string   `json:"contentType" binding:"required"`
	Tone        string   `json:"tone"`
	Keywords    []string `json:"keywords"`
	Language    string   `json:"language"`
	Length      string   `json:"length"`
	Style       string   `json:"style"`
}

// GenerateContent - إنشاء محتوى باستخدام الذكاء الاصطناعي
// @Summary إنشاء محتوى باستخدام الذكاء الاصطناعي
// @Description إنشاء محتوى باستخدام الذكاء الاصطناعي (للمشرفين فقط)
// @Tags Content
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body GenerateContentRequest true "بيانات إنشاء المحتوى"
// @Success 200 {object} utils.Response
// @Router /api/v1/content/generate [post]
func (h *ContentHandler) GenerateContent(c *gin.Context) {
	var req GenerateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	// الحصول على المستخدم من السياق
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	// التحقق من صلاحيات المشرف
	if !h.authService.IsAdmin(c, userID.(string)) {
		utils.ErrorResponse(c, http.StatusForbidden, "غير مصرح", "FORBIDDEN")
		return
	}

	// تطبيق rate limiting
	if !middleware.CheckRateLimit(c, "content_generate", 25, 10*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	// إنشاء المحتوى باستخدام الذكاء الاصطناعي
	content, err := h.contentService.GenerateContent(c, services.GenerateContentParams{
		Topic:       req.Topic,
		Platform:    req.Platform,
		ContentType: req.ContentType,
		Tone:        req.Tone,
		Keywords:    req.Keywords,
		Language:    req.Language,
		Length:      req.Length,
		Style:       req.Style,
		UserID:      userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في إنشاء المحتوى", "GENERATION_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم إنشاء المحتوى بنجاح", content)
}

// BatchGenerateContentRequest - طلب إنشاء محتوى جماعي
type BatchGenerateContentRequest struct {
	Topics    []string               `json:"topics" binding:"required"`
	Platforms []string               `json:"platforms" binding:"required"`
	Schedule  map[string]interface{} `json:"schedule"`
	ContentPlan map[string]interface{} `json:"contentPlan"`
}

// BatchGenerateContent - إنشاء مجموعة محتوى باستخدام الذكاء الاصطناعي
// @Summary إنشاء مجموعة محتوى باستخدام الذكاء الاصطناعي
// @Description إنشاء مجموعة محتوى باستخدام الذكاء الاصطناعي (للمشرفين فقط)
// @Tags Content
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body BatchGenerateContentRequest true "بيانات إنشاء المحتوى الجماعي"
// @Success 200 {object} utils.Response
// @Router /api/v1/content/batch-generate [post]
func (h *ContentHandler) BatchGenerateContent(c *gin.Context) {
	var req BatchGenerateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !h.authService.IsAdmin(c, userID.(string)) {
		utils.ErrorResponse(c, http.StatusForbidden, "غير مصرح", "FORBIDDEN")
		return
	}

	if !middleware.CheckRateLimit(c, "batch_generate", 5, 15*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	batchContent, err := h.contentService.BatchGenerateContent(c, services.BatchGenerateContentParams{
		Topics:      req.Topics,
		Platforms:   req.Platforms,
		Schedule:    req.Schedule,
		ContentPlan: req.ContentPlan,
		UserID:      userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في إنشاء المحتوى الجماعي", "BATCH_GENERATION_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم إنشاء المحتوى الجماعي بنجاح", batchContent)
}

// GetContent - الحصول على المحتوى المُنشأ
// @Summary الحصول على المحتوى المُنشأ
// @Description الحصول على المحتوى المُنشأ (للمشرفين فقط)
// @Tags Content
// @Security BearerAuth
// @Produce json
// @Param page query int false "الصفحة" default(1)
// @Param limit query int false "الحد" default(20)
// @Param platform query string false "المنصة"
// @Param status query string false "الحالة" default(all)
// @Param sortBy query string false "ترتيب حسب" default(createdAt)
// @Param sortOrder query string false "اتجاه الترتيب" default(desc)
// @Success 200 {object} utils.Response
// @Router /api/v1/content [get]
func (h *ContentHandler) GetContent(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !h.authService.IsAdmin(c, userID.(string)) {
		utils.ErrorResponse(c, http.StatusForbidden, "غير مصرح", "FORBIDDEN")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	platform := c.Query("platform")
	status := c.DefaultQuery("status", "all")
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	content, pagination, err := h.contentService.GetContent(c, services.GetContentParams{
		Page:      page,
		Limit:     limit,
		Platform:  platform,
		Status:    status,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		UserID:    userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب المحتوى", "FETCH_FAILED")
		return
	}

	response := map[string]interface{}{
		"content":    content,
		"pagination": pagination,
		"filters": map[string]interface{}{
			"platform":  platform,
			"status":    status,
			"sortBy":    sortBy,
			"sortOrder": sortOrder,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب المحتوى بنجاح", response)
}

// GetContentByID - الحصول على محتوى محدد
// @Summary الحصول على محتوى محدد
// @Description الحصول على محتوى محدد (للمشرفين فقط)
// @Tags Content
// @Security BearerAuth
// @Produce json
// @Param id path string true "معرف المحتوى"
// @Success 200 {object} utils.Response
// @Router /api/v1/content/{id} [get]
func (h *ContentHandler) GetContentByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !h.authService.IsAdmin(c, userID.(string)) {
		utils.ErrorResponse(c, http.StatusForbidden, "غير مصرح", "FORBIDDEN")
		return
	}

	contentID := c.Param("id")
	content, err := h.contentService.GetContentByID(c, contentID, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "المحتوى غير موجود", "CONTENT_NOT_FOUND")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب المحتوى بنجاح", content)
}

// UpdateContentRequest - طلب تحديث المحتوى
type UpdateContentRequest struct {
	Content   string   `json:"content"`
	Platform  string   `json:"platform"`
	Status    string   `json:"status"`
	Keywords  []string `json:"keywords"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// UpdateContent - تحديث محتوى محدد
// @Summary تحديث محتوى محدد
// @Description تحديث محتوى محدد (للمشرفين فقط)
// @Tags Content
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "معرف المحتوى"
// @Param input body UpdateContentRequest true "بيانات التحديث"
// @Success 200 {object} utils.Response
// @Router /api/v1/content/{id} [put]
func (h *ContentHandler) UpdateContent(c *gin.Context) {
	var req UpdateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !h.authService.IsAdmin(c, userID.(string)) {
		utils.ErrorResponse(c, http.StatusForbidden, "غير مصرح", "FORBIDDEN")
		return
	}

	contentID := c.Param("id")
	updatedContent, err := h.contentService.UpdateContent(c, contentID, services.UpdateContentParams{
		Content:  req.Content,
		Platform: req.Platform,
		Status:   req.Status,
		Keywords: req.Keywords,
		Metadata: req.Metadata,
		UserID:   userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحديث المحتوى", "UPDATE_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحديث المحتوى بنجاح", updatedContent)
}

// DeleteContent - حذف محتوى محدد
// @Summary حذف محتوى محدد
// @Description حذف محتوى محدد (للمشرفين فقط)
// @Tags Content
// @Security BearerAuth
// @Produce json
// @Param id path string true "معرف المحتوى"
// @Success 200 {object} utils.Response
// @Router /api/v1/content/{id} [delete]
func (h *ContentHandler) DeleteContent(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !h.authService.IsAdmin(c, userID.(string)) {
		utils.ErrorResponse(c, http.StatusForbidden, "غير مصرح", "FORBIDDEN")
		return
	}

	contentID := c.Param("id")
	err := h.contentService.DeleteContent(c, contentID, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في حذف المحتوى", "DELETE_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم حذف المحتوى بنجاح", nil)
}

// AnalyzeContentRequest - طلب تحليل المحتوى
type AnalyzeContentRequest struct {
	AnalysisType string `json:"analysisType"`
}

// AnalyzeContent - تحليل محتوى باستخدام الذكاء الاصطناعي
// @Summary تحليل محتوى باستخدام الذكاء الاصطناعي
// @Description تحليل محتوى باستخدام الذكاء الاصطناعي (للمشرفين فقط)
// @Tags Content
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "معرف المحتوى"
// @Param input body AnalyzeContentRequest true "بيانات التحليل"
// @Success 200 {object} utils.Response
// @Router /api/v1/content/{id}/analyze [post]
func (h *ContentHandler) AnalyzeContent(c *gin.Context) {
	var req AnalyzeContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !h.authService.IsAdmin(c, userID.(string)) {
		utils.ErrorResponse(c, http.StatusForbidden, "غير مصرح", "FORBIDDEN")
		return
	}

	if !middleware.CheckRateLimit(c, "content_analyze", 20, 5*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	contentID := c.Param("id")
	analysis, err := h.contentService.AnalyzeContent(c, contentID, req.AnalysisType, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحليل المحتوى", "ANALYSIS_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحليل المحتوى بنجاح", analysis)
}

// OptimizeContentRequest - طلب تحسين المحتوى
type OptimizeContentRequest struct {
	Content  string   `json:"content" binding:"required"`
	Platform string   `json:"platform" binding:"required"`
	Goals    []string `json:"goals"`
}

// OptimizeContent - تحسين محتوى موجود باستخدام الذكاء الاصطناعي
// @Summary تحسين محتوى موجود باستخدام الذكاء الاصطناعي
// @Description تحسين محتوى موجود باستخدام الذكاء الاصطناعي (للمشرفين فقط)
// @Tags Content
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body OptimizeContentRequest true "بيانات التحسين"
// @Success 200 {object} utils.Response
// @Router /api/v1/content/optimize [post]
func (h *ContentHandler) OptimizeContent(c *gin.Context) {
	var req OptimizeContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !h.authService.IsAdmin(c, userID.(string)) {
		utils.ErrorResponse(c, http.StatusForbidden, "غير مصرح", "FORBIDDEN")
		return
	}

	if !middleware.CheckRateLimit(c, "content_optimize", 15, 10*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	optimization, err := h.contentService.OptimizeContent(c, services.OptimizeContentParams{
		Content:  req.Content,
		Platform: req.Platform,
		Goals:    req.Goals,
		UserID:   userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحسين المحتوى", "OPTIMIZATION_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحسين المحتوى بنجاح", optimization)
}

// GetContentPerformance - الحصول على أداء المحتوى
// @Summary الحصول على أداء المحتوى
// @Description الحصول على أداء المحتوى (للمشرفين فقط)
// @Tags Content
// @Security BearerAuth
// @Produce json
// @Param id path string true "معرف المحتوى"
// @Param timeframe query string false "الفترة الزمنية" default(30d)
// @Success 200 {object} utils.Response
// @Router /api/v1/content/{id}/performance [get]
func (h *ContentHandler) GetContentPerformance(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !h.authService.IsAdmin(c, userID.(string)) {
		utils.ErrorResponse(c, http.StatusForbidden, "غير مصرح", "FORBIDDEN")
		return
	}

	contentID := c.Param("id")
	timeframe := c.DefaultQuery("timeframe", "30d")

	performance, err := h.contentService.GetContentPerformance(c, contentID, timeframe, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب أداء المحتوى", "PERFORMANCE_FETCH_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب أداء المحتوى بنجاح", performance)
}