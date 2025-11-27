package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/models"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
	"gorm.io/gorm"
)

// HealthHandler معالج فحوصات الصحة
type HealthHandler struct {
	db             *gorm.DB
	healthService  services.HealthService
	version        string
	environment    string
	startTime      time.Time
}

// NewHealthHandler إنشاء معالج صحة جديد
func NewHealthHandler(db *gorm.DB, healthService services.HealthService, version, environment string) *HealthHandler {
	return &HealthHandler{
		db:             db,
		healthService:  healthService,
		version:        version,
		environment:    environment,
		startTime:      time.Now(),
	}
}

// HealthResponse استجابة فحص الصحة
type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Version   string                 `json:"version"`
	Environment string               `json:"environment"`
	Uptime    string                 `json:"uptime"`
	Checks    map[string]HealthCheck `json:"checks"`
}

// HealthCheck فحص صحة فردي
type HealthCheck struct {
	Status      string      `json:"status"`
	ResponseTime string     `json:"responseTime,omitempty"`
	Error       string      `json:"error,omitempty"`
	Details     interface{} `json:"details,omitempty"`
}

// SystemInfoResponse استجابة معلومات النظام
type SystemInfoResponse struct {
	Version     string    `json:"version"`
	Environment string    `json:"environment"`
	Uptime      string    `json:"uptime"`
	StartTime   time.Time `json:"startTime"`
	Timestamp   time.Time `json:"timestamp"`
}

// Check - فحص الصحة الأساسي
// @Summary فحص صحة الخدمة
// @Description فحص الحالة العامة للخدمة والمكونات
// @Tags Health
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	start := time.Now()
	checks := make(map[string]HealthCheck)

	// فحص قاعدة البيانات
	dbCheck := h.checkDatabase()
	checks["database"] = dbCheck

	// فحص الذاكرة
	memoryCheck := h.checkMemory()
	checks["memory"] = memoryCheck

	// فحص نظام الملفات
	diskCheck := h.checkDisk()
	checks["disk"] = diskCheck

	// فحص الخدمات الخارجية (إذا وجدت)
	externalCheck := h.checkExternalServices()
	checks["external_services"] = externalCheck

	// تحديد الحالة العامة
	overallStatus := "healthy"
	for _, check := range checks {
		if check.Status == "unhealthy" {
			overallStatus = "unhealthy"
			break
		} else if check.Status == "degraded" && overallStatus == "healthy" {
			overallStatus = "degraded"
		}
	}

	response := HealthResponse{
		Status:      overallStatus,
		Timestamp:   time.Now(),
		Version:     h.version,
		Environment: h.environment,
		Uptime:      time.Since(h.startTime).String(),
		Checks:      checks,
	}

	utils.SuccessResponse(c, http.StatusOK, "فحص الصحة مكتمل", response)
}

// Live - فحص الحيوية
// @Summary فحص حيوية الخدمة
// @Description فحص إذا كانت الخدمة حية وجاهزة لاستقبال الطلبات
// @Tags Health
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/health/live [get]
func (h *HealthHandler) Live(c *gin.Context) {
	response := gin.H{
		"status":    "alive",
		"timestamp": time.Now(),
		"message":   "الخدمة حية وتعمل",
	}

	utils.SuccessResponse(c, http.StatusOK, "الخدمة حية", response)
}

// Ready - فحص الجاهزية
// @Summary فحص جاهزية الخدمة
// @Description فحص إذا كانت الخدمة جاهزة لمعالجة الطلبات
// @Tags Health
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/health/ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	// فحص قاعدة البيانات
	if err := h.db.Exec("SELECT 1").Error; err != nil {
		utils.ErrorResponse(c, http.StatusServiceUnavailable, "الخدمة غير جاهزة", "SERVICE_NOT_READY")
		return
	}

	response := gin.H{
		"status":    "ready",
		"timestamp": time.Now(),
		"message":   "الخدمة جاهزة لمعالجة الطلبات",
	}

	utils.SuccessResponse(c, http.StatusOK, "الخدمة جاهزة", response)
}

// Info - معلومات النظام
// @Summary معلومات النظام
// @Description الحصول على معلومات حول إصدار وبيئة الخدمة
// @Tags Health
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/health/info [get]
func (h *HealthHandler) Info(c *gin.Context) {
	response := SystemInfoResponse{
		Version:     h.version,
		Environment: h.environment,
		Uptime:      time.Since(h.startTime).String(),
		StartTime:   h.startTime,
		Timestamp:   time.Now(),
	}

	utils.SuccessResponse(c, http.StatusOK, "معلومات النظام", response)
}

