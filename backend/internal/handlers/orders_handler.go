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

type OrdersHandler struct {
	ordersService services.OrdersService
	authService   services.AuthService
}

func NewOrdersHandler(ordersService services.OrdersService, authService services.AuthService) *OrdersHandler {
	return &OrdersHandler{
		ordersService: ordersService,
		authService:   authService,
	}
}

// GetOrders - الحصول على جميع الطلبات (للمسؤولين) أو طلبات المستخدم الحالي
// @Summary الحصول على جميع الطلبات (للمسؤولين) أو طلبات المستخدم الحالي
// @Description الحصول على جميع الطلبات (للمسؤولين) أو طلبات المستخدم الحالي
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Param page query int false "الصفحة" default(1)
// @Param limit query int false "الحد" default(10)
// @Param status query string false "حالة الطلب"
// @Param sortBy query string false "ترتيب حسب" default(createdAt)
// @Param sortOrder query string false "اتجاه الترتيب" default(desc)
// @Param startDate query string false "تاريخ البدء"
// @Param endDate query string false "تاريخ الانتهاء"
// @Success 200 {object} utils.Response
// @Router /api/v1/orders [get]
func (h *OrdersHandler) GetOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	userRole, _ := c.Get("userRole")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "desc")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	params := services.GetOrdersParams{
		UserID:    userID.(string),
		UserRole:  userRole.(string),
		Page:      page,
		Limit:     limit,
		Status:    status,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		StartDate: startDate,
		EndDate:   endDate,
	}

	orders, pagination, err := h.ordersService.GetOrders(c, params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب الطلبات", "FETCH_ORDERS_FAILED")
		return
	}

	response := map[string]interface{}{
		"orders":     orders,
		"pagination": pagination,
		"filters": map[string]interface{}{
			"status":    status,
			"sortBy":    sortBy,
			"sortOrder": sortOrder,
			"startDate": startDate,
			"endDate":   endDate,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب الطلبات بنجاح", response)
}

// CreateOrderRequest - طلب إنشاء طلب جديد
type CreateOrderRequest struct {
	Items       []models.OrderItem         `json:"items" binding:"required"`
	TotalAmount float64                    `json:"totalAmount" binding:"required"`
	ShippingAddress models.ShippingAddress `json:"shippingAddress" binding:"required"`
	PaymentMethod  string                  `json:"paymentMethod" binding:"required"`
	Notes        string                    `json:"notes"`
}

// CreateOrder - إنشاء طلب جديد
// @Summary إنشاء طلب جديد
// @Description إنشاء طلب جديد
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body CreateOrderRequest true "بيانات الطلب"
// @Success 201 {object} utils.Response
// @Router /api/v1/orders [post]
func (h *OrdersHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	if !middleware.CheckRateLimit(c, "orders_create", 20, time.Hour) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "تم تجاوز الحد المسموح", "RATE_LIMIT_EXCEEDED")
		return
	}

	order, err := h.ordersService.CreateOrder(c, services.CreateOrderParams{
		UserID:          userID.(string),
		Items:           req.Items,
		TotalAmount:     req.TotalAmount,
		ShippingAddress: req.ShippingAddress,
		PaymentMethod:   req.PaymentMethod,
		Notes:           req.Notes,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في إنشاء الطلب", "CREATE_ORDER_FAILED")
		return
	}

	response := map[string]interface{}{
		"order":     order,
		"orderId":   order.ID,
		"nextSteps": []string{"انتظار التأكيد", "إتمام الدفع"},
	}

	utils.SuccessResponse(c, http.StatusCreated, "تم إنشاء الطلب بنجاح", response)
}

// GetOrderByID - الحصول على تفاصيل طلب محدد
// @Summary الحصول على تفاصيل طلب محدد
// @Description الحصول على تفاصيل طلب محدد
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Param orderId path string true "معرف الطلب"
// @Success 200 {object} utils.Response
// @Router /api/v1/orders/{orderId} [get]
func (h *OrdersHandler) GetOrderByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	userRole, _ := c.Get("userRole")
	orderID := c.Param("orderId")

	order, err := h.ordersService.GetOrderByID(c, orderID, userID.(string), userRole.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "الطلب غير موجود", "ORDER_NOT_FOUND")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب تفاصيل الطلب بنجاح", order)
}

