package models

import (
	"time"
)

// ================================
// النماذج الأساسية (Core Models)
// ================================

// User نموذج المستخدم
type User struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	Email         string    `json:"email" gorm:"not null;uniqueIndex"`
	Username      string    `json:"username" gorm:"not null;uniqueIndex"`
	Password      string    `json:"-" gorm:"not null"`
	FirstName     string    `json:"first_name" gorm:"not null"`
	LastName      string    `json:"last_name" gorm:"not null"`
	Phone         string    `json:"phone,omitempty"`
	Avatar        string    `json:"avatar,omitempty"`
	Role          string    `json:"role" gorm:"default:'user'"`
	Status        string    `json:"status" gorm:"default:'active'"`
	EmailVerified bool      `json:"email_verified" gorm:"default:false"`
	LastLogin     time.Time `json:"last_login,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Service نموذج الخدمة
type Service struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	Price       float64   `json:"price" gorm:"not null"`
	Duration    int       `json:"duration" gorm:"not null"` // بالمدة بالدقائق
	CategoryID  string    `json:"category_id" gorm:"not null;index"`
	ProviderID  string    `json:"provider_id" gorm:"not null;index"`
	Images      []string  `json:"images" gorm:"type:json;serializer:json"`
	Tags        []string  `json:"tags" gorm:"type:json;serializer:json"`
	Rating      float64   `json:"rating" gorm:"default:0"`
	TotalOrders int       `json:"total_orders" gorm:"default:0"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	IsFeatured  bool      `json:"is_featured" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Content نموذج المحتوى
type Content struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Content     string    `json:"content" gorm:"type:text;not null"`
	Type        string    `json:"type" gorm:"not null;index"` // article, blog, news, etc.
	AuthorID    string    `json:"author_id" gorm:"not null;index"`
	Slug        string    `json:"slug" gorm:"not null;uniqueIndex"`
	Image       string    `json:"image,omitempty"`
	Tags        []string  `json:"tags" gorm:"type:json;serializer:json"`
	IsPublished bool      `json:"is_published" gorm:"default:false"`
	Views       int       `json:"views" gorm:"default:0"`
	Likes       int       `json:"likes" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PublishedAt time.Time `json:"published_at,omitempty"`
}

