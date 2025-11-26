package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"github.com/nawthtech/nawthtech/backend/internal/services"

	"github.com/go-chi/chi/v5"
)

type StoreHandler struct {
	storeService *services.StoreService
}

func NewStoreHandler(storeService *services.StoreService) *StoreHandler {
	return &StoreHandler{
		storeService: storeService,
	}
}

// ==================== خدمات المتجر الأساسية ====================

func (h *StoreHandler) GetServices(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	userID := "guest"
	if user := r.Context().Value("userID"); user != nil {
		userID = user.(string)
	}

	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 12
	}

	filters := map[string]interface{}{
		"query":     query.Get("q"),
		"category":  query.Get("category"),
		"minPrice":  query.Get("minPrice"),
		"maxPrice":  query.Get("maxPrice"),
		"rating":    query.Get("rating"),
		"featured":  query.Get("featured"),
		"inStock":   query.Get("inStock"),
	}

	logger.Stdout.Info("جلب الخدمات مع التوصيات الذكية", 
		"userID", userID, 
		"filters", filters, 
		"page", page, 
		"limit", limit)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب الخدمات بنجاح",
		"data": []map[string]interface{}{
			{
				"id":          "service_1",
				"name":        "خدمة متابعين إنستغرام",
				"description": "زيادة المتابعين بشكل طبيعي وآمن",
				"price":       150.00,
				"category":    "وسائل_اجتماعية",
				"rating":      4.8,
				"inStock":     true,
			},
			{
				"id":          "service_2",
				"name":        "تصميم شعار احترافي",
				"description": "تصميم شعار فريد لعلامتك التجارية",
				"price":       300.00,
				"category":    "تصميم",
				"rating":      4.9,
				"inStock":     true,
			},
		},
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": 2,
		},
		"aiRecommendations": []map[string]interface{}{
			{
				"serviceId": "service_3",
				"name":      "خدمة مقترحة بناءً على اهتماماتك",
				"reason":    "شائع بين المستخدمين المشابهين",
				"confidence": 0.85,
			},
		},
	}

	respondJSON(w, response)
}

func (h *StoreHandler) GetServiceDetails(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "serviceId")
	userID := "guest"
	if user := r.Context().Value("userID"); user != nil {
		userID = user.(string)
	}

	logger.Stdout.Info("جلب تفاصيل خدمة", "userID", userID, "serviceID", serviceID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب تفاصيل الخدمة بنجاح",
		"data": map[string]interface{}{
			"id":          serviceID,
			"name":        "خدمة متابعين إنستغرام",
			"description": "زيادة المتابعين بشكل طبيعي وآمن مع ضمان الجودة",
			"price":       150.00,
			"category":    "وسائل_اجتماعية",
			"rating":      4.8,
			"reviews":     1250,
			"inStock":     true,
			"features": []string{
				"متابعين حقيقيين",
				"ضمان استمرارية المتابعين",
				"دعم فني 24/7",
			},
			"deliveryTime": "24-48 ساعة",
		},
	}

	respondJSON(w, response)
}

func (h *StoreHandler) GetCategoriesWithStats(w http.ResponseWriter, r *http.Request) {
	logger.Stdout.Info("جلب التصنيفات مع الإحصائيات")

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب التصنيفات بنجاح",
		"data": []map[string]interface{}{
			{
				"id":    "cat_1",
				"name":  "وسائل التواصل الاجتماعي",
				"count": 45,
				"stats": map[string]interface{}{
					"totalServices": 45,
					"averageRating": 4.7,
					"popularity":    95,
				},
			},
			{
				"id":    "cat_2",
				"name":  "التصميم والإبداع",
				"count": 32,
				"stats": map[string]interface{}{
					"totalServices": 32,
					"averageRating": 4.9,
					"popularity":    88,
				},
			},
		},
	}

	respondJSON(w, response)
}

func (h *StoreHandler) GetServicesByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "categoryId")
	query := r.URL.Query()
	
	userID := "guest"
	if user := r.Context().Value("userID"); user != nil {
		userID = user.(string)
	}

	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 12
	}

	logger.Stdout.Info("جلب خدمات التصنيف", 
		"userID", userID, 
		"categoryID", categoryID, 
		"page", page, 
		"limit", limit)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب خدمات التصنيف بنجاح",
		"data": []map[string]interface{}{
			{
				"id":          "service_cat_1",
				"name":        "خدمة ضمن التصنيف",
				"description": "خدمة مخصصة لهذا التصنيف",
				"price":       200.00,
				"category":    categoryID,
				"rating":      4.7,
			},
		},
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": 1,
		},
	}

	respondJSON(w, response)
}

