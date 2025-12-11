package models

import "time"

// ================================
// User
// ================================

type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	PasswordHash  string    `json:"-"` // لن يتم إرجاعه في JSON
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Phone         string    `json:"phone,omitempty"`
	Avatar        string    `json:"avatar,omitempty"`
	Role          string    `json:"role"`
	Status        string    `json:"status"` // active / inactive / banned / deleted
	EmailVerified bool      `json:"email_verified"`
	Settings      string    `json:"settings,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LastLogin     time.Time `json:"last_login,omitempty"`
}

// ================================
// Service
// ================================

type Service struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Duration    int       `json:"duration"`
	CategoryID  string    `json:"category_id"`
	ProviderID  string    `json:"provider_id"`
	Images      []string  `json:"images"`
	Tags        []string  `json:"tags"`
	IsActive    bool      `json:"is_active"`
	IsFeatured  bool      `json:"is_featured"`
	Rating      float64   `json:"rating"`
	ReviewCount int       `json:"review_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ================================
// Category
// ================================

type Category struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Image     string    `json:"image,omitempty"`
	Description string  `json:"description,omitempty"`
	ParentID   string   `json:"parent_id,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ================================
// Order
// ================================

type Order struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ServiceID string    `json:"service_id"`
	Status    string    `json:"status"` // pending, completed, cancelled
	Amount    float64   `json:"amount"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ================================
// Payment
// ================================

type Payment struct {
	ID            string    `json:"id"`
	OrderID       string    `json:"order_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"` // pending, processing, completed, failed, refunded
	PaymentMethod string    `json:"payment_method,omitempty"`
	TransactionID string    `json:"transaction_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ================================
// File
// ================================

type File struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Size      int64     `json:"size,omitempty"`
	Type      string    `json:"type,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ================================
// SystemLog
// ================================

type SystemLog struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id,omitempty"`
	Level     string    `json:"level"` // debug, info, warn, error
	Action    string    `json:"action"`
	Resource  string    `json:"resource,omitempty"`
	Details   string    `json:"details,omitempty"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ================================
// Notification
// ================================

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"` // info, success, warning, error
	IsRead    bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

// ================================
// Cart
// ================================

type Cart struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Items     []CartItem `json:"items"`
	Total     float64    `json:"total"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CartItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// ================================
// AuthToken (للتوافق مع services.go إذا لزم)
// ================================

type AuthToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// ================================
// Stats (للتوافق مع services.go)
// ================================

type UserStats struct {
	TotalOrders   int     `json:"total_orders"`
	TotalSpent    float64 `json:"total_spent"`
	ActiveSince   string  `json:"active_since"`
	ServicesCount int     `json:"services_count"`
}

type OrderStats struct {
	TotalOrders   int     `json:"total_orders"`
	PendingOrders int     `json:"pending_orders"`
	Completed     int     `json:"completed_orders"`
	Cancelled     int     `json:"cancelled_orders"`
	TotalRevenue  float64 `json:"total_revenue"`
	AvgOrderValue float64 `json:"avg_order_value"`
}

type DashboardStats struct {
	TotalUsers      int64   `json:"total_users"`
	TotalServices   int64   `json:"total_services"`
	TotalOrders     int64   `json:"total_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
	ActiveUsers     int64   `json:"active_users"`
	PendingOrders   int64   `json:"pending_orders"`
	CompletedOrders int64   `json:"completed_orders"`
}

// ================================
// Request/Response Types (للتوافق مع services.go)
// ================================

type AuthRegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Phone     string `json:"phone,omitempty"`
	Password  string `json:"password"`
}

type AuthLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User         *User     `json:"user"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type UserUpdateRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
}

type ServiceCreateRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Duration    int      `json:"duration"`
	CategoryID  string   `json:"category_id"`
	ProviderID  string   `json:"provider_id"`
	Images      []string `json:"images"`
	Tags        []string `json:"tags"`
}

type ServiceUpdateRequest struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Price       float64  `json:"price,omitempty"`
	Duration    int      `json:"duration,omitempty"`
	CategoryID  string   `json:"category_id,omitempty"`
	Images      []string `json:"images,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	IsActive    bool     `json:"is_active"`
	IsFeatured  bool     `json:"is_featured"`
}

