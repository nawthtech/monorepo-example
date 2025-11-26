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

type NotificationHandler struct {
	notificationService services.NotificationService
	authService         services.AuthService
}

func NewNotificationHandler(notificationService services.NotificationService, authService services.AuthService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		authService:         authService,
	}
}

// GetNotifications - الحصول على الإشعارات مع التصنيف الذكي
// @Summary الحصول على الإشعارات مع التصنيف الذكي
// @Description الحصول على الإشعارات مع التصنيف الذكي
// @Tags Notifications
// @Security BearerAuth
// @Produce json
// @Param page query int false "الصفحة" default(1)
// @Param limit query int false "الحد" default(20)
// @Param type query string false "النوع" default(all)
// @Param status query string false "الحالة" default(unread)
// @Param priority query string false "الأولوية" default(all)
// @Success 200 {object} utils.Response
// @Router /api/v1/notifications [get]
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "notifications_get", 60, 5*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	notificationType := c.DefaultQuery("type", "all")
	status := c.DefaultQuery("status", "unread")
	priority := c.DefaultQuery("priority", "all")

	notifications, pagination, err := h.notificationService.GetNotifications(c, services.GetNotificationsParams{
		UserID:   userID.(string),
		Page:     page,
		Limit:    limit,
		Type:     notificationType,
		Status:   status,
		Priority: priority,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب الإشعارات", "FETCH_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب الإشعارات بنجاح", map[string]interface{}{
		"notifications": notifications,
		"pagination":    pagination,
		"filters": map[string]interface{}{
			"type":     notificationType,
			"status":   status,
			"priority": priority,
		},
	})
}

// GetNotificationStats - الحصول على إحصائيات الإشعارات مع تحليل الذكاء الاصطناعي
// @Summary الحصول على إحصائيات الإشعارات مع تحليل الذكاء الاصطناعي
// @Description الحصول على إحصائيات الإشعارات مع تحليل الذكاء الاصطناعي
// @Tags Notifications
// @Security BearerAuth
// @Produce json
// @Param timeframe query string false "الفترة الزمنية" default(30d)
// @Success 200 {object} utils.Response
// @Router /api/v1/notifications/stats [get]
func (h *NotificationHandler) GetNotificationStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "notifications_stats", 30, 10*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	timeframe := c.DefaultQuery("timeframe", "30d")

	stats, err := h.notificationService.GetNotificationStats(c, userID.(string), timeframe)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب إحصائيات الإشعارات", "STATS_FETCH_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب إحصائيات الإشعارات بنجاح", stats)
}

// MarkAsReadRequest - طلب تعليم إشعار كمقروء
type MarkAsReadRequest struct {
	NotificationID string `json:"notificationId"`
}

// MarkAsRead - تعليم إشعار كمقروء مع تحليل السلوك
// @Summary تعليم إشعار كمقروء مع تحليل السلوك
// @Description تعليم إشعار كمقروء مع تحليل السلوك
// @Tags Notifications
// @Security BearerAuth
// @Produce json
// @Param id path string true "معرف الإشعار"
// @Success 200 {object} utils.Response
// @Router /api/v1/notifications/{id}/read [put]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "notifications_mark_read", 100, 5*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	notificationID := c.Param("id")

	result, err := h.notificationService.MarkAsRead(c, notificationID, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "الإشعار غير موجود", "NOTIFICATION_NOT_FOUND")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تعليم الإشعار كمقروء", result)
}

// MarkAllAsReadRequest - طلب تعليم جميع الإشعارات كمقروءة
type MarkAllAsReadRequest struct {
	Type string `json:"type"`
}

// MarkAllAsRead - تعليم جميع الإشعارات كمقروءة
// @Summary تعليم جميع الإشعارات كمقروءة
// @Description تعليم جميع الإشعارات كمقروءة
// @Tags Notifications
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body MarkAllAsReadRequest false "بيانات التعليم"
// @Success 200 {object} utils.Response
// @Router /api/v1/notifications/read-all [put]
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	var req MarkAllAsReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "notifications_mark_all_read", 10, 10*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	result, err := h.notificationService.MarkAllAsRead(c, userID.(string), req.Type)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تعليم الإشعارات كمقروءة", "MARK_ALL_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تعليم جميع الإشعارات كمقروءة", result)
}