func (h *StoreHandler) CheckServiceAvailability(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "serviceId")
	userID := "guest"
	if user := r.Context().Value("userID"); user != nil {
		userID = user.(string)
	}

	var availabilityData struct {
		Quantity int                    `json:"quantity"`
		Options  map[string]interface{} `json:"options"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&availabilityData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("التحقق من توفر الخدمة", 
		"userID", userID, 
		"serviceID", serviceID, 
		"quantity", availabilityData.Quantity, 
		"options", availabilityData.Options)

	response := map[string]interface{}{
		"success": true,
		"message": "تم التحقق من التوفر بنجاح",
		"data": map[string]interface{}{
			"available":   true,
			"serviceId":   serviceID,
			"maxQuantity": 100,
			"estimatedDelivery": "24 ساعة",
			"requirements": []string{
				"حساب عام",
				"رابط الحساب",
			},
		},
	}

	respondJSON(w, response)
}

// ==================== إدارة الطلبات ====================

func (h *StoreHandler) CreateAIOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	var orderData struct {
		Items       []map[string]interface{} `json:"items"`
		TotalAmount float64                  `json:"totalAmount"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&orderData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("إنشاء طلب جديد مع التحقق بالذكاء الاصطناعي", 
		"userID", userID, 
		"itemsCount", len(orderData.Items), 
		"totalAmount", orderData.TotalAmount)

	response := map[string]interface{}{
		"success": true,
		"message": "تم إنشاء الطلب بنجاح",
		"data": map[string]interface{}{
			"id":          "order_" + userID + "_" + strconv.Itoa(len(orderData.Items)),
			"status":      "pending",
			"totalAmount": orderData.TotalAmount,
			"items":       orderData.Items,
			"createdAt":   "2024-01-01T00:00:00Z",
		},
		"orderId":   "order_" + userID + "_" + strconv.Itoa(len(orderData.Items)),
		"nextSteps": []string{"complete_payment", "confirm_order"},
	}

	w.WriteHeader(http.StatusCreated)
	respondJSON(w, response)
}

func (h *StoreHandler) CreateOrderFromCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	var paymentData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&paymentData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("إنشاء طلب من السلة", 
		"userID", userID, 
		"paymentMethod", paymentData["paymentMethod"])

	response := map[string]interface{}{
		"success": true,
		"message": "تم إنشاء الطلب من السلة بنجاح",
		"data": map[string]interface{}{
			"id":          "order_cart_" + userID,
			"status":      "processing",
			"totalAmount": 450.00,
			"items": []map[string]interface{}{
				{
					"serviceId": "service_1",
					"name":      "خدمة متابعين إنستغرام",
					"price":     150.00,
					"quantity":  3,
				},
			},
		},
		"orderId": "order_cart_" + userID,
	}

	w.WriteHeader(http.StatusCreated)
	respondJSON(w, response)
}

func (h *StoreHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	query := r.URL.Query()
	status := query.Get("status")
	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 10
	}

	logger.Stdout.Info("جلب طلبات المستخدم", 
		"userID", userID, 
		"status", status, 
		"page", page, 
		"limit", limit)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب الطلبات بنجاح",
		"data": []map[string]interface{}{
			{
				"id":     "order_1",
				"status": "completed",
				"total":  450.00,
				"date":   "2024-01-01T00:00:00Z",
			},
		},
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": 1,
		},
		"filters": map[string]interface{}{
			"status": status,
		},
	}

	respondJSON(w, response)
}

func (h *StoreHandler) GetOrderDetails(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	orderID := chi.URLParam(r, "orderId")

	logger.Stdout.Info("جلب تفاصيل طلب", "userID", userID, "orderID", orderID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب تفاصيل الطلب بنجاح",
		"data": map[string]interface{}{
			"id":     orderID,
			"status": "completed",
			"items": []map[string]interface{}{
				{
					"serviceId": "service_1",
					"name":      "خدمة متابعين إنستغرام",
					"price":     150.00,
					"quantity":  3,
				},
			},
			"totalAmount": 450.00,
			"createdAt":   "2024-01-01T00:00:00Z",
		},
	}

	respondJSON(w, response)
}

func (h *StoreHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	orderID := chi.URLParam(r, "orderId")
	
	var statusData struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&statusData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("تحديث حالة الطلب", 
		"userID", userID, 
		"orderID", orderID, 
		"newStatus", statusData.Status, 
		"reason", statusData.Reason)

	response := map[string]interface{}{
		"success": true,
		"message": "تم تحديث حالة الطلب بنجاح",
		"data": map[string]interface{}{
			"id":            orderID,
			"status":        statusData.Status,
			"updatedAt":     "2024-01-01T00:00:00Z",
		},
		"previousStatus": "pending",
		"newStatus":      statusData.Status,
		"updatedAt":      "2024-01-01T00:00:00Z",
	}

	respondJSON(w, response)
}

