package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
)

// CacheHandler معالج التخزين المؤقت
type CacheHandler struct {
	cacheService services.CacheService
}

// NewCacheHandler إنشاء معالج تخزين مؤقت جديد
func NewCacheHandler(cacheService services.CacheService) *CacheHandler {
	return &CacheHandler{
		cacheService: cacheService,
	}
}

// SetRequest طلب تعيين قيمة في التخزين المؤقت
type SetRequest struct {
	Key   string      `json:"key" binding:"required"`
	Value interface{} `json:"value" binding:"required"`
	TTL   int64       `json:"ttl"` // بالثواني
}

// Set - تعيين قيمة في التخزين المؤقت
// @Summary تعيين قيمة في التخزين المؤقت
// @Description تخزين قيمة في التخزين المؤقت باستخدام المفتاح المحدد
// @Tags Cache
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body SetRequest true "بيانات التخزين"
// @Success 200 {object} utils.Response
// @Router /api/v1/cache [post]
func (h *CacheHandler) Set(c *gin.Context) {
	var req SetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	ttl := time.Duration(req.TTL) * time.Second
	if req.TTL == 0 {
		ttl = 0 // استخدام القيمة الافتراضية
	}

	err := h.cacheService.Set(c.Request.Context(), req.Key, req.Value, ttl)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تخزين البيانات", "CACHE_SET_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تخزين البيانات بنجاح", gin.H{
		"key": req.Key,
		"ttl": req.TTL,
	})
}

// Get - الحصول على قيمة من التخزين المؤقت
// @Summary الحصول على قيمة من التخزين المؤقت
// @Description جلب قيمة من التخزين المؤقت باستخدام المفتاح المحدد
// @Tags Cache
// @Security BearerAuth
// @Produce json
// @Param key query string true "المفتاح"
// @Success 200 {object} utils.Response
// @Router /api/v1/cache [get]
func (h *CacheHandler) Get(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "المفتاح مطلوب", "KEY_REQUIRED")
		return
	}

	value, err := h.cacheService.Get(c.Request.Context(), key)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب البيانات", "CACHE_GET_FAILED")
		return
	}

	if value == nil {
		utils.SuccessResponse(c, http.StatusOK, "المفتاح غير موجود", gin.H{
			"key":   key,
			"value": nil,
			"found": false,
		})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب البيانات بنجاح", gin.H{
		"key":   key,
		"value": value,
		"found": true,
	})
}

// Delete - حذف قيمة من التخزين المؤقت
// @Summary حذف قيمة من التخزين المؤقت
// @Description حذف قيمة من التخزين المؤقت باستخدام المفتاح المحدد
// @Tags Cache
// @Security BearerAuth
// @Produce json
// @Param key query string true "المفتاح"
// @Success 200 {object} utils.Response
// @Router /api/v1/cache [delete]
func (h *CacheHandler) Delete(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "المفتاح مطلوب", "KEY_REQUIRED")
		return
	}

	err := h.cacheService.Delete(c.Request.Context(), key)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في حذف البيانات", "CACHE_DELETE_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم حذف البيانات بنجاح", gin.H{
		"key": key,
	})
}

// Exists - التحقق من وجود مفتاح في التخزين المؤقت
// @Summary التحقق من وجود مفتاح
// @Description التحقق من وجود مفتاح في التخزين المؤقت
// @Tags Cache
// @Security BearerAuth
// @Produce json
// @Param key query string true "المفتاح"
// @Success 200 {object} utils.Response
// @Router /api/v1/cache/exists [get]
func (h *CacheHandler) Exists(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "المفتاح مطلوب", "KEY_REQUIRED")
		return
	}

	exists, err := h.cacheService.Exists(c.Request.Context(), key)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في التحقق من وجود المفتاح", "CACHE_EXISTS_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم التحقق من وجود المفتاح", gin.H{
		"key":    key,
		"exists": exists,
	})
}

// TTL - الحصول على وقت انتهاء الصلاحية المتبقي
// @Summary وقت انتهاء الصلاحية المتبقي
// @Description الحصول على الوقت المتبقي لانتهاء صلاحية المفتاح
// @Tags Cache
// @Security BearerAuth
// @Produce json
// @Param key query string true "المفتاح"
// @Success 200 {object} utils.Response
// @Router /api/v1/cache/ttl [get]
func (h *CacheHandler) TTL(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "المفتاح مطلوب", "KEY_REQUIRED")
		return
	}

	ttl, err := h.cacheService.TTL(c.Request.Context(), key)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في الحصول على وقت الصلاحية", "CACHE_TTL_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم الحصول على وقت الصلاحية", gin.H{
		"key": key,
		"ttl": ttl.Seconds(), // بالثواني
	})
}

