package handlers

import (
	"github.com/gin-gonic/gin"
)

func RegisterOrdersRoutes(router *gin.RouterGroup, ordersHandler *OrdersHandler, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc, sellerMiddleware gin.HandlerFunc) {
	ordersRoutes := router.Group("/orders")
	ordersRoutes.Use(authMiddleware)
	
	// ==================== الطلبات الأساسية ====================
	ordersRoutes.GET("", ordersHandler.GetOrders)
	ordersRoutes.POST("", ordersHandler.CreateOrder)
	ordersRoutes.GET("/:orderId", ordersHandler.GetOrderByID)
	
	// ==================== إدارة حالة الطلب ====================
	ordersRoutes.PATCH("/:orderId/status", ordersHandler.UpdateOrderStatus)
	ordersRoutes.POST("/:orderId/cancel", ordersHandler.CancelOrder)
	ordersRoutes.POST("/:orderId/refund", ordersHandler.RequestRefund)
	
	// ==================== إدارة الشحن والتسليم ====================
	ordersRoutes.GET("/:orderId/tracking", ordersHandler.GetOrderTracking)
	
	// ==================== الفواتير والمستندات ====================
	ordersRoutes.GET("/:orderId/invoice", ordersHandler.GenerateInvoice)
	ordersRoutes.GET("/:orderId/receipt", ordersHandler.GetOrderReceipt)
	
	// ==================== البحث والتصفية ====================
	ordersRoutes.GET("/search", ordersHandler.SearchOrders)
	
	// ==================== إحصائيات وتقارير ====================
	ordersRoutes.GET("/stats/overview", ordersHandler.GetOrdersStats)
	ordersRoutes.GET("/stats/revenue", ordersHandler.GetRevenueStats)
	
	// ==================== إدارة الطلبات للمسؤولين والبائعين ====================
	adminOrdersRoutes := ordersRoutes.Group("")
	adminOrdersRoutes.Use(sellerMiddleware)
	adminOrdersRoutes.POST("/:orderId/ship", ordersHandler.UpdateShippingInfo)
	
	// ==================== إدارة الطلبات للمسؤولين فقط ====================
	adminOnlyRoutes := ordersRoutes.Group("/admin")
	adminOnlyRoutes.Use(adminMiddleware)
	adminOnlyRoutes.GET("/pending", ordersHandler.GetPendingOrders)
	adminOnlyRoutes.POST("/:orderId/approve", ordersHandler.ApproveOrder)
	adminOnlyRoutes.POST("/:orderId/reject", ordersHandler.RejectOrder)
}