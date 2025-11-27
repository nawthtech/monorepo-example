package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/services"
	"github.com/nawthtech/nawthtech/backend/internal/utils"
)

type ServicesHandler struct {
	servicesService services.ServicesService
	authService     services.AuthService
}

func NewServicesHandler(servicesService services.ServicesService, authService services.AuthService) *ServicesHandler {
	return &ServicesHandler{
		servicesService: servicesService,
		authService:     authService,
	}
}

// GetServices - الحصول على جميع الخدمات
// @Summary الحصول على جميع الخدمات
// @Description الحصول على جميع الخدمات
// @Tags Services
// @Produce json
// @Param page query int false "الصفحة" default(1)
// @Param limit query int false "الحد" default(20)
// @Param category query string false "الفئة"
// @Param sortBy query string false "ترتيب حسب" default(createdAt)
// @Param sortOrder query string false "اتجاه الترتيب" default(desc)
// @Param minPrice query number false "أقل سعر"
// @Param maxPrice query number false "أعلى سعر"
// @Success 200 {object} utils.Response
// @Router /api/v1/services [get]
func (h *ServicesHandler) GetServices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	category := c.Query("category")
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "desc")
	minPrice, _ := strconv.ParseFloat(c.Query("minPrice"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("maxPrice"), 64)

	services, pagination, err := h.servicesService.GetServices(c, services.GetServicesParams{
		Page:      page,
		Limit:     limit,
		Category:  category,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		MinPrice:  minPrice,
		MaxPrice:  maxPrice,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب الخدمات", "SERVICES_FETCH_FAILED")
		return
	}

	response := map[string]interface{}{
		"services":   services,
		"pagination": pagination,
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب الخدمات بنجاح", response)
}

// SearchServices - البحث في الخدمات
// @Summary البحث في الخدمات
// @Description البحث في الخدمات
// @Tags Services
// @Produce json
// @Param q query string true "نص البحث"
// @Param page query int false "الصفحة" default(1)
// @Param limit query int false "الحد" default(20)
// @Param category query string false "الفئة"
// @Success 200 {object} utils.Response
// @Router /api/v1/services/search [get]
func (h *ServicesHandler) SearchServices(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "نص البحث مطلوب", "SEARCH_QUERY_REQUIRED")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	category := c.Query("category")

	results, pagination, err := h.servicesService.SearchServices(c, services.SearchServicesParams{
		Query:    query,
		Page:     page,
		Limit:    limit,
		Category: category,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في البحث في الخدمات", "SEARCH_FAILED")
		return
	}

	response := map[string]interface{}{
		"results":    results,
		"pagination": pagination,
		"query":      query,
	}

	utils.SuccessResponse(c, http.StatusOK, "تم البحث في الخدمات بنجاح", response)
}

// GetFeaturedServices - الحصول على الخدمات المميزة
// @Summary الحصول على الخدمات المميزة
// @Description الحصول على الخدمات المميزة
// @Tags Services
// @Produce json
// @Param limit query int false "الحد" default(10)
// @Success 200 {object} utils.Response
// @Router /api/v1/services/featured [get]
func (h *ServicesHandler) GetFeaturedServices(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	services, err := h.servicesService.GetFeaturedServices(c, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب الخدمات المميزة", "FEATURED_SERVICES_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب الخدمات المميزة بنجاح", services)
}

// GetServiceDetails - الحصول على تفاصيل الخدمة
// @Summary الحصول على تفاصيل الخدمة
// @Description الحصول على تفاصيل الخدمة
// @Tags Services
// @Produce json
// @Param serviceId path string true "معرف الخدمة"
// @Success 200 {object} utils.Response
// @Router /api/v1/services/{serviceId} [get]
func (h *ServicesHandler) GetServiceDetails(c *gin.Context) {
	serviceID := c.Param("serviceId")

	service, err := h.servicesService.GetServiceDetails(c, serviceID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "الخدمة غير موجودة", "SERVICE_NOT_FOUND")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب تفاصيل الخدمة بنجاح", service)
}

// GetRecommendedServices - الحصول على الخدمات الموصى بها
// @Summary الحصول على الخدمات الموصى بها
// @Description الحصول على الخدمات الموصى بها
// @Tags Services
// @Produce json
// @Param serviceId path string true "معرف الخدمة"
// @Param limit query int false "الحد" default(5)
// @Success 200 {object} utils.Response
// @Router /api/v1/services/{serviceId}/recommended [get]
func (h *ServicesHandler) GetRecommendedServices(c *gin.Context) {
	serviceID := c.Param("serviceId")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	services, err := h.servicesService.GetRecommendedServices(c, serviceID, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب الخدمات الموصى بها", "RECOMMENDED_SERVICES_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب الخدمات الموصى بها بنجاح", services)
}

// GetSellerServices - الحصول على خدمات البائع
// @Summary الحصول على خدمات البائع
// @Description الحصول على خدمات البائع
// @Tags Services
// @Produce json
// @Param sellerId path string true "معرف البائع"
// @Param page query int false "الصفحة" default(1)
// @Param limit query int false "الحد" default(20)
// @Param status query string false "الحالة"
// @Success 200 {object} utils.Response
// @Router /api/v1/services/seller/{sellerId} [get]
func (h *ServicesHandler) GetSellerServices(c *gin.Context) {
	sellerID := c.Param("sellerId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	services, pagination, err := h.servicesService.GetSellerServices(c, services.GetSellerServicesParams{
		SellerID: sellerID,
		Page:     page,
		Limit:    limit,
		Status:   status,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب خدمات البائع", "SELLER_SERVICES_FAILED")
		return
	}

	response := map[string]interface{}{
		"services":   services,
		"pagination": pagination,
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب خدمات البائع بنجاح", response)
}

// GetAllCategories جلب جميع الفئات
// @Summary جلب جميع الفئات
// @Description جلب جميع الفئات المتاحة
// @Tags Services
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/services/categories [get]
func (h *ServicesHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.servicesService.GetAllCategories(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب الفئات", "CATEGORIES_FETCH_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب الفئات بنجاح", categories)
}

// GetPopularTags جلب الوسوم الشائعة
// @Summary جلب الوسوم الشائعة
// @Description جلب الوسوم الأكثر شيوعاً
// @Tags Services
// @Produce json
// @Param limit query int false "الحد" default(10)
// @Success 200 {object} utils.Response
// @Router /api/v1/services/tags/popular [get]
func (h *ServicesHandler) GetPopularTags(c *gin.Context) {
	limit := 10
	if limitParam, err := strconv.Atoi(c.DefaultQuery("limit", "10")); err == nil {
		limit = limitParam
	}

	tags, err := h.servicesService.GetPopularTags(c.Request.Context(), limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب الوسوم", "TAGS_FETCH_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب الوسوم بنجاح", tags)
}

// GetPopularServices جلب الخدمات الشائعة
// @Summary جلب الخدمات الشائعة
// @Description جلب الخدمات الأكثر شيوعاً
// @Tags Services
// @Produce json
// @Param limit query int false "الحد" default(10)
// @Param category query string false "الفئة"
// @Success 200 {object} utils.Response
// @Router /api/v1/services/popular [get]
func (h *ServicesHandler) GetPopularServices(c *gin.Context) {
	limit := 10
	if limitParam, err := strconv.Atoi(c.DefaultQuery("limit", "10")); err == nil {
		limit = limitParam
	}
	category := c.Query("category")

	services, err := h.servicesService.GetPopularServices(c.Request.Context(), limit, category)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب الخدمات الشائعة", "POPULAR_SERVICES_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب الخدمات الشائعة بنجاح", services)
}

// GetServicesByCategory جلب الخدمات حسب الفئة
// @Summary جلب الخدمات حسب الفئة
// @Description جلب الخدمات بناءً على الفئة المحددة
// @Tags Services
// @Produce json
// @Param category path string true "الفئة"
// @Param limit query int false "الحد" default(10)
// @Success 200 {object} utils.Response
// @Router /api/v1/services/category/{category} [get]
func (h *ServicesHandler) GetServicesByCategory(c *gin.Context) {
	category := c.Param("category")
	limit := 10
	if limitParam, err := strconv.Atoi(c.DefaultQuery("limit", "10")); err == nil {
		limit = limitParam
	}

	services, err := h.servicesService.GetServicesByCategory(c.Request.Context(), category, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب الخدمات", "CATEGORY_SERVICES_FAILED")
		return
	}

	response := map[string]interface{}{
		"services": services,
		"category": category,
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب الخدمات بنجاح", response)
}

// GetServicesByTag جلب الخدمات حسب الوسم
// @Summary جلب الخدمات حسب الوسم
// @Description جلب الخدمات بناءً على الوسم المحدد
// @Tags Services
// @Produce json
// @Param tag path string true "الوسم"
// @Param page query int false "الصفحة" default(1)
// @Param limit query int false "الحد" default(10)
// @Success 200 {object} utils.Response
// @Router /api/v1/services/tag/{tag} [get]
func (h *ServicesHandler) GetServicesByTag(c *gin.Context) {
	tag := c.Param("tag")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	services, pagination, err := h.servicesService.GetServicesByTag(c.Request.Context(), tag, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب الخدمات", "TAG_SERVICES_FAILED")
		return
	}

	response := map[string]interface{}{
		"services":   services,
		"tag":        tag,
		"pagination": pagination,
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب الخدمات بنجاح", response)
}

// GetSimilarServices جلب خدمات مشابهة
// @Summary جلب خدمات مشابهة
// @Description جلب خدمات مشابهة للخدمة المحددة
// @Tags Services
// @Produce json
// @Param serviceId path string true "معرف الخدمة"
// @Param limit query int false "الحد" default(5)
// @Success 200 {object} utils.Response
// @Router /api/v1/services/{serviceId}/similar [get]
func (h *ServicesHandler) GetSimilarServices(c *gin.Context) {
	serviceID := c.Param("serviceId")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	services, err := h.servicesService.GetSimilarServices(c.Request.Context(), serviceID, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب الخدمات المشابهة", "SIMILAR_SERVICES_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب الخدمات المشابهة بنجاح", services)
}

// GetServiceRatings جلب تقييمات الخدمة
// @Summary جلب تقييمات الخدمة
// @Description جلب جميع التقييمات الخاصة بالخدمة
// @Tags Services
// @Produce json
// @Param serviceId path string true "معرف الخدمة"
// @Param page query int false "الصفحة" default(1)
// @Param limit query int false "الحد" default(10)
// @Success 200 {object} utils.Response
// @Router /api/v1/services/{serviceId}/ratings [get]
func (h *ServicesHandler) GetServiceRatings(c *gin.Context) {
	serviceID := c.Param("serviceId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	ratings, pagination, err := h.servicesService.GetServiceRatings(c.Request.Context(), serviceID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب التقييمات", "RATINGS_FETCH_FAILED")
		return
	}

	response := map[string]interface{}{
		"ratings":    ratings,
		"pagination": pagination,
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب التقييمات بنجاح", response)
}

// GetMyServices جلب خدماتي (للبائع)
// @Summary جلب خدماتي
// @Description جلب جميع خدمات البائع الحالي
// @Tags Services
// @Security BearerAuth
// @Produce json
// @Param page query int false "الصفحة" default(1)
// @Param limit query int false "الحد" default(10)
// @Param status query string false "الحالة"
// @Success 200 {object} utils.Response
// @Router /api/v1/services/my/services [get]
func (h *ServicesHandler) GetMyServices(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	var params services.GetSellerServicesParams
	params.SellerID = userID.(string)
	params.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))
	params.Status = c.Query("status")

	services, pagination, err := h.servicesService.GetSellerServices(c.Request.Context(), params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب الخدمات", "MY_SERVICES_FAILED")
		return
	}

	response := map[string]interface{}{
		"services":   services,
		"pagination": pagination,
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب الخدمات بنجاح", response)
}

// GetServicesStatusCount جلب عدد الخدمات حسب الحالة
// @Summary جلب عدد الخدمات حسب الحالة
// @Description جلب عدد الخدمات مصنفة حسب حالتها للبائع الحالي
// @Tags Services
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/services/my/stats/status [get]
func (h *ServicesHandler) GetServicesStatusCount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	counts, err := h.servicesService.CountServicesByStatus(c.Request.Context(), userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب الإحصائيات", "STATUS_COUNT_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب الإحصائيات بنجاح", counts)
}

// CheckAvailabilityRequest - طلب التحقق من توفر الخدمة
type CheckAvailabilityRequest struct {
	Date   string `json:"date" binding:"required"`
	Time   string `json:"time" binding:"required"`
	Guests int    `json:"guests"`
}

// CheckAvailability - التحقق من توفر الخدمة
// @Summary التحقق من توفر الخدمة
// @Description التحقق من توفر الخدمة
// @Tags Services
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param serviceId path string true "معرف الخدمة"
// @Param input body CheckAvailabilityRequest true "بيانات التوفر"
// @Success 200 {object} utils.Response
// @Router /api/v1/services/{serviceId}/check-availability [post]
func (h *ServicesHandler) CheckAvailability(c *gin.Context) {
	var req CheckAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	serviceID := c.Param("serviceId")

	availability, err := h.servicesService.CheckAvailability(c, services.CheckAvailabilityParams{
		ServiceID: serviceID,
		Date:      req.Date,
		Time:      req.Time,
		Guests:    req.Guests,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في التحقق من التوفر", "AVAILABILITY_CHECK_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم التحقق من التوفر بنجاح", availability)
}

// AddRatingRequest - طلب إضافة تقييم
type AddRatingRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"required"`
}

// AddRating - إضافة تقييم
// @Summary إضافة تقييم
// @Description إضافة تقييم للخدمة
// @Tags Services
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param serviceId path string true "معرف الخدمة"
// @Param input body AddRatingRequest true "بيانات التقييم"
// @Success 200 {object} utils.Response
// @Router /api/v1/services/{serviceId}/ratings [post]
func (h *ServicesHandler) AddRating(c *gin.Context) {
	var req AddRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	serviceID := c.Param("serviceId")

	rating, err := h.servicesService.AddRating(c, services.AddRatingParams{
		ServiceID: serviceID,
		UserID:    userID.(string),
		Rating:    req.Rating,
		Comment:   req.Comment,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في إضافة التقييم", "RATING_ADD_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم إضافة التقييم بنجاح", rating)
}

// CreateServiceRequest - طلب إنشاء خدمة جديدة
type CreateServiceRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Category    string   `json:"category" binding:"required"`
	Price       float64  `json:"price" binding:"required"`
	Duration    int      `json:"duration" binding:"required"`
	Images      []string `json:"images"`
	Features    []string `json:"features"`
	Tags        []string `json:"tags"`
}

// CreateService - إنشاء خدمة جديدة
// @Summary إنشاء خدمة جديدة
// @Description إنشاء خدمة جديدة (للبائعين والمسؤولين فقط)
// @Tags Services
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body CreateServiceRequest true "بيانات الخدمة"
// @Success 201 {object} utils.Response
// @Router /api/v1/services [post]
func (h *ServicesHandler) CreateService(c *gin.Context) {
	var req CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	service, err := h.servicesService.CreateService(c, services.CreateServiceParams{
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Price:       req.Price,
		Duration:    req.Duration,
		Images:      req.Images,
		Features:    req.Features,
		Tags:        req.Tags,
		SellerID:    userID.(string),
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في إنشاء الخدمة", "SERVICE_CREATION_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "تم إنشاء الخدمة بنجاح", service)
}

// UpdateServiceRequest - طلب تحديث الخدمة
type UpdateServiceRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Price       float64  `json:"price"`
	Duration    int      `json:"duration"`
	Images      []string `json:"images"`
	Features    []string `json:"features"`
	Tags        []string `json:"tags"`
}

// UpdateService - تحديث الخدمة
// @Summary تحديث الخدمة
// @Description تحديث الخدمة (للبائعين والمسؤولين فقط)
// @Tags Services
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param serviceId path string true "معرف الخدمة"
// @Param input body UpdateServiceRequest true "بيانات التحديث"
// @Success 200 {object} utils.Response
// @Router /api/v1/services/{serviceId} [put]
func (h *ServicesHandler) UpdateService(c *gin.Context) {
	var req UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	serviceID := c.Param("serviceId")

	updatedService, err := h.servicesService.UpdateService(c, services.UpdateServiceParams{
		ServiceID:   serviceID,
		SellerID:    userID.(string),
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Price:       req.Price,
		Duration:    req.Duration,
		Images:      req.Images,
		Features:    req.Features,
		Tags:        req.Tags,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحديث الخدمة", "SERVICE_UPDATE_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحديث الخدمة بنجاح", updatedService)
}

// UpdateServiceStatusRequest - طلب تحديث حالة الخدمة
type UpdateServiceStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateServiceStatus - تحديث حالة الخدمة
// @Summary تحديث حالة الخدمة
// @Description تحديث حالة الخدمة (للبائعين والمسؤولين فقط)
// @Tags Services
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param serviceId path string true "معرف الخدمة"
// @Param input body UpdateServiceStatusRequest true "بيانات الحالة"
// @Success 200 {object} utils.Response
// @Router /api/v1/services/{serviceId}/status [patch]
func (h *ServicesHandler) UpdateServiceStatus(c *gin.Context) {
	var req UpdateServiceStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	serviceID := c.Param("serviceId")

	updatedService, err := h.servicesService.UpdateServiceStatus(c, services.UpdateServiceStatusParams{
		ServiceID: serviceID,
		SellerID:  userID.(string),
		Status:    req.Status,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحديث حالة الخدمة", "STATUS_UPDATE_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحديث حالة الخدمة بنجاح", updatedService)
}

// DeleteService - حذف الخدمة
// @Summary حذف الخدمة
// @Description حذف الخدمة (للبائعين والمسؤولين فقط)
// @Tags Services
// @Security BearerAuth
// @Produce json
// @Param serviceId path string true "معرف الخدمة"
// @Success 200 {object} utils.Response
// @Router /api/v1/services/{serviceId} [delete]
func (h *ServicesHandler) DeleteService(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	serviceID := c.Param("serviceId")

	err := h.servicesService.DeleteService(c, serviceID, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في حذف الخدمة", "SERVICE_DELETE_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم حذف الخدمة بنجاح", nil)
}

// GetServicesStats - الحصول على إحصائيات الخدمات
// @Summary الحصول على إحصائيات الخدمات
// @Description الحصول على إحصائيات الخدمات (للبائعين والمسؤولين فقط)
// @Tags Services
// @Security BearerAuth
// @Produce json
// @Param timeframe query string false "الفترة الزمنية" default(30d)
// @Success 200 {object} utils.Response
// @Router /api/v1/services/my/stats [get]
func (h *ServicesHandler) GetServicesStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	timeframe := c.DefaultQuery("timeframe", "30d")

	stats, err := h.servicesService.GetServicesStats(c, userID.(string), timeframe)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب إحصائيات الخدمات", "STATS_FETCH_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب إحصائيات الخدمات بنجاح", stats)
}

// AdvancedSearchParams - معاملات البحث المتقدم
type AdvancedSearchParams struct {
	Query     string   `json:"query"`
	Category  string   `json:"category"`
	Tags      []string `json:"tags"`
	MinPrice  float64  `json:"min_price"`
	MaxPrice  float64  `json:"max_price"`
	MinRating float64  `json:"min_rating"`
	SellerID  string   `json:"seller_id"`
	Status    string   `json:"status"`
	SortBy    string   `json:"sort_by"`
	SortOrder string   `json:"sort_order"`
	Page      int      `json:"page"`
	Limit     int      `json:"limit"`
}

// AdvancedSearch بحث متقدم
// @Summary بحث متقدم في الخدمات
// @Description بحث متقدم في الخدمات مع فلاتر متعددة
// @Tags Services
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body AdvancedSearchParams true "معاملات البحث"
// @Success 200 {object} utils.Response
// @Router /api/v1/services/search/advanced [post]
func (h *ServicesHandler) AdvancedSearch(c *gin.Context) {
	var params services.AdvancedSearchParams
	
	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	// تعيين القيم الافتراضية
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	services, pagination, err := h.servicesService.AdvancedSearch(c.Request.Context(), params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في البحث", "ADVANCED_SEARCH_FAILED")
		return
	}

	response := map[string]interface{}{
		"services":   services,
		"pagination": pagination,
	}

	utils.SuccessResponse(c, http.StatusOK, "تم البحث بنجاح", response)
}

// GetAllServices - الحصول على جميع الخدمات (للمسؤولين)
// @Summary الحصول على جميع الخدمات
// @Description الحصول على جميع الخدمات في النظام (للمسؤولين فقط)
// @Tags Services-Admin
// @Security BearerAuth
// @Produce json
// @Param page query int false "الصفحة" default(1)
// @Param limit query int false "الحد" default(20)
// @Param status query string false "الحالة"
// @Success 200 {object} utils.Response
// @Router /api/v1/services/admin [get]
func (h *ServicesHandler) GetAllServices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	services, pagination, err := h.servicesService.GetAllServices(c.Request.Context(), page, limit, status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب الخدمات", "ADMIN_SERVICES_FETCH_FAILED")
		return
	}

	response := map[string]interface{}{
		"services":   services,
		"pagination": pagination,
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب الخدمات بنجاح", response)
}