// IncrementRequest طلب زيادة قيمة رقمية
type IncrementRequest struct {
	Key   string `json:"key" binding:"required"`
	Value int64  `json:"value" binding:"required"`
}

// Increment - زيادة قيمة رقمية
// @Summary زيادة قيمة رقمية
// @Description زيادة قيمة رقمية في التخزين المؤقت
// @Tags Cache
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body IncrementRequest true "بيانات الزيادة"
// @Success 200 {object} utils.Response
// @Router /api/v1/cache/increment [post]
func (h *CacheHandler) Increment(c *gin.Context) {
	var req IncrementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	result, err := h.cacheService.Increment(c.Request.Context(), req.Key, req.Value)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في زيادة القيمة", "CACHE_INCREMENT_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم زيادة القيمة بنجاح", gin.H{
		"key":   req.Key,
		"value": result,
	})
}

// LPushRequest طلب إضافة إلى القائمة
type LPushRequest struct {
	Key    string        `json:"key" binding:"required"`
	Values []interface{} `json:"values" binding:"required"`
}

// LPush - إضافة قيم إلى القائمة
// @Summary إضافة إلى القائمة
// @Description إضافة قيم إلى قائمة في التخزين المؤقت
// @Tags Cache
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body LPushRequest true "بيانات القائمة"
// @Success 200 {object} utils.Response
// @Router /api/v1/cache/lpush [post]
func (h *CacheHandler) LPush(c *gin.Context) {
	var req LPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	err := h.cacheService.LPush(c.Request.Context(), req.Key, req.Values...)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في إضافة البيانات إلى القائمة", "CACHE_LPUSH_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم إضافة البيانات إلى القائمة بنجاح", gin.H{
		"key":    req.Key,
		"values": req.Values,
	})
}

// LRange - جلب بيانات من القائمة
// @Summary جلب بيانات من القائمة
// @Description جلب بيانات من قائمة في التخزين المؤقت
// @Tags Cache
// @Security BearerAuth
// @Produce json
// @Param key query string true "المفتاح"
// @Param start query int false "بداية النطاق" default(0)
// @Param stop query int false "نهاية النطاق" default(-1)
// @Success 200 {object} utils.Response
// @Router /api/v1/cache/lrange [get]
func (h *CacheHandler) LRange(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "المفتاح مطلوب", "KEY_REQUIRED")
		return
	}

	start, _ := strconv.ParseInt(c.DefaultQuery("start", "0"), 10, 64)
	stop, _ := strconv.ParseInt(c.DefaultQuery("stop", "-1"), 10, 64)

	values, err := h.cacheService.LRange(c.Request.Context(), key, start, stop)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب البيانات من القائمة", "CACHE_LRANGE_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب البيانات من القائمة بنجاح", gin.H{
		"key":    key,
		"start":  start,
		"stop":   stop,
		"values": values,
		"count":  len(values),
	})
}

// HSetRequest طلب تعيين قيمة في الهاش
type HSetRequest struct {
	Key   string      `json:"key" binding:"required"`
	Field string      `json:"field" binding:"required"`
	Value interface{} `json:"value" binding:"required"`
}

// HSet - تعيين قيمة في الهاش
// @Summary تعيين قيمة في الهاش
// @Description تخزين قيمة في حقل هاش في التخزين المؤقت
// @Tags Cache
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body HSetRequest true "بيانات الهاش"
// @Success 200 {object} utils.Response
// @Router /api/v1/cache/hset [post]
func (h *CacheHandler) HSet(c *gin.Context) {
	var req HSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	err := h.cacheService.HSet(c.Request.Context(), req.Key, req.Field, req.Value)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تخزين البيانات في الهاش", "CACHE_HSET_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تخزين البيانات في الهاش بنجاح", gin.H{
		"key":   req.Key,
		"field": req.Field,
	})
}

