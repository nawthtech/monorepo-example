package models

import "time"

// UserListResponse استجابة قائمة المستخدمين
type UserListResponse struct {
	Users []User `json:"users"`
	Total int64  `json:"total"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

// UserDetails تفاصيل المستخدم
type UserDetails struct {
	User      *User     `json:"user"`
	Stats     *UserStats `json:"stats"`
	LastLogin time.Time `json:"last_login"`
}

// ServicesReport تقرير الخدمات
type ServicesReport struct {
	Timeframe string                   `json:"timeframe"`
	Summary   map[string]interface{}   `json:"summary"`
	Metrics   []map[string]interface{} `json:"metrics"`
}

// FinancialReport التقرير المالي
type FinancialReport struct {
	Timeframe string               `json:"timeframe"`
	Revenue   float64              `json:"revenue"`
	Expenses  float64              `json:"expenses"`
	Profit    float64              `json:"profit"`
	Growth    string               `json:"growth"`
	Breakdown []RevenueBreakdown   `json:"breakdown"`
}

// RevenueBreakdown تفصيل الإيرادات
type RevenueBreakdown struct {
	Category   string  `json:"category"`
	Amount     float64 `json:"amount"`
	Percentage float64 `json:"percentage"`
}

// PlatformAnalytics تحليلات المنصة
type PlatformAnalytics struct {
	TotalUsers      int     `json:"total_users"`
	ActiveUsers     int     `json:"active_users"`
	TotalServices   int     `json:"total_services"`
	ActiveServices  int     `json:"active_services"`
	TotalOrders     int     `json:"total_orders"`
	CompletedOrders int     `json:"completed_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
	AverageRating   float64 `json:"average_rating"`
	GrowthRate      string  `json:"growth_rate"`
	RetentionRate   string  `json:"retention_rate"`
	ConversionRate  string  `json:"conversion_rate"`
}

// ServiceStats إحصائيات الخدمات
type ServiceStats struct {
	TotalServices     int     `json:"total_services"`
	ActiveServices    int     `json:"active_services"`
	InactiveServices  int     `json:"inactive_services"`
	SuspendedServices int     `json:"suspended_services"`
	TotalRevenue      float64 `json:"total_revenue"`
	AverageRating     float64 `json:"average_rating"`
	TotalOrders       int     `json:"total_orders"`
	PopularCategory   string  `json:"popular_category"`
}

// UserStats إحصائيات المستخدم
type UserStats struct {
	UserID          string  `json:"user_id"`
	TotalServices   int     `json:"total_services"`
	ActiveServices  int     `json:"active_services"`
	TotalOrders     int     `json:"total_orders"`
	CompletedOrders int     `json:"completed_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
	AverageRating   float64 `json:"average_rating"`
	TotalReviews    int     `json:"total_reviews"`
}