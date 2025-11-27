package models

import (
	"time"
)

// Service نموذج الخدمة
type Service struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Title        string    `json:"title" gorm:"not null"`
	Description  string    `json:"description" gorm:"type:text"`
	Category     string    `json:"category" gorm:"not null"`
	Price        float64   `json:"price" gorm:"not null"`
	Duration     int       `json:"duration"` // المدة بالأيام
	Rating       float64   `json:"rating" gorm:"default:0"`
	TotalOrders  int       `json:"total_orders" gorm:"default:0"`
	TotalReviews int       `json:"total_reviews" gorm:"default:0"`
	Status       string    `json:"status" gorm:"default:'active'"` // active, inactive, suspended
	Featured     bool      `json:"featured" gorm:"default:false"`
	SellerID     string    `json:"seller_id" gorm:"not null;index"`
	Images       []string  `json:"images" gorm:"type:json"`
	Features     []string  `json:"features" gorm:"type:json"`
	Tags         []string  `json:"tags" gorm:"type:json"`
	Reviews      []Review  `json:"reviews,omitempty" gorm:"foreignKey:ServiceID"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Review نموذج التقييم
type Review struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	ServiceID string    `json:"service_id" gorm:"not null;index"`
	UserID    string    `json:"user_id" gorm:"not null"`
	UserName  string    `json:"user_name"`
	Rating    int       `json:"rating" gorm:"check:rating>=1 AND rating<=5"`
	Comment   string    `json:"comment" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
}

// Rating نموذج التقييم (مبسط)
type Rating struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	ServiceID string    `json:"service_id" gorm:"not null;index"`
	UserID    string    `json:"user_id" gorm:"not null"`
	Rating    int       `json:"rating" gorm:"check:rating>=1 AND rating<=5"`
	Comment   string    `json:"comment" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
}

// Availability نموذج التوفر
type Availability struct {
	Available      bool     `json:"available"`
	ServiceID      string   `json:"service_id"`
	Date           string   `json:"date"`
	Time           string   `json:"time"`
	Guests         int      `json:"guests"`
	Message        string   `json:"message"`
	SuggestedTimes []string `json:"suggested_times"`
	CheckedAt      time.Time `json:"checked_at"`
}

// ServicesStats نموذج إحصائيات الخدمات
type ServicesStats struct {
	UserID         string    `json:"user_id"`
	Timeframe      string    `json:"timeframe"`
	TotalServices  int       `json:"total_services"`
	ActiveServices int       `json:"active_services"`
	TotalOrders    int       `json:"total_orders"`
	AverageRating  float64   `json:"average_rating"`
	Revenue        float64   `json:"revenue,omitempty"`
	CalculatedAt   time.Time `json:"calculated_at"`
}

// TimeSlot نموذج الفترة الزمنية
type TimeSlot struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	ServiceID string    `json:"service_id" gorm:"not null;index"`
	Date      string    `json:"date" gorm:"not null"` // YYYY-MM-DD
	StartTime string    `json:"start_time" gorm:"not null"` // HH:MM
	EndTime   string    `json:"end_time" gorm:"not null"`   // HH:MM
	Available bool      `json:"available" gorm:"default:true"`
	Booked    bool      `json:"booked" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ServicesGrowth نموذج نمو الخدمات
type ServicesGrowth struct {
	SellerID      string    `json:"seller_id"`
	Period        string    `json:"period"`
	NewServices   int       `json:"new_services"`
	TotalOrders   int       `json:"total_orders"`
	Revenue       float64   `json:"revenue"`
	GrowthRate    float64   `json:"growth_rate"`
	CalculatedAt  time.Time `json:"calculated_at"`
}

