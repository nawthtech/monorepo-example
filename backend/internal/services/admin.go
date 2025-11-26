package services

import (
	"nawthtech/backend/internal/models"
	"time"
)

type AdminService struct {
	// يمكن إضافة حقول مثل قاعدة البيانات هنا
}

func NewAdminService() *AdminService {
	return &AdminService{}
}

func (s *AdminService) GetDashboardData(timeRange string) (*models.DashboardData, error) {
	// في الواقع، ستجلب هذه البيانات من قاعدة البيانات
	// لكن حالياً سنعيد بيانات تجريبية مشابهة للبيانات في React

	stats := models.DashboardStats{
		TotalUsers:     1250,
		TotalOrders:    543,
		TotalRevenue:   125430,
		ActiveServices: 28,
		PendingOrders:  12,
		SupportTickets: 8,
		ConversionRate: 4.2,
		BounceRate:     32.1,
		StoreVisits:    3450,
		NewCustomers:   89,
	}

	storeMetrics := models.StoreMetrics{
		TotalProducts:         45,
		LowStockItems:         3,
		StoreRevenue:          89450,
		StoreOrders:           432,
		AverageOrderValue:     207,
		TopSellingCategory:    "خدمات الوسائل الاجتماعية",
		CustomerSatisfaction:  4.8,
		ReturnRate:            1.2,
	}

	recentOrders := []models.Order{
		{
			ID:       "ORD-001",
			User:     "أحمد محمد",
			Service:  "متابعين إنستغرام - 1000 متابع",
			Amount:   150,
			Status:   "completed",
			Date:     "2024-01-15",
			Type:     "store",
			Category: "وسائل اجتماعية",
		},
		{
			ID:       "ORD-002",
			User:     "سارة أحمد",
			Service:  "مشاهدات تيك توك - 5000 مشاهدة",
			Amount:   75,
			Status:   "processing",
			Date:     "2024-01-15",
			Type:     "store",
			Category: "تيك توك",
		},
		{
			ID:       "ORD-003",
			User:     "خالد عبدالله",
			Service:  "إعجابات يوتيوب - 2000 إعجاب",
			Amount:   200,
			Status:   "pending",
			Date:     "2024-01-14",
			Type:     "store",
			Category: "يوتيوب",
		},
		{
			ID:       "ORD-004",
			User:     "فاطمة علي",
			Service:  "استشارة تحليل أداء",
			Amount:   300,
			Status:   "completed",
			Date:     "2024-01-14",
			Type:     "consultation",
			Category: "استشارات",
		},
		{
			ID:       "ORD-005",
			User:     "محمد إبراهيم",
			Service:  "حزمة متابعين إنستغرام - 5000 متابع",
			Amount:   450,
			Status:   "completed",
			Date:     "2024-01-13",
			Type:     "store",
			Category: "وسائل اجتماعية",
		},
		{
			ID:       "ORD-006",
			User:     "نورة الكندري",
			Service:  "تحليل محتوى متقدم",
			Amount:   180,
			Status:   "processing",
			Date:     "2024-01-13",
			Type:     "ai_service",
			Category: "ذكاء اصطناعي",
		},
	}

	userActivity := []models.UserActivity{
		{
			User:    "أحمد محمد",
			Action:  "شراء من المتجر",
			Service: "متابعين إنستغرام",
			Time:    "منذ 5 دقائق",
			IP:      "192.168.1.100",
			Type:    "purchase",
		},
		{
			User:    "سارة أحمد",
			Action:  "إضافة خدمة إلى السلة",
			Service: "مشاهدات تيك توك",
			Time:    "منذ 15 دقيقة",
			IP:      "192.168.1.101",
			Type:    "cart",
		},
		{
			User:   "خالد عبدالله",
			Action: "تسجيل دخول",
			Time:   "منذ 30 دقيقة",
			IP:     "192.168.1.102",
			Type:   "login",
		},
		{
			User:    "فاطمة علي",
			Action:  "إكمال عملية دفع",
			Service: "استشارة تحليل أداء",
			Time:    "منذ 45 دقيقة",
			IP:      "192.168.1.103",
			Type:    "payment",
		},
		{
			User:    "محمد إبراهيم",
			Action:  "تقييم خدمة",
			Service: "حزمة متابعين إنستغرام",
			Time:    "منذ ساعة",
			IP:      "192.168.1.104",
			Type:    "review",
		},
		{
			User:    "لينا السعد",
			Action:  "طلب خدمة مخصصة",
			Service: "تحليل استراتيجي",
			Time:    "منذ ساعتين",
			IP:      "192.168.1.105",
			Type:    "custom_order",
		},
	}

	systemAlerts := []models.SystemAlert{
		{
			Type:      "warning",
			Title:     "خدمات منخفضة المخزون",
			Message:   "هناك 3 خدمات تحتاج لتجديد المخزون",
			Priority:  "medium",
			Action:    "معالجة",
			Timestamp: time.Now(),
		},
		{
			Type:      "info",
			Title:     "أداء المتجر ممتاز",
			Message:   "معدل التحويل في المتجر أعلى من المتوسط",
			Priority:  "low",
			Action:    "عرض التقرير",
			Timestamp: time.Now(),
		},
		{
			Type:      "success",
			Title:     "جميع الخدمات نشطة",
			Message:   "جميع خدمات المتجر تعمل بشكل طبيعي",
			Priority:  "low",
			Timestamp: time.Now(),
		},
	}

	return &models.DashboardData{
		Stats:        stats,
		StoreMetrics: storeMetrics,
		RecentOrders: recentOrders,
		UserActivity: userActivity,
		SystemAlerts: systemAlerts,
	}, nil
}

