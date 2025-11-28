package models

import (
	"time"
)

// ================================
// النماذج الأساسية (Core Models)
// ================================

// User, Service, Content, Notification, Review (كما في السابق)

// Cart عربة التسوق
type Cart struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	UserID      string     `json:"user_id" gorm:"not null;index"`
	Items       []CartItem `json:"items" gorm:"type:json;serializer:json"`
	TotalAmount float64    `json:"total_amount" gorm:"default:0"`
	Discount    float64    `json:"discount" gorm:"default:0"`
	Tax         float64    `json:"tax" gorm:"default:0"`
	Shipping    float64    `json:"shipping" gorm:"default:0"`
	FinalAmount float64    `json:"final_amount" gorm:"default:0"`
	CouponCode  string     `json:"coupon_code,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CartItem عنصر في عربة التسوق
type CartItem struct {
	ID          string  `json:"id"`
	ServiceID   string  `json:"service_id"`
	ServiceName string  `json:"service_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Image       string  `json:"image,omitempty"`
}

// Category فئة
type Category struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;uniqueIndex"`
	Description string    `json:"description" gorm:"type:text"`
	Slug        string    `json:"slug" gorm:"not null;uniqueIndex"`
	ParentID    string    `json:"parent_id,omitempty" gorm:"index"`
	Icon        string    `json:"icon,omitempty"`
	Image       string    `json:"image,omitempty"`
	Color       string    `json:"color,omitempty"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	ServiceCount int      `json:"service_count,omitempty" gorm:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Store متجر
type Store struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Slug        string    `json:"slug" gorm:"not null;uniqueIndex"`
	Description string    `json:"description" gorm:"type:text"`
	OwnerID     string    `json:"owner_id" gorm:"not null;index"`
	Banner      string    `json:"banner,omitempty"`
	Logo        string    `json:"logo,omitempty"`
	ContactEmail string   `json:"contact_email,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	Address     string    `json:"address,omitempty"`
	Rating      float64   `json:"rating" gorm:"default:0"`
	TotalSales  int       `json:"total_sales" gorm:"default:0"`
	IsVerified  bool      `json:"is_verified" gorm:"default:false"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Strategy استراتيجية
type Strategy struct {
	ID          string                 `json:"id" gorm:"primaryKey"`
	Name        string                 `json:"name" gorm:"not null"`
	Description string                 `json:"description" gorm:"type:text"`
	Type        string                 `json:"type" gorm:"not null;index"`
	Parameters  map[string]interface{} `json:"parameters" gorm:"type:json;serializer:json"`
	Rules       []StrategyRule         `json:"rules" gorm:"type:json;serializer:json"`
	IsActive    bool                   `json:"is_active" gorm:"default:true"`
	CreatedBy   string                 `json:"created_by" gorm:"not null;index"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// StrategyRule قاعدة استراتيجية
type StrategyRule struct {
	ID        string      `json:"id"`
	Condition string      `json:"condition"`
	Action    string      `json:"action"`
	Value     interface{} `json:"value"`
	Priority  int         `json:"priority"`
}

// File ملف
type File struct {
	ID          string            `json:"id" gorm:"primaryKey"`
	Filename    string            `json:"filename" gorm:"not null"`
	OriginalName string           `json:"original_name" gorm:"not null"`
	Path        string            `json:"path" gorm:"not null"`
	URL         string            `json:"url" gorm:"not null"`
	Size        int64             `json:"size" gorm:"not null"`
	MimeType    string            `json:"mime_type" gorm:"not null"`
	Extension   string            `json:"extension" gorm:"not null"`
	Metadata    map[string]string `json:"metadata" gorm:"type:json;serializer:json"`
	UserID      string            `json:"user_id" gorm:"not null;index"`
	IsPublic    bool              `json:"is_public" gorm:"default:false"`
	UploadedAt  time.Time         `json:"uploaded_at"`
}

// ================================
// نماذج المعاملات (Transaction Models)
// ================================

// Order طلب (محدث)
type Order struct {
	ID           string      `json:"id" gorm:"primaryKey"`
	UserID       string      `json:"user_id" gorm:"not null;index"`
	ServiceID    string      `json:"service_id" gorm:"not null;index"`
	SellerID     string      `json:"seller_id" `