// DeleteNotification - حذف إشعار محدد
// @Summary حذف إشعار محدد
// @Description حذف إشعار محدد
// @Tags Notifications
// @Security BearerAuth
// @Produce json
// @Param id path string true "معرف الإشعار"
// @Success 200 {object} utils.Response
// @Router /api/v1/notifications/{id} [delete]
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "notifications_delete", 50, 5*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	notificationID := c.Param("id")

	err := h.notificationService.DeleteNotification(c, notificationID, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "الإشعار غير موجود", "NOTIFICATION_NOT_FOUND")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم حذف الإشعار بنجاح", nil)
}

// DeleteReadNotifications - حذف جميع الإشعارات المقروءة
// @Summary حذف جميع الإشعارات المقروءة
// @Description حذف جميع الإشعارات المقروءة
// @Tags Notifications
// @Security BearerAuth
// @Produce json
// @Param type query string false "نوع الإشعارات"
// @Success 200 {object} utils.Response
// @Router /api/v1/notifications [delete]
func (h *NotificationHandler) DeleteReadNotifications(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "notifications_delete_read", 5, 10*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	notificationType := c.Query("type")

	result, err := h.notificationService.DeleteReadNotifications(c, userID.(string), notificationType)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في حذف الإشعارات المقروءة", "DELETE_READ_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم حذف الإشعارات المقروءة بنجاح", result)
}

// CreateSmartNotificationsRequest - طلب إنشاء إشعارات ذكية
type CreateSmartNotificationsRequest struct {
	TargetUsers  []string               `json:"targetUsers"`
	Template     string                 `json:"template" binding:"required"`
	Data         map[string]interface{} `json:"data"`
	Triggers     []string               `json:"triggers"`
	Optimization bool                   `json:"optimization"`
}