// UpdateOrderStatusRequest - طلب تحديث حالة الطلب
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
	Reason string `json:"reason"`
}

// UpdateOrderStatus - تحديث حالة الطلب
// @Summary تحديث حالة الطلب
// @Description تحديث حالة الطلب
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param orderId path string true "معرف الطلب"
// @Param input body UpdateOrderStatusRequest true "بيانات تحديث الحالة"
// @Success 200 {object} utils.Response
// @Router /api/v1/orders/{orderId}/status [patch]
func (h *OrdersHandler) UpdateOrderStatus(c *gin.Context) {
	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	orderID := c.Param("orderId")

	order, err := h.ordersService.UpdateOrderStatus(c, services.UpdateOrderStatusParams{
		OrderID: orderID,
		UserID:  userID.(string),
		Status:  req.Status,
		Reason:  req.Reason,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحديث حالة الطلب", "UPDATE_STATUS_FAILED")
		return
	}

	response := map[string]interface{}{
		"order":         order,
		"previousStatus": order.PreviousStatus,
		"newStatus":     order.Status,
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحديث حالة الطلب بنجاح", response)
}

// CancelOrderRequest - طلب إلغاء الطلب
type CancelOrderRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// CancelOrder - إلغاء الطلب
// @Summary إلغاء الطلب
// @Description إلغاء الطلب
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param orderId path string true "معرف الطلب"
// @Param input body CancelOrderRequest true "بيانات الإلغاء"
// @Success 200 {object} utils.Response
// @Router /api/v1/orders/{orderId}/cancel [post]
func (h *OrdersHandler) CancelOrder(c *gin.Context) {
	var req CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	orderID := c.Param("orderId")

	result, err := h.ordersService.CancelOrder(c, services.CancelOrderParams{
		OrderID: orderID,
		UserID:  userID.(string),
		Reason:  req.Reason,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في إلغاء الطلب", "CANCEL_ORDER_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم إلغاء الطلب بنجاح", result)
}

// RequestRefundRequest - طلب استرداد أموال
type RequestRefundRequest struct {
	Reason string  `json:"reason" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

// RequestRefund - طلب استرداد أموال
// @Summary طلب استرداد أموال
// @Description طلب استرداد أموال
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param orderId path string true "معرف الطلب"
// @Param input body RequestRefundRequest true "بيانات الاسترداد"
// @Success 200 {object} utils.Response
// @Router /api/v1/orders/{orderId}/refund [post]
func (h *OrdersHandler) RequestRefund(c *gin.Context) {
	var req RequestRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	orderID := c.Param("orderId")

	result, err := h.ordersService.RequestRefund(c, services.RequestRefundParams{
		OrderID: orderID,
		UserID:  userID.(string),
		Reason:  req.Reason,
		Amount:  req.Amount,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تقديم طلب الاسترداد", "REFUND_REQUEST_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تقديم طلب الاسترداد بنجاح", result)
}

// GetOrderTracking - الحصول على معلومات تتبع الشحن
// @Summary الحصول على معلومات تتبع الشحن
// @Description الحصول على معلومات تتبع الشحن
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Param orderId path string true "معرف الطلب"
// @Success 200 {object} utils.Response
// @Router /api/v1/orders/{orderId}/tracking [get]
func (h *OrdersHandler) GetOrderTracking(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	orderID := c.Param("orderId")

	tracking, err := h.ordersService.GetOrderTracking(c, orderID, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "معلومات التتبع غير متوفرة", "TRACKING_NOT_FOUND")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب معلومات التتبع بنجاح", tracking)
}

// UpdateShippingInfoRequest - طلب تحديث معلومات الشحن
type UpdateShippingInfoRequest struct {
	TrackingNumber   string `json:"trackingNumber" binding:"required"`
	Carrier          string `json:"carrier" binding:"required"`
	EstimatedDelivery string `json:"estimatedDelivery"`
}

// UpdateShippingInfo - تحديث معلومات الشحن (للبائعين والمسؤولين)
// @Summary تحديث معلومات الشحن (للبائعين والمسؤولين)
// @Description تحديث معلومات الشحن (للبائعين والمسؤولين)
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param orderId path string true "معرف الطلب"
// @Param input body UpdateShippingInfoRequest true "بيانات الشحن"
// @Success 200 {object} utils.Response
// @Router /api/v1/orders/{orderId}/ship [post]
func (h *OrdersHandler) UpdateShippingInfo(c *gin.Context) {
	var req UpdateShippingInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	orderID := c.Param("orderId")

	result, err := h.ordersService.UpdateShippingInfo(c, services.UpdateShippingInfoParams{
		OrderID:           orderID,
		UserID:            userID.(string),
		TrackingNumber:    req.TrackingNumber,
		Carrier:           req.Carrier,
		EstimatedDelivery: req.EstimatedDelivery,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في تحديث معلومات الشحن", "UPDATE_SHIPPING_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم تحديث معلومات الشحن بنجاح", result)
}

// GenerateInvoice - الحصول على فاتورة الطلب
// @Summary الحصول على فاتورة الطلب
// @Description الحصول على فاتورة الطلب
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Param orderId path string true "معرف الطلب"
// @Param format query string false "صيغة الفاتورة" default(pdf)
// @Success 200 {object} utils.Response
// @Router /api/v1/orders/{orderId}/invoice [get]
func (h *OrdersHandler) GenerateInvoice(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	orderID := c.Param("orderId")
	format := c.DefaultQuery("format", "pdf")

	invoice, err := h.ordersService.GenerateInvoice(c, orderID, userID.(string), format)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في إنشاء الفاتورة", "INVOICE_GENERATION_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم إنشاء الفاتورة بنجاح", invoice)
}

// GetOrderReceipt - الحصول على إيصال الطلب
// @Summary الحصول على إيصال الطلب
// @Description الحصول على إيصال الطلب
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Param orderId path string true "معرف الطلب"
// @Success 200 {object} utils.Response
// @Router /api/v1/orders/{orderId}/receipt [get]
func (h *OrdersHandler) GetOrderReceipt(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	orderID := c.Param("orderId")

	receipt, err := h.ordersService.GetOrderReceipt(c, orderID, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "الإيصال غير متوفر", "RECEIPT_NOT_FOUND")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب الإيصال بنجاح", receipt)
}

// SearchOrders - البحث في الطلبات
// @Summary البحث في الطلبات
// @Description البحث في الطلبات
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Param q query string false "نص البحث"
// @Param status query string false "حالة الطلب"
// @Param startDate query string false "تاريخ البدء"
// @Param endDate query string false "تاريخ الانتهاء"
// @Param page query int false "الصفحة" default(1)
// @Param limit query int false "الحد" default(10)
// @Success 200 {object} utils.Response
// @Router /api/v1/orders/search [get]
func (h *OrdersHandler) SearchOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	userRole, _ := c.Get("userRole")
	query := c.Query("q")
	status := c.Query("status")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.ordersService.SearchOrders(c, services.SearchOrdersParams{
		UserID:    userID.(string),
		UserRole:  userRole.(string),
		Query:     query,
		Status:    status,
		StartDate: startDate,
		EndDate:   endDate,
		Page:      page,
		Limit:     limit,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في البحث في الطلبات", "SEARCH_ORDERS_FAILED")
		return
	}

	response := map[string]interface{}{
		"orders":     result.Orders,
		"pagination": result.Pagination,
		"searchQuery": query,
	}

	utils.SuccessResponse(c, http.StatusOK, "تم البحث في الطلبات بنجاح", response)
}

// GetOrdersStats - الحصول على إحصائيات الطلبات الشاملة
// @Summary الحصول على إحصائيات الطلبات الشاملة
// @Description الحصول على إحصائيات الطلبات الشاملة (للمسؤولين والبائعين)
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Param period query string false "الفترة" default(30d)
// @Param type query string false "نوع التقرير" default(overview)
// @Success 200 {object} utils.Response
// @Router /api/v1/orders/stats/overview [get]
func (h *OrdersHandler) GetOrdersStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	userRole, _ := c.Get("userRole")
	period := c.DefaultQuery("period", "30d")
	reportType := c.DefaultQuery("type", "overview")

	stats, err := h.ordersService.GetOrdersStats(c, services.GetOrdersStatsParams{
		UserID: userID.(string),
		UserRole: userRole.(string),
		Period: period,
		Type:   reportType,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب إحصائيات الطلبات", "STATS_FETCH_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب إحصائيات الطلبات بنجاح", stats)
}

// GetRevenueStats - الحصول على إحصائيات الإيرادات
// @Summary الحصول على إحصائيات الإيرادات
// @Description الحصول على إحصائيات الإيرادات (للمسؤولين والبائعين)
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Param period query string false "الفترة" default(30d)
// @Param groupBy query string false "التجميع" default(day)
// @Success 200 {object} utils.Response
// @Router /api/v1/orders/stats/revenue [get]
func (h *OrdersHandler) GetRevenueStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	userRole, _ := c.Get("userRole")
	period := c.DefaultQuery("period", "30d")
	groupBy := c.DefaultQuery("groupBy", "day")

	revenue, err := h.ordersService.GetRevenueStats(c, services.GetRevenueStatsParams{
		UserID:  userID.(string),
		UserRole: userRole.(string),
		Period:  period,
		GroupBy: groupBy,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب إحصائيات الإيرادات", "REVENUE_STATS_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب إحصائيات الإيرادات بنجاح", revenue)
}

// GetPendingOrders - الحصول على الطلبات قيد الانتظار (للمسؤولين)
// @Summary الحصول على الطلبات قيد الانتظار (للمسؤولين)
// @Description الحصول على الطلبات قيد الانتظار (للمسؤولين)
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Param page query int false "الصفحة" default(1)
// @Param limit query int false "الحد" default(20)
// @Success 200 {object} utils.Response
// @Router /api/v1/orders/admin/pending [get]
func (h *OrdersHandler) GetPendingOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.ordersService.GetPendingOrders(c, services.GetPendingOrdersParams{
		AdminID: userID.(string),
		Page:    page,
		Limit:   limit,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في جلب الطلبات قيد الانتظار", "PENDING_ORDERS_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم جلب الطلبات قيد الانتظار بنجاح", result)
}

// ApproveOrderRequest - طلب الموافقة على طلب
type ApproveOrderRequest struct {
	Notes string `json:"notes"`
}

// ApproveOrder - الموافقة على طلب (للمسؤولين)
// @Summary الموافقة على طلب (للمسؤولين)
// @Description الموافقة على طلب (للمسؤولين)
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param orderId path string true "معرف الطلب"
// @Param input body ApproveOrderRequest true "ملاحظات الموافقة"
// @Success 200 {object} utils.Response
// @Router /api/v1/orders/admin/{orderId}/approve [post]
func (h *OrdersHandler) ApproveOrder(c *gin.Context) {
	var req ApproveOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	orderID := c.Param("orderId")

	result, err := h.ordersService.ApproveOrder(c, services.ApproveOrderParams{
		OrderID: orderID,
		AdminID: userID.(string),
		Notes:   req.Notes,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في الموافقة على الطلب", "APPROVE_ORDER_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تمت الموافقة على الطلب بنجاح", result)
}

// RejectOrderRequest - طلب رفض طلب
type RejectOrderRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// RejectOrder - رفض طلب (للمسؤولين)
// @Summary رفض طلب (للمسؤولين)
// @Description رفض طلب (للمسؤولين)
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param orderId path string true "معرف الطلب"
// @Param input body RejectOrderRequest true "سبب الرفض"
// @Success 200 {object} utils.Response
// @Router /api/v1/orders/admin/{orderId}/reject [post]
func (h *OrdersHandler) RejectOrder(c *gin.Context) {
	var req RejectOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "بيانات غير صالحة", "INVALID_INPUT")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "غير مصرح", "UNAUTHORIZED")
		return
	}

	orderID := c.Param("orderId")

	result, err := h.ordersService.RejectOrder(c, services.RejectOrderParams{
		OrderID: orderID,
		AdminID: userID.(string),
		Reason:  req.Reason,
	})

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "فشل في رفض الطلب", "REJECT_ORDER_FAILED")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "تم رفض الطلب بنجاح", result)
}