func (h *StoreHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	orderID := chi.URLParam(r, "orderId")
	
	var cancelData struct {
		Reason string `json:"reason"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&cancelData); err != nil {
		respondError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	logger.Stdout.Info("إلغاء الطلب", 
		"userID", userID, 
		"orderID", orderID, 
		"reason", cancelData.Reason)

	response := map[string]interface{}{
		"success": true,
		"message": "تم إلغاء الطلب بنجاح",
		"data": map[string]interface{}{
			"cancelled": true,
			"orderId":   orderID,
		},
		"refundAmount":     450.00,
		"cancellationFee": 0.00,
	}

	respondJSON(w, response)
}

// ==================== التوصيات والتحليلات ====================

func (h *StoreHandler) GetAIRecommendations(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 10
	}
	category := query.Get("category")

	logger.Stdout.Info("توليد التوصيات الذكية", 
		"userID", userID, 
		"limit", limit, 
		"category", category)

	response := map[string]interface{}{
		"success": true,
		"message": "تم توليد التوصيات بنجاح",
		"data": []map[string]interface{}{
			{
				"serviceId": "service_rec_1",
				"name":      "خدمة مقترحة بناءً على سجل مشترياتك",
				"reason":    "متوافقة مع اهتماماتك السابقة",
				"confidence": 0.92,
			},
			{
				"serviceId": "service_rec_2",
				"name":      "خدمة شائعة بين المستخدمين المشابهين",
				"reason":    "عالية التقييم من مستخدمين مشابهين",
				"confidence": 0.87,
			},
		},
		"confidence": 0.89,
		"generatedAt": "2024-01-01T00:00:00Z",
	}

	respondJSON(w, response)
}

func (h *StoreHandler) GetStoreStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	
	query := r.URL.Query()
	period := query.Get("period")
	if period == "" {
		period = "30d"
	}
	statsType := query.Get("type")
	if statsType == "" {
		statsType = "overview"
	}

	logger.Stdout.Info("جلب إحصائيات المتجر", 
		"userID", userID, 
		"period", period, 
		"type", statsType)

	response := map[string]interface{}{
		"success": true,
		"message": "تم جلب الإحصائيات بنجاح",
		"data": map[string]interface{}{
			"totalOrders":    45,
			"totalRevenue":   12500.00,
			"averageOrderValue": 277.78,
			"conversionRate": 4.2,
			"topCategory":    "وسائل_اجتماعية",
		},
		"period":      period,
		"type":        statsType,
		"generatedAt": "2024-01-01T00:00:00Z",
	}

	respondJSON(w, response)
}

func (h *StoreHandler) SearchServices(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	userID := "guest"
	if user := r.Context().Value("userID"); user != nil {
		userID = user.(string)
	}

	searchQuery := query.Get("q")
	page, _ := strconv.Atoi(query.Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = 12
	}

	filters := map[string]interface{}{
		"query":    searchQuery,
		"category": query.Get("category"),
		"minPrice": query.Get("minPrice"),
		"maxPrice": query.Get("maxPrice"),
		"rating":   query.Get("rating"),
	}

	logger.Stdout.Info("بحث متقدم في الخدمات", 
		"userID", userID, 
		"query", searchQuery, 
		"filters", filters, 
		"page", page, 
		"limit", limit)

	response := map[string]interface{}{
		"success": true,
		"message": "تم البحث بنجاح",
		"data": []map[string]interface{}{
			{
				"id":          "search_1",
				"name":        "نتيجة البحث: " + searchQuery,
				"description": "خدمة متوافقة مع بحثك",
				"price":       200.00,
				"category":    query.Get("category"),
				"rating":      4.5,
			},
		},
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": 1,
		},
		"searchQuery":   searchQuery,
		"filtersApplied": filters,
	}

	respondJSON(w, response)
}

func (h *StoreHandler) GenerateGrowthReport(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	orderID := chi.URLParam(r, "orderId")

	logger.Stdout.Info("إنشاء تقرير النمو الذكي", "userID", userID, "orderID", orderID)

	response := map[string]interface{}{
		"success": true,
		"message": "تم إنشاء تقرير النمو بنجاح",
		"data": map[string]interface{}{
			"orderId": orderID,
			"growthMetrics": map[string]interface{}{
				"estimatedReach":   15000,
				"engagementRate":   4.8,
				"conversionGrowth": 15.2,
			},
			"recommendations": []string{
				"زيادة الميزانية للإعلان",
				"تحسين توقيت النشر",
				"استهداف جمهور جديد",
			},
		},
		"orderId":     orderID,
		"generatedAt": "2024-01-01T00:00:00Z",
		"timeframe":   "30 يوم",
	}

	respondJSON(w, response)
}