// CreateSmartNotifications - إنشاء إشعارات ذكية باستخدام الذكاء الاصطناعي
// @Summary إنشاء إشعارات ذكية باستخدام الذكاء الاصطناعي
// @Description إنشاء إشعارات ذكية باستخدام الذكاء الاصطناعي (للمشرفين فقط)
// @Tags Notifications
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body CreateSmartNotificationsRequest true "بيانات الإشعارات الذكية"
// @Success 200 {object} utils.Response
// @Router /api/v1/notifications/smart [post]
func (h *NotificationHandler) CreateSmartNotifications(c *gin.Context) {
	var req CreateSmartNotificationsRequest
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

	if !middleware.CheckRateLimit(c, "notifications_smart", 20, 10*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	result, err := h.notificationService.CreateSmartNotifications(c, services.CreateSmartNotificationsParams{
		TargetUsers:  req.TargetUsers,
		Template:     req.Template,
		Data:         req.Data,
		Triggers:     req.Triggers,
		Optimization: req.Optimization,
		CreatedBy:    userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في إنشاء الإشعارات الذكية", "SMART_NOTIFICATIONS_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم إنشاء الإشعارات الذكية بنجاح", result)
}

// GetAIRecommendationsRequest - طلب الحصول على توصيات الإشعارات
type GetAIRecommendationsRequest struct {
	Category          string `json:"category"`
	MaxRecommendations int    `json:"maxRecommendations"`
}

// GetAIRecommendations - الحصول على توصيات الإشعارات باستخدام الذكاء الاصطناعي
// @Summary الحصول على توصيات الإشعارات باستخدام الذكاء الاصطناعي
// @Description الحصول على توصيات الإشعارات باستخدام الذكاء الاصطناعي
// @Tags Notifications
// @Security BearerAuth
// @Produce json
// @Param category query string false "الفئة" default(all)
// @Param maxRecommendations query int false "الحد الأقصى للتوصيات" default(5)
// @Success 200 {object} utils.Response
// @Router /api/v1/notifications/ai-recommendations [get]
func (h *NotificationHandler) GetAIRecommendations(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "notifications_recommendations", 15, 15*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	category := c.DefaultQuery("category", "all")
	maxRec, _ := strconv.Atoi(c.DefaultQuery("maxRecommendations", "5"))

	recommendations, err := h.notificationService.GetAIRecommendations(c, services.GetAIRecommendationsParams{
		UserID:            userID.(string),
		Category:          category,
		MaxRecommendations: maxRec,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في توليد التوصيات", "RECOMMENDATIONS_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم توليد توصيات الإشعارات بنجاح", recommendations)
}

// GetPreferences - الحصول على تفضيلات الإشعارات
// @Summary الحصول على تفضيلات الإشعارات
// @Description الحصول على تفضيلات الإشعارات
// @Tags Notifications
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/notifications/preferences [get]
func (h *NotificationHandler) GetPreferences(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	preferences, err := h.notificationService.GetPreferences(c, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب التفضيلات", "PREFERENCES_FETCH_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب تفضيلات الإشعارات بنجاح", preferences)
}

// UpdatePreferencesRequest - طلب تحديث تفضيلات الإشعارات
type UpdatePreferencesRequest struct {
	EmailEnabled  *bool    `json:"emailEnabled"`
	PushEnabled   *bool    `json:"pushEnabled"`
	SMSEnabled    *bool    `json:"smsEnabled"`
	AllowedTypes  []string `json:"allowedTypes"`
	QuietHours    []string `json:"quietHours"`
	Language      string   `json:"language"`
}

// UpdatePreferences - تحديث تفضيلات الإشعارات
// @Summary تحديث تفضيلات الإشعارات
// @Description تحديث تفضيلات الإشعارات
// @Tags Notifications
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body UpdatePreferencesRequest true "بيانات التحديث"
// @Success 200 {object} utils.Response
// @Router /api/v1/notifications/preferences [put]
func (h *NotificationHandler) UpdatePreferences(c *gin.Context) {
	var req UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	updatedPreferences, err := h.notificationService.UpdatePreferences(c, userID.(string), services.UpdatePreferencesParams{
		EmailEnabled: req.EmailEnabled,
		PushEnabled:  req.PushEnabled,
		SMSEnabled:   req.SMSEnabled,
		AllowedTypes: req.AllowedTypes,
		QuietHours:   req.QuietHours,
		Language:     req.Language,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحديث التفضيلات", "PREFERENCES_UPDATE_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحديث تفضيلات الإشعارات بنجاح", updatedPreferences)
}

// CreateSystemNotificationRequest - طلب إنشاء إشعار نظام
type CreateSystemNotificationRequest struct {
	Title       string                 `json:"title" binding:"required"`
	Message     string                 `json:"message" binding:"required"`
	Type        string                 `json:"type"`
	Priority    string                 `json:"priority"`
	TargetUsers string                 `json:"targetUsers"`
	ActionURL   string                 `json:"actionUrl"`
	ExpiresAt   string                 `json:"expiresAt"`
}

// CreateSystemNotification - إرسال إشعار نظام للمستخدمين
// @Summary إرسال إشعار نظام للمستخدمين
// @Description إرسال إشعار نظام للمستخدمين (للمشرفين فقط)
// @Tags Notifications
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body CreateSystemNotificationRequest true "بيانات إشعار النظام"
// @Success 200 {object} utils.Response
// @Router /api/v1/notifications/system [post]
func (h *NotificationHandler) CreateSystemNotification(c *gin.Context) {
	var req CreateSystemNotificationRequest
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

	if !middleware.CheckRateLimit(c, "notifications_system", 10, 5*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	systemNotification, err := h.notificationService.CreateSystemNotification(c, services.CreateSystemNotificationParams{
		Title:       req.Title,
		Message:     req.Message,
		Type:        req.Type,
		Priority:    req.Priority,
		TargetUsers: req.TargetUsers,
		ActionURL:   req.ActionURL,
		ExpiresAt:   req.ExpiresAt,
		CreatedBy:   userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في إنشاء إشعار النظام", "SYSTEM_NOTIFICATION_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم إرسال إشعار النظام بنجاح", systemNotification)
}