type CategoryCreateRequest struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Image string `json:"image,omitempty"`
}

type CategoryUpdateRequest struct {
	Name     string `json:"name,omitempty"`
	Slug     string `json:"slug,omitempty"`
	Image    string `json:"image,omitempty"`
	IsActive bool   `json:"is_active"`
}

type OrderCreateRequest struct {
	UserID    string  `json:"user_id"`
	ServiceID string  `json:"service_id"`
	Amount    float64 `json:"amount"`
	Notes     string  `json:"notes,omitempty"`
}

type PaymentIntentRequest struct {
	OrderID   string  `json:"order_id"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Customer  string  `json:"customer,omitempty"`
	ReturnURL string  `json:"return_url,omitempty"`
}

type UploadRequest struct {
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
	FileSize int64  `json:"file_size"`
}

type UploadResult struct {
	URL      string    `json:"url"`
	FileName string    `json:"file_name"`
	FileType string    `json:"file_type"`
	FileSize int64     `json:"file_size"`
	Uploaded time.Time `json:"uploaded"`
}

type NotificationCreateRequest struct {
	UserID  string `json:"user_id"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

// ================================
// Query Params (للتوافق مع services.go)
// ================================

type UserQueryParams struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Role  string `json:"role,omitempty"`
	Email string `json:"email,omitempty"`
}

type ServiceQueryParams struct {
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
	CategoryID string  `json:"category_id,omitempty"`
	ProviderID string  `json:"provider_id,omitempty"`
	MinPrice   float64 `json:"min_price,omitempty"`
	MaxPrice   float64 `json:"max_price,omitempty"`
	IsActive   bool    `json:"is_active"`
	IsFeatured bool    `json:"is_featured"`
	Search     string  `json:"search,omitempty"`
}

type CategoryQueryParams struct {
	Page     int  `json:"page"`
	Limit    int  `json:"limit"`
	IsActive bool `json:"is_active"`
}

type OrderQueryParams struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Status string `json:"status,omitempty"`
	UserID string `json:"user_id,omitempty"`
}

type PaymentQueryParams struct {
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
	Status   string    `json:"status,omitempty"`
	UserID   string    `json:"user_id,omitempty"`
	OrderID  string    `json:"order_id,omitempty"`
	FromDate time.Time `json:"from_date,omitempty"`
	ToDate   time.Time `json:"to_date,omitempty"`
}

type NotificationQueryParams struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	UserID string `json:"user_id,omitempty"`
	Type   string `json:"type,omitempty"`
	Read   *bool  `json:"read,omitempty"`
}

type SystemLogQuery struct {
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
	Level    string    `json:"level,omitempty"`
	UserID   string    `json:"user_id,omitempty"`
	FromDate time.Time `json:"from_date,omitempty"`
	ToDate   time.Time `json:"to_date,omitempty"`
}

// ================================
// Validation Types
// ================================

type PaymentValidation struct {
	Valid    bool                   `json:"valid"`
	Reason   string                 `json:"reason,omitempty"`
	Payment  *Payment               `json:"payment,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type PaymentResult struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	PaymentID  string      `json:"payment_id,omitempty"`
	OrderID    string      `json:"order_id,omitempty"`
	Amount     float64     `json:"amount,omitempty"`
	Currency   string      `json:"currency,omitempty"`
	Status     string      `json:"status,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
}

type PaymentIntent struct {
	ID           string                 `json:"id"`
	ClientSecret string                 `json:"client_secret"`
	Status       string                 `json:"status"`
	Amount       float64                `json:"amount"`
	Currency     string                 `json:"currency"`
	CreatedAt    time.Time              `json:"created_at"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

type CategoryNode struct {
	Category  *Category    `json:"category"`
	Children  []CategoryNode `json:"children"`
	Services  int          `json:"services_count"`
}

// ================================
// Token Claims
// ================================

type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
}

// ================================
// Change Password
// ================================

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}