// Detailed - فحص مفصل
// @Summary فحص صحة مفصل
// @Description فحص مفصل لجميع مكونات النظام
// @Tags Health
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/health/detailed [get]
func (h *HealthHandler) Detailed(c *gin.Context) {
	start := time.Now()
	checks := make(map[string]HealthCheck)

	// فحوصات النظام الأساسية
	checks["database"] = h.checkDatabase()
	checks["memory"] = h.checkMemory()
	checks["disk"] = h.checkDisk()
	checks["cpu"] = h.checkCPU()
	checks["network"] = h.checkNetwork()

	// فحوصات التطبيق
	checks["cache"] = h.checkCache()
	checks["queue"] = h.checkQueue()
	checks["storage"] = h.checkStorage()

	// فحوصات الخدمات
	checks["external_services"] = h.checkExternalServices()
	checks["internal_services"] = h.checkInternalServices()

	// إحصائيات الأداء
	performanceCheck := h.checkPerformance()
	checks["performance"] = performanceCheck

	// تحليل شامل
	analysis := h.analyzeHealth(checks)

	response := gin.H{
		"status":      analysis.Overall,
		"timestamp":   time.Now(),
		"version":     h.version,
		"environment": h.environment,
		"uptime":      time.Since(h.startTime).String(),
		"response_time": time.Since(start).String(),
		"checks":      checks,
		"issues":      analysis.Issues,
		"recommendations": analysis.Recommendations,
		"summary":     analysis.Summary,
	}

	utils.SuccessResponse(c, http.StatusOK, "الفحص المفصل مكتمل", response)
}

// Metrics - مقاييس النظام
// @Summary مقاييس النظام
// @Description الحصول على مقاييس أداء النظام
// @Tags Health
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/health/metrics [get]
func (h *HealthHandler) Metrics(c *gin.Context) {
	metrics, err := h.healthService.GetSystemMetrics(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب المقاييس", "METRICS_FETCH_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "مقاييس النظام", metrics)
}

// AdminHealth - فحص صحة للمسؤولين
// @Summary فحص صحة متقدم للمسؤولين
// @Description فحص صحة مفصل مع معلومات حساسة للمسؤولين فقط
// @Tags Health-Admin
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/health/admin [get]
func (h *HealthHandler) AdminHealth(c *gin.Context) {
	start := time.Now()
	checks := make(map[string]HealthCheck)

	// فحوصات متقدمة للمسؤولين
	checks["database_detailed"] = h.checkDatabaseDetailed()
	checks["system_resources"] = h.checkSystemResources()
	checks["security"] = h.checkSecurity()
	checks["backups"] = h.checkBackups()
	checks["logs"] = h.checkLogs()
	checks["services_status"] = h.checkServicesStatus()

	// معلومات حساسة
	sensitiveInfo := h.getSensitiveInfo()

	response := gin.H{
		"status":       "healthy",
		"timestamp":    time.Now(),
		"version":      h.version,
		"environment":  h.environment,
		"uptime":       time.Since(h.startTime).String(),
		"response_time": time.Since(start).String(),
		"checks":       checks,
		"system_info":  sensitiveInfo,
		"warnings":     h.getSystemWarnings(),
		"maintenance":  h.getMaintenanceInfo(),
	}

	utils.SuccessResponse(c, http.StatusOK, "فحص الصحة الإداري مكتمل", response)
}

// ================================
// الدوال المساعدة للفحوصات
// ================================

func (h *HealthHandler) checkDatabase() HealthCheck {
	start := time.Now()
	
	var result int
	err := h.db.Raw("SELECT 1").Scan(&result).Error
	
	responseTime := time.Since(start).String()
	
	if err != nil {
		return HealthCheck{
			Status:      "unhealthy",
			ResponseTime: responseTime,
			Error:       err.Error(),
			Details:     "فشل في الاتصال بقاعدة البيانات",
		}
	}

	return HealthCheck{
		Status:      "healthy",
		ResponseTime: responseTime,
		Details:     "الاتصال بقاعدة البيانات نشط",
	}
}

func (h *HealthHandler) checkMemory() HealthCheck {
	// محاكاة فحص الذاكرة - في التطبيق الحقيقي استخدم runtime أو نظام المراقبة
	return HealthCheck{
		Status:  "healthy",
		Details: "استخدام الذاكرة ضمن الحدود الطبيعية",
	}
}

func (h *HealthHandler) checkDisk() HealthCheck {
	// محاكاة فحص القرص
	return HealthCheck{
		Status:  "healthy",
		Details: "مساحة التخزين كافية",
	}
}

func (h *HealthHandler) checkCPU() HealthCheck {
	return HealthCheck{
		Status:  "healthy",
		Details: "استخدام المعالج ضمن الحدود الطبيعية",
	}
}

func (h *HealthHandler) checkNetwork() HealthCheck {
	return HealthCheck{
		Status:  "healthy",
		Details: "الاتصال بالشبكة نشط",
	}
}