// HGet - جلب قيمة من الهاش
// @Summary جلب قيمة من الهاش
// @Description جلب قيمة من حقل هاش في التخزين المؤقت
// @Tags Cache
// @Security BearerAuth
// @Produce json
// @Param key query string true "المفتاح"
// @Param field query string true "الحقل"
// @Success 200 {object} utils.Response
// @Router /api/v1/cache/hget [get]
func (h *CacheHandler) HGet(c *gin.Context) {
	key := c.Query("key")
	field := c.Query("field")
	if key == "" || field == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "المفتاح والحقل مطلوبان", "KEY_FIELD_REQUIRED")
		return
	}

	value, err := h.cacheService.HGet(c.Request.Context(), key, field)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب البيانات من الهاش", "CACHE_HGET_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب البيانات من الهاش بنجاح", gin.H{
		"key":   key,
		"field": field,
		"value": value,
	})
}

// HGetAll - جلب جميع بيانات الهاش
// @Summary جلب جميع بيانات الهاش
// @Description جلب جميع حقول الهاش في التخزين المؤقت
// @Tags Cache
// @Security BearerAuth
// @Produce json
// @Param key query string true "المفتاح"
// @Success 200 {object} utils.Response
// @Router /api/v1/cache/hgetall [get]
func (h *CacheHandler) HGetAll(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "المفتاح مطلوب", "KEY_REQUIRED")
		return
	}

	hash, err := h.cacheService.HGetAll(c.Request.Context(), key)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب بيانات الهاش", "CACHE_HGETALL_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب بيانات الهاش بنجاح", gin.H{
		"key":  key,
		"hash": hash,
		"size": len(hash),
	})
}

// Keys - البحث عن المفاتيح
// @Summary البحث عن المفاتيح
// @Description البحث عن المفاتيح باستخدام النمط المحدد
// @Tags Cache
// @Security BearerAuth
// @Produce json
// @Param pattern query string true "نمط البحث"
// @Success 200 {object} utils.Response
// @Router /api/v1/cache/keys [get]
func (h *CacheHandler) Keys(c *gin.Context) {
	pattern := c.Query("pattern")
	if pattern == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "نمط البحث مطلوب", "PATTERN_REQUIRED")
		return
	}

	keys, err := h.cacheService.Keys(c.Request.Context(), pattern)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في البحث عن المفاتيح", "CACHE_KEYS_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم البحث عن المفاتيح بنجاح", gin.H{
		"pattern": pattern,
		"keys":    keys,
		"count":   len(keys),
	})
}

// Flush - مسح جميع البيانات
// @Summary مسح جميع البيانات
// @Description مسح جميع البيانات من التخزين المؤقت
// @Tags Cache-Admin
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/cache/admin/flush [delete]
func (h *CacheHandler) Flush(c *gin.Context) {
	err := h.cacheService.Flush(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في مسح البيانات", "CACHE_FLUSH_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم مسح جميع البيانات بنجاح", nil)
}

// FlushPattern - مسح البيانات باستخدام النمط
// @Summary مسح البيانات باستخدام النمط
// @Description مسح البيانات من التخزين المؤقت باستخدام النمط المحدد
// @Tags Cache-Admin
// @Security BearerAuth
// @Produce json
// @Param pattern query string true "نمط المسح"
// @Success 200 {object} utils.Response
// @Router /api/v1/cache/admin/flush-pattern [delete]
func (h *CacheHandler) FlushPattern(c *gin.Context) {
	pattern := c.Query("pattern")
	if pattern == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "نمط المسح مطلوب", "PATTERN_REQUIRED")
		return
	}

	deletedCount, err := h.cacheService.FlushPattern(c.Request.Context(), pattern)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في مسح البيانات", "CACHE_FLUSH_PATTERN_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم مسح البيانات بنجاح", gin.H{
		"pattern":      pattern,
		"deletedCount": deletedCount,
	})
}

// Stats - إحصائيات التخزين المؤقت
// @Summary إحصائيات التخزين المؤقت
// @Description الحصول على إحصائيات التخزين المؤقت
// @Tags Cache-Admin
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/cache/admin/stats [get]
func (h *CacheHandler) Stats(c *gin.Context) {
	stats, err := h.cacheService.GetStats(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب الإحصائيات", "CACHE_STATS_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب الإحصائيات بنجاح", stats)
}

// Health - فحص صحة التخزين المؤقت
// @Summary فحص صحة التخزين المؤقت
// @Description فحص صحة خدمة التخزين المؤقت
// @Tags Cache-Admin
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/cache/admin/health [get]
func (h *CacheHandler) Health(c *gin.Context) {
	health, err := h.cacheService.HealthCheck(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في فحص الصحة", "CACHE_HEALTH_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم فحص الصحة بنجاح", health)
}