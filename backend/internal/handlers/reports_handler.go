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

type ReportsHandler struct {
	reportsService services.ReportsService
	authService    services.AuthService
}

func NewReportsHandler(reportsService services.ReportsService, authService services.AuthService) *ReportsHandler {
	return &ReportsHandler{
		reportsService: reportsService,
		authService:    authService,
	}
}

// GenerateReportRequest - طلب إنشاء تقرير تلقائي
type GenerateReportRequest struct {
	Type              string   `json:"type" binding:"required"`
	Timeframe         string   `json:"timeframe" binding:"required"`
	Platforms         []string `json:"platforms"`
	Metrics           []string `json:"metrics"`
	IncludeInsights   bool     `json:"includeInsights"`
	IncludePredictions bool    `json:"includePredictions"`
	Format            string   `json:"format"`
}

// GenerateReport - إنشاء تقرير تلقائي باستخدام الذكاء الاصطناعي
// @Summary إنشاء تقرير تلقائي باستخدام الذكاء الاصطناعي
// @Description إنشاء تقرير تلقائي باستخدام الذكاء الاصطناعي (للمشرفين فقط)
// @Tags Reports
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body GenerateReportRequest true "بيانات التقرير"
// @Success 200 {object} utils.Response
// @Router /api/v1/reports/generate [post]
func (h *ReportsHandler) GenerateReport(c *gin.Context) {
	var req GenerateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "reports_generate", 10, 15*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	report, err := h.reportsService.GenerateReport(c, services.GenerateReportParams{
		Type:              req.Type,
		Timeframe:         req.Timeframe,
		Platforms:         req.Platforms,
		Metrics:           req.Metrics,
		IncludeInsights:   req.IncludeInsights,
		IncludePredictions: req.IncludePredictions,
		Format:            req.Format,
		UserID:            userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في إنشاء التقرير", "REPORT_GENERATION_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم إنشاء التقرير بنجاح", report)
}

// GenerateComparisonReportRequest - طلب إنشاء تقرير مقارنة
type GenerateComparisonReportRequest struct {
	Periods    []string `json:"periods" binding:"required"`
	Platforms  []string `json:"platforms"`
	Metrics    []string `json:"metrics"`
	FocusAreas []string `json:"focusAreas"`
}

// GenerateComparisonReport - إنشاء تقرير مقارنة باستخدام الذكاء الاصطناعي
// @Summary إنشاء تقرير مقارنة باستخدام الذكاء الاصطناعي
// @Description إنشاء تقرير مقارنة باستخدام الذكاء الاصطناعي (للمشرفين فقط)
// @Tags Reports
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body GenerateComparisonReportRequest true "بيانات المقارنة"
// @Success 200 {object} utils.Response
// @Router /api/v1/reports/compare [post]
func (h *ReportsHandler) GenerateComparisonReport(c *gin.Context) {
	var req GenerateComparisonReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if len(req.Periods) < 2 {
		utils.ErrorResponse(c, http.StatusBadRequest, "يجب توفير فترتين على الأقل للمقارنة", "INSUFFICIENT_PERIODS")
		return
	}

	if !middleware.CheckRateLimit(c, "reports_compare", 5, 30*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	comparisonReport, err := h.reportsService.GenerateComparisonReport(c, services.GenerateComparisonReportParams{
		Periods:    req.Periods,
		Platforms:  req.Platforms,
		Metrics:    req.Metrics,
		FocusAreas: req.FocusAreas,
		UserID:     userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في إنشاء تقرير المقارنة", "COMPARISON_REPORT_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم إنشاء تقرير المقارنة بنجاح", comparisonReport)
}

// GetReports - الحصول على التقارير المحفوظة
// @Summary الحصول على التقارير المحفوظة
// @Description الحصول على التقارير المحفوظة (للمشرفين فقط)
// @Tags Reports
// @Security BearerAuth
// @Produce json
// @Param page query int false "الصفحة" default(1)
// @Param limit query int false "الحد" default(20)
// @Param type query string false "نوع التقرير"
// @Param status query string false "الحالة" default(all)
// @Param sortBy query string false "ترتيب حسب" default(createdAt)
// @Param sortOrder query string false "اتجاه الترتيب" default(desc)
// @Success 200 {object} utils.Response
// @Router /api/v1/reports [get]
func (h *ReportsHandler) GetReports(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	reportType := c.Query("type")
	status := c.DefaultQuery("status", "all")
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	reports, pagination, err := h.reportsService.GetReports(c, services.GetReportsParams{
		Page:      page,
		Limit:     limit,
		Type:      reportType,
		Status:    status,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		UserID:    userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب التقارير", "REPORTS_FETCH_FAILED")
		return
	}

	response := map[string]interface{}{
		"reports":    reports,
		"pagination": pagination,
		"filters": map[string]interface{}{
			"type":      reportType,
			"status":    status,
			"sortBy":    sortBy,
			"sortOrder": sortOrder,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب التقارير بنجاح", response)
}

// GetReportByID - الحصول على تقرير محدد
// @Summary الحصول على تقرير محدد
// @Description الحصول على تقرير محدد (للمشرفين فقط)
// @Tags Reports
// @Security BearerAuth
// @Produce json
// @Param id path string true "معرف التقرير"
// @Success 200 {object} utils.Response
// @Router /api/v1/reports/{id} [get]
func (h *ReportsHandler) GetReportByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	reportID := c.Param("id")

	report, err := h.reportsService.GetReportByID(c, reportID, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "التقرير غير موجود", "REPORT_NOT_FOUND")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب التقرير بنجاح", report)
}

// UpdateReportRequest - طلب تحديث تقرير
type UpdateReportRequest struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// UpdateReport - تحديث تقرير محدد
// @Summary تحديث تقرير محدد
// @Description تحديث تقرير محدد (للمشرفين فقط)
// @Tags Reports
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "معرف التقرير"
// @Param input body UpdateReportRequest true "بيانات التحديث"
// @Success 200 {object} utils.Response
// @Router /api/v1/reports/{id} [put]
func (h *ReportsHandler) UpdateReport(c *gin.Context) {
	var req UpdateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	reportID := c.Param("id")

	updatedReport, err := h.reportsService.UpdateReport(c, services.UpdateReportParams{
		ReportID:    reportID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Metadata:    req.Metadata,
		UserID:      userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحديث التقرير", "REPORT_UPDATE_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحديث التقرير بنجاح", updatedReport)
}

// DeleteReport - حذف تقرير محدد
// @Summary حذف تقرير محدد
// @Description حذف تقرير محدد (للمشرفين فقط)
// @Tags Reports
// @Security BearerAuth
// @Produce json
// @Param id path string true "معرف التقرير"
// @Success 200 {object} utils.Response
// @Router /api/v1/reports/{id} [delete]
func (h *ReportsHandler) DeleteReport(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	reportID := c.Param("id")

	err := h.reportsService.DeleteReport(c, reportID, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في حذف التقرير", "REPORT_DELETE_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم حذف التقرير بنجاح", nil)
}

// AnalyzeReportRequest - طلب تحليل تقرير
type AnalyzeReportRequest struct {
	AnalysisType string `json:"analysisType"`
}

// AnalyzeReport - تحليل تقرير باستخدام الذكاء الاصطناعي
// @Summary تحليل تقرير باستخدام الذكاء الاصطناعي
// @Description تحليل تقرير باستخدام الذكاء الاصطناعي (للمشرفين فقط)
// @Tags Reports
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "معرف التقرير"
// @Param input body AnalyzeReportRequest true "بيانات التحليل"
// @Success 200 {object} utils.Response
// @Router /api/v1/reports/{id}/analyze [post]
func (h *ReportsHandler) AnalyzeReport(c *gin.Context) {
	var req AnalyzeReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "reports_analyze", 15, 10*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	reportID := c.Param("id")

	analysis, err := h.reportsService.AnalyzeReport(c, reportID, req.AnalysisType, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحليل التقرير", "REPORT_ANALYSIS_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحليل التقرير بنجاح", analysis)
}

// ExportReport - تصدير تقرير بتنسيقات مختلفة
// @Summary تصدير تقرير بتنسيقات مختلفة
// @Description تصدير تقرير بتنسيقات مختلفة (للمشرفين فقط)
// @Tags Reports
// @Security BearerAuth
// @Produce json
// @Param id path string true "معرف التقرير"
// @Param format query string false "التنسيق" default(pdf)
// @Param includeCharts query bool false "تضمين الرسوم البيانية" default(true)
// @Success 200 {object} utils.Response
// @Router /api/v1/reports/{id}/export [get]
func (h *ReportsHandler) ExportReport(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	reportID := c.Param("id")
	format := c.DefaultQuery("format", "pdf")
	includeCharts := c.DefaultQuery("includeCharts", "true") == "true"

	exportResult, err := h.reportsService.ExportReport(c, services.ExportReportParams{
		ReportID:     reportID,
		Format:       format,
		IncludeCharts: includeCharts,
		UserID:       userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تصدير التقرير", "REPORT_EXPORT_FAILED")
		return
	}

	// إعداد الاستجابة للتحميل
	c.Header("Content-Type", exportResult.ContentType)
	c.Header("Content-Disposition", "attachment; filename="+exportResult.FileName)
	c.Header("Content-Length", strconv.Itoa(len(exportResult.Data)))

	c.Data(http.StatusOK, exportResult.ContentType, exportResult.Data)
}

// GetDashboardPerformance - الحصول على تقرير أداء اللوحة الرئيسية
// @Summary الحصول على تقرير أداء اللوحة الرئيسية
// @Description الحصول على تقرير أداء اللوحة الرئيسية (للمشرفين فقط)
// @Tags Reports
// @Security BearerAuth
// @Produce json
// @Param timeframe query string false "الفترة الزمنية" default(7d)
// @Param platforms query string false "المنصات" default(all)
// @Success 200 {object} utils.Response
// @Router /api/v1/reports/performance/dashboard [get]
func (h *ReportsHandler) GetDashboardPerformance(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "reports_dashboard", 30, 5*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	timeframe := c.DefaultQuery("timeframe", "7d")
	platforms := c.Query("platforms")

	dashboardReport, err := h.reportsService.GetDashboardPerformance(c, services.GetDashboardPerformanceParams{
		Timeframe: timeframe,
		Platforms: platforms,
		UserID:    userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب تقرير الأداء", "DASHBOARD_REPORT_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب تقرير الأداء بنجاح", dashboardReport)
}

// GenerateCustomPerformanceReportRequest - طلب إنشاء تقرير أداء مخصص
type GenerateCustomPerformanceReportRequest struct {
	Name          string                 `json:"name" binding:"required"`
	Metrics       []string               `json:"metrics" binding:"required"`
	Dimensions    []string               `json:"dimensions"`
	Filters       map[string]interface{} `json:"filters"`
	Timeframe     string                 `json:"timeframe" binding:"required"`
	Platforms     []string               `json:"platforms"`
	Visualization string                 `json:"visualization"`
}

// GenerateCustomPerformanceReport - إنشاء تقرير أداء مخصص
// @Summary إنشاء تقرير أداء مخصص
// @Description إنشاء تقرير أداء مخصص (للمشرفين فقط)
// @Tags Reports
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body GenerateCustomPerformanceReportRequest true "بيانات التقرير المخصص"
// @Success 200 {object} utils.Response
// @Router /api/v1/reports/performance/custom [post]
func (h *ReportsHandler) GenerateCustomPerformanceReport(c *gin.Context) {
	var req GenerateCustomPerformanceReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "reports_custom", 10, 10*time.Minute) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	customReport, err := h.reportsService.GenerateCustomPerformanceReport(c, services.GenerateCustomPerformanceReportParams{
		Name:          req.Name,
		Metrics:       req.Metrics,
		Dimensions:    req.Dimensions,
		Filters:       req.Filters,
		Timeframe:     req.Timeframe,
		Platforms:     req.Platforms,
		Visualization: req.Visualization,
		UserID:        userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في إنشاء التقرير المخصص", "CUSTOM_REPORT_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم إنشاء التقرير المخصص بنجاح", customReport)
}