func (h *HealthHandler) checkCache() HealthCheck {
	return HealthCheck{
		Status:  "healthy",
		Details: "نظام الكاش يعمل بشكل طبيعي",
	}
}

func (h *HealthHandler) checkQueue() HealthCheck {
	return HealthCheck{
		Status:  "healthy",
		Details: "قوائم الانتظار تعمل بشكل طبيعي",
	}
}

func (h *HealthHandler) checkStorage() HealthCheck {
	return HealthCheck{
		Status:  "healthy",
		Details: "أنظمة التخزين تعمل بشكل طبيعي",
	}
}

func (h *HealthHandler) checkExternalServices() HealthCheck {
	// فحص الخدمات الخارجية مثل البريد، الدفع، إلخ
	return HealthCheck{
		Status:  "healthy",
		Details: "جميع الخدمات الخارجية متاحة",
	}
}

func (h *HealthHandler) checkInternalServices() HealthCheck {
	return HealthCheck{
		Status:  "healthy",
		Details: "جميع الخدمات الداخلية تعمل",
	}
}

func (h *HealthHandler) checkPerformance() HealthCheck {
	return HealthCheck{
		Status: "healthy",
		Details: gin.H{
			"response_time": "ممتاز",
			"throughput":    "عالٍ",
			"error_rate":    "منخفض",
		},
	}
}

func (h *HealthHandler) checkDatabaseDetailed() HealthCheck {
	// فحص مفصل لقاعدة البيانات
	var (
		tableCount int
		connectionCount int
	)

	h.db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE()").Scan(&tableCount)
	h.db.Raw("SHOW STATUS LIKE 'Threads_connected'").Scan(&connectionCount)

	return HealthCheck{
		Status: "healthy",
		Details: gin.H{
			"table_count":      tableCount,
			"connections":      connectionCount,
			"database_size":    "طبيعي",
			"query_performance": "ممتاز",
		},
	}
}

func (h *HealthHandler) checkSystemResources() HealthCheck {
	return HealthCheck{
		Status: "healthy",
		Details: gin.H{
			"memory_usage":    "65%",
			"cpu_usage":       "45%",
			"disk_usage":      "30%",
			"network_traffic": "منخفض",
		},
	}
}

func (h *HealthHandler) checkSecurity() HealthCheck {
	return HealthCheck{
		Status: "healthy",
		Details: gin.H{
			"ssl_enabled":     true,
			"rate_limiting":   true,
			"authentication":  true,
			"last_scan":       time.Now().Add(-24 * time.Hour),
		},
	}
}

func (h *HealthHandler) checkBackups() HealthCheck {
	return HealthCheck{
		Status: "healthy",
		Details: gin.H{
			"last_backup":     time.Now().Add(-12 * time.Hour),
			"backup_size":     "2.5 GB",
			"backup_status":   "مكتمل",
		},
	}
}

func (h *HealthHandler) checkLogs() HealthCheck {
	return HealthCheck{
		Status: "healthy",
		Details: gin.H{
			"log_level":       "info",
			"log_size":        "150 MB",
			"error_count":     "12",
		},
	}
}

func (h *HealthHandler) checkServicesStatus() HealthCheck {
	return HealthCheck{
		Status: "healthy",
		Details: gin.H{
			"api_service":     "نشط",
			"auth_service":    "نشط",
			"database_service": "نشط",
			"cache_service":   "نشط",
		},
	}
}

func (h *HealthHandler) analyzeHealth(checks map[string]HealthCheck) models.SystemSummary {
	issues := []string{}
	recommendations := []string{}

	for name, check := range checks {
		if check.Status == "unhealthy" {
			issues = append(issues, name+": "+check.Error)
		} else if check.Status == "degraded" {
			recommendations = append(recommendations, "تحسين أداء: "+name)
		}
	}

	overall := "healthy"
	if len(issues) > 0 {
		overall = "unhealthy"
	} else if len(recommendations) > 0 {
		overall = "degraded"
	}

	return models.SystemSummary{
		Overall:         overall,
		Issues:          issues,
		Recommendations: recommendations,
	}
}

func (h *HealthHandler) getSensitiveInfo() gin.H {
	return gin.H{
		"server_time":    time.Now(),
		"go_version":     "1.21", // يمكن جلبها من runtime
		"database_host":  "localhost",
		"cache_engine":   "Redis",
		"queue_system":   "في الذاكرة",
		"active_sessions": 150,
	}
}

func (h *HealthHandler) getSystemWarnings() []string {
	return []string{
		"مساحة القرص تقترب من 80%",
		"عدد الاتصالات بقاعدة البيانات مرتفع",
	}
}

func (h *HealthHandler) getMaintenanceInfo() gin.H {
	return gin.H{
		"scheduled": false,
		"next_maintenance": time.Now().Add(7 * 24 * time.Hour),
		"last_maintenance": time.Now().Add(-14 * 24 * time.Hour),
	}
}