func (s *AdminService) GetStoreMetrics(timeRange string) (*models.StoreMetrics, error) {
	metrics := models.StoreMetrics{
		TotalProducts:         45,
		LowStockItems:         3,
		StoreRevenue:          89450,
		StoreOrders:           432,
		AverageOrderValue:     207,
		TopSellingCategory:    "خدمات الوسائل الاجتماعية",
		CustomerSatisfaction:  4.8,
		ReturnRate:            1.2,
	}
	return &metrics, nil
}

func (s *AdminService) GetRecentOrders(limit int) ([]models.Order, error) {
	// محاكاة جلب البيانات من قاعدة البيانات
	orders := []models.Order{
		{
			ID:       "ORD-001",
			User:     "أحمد محمد",
			Service:  "متابعين إنستغرام - 1000 متابع",
			Amount:   150,
			Status:   "completed",
			Date:     "2024-01-15",
			Type:     "store",
			Category: "وسائل اجتماعية",
		},
		// ... إضافة باقي الطلبات
	}
	return orders, nil
}

func (s *AdminService) GetUserActivity(limit int) ([]models.UserActivity, error) {
	// محاكاة جلب البيانات من قاعدة البيانات
	activity := []models.UserActivity{
		{
			User:    "أحمد محمد",
			Action:  "شراء من المتجر",
			Service: "متابعين إنستغرام",
			Time:    "منذ 5 دقائق",
			IP:      "192.168.1.100",
			Type:    "purchase",
		},
		// ... إضافة باقي النشاطات
	}
	return activity, nil
}

func (s *AdminService) GetSystemAlerts() ([]models.SystemAlert, error) {
	alerts := []models.SystemAlert{
		{
			Type:      "warning",
			Title:     "خدمات منخفضة المخزون",
			Message:   "هناك 3 خدمات تحتاج لتجديد المخزون",
			Priority:  "medium",
			Action:    "معالجة",
			Timestamp: time.Now(),
		},
	}
	return alerts, nil
}

func (s *AdminService) ExportReport(reportType, timeRange string) (interface{}, error) {
	// محاكاة تصدير التقرير
	return map[string]interface{}{
		"type":      reportType,
		"timeRange": timeRange,
		"data":      "بيانات التقرير...",
	}, nil
}

func (s *AdminService) UpdateOrderStatus(orderID, status string) error {
	// محاكاة تحديث حالة الطلب في قاعدة البيانات
	return nil
}
