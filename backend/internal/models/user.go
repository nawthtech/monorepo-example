package models

import (
	"time"
)

// User نموذج المستخدم
type User struct {
	ID              string    `json:"id" gorm:"primaryKey"`
	Email           string    `json:"email" gorm:"uniqueIndex;not null"`
	Username        string    `json:"username" gorm:"uniqueIndex;not null"`
	FirstName       string    `json:"first_name" gorm:"not null"`
	LastName        string    `json:"last_name" gorm:"not null"`
	Phone           string    `json:"phone,omitempty"`
	Avatar          string    `json:"avatar,omitempty"`
	Role            string    `json:"role" gorm:"default:'user';index"` // user, seller, admin, super_admin
	Status          string    `json:"status" gorm:"default:'active';index"` // active, inactive, suspended, pending
	EmailVerified   bool      `json:"email_verified" gorm:"default:false"`
	PhoneVerified   bool      `json:"phone_verified" gorm:"default:false"`
	LastLogin       time.Time `json:"last_login,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// UserProfile ملف المستخدم الشخصي
type UserProfile struct {
	UserID      string `json:"user_id" gorm:"primaryKey"`
	Bio         string `json:"bio,omitempty" gorm:"type:text"`
	Website     string `json:"website,omitempty"`
	Location    string `json:"location,omitempty"`
	SocialLinks map[string]string `json:"social_links,omitempty" gorm:"type:json;serializer:json"`
	Skills      []string `json:"skills,omitempty" gorm:"type:json;serializer:json"`
	Languages   []string `json:"languages,omitempty" gorm:"type:json;serializer:json"`
	Experience  string `json:"experience,omitempty" gorm:"type:text"`
}

// UserStats إحصائيات المستخدم
type UserStats struct {
	UserID          string  `json:"user_id" gorm:"primaryKey"`
	TotalServices   int     `json:"total_services" gorm:"default:0"`
	ActiveServices  int     `json:"active_services" gorm:"default:0"`
	TotalOrders     int     `json:"total_orders" gorm:"default:0"`
	CompletedOrders int     `json:"completed_orders" gorm:"default:0"`
	TotalRevenue    float64 `json:"total_revenue" gorm:"default:0"`
	AverageRating   float64 `json:"average_rating" gorm:"default:0"`
	TotalReviews    int     `json:"total_reviews" gorm:"default:0"`
}

// SellerInfo معلومات البائع
type SellerInfo struct {
	UserID        string    `json:"user_id" gorm:"primaryKey"`
	StoreName     string    `json:"store_name" gorm:"not null"`
	StoreSlug     string    `json:"store_slug" gorm:"uniqueIndex;not null"`
	Description   string    `json:"description,omitempty" gorm:"type:text"`
	BannerImage   string    `json:"banner_image,omitempty"`
	Verified      bool      `json:"verified" gorm:"default:false"`
	Level         string    `json:"level" gorm:"default:'basic'"` // basic, pro, premium
	ResponseRate  float64   `json:"response_rate" gorm:"default:0"`
	ResponseTime  int       `json:"response_time" gorm:"default:0"` // بالدقائق
	JoinDate      time.Time `json:"join_date"`
	TotalSales    int       `json:"total_sales" gorm:"default:0"`
}

// UserCreateRequest طلب إنشاء مستخدم
type UserCreateRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required,min=3,max=50"`
	FirstName string `json:"first_name" binding:"required,min=2,max=50"`
	LastName  string `json:"last_name" binding:"required,min=2,max=50"`
	Password  string `json:"password" binding:"required,min=6"`
	Phone     string `json:"phone,omitempty"`
}

// UserUpdateRequest طلب تحديث مستخدم
type UserUpdateRequest struct {
	FirstName string `json:"first_name,omitempty" binding:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name,omitempty" binding:"omitempty,min=2,max=50"`
	Phone     string `json:"phone,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
}

// UserLoginRequest طلب تسجيل الدخول
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserChangePasswordRequest طلب تغيير كلمة المرور
type UserChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// UserResetPasswordRequest طلب إعادة تعيين كلمة المرور
type UserResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// UserVerifyEmailRequest طلب التحقق من البريد الإلكتروني
type UserVerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}