// Notification نموذج الإشعار
type Notification struct {
	ID        string                 `json:"id" gorm:"primaryKey"`
	UserID    string                 `json:"user_id" gorm:"not null;index"`
	Title     string                 `json:"title" gorm:"not null"`
	Message   string                 `json:"message" gorm:"not null"`
	Type      string                 `json:"type" gorm:"not null;index"` // info, success, warning, error
	Data      map[string]interface{} `json:"data" gorm:"type:json;serializer:json"`
	IsRead    bool                   `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time              `json:"created_at"`
	ReadAt    time.Time              `json:"read_at,omitempty"`
}

// Review نموذج التقييم
type Review struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	UserID     string    `json:"user_id" gorm:"not null;index"`
	ServiceID  string    `json:"service_id" gorm:"not null;index"`
	OrderID    string    `json:"order_id" gorm:"not null;index"`
	Rating     int       `json:"rating" gorm:"not null;check:rating>=1 AND rating<=5"`
	Title      string    `json:"title,omitempty"`
	Comment    string    `json:"comment,omitempty" gorm:"type:text"`
	IsVerified bool      `json:"is_verified" gorm:"default:false"`
	IsHelpful  int       `json:"is_helpful" gorm:"default:0"`
	IsReported bool      `json:"is_reported" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

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
	ID          string    `json:"id"`
	ServiceID   string    `json:"service_id"`
	ServiceName string    `json:"service_name"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	Image       string    `json:"image,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Category فئة
type Category struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"not null;uniqueIndex"`
	Description  string    `json:"description" gorm:"type:text"`
	Slug         string    `json:"slug" gorm:"not null;uniqueIndex"`
	ParentID     string    `json:"parent_id,omitempty" gorm:"index"`
	Icon         string    `json:"icon,omitempty"`
	Image        string    `json:"image,omitempty"`
	Color        string    `json:"color,omitempty"`
	SortOrder    int       `json:"sort_order" gorm:"default:0"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	ServiceCount int       `json:"service_count,omitempty" gorm:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Store متجر
type Store struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"not null"`
	Slug         string    `json:"slug" gorm:"not null;uniqueIndex"`
	Description  string    `json:"description" gorm:"type:text"`
	OwnerID      string    `json:"owner_id" gorm:"not null;index"`
	Banner       string    `json:"banner,omitempty"`
	Logo         string    `json:"logo,omitempty"`
	ContactEmail string    `json:"contact_email,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	Address      string    `json:"address,omitempty"`
	Rating       float64   `json:"rating" gorm:"default:0"`
	TotalSales   int       `json:"total_sales" gorm:"default:0"`
	IsVerified   bool      `json:"is_verified" gorm:"default:false"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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
	ID           string            `json:"id" gorm:"primaryKey"`
	Filename     string            `json:"filename" gorm:"not null"`
	OriginalName string            `json:"original_name" gorm:"not null"`
	Path         string            `json:"path" gorm:"not null"`
	URL          string            `json:"url" gorm:"not null"`
	Size         int64             `json:"size" gorm:"not null"`
	MimeType     string            `json:"mime_type" gorm:"not null"`
	Extension    string            `json:"extension" gorm:"not null"`
	ContentType  string            `json:"content_type" gorm:"not null"` // أضيف هذا الحقل
	Metadata     map[string]string `json:"metadata" gorm:"type:json;serializer:json"`
	UserID       string            `json:"user_id" gorm:"not null;index"`
	IsPublic     bool              `json:"is_public" gorm:"default:false"`
	UploadedAt   time.Time         `json:"uploaded_at"`
	CreatedAt    time.Time         `json:"created_at"` // أضيف هذا الحقل
	UpdatedAt    time.Time         `json:"updated_at"` // أضيف هذا الحقل
}

// ================================
// نماذج المعاملات (Transaction Models)
// ================================

// Order طلب
type Order struct {
	ID            string       `json:"id" gorm:"primaryKey"`
	UserID        string       `json:"user_id" gorm:"not null;index"`
	ServiceID     string       `json:"service_id,omitempty" gorm:"index"` // جعلته اختياريًا
	SellerID      string       `json:"seller_id" gorm:"not null;index"`
	Items         []OrderItem  `json:"items" gorm:"type:json;serializer:json"`
	Status        string       `json:"status" gorm:"not null;index"` // pending, confirmed, processing, shipped, delivered, cancelled, refunded
	TotalAmount   float64      `json:"total_amount" gorm:"not null"`
	Discount      float64      `json:"discount" gorm:"default:0"`
	Tax           float64      `json:"tax" gorm:"default:0"`
	Shipping      float64      `json:"shipping" gorm:"default:0"`
	FinalAmount   float64      `json:"final_amount" gorm:"not null"`
	PaymentStatus string       `json:"payment_status" gorm:"not null;index"` // pending, paid, failed, refunded
	PaymentMethod string       `json:"payment_method" gorm:"not null"`
	ShippingInfo  ShippingInfo `json:"shipping_info" gorm:"type:json;serializer:json"`
	CustomerNotes string       `json:"customer_notes,omitempty" gorm:"type:text"`
	CancelledAt   time.Time    `json:"cancelled_at,omitempty"`
	DeliveredAt   time.Time    `json:"delivered_at,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// OrderItem عنصر في الطلب
type OrderItem struct {
	ID          string    `json:"id"`
	ServiceID   string    `json:"service_id"`
	ServiceName string    `json:"service_name"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	Image       string    `json:"image,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ShippingInfo معلومات الشحن
type ShippingInfo struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Address        string `json:"address"`
	City           string `json:"city"`
	Country        string `json:"country"`
	PostalCode     string `json:"postal_code"`
	ShippingMethod string `json:"shipping_method"`
	TrackingNumber string `json:"tracking_number,omitempty"`
}

// Payment نموذج الدفع
type Payment struct {
	ID            string                 `json:"id" gorm:"primaryKey"`
	OrderID       string                 `json:"order_id" gorm:"not null;index"`
	UserID        string                 `json:"user_id" gorm:"not null;index"`
	Amount        float64                `json:"amount" gorm:"not null"`
	Currency      string                 `json:"currency" gorm:"not null;default:'USD'"`
	Status        string                 `json:"status" gorm:"not null;index"` // pending, succeeded, failed, refunded
	PaymentMethod string                 `json:"payment_method" gorm:"not null"`
	PaymentIntent string                 `json:"payment_intent,omitempty"` // معرف الدفع من مزود الدفع
	Metadata      map[string]interface{} `json:"metadata" gorm:"type:json;serializer:json"`
	FailureReason string                 `json:"failure_reason,omitempty"`
	PaidAt        time.Time              `json:"paid_at,omitempty"`
	RefundedAt    time.Time              `json:"refunded_at,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// ================================
// نماذج التحليل والإحصائيات (Analytics Models)
// ================================

// Analytics نموذج التحليلات
type Analytics struct {
	ID        string                 `json:"id" gorm:"primaryKey"`
	Type      string                 `json:"type" gorm:"not null;index"` // page_view, event, conversion, etc.
	UserID    string                 `json:"user_id,omitempty" gorm:"index"`
	SessionID string                 `json:"session_id" gorm:"not null"`
	Page      string                 `json:"page,omitempty"`
	Action    string                 `json:"action,omitempty"`
	Data      map[string]interface{} `json:"data" gorm:"type:json;serializer:json"`
	IPAddress string                 `json:"ip_address,omitempty"`
	UserAgent string                 `json:"user_agent,omitempty"`
	Referrer  string                 `json:"referrer,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

// ================================
// نماذج النظام (System Models)
// ================================

// SystemLog سجل النظام
type SystemLog struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Level     string    `json:"level" gorm:"not null;index"` // info, warning, error, debug
	Message   string    `json:"message" gorm:"not null"`
	Module    string    `json:"module" gorm:"not null;index"`
	UserID    string    `json:"user_id,omitempty" gorm:"index"`
	IPAddress string    `json:"ip_address,omitempty"`
	Data      string    `json:"data,omitempty" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
}

// Setting إعدادات النظام
type Setting struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Key         string    `json:"key" gorm:"not null;uniqueIndex"`
	Value       string    `json:"value" gorm:"type:text;not null"`
	Type        string    `json:"type" gorm:"not null"` // string, number, boolean, json
	Description string    `json:"description,omitempty"`
	Category    string    `json:"category" gorm:"not null;index"`
	IsPublic    bool      `json:"is_public" gorm:"default:false"`
	UpdatedBy   string    `json:"updated_by" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ================================
// نماذج إضافية (Additional Models)
// ================================

// Coupon نموذج الكوبون
type Coupon struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Code         string    `json:"code" gorm:"not null;uniqueIndex"`
	Description  string    `json:"description,omitempty"`
	DiscountType string    `json:"discount_type" gorm:"not null"` // percentage, fixed
	DiscountValue float64  `json:"discount_value" gorm:"not null"`
	MinAmount    float64   `json:"min_amount,omitempty"`
	MaxDiscount  float64   `json:"max_discount,omitempty"`
	UsageLimit   int       `json:"usage_limit,omitempty"`
	UsedCount    int       `json:"used_count" gorm:"default:0"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Wishlist قائمة الرغبات
type Wishlist struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null;index"`
	ServiceID string    `json:"service_id" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
}

// Subscription نموذج الاشتراك
type Subscription struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	UserID      string    `json:"user_id" gorm:"not null;index"`
	PlanID      string    `json:"plan_id" gorm:"not null;index"`
	Status      string    `json:"status" gorm:"not null;index"` // active, cancelled, expired
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	RenewalDate time.Time `json:"renewal_date,omitempty"`
	CancelledAt time.Time `json:"cancelled_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Session نموذج الجلسة
type Session struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null;index"`
	Token     string    `json:"token" gorm:"not null;uniqueIndex"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}