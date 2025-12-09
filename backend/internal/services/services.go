package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/models"
)

// ================================
// هياكل المعاملات المحدثة
// ================================

type (
	ReviewQueryParams struct {
		Page   int    `json:"page"`
		Limit  int    `json:"limit"`
		Rating int    `json:"rating"`
		SortBy string `json:"sort_by"`
	}
)

// ================================
// الواجهات الرئيسية
// ================================

type (
	AuthService interface {
		Register(ctx context.Context, req AuthRegisterRequest) (*AuthResponse, error)
		Login(ctx context.Context, req AuthLoginRequest) (*AuthResponse, error)
		Logout(ctx context.Context, token string) error
		RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)
		VerifyToken(ctx context.Context, token string) (*TokenClaims, error)
		ForgotPassword(ctx context.Context, email string) error
		ResetPassword(ctx context.Context, token string, newPassword string) error
		ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error
	}

	UserService interface {
		GetProfile(ctx context.Context, userID string) (*models.User, error)
		UpdateProfile(ctx context.Context, userID string, req UserUpdateRequest) (*models.User, error)
		UpdateAvatar(ctx context.Context, userID string, avatarURL string) error
		DeleteAccount(ctx context.Context, userID string) error
		SearchUsers(ctx context.Context, query string, params UserQueryParams) ([]models.User, error)
		GetUserStats(ctx context.Context, userID string) (*UserStats, error)
	}

	ServiceService interface {
		CreateService(ctx context.Context, req ServiceCreateRequest) (*models.Service, error)
		GetServiceByID(ctx context.Context, serviceID string) (*models.Service, error)
		UpdateService(ctx context.Context, serviceID string, req ServiceUpdateRequest) (*models.Service, error)
		DeleteService(ctx context.Context, serviceID string) error
		GetServices(ctx context.Context, params ServiceQueryParams) ([]models.Service, error)
		SearchServices(ctx context.Context, query string, params ServiceQueryParams) ([]models.Service, error)
		GetFeaturedServices(ctx context.Context) ([]models.Service, error)
		GetSimilarServices(ctx context.Context, serviceID string) ([]models.Service, error)
	}

	CategoryService interface {
		GetCategories(ctx context.Context, params CategoryQueryParams) ([]models.Category, error)
		GetCategoryByID(ctx context.Context, categoryID string) (*models.Category, error)
		CreateCategory(ctx context.Context, req CategoryCreateRequest) (*models.Category, error)
		UpdateCategory(ctx context.Context, categoryID string, req CategoryUpdateRequest) (*models.Category, error)
		DeleteCategory(ctx context.Context, categoryID string) error
		GetCategoryTree(ctx context.Context) ([]CategoryNode, error)
	}

	OrderService interface {
		CreateOrder(ctx context.Context, req OrderCreateRequest) (*models.Order, error)
		GetOrderByID(ctx context.Context, orderID string) (*models.Order, error)
		GetUserOrders(ctx context.Context, userID string, params OrderQueryParams) ([]models.Order, error)
		UpdateOrderStatus(ctx context.Context, orderID string, status string, notes string) (*models.Order, error)
		CancelOrder(ctx context.Context, orderID string, reason string) (*models.Order, error)
		GetOrderStats(ctx context.Context, timeframe string) (*OrderStats, error)
	}

	PaymentService interface {
		CreatePaymentIntent(ctx context.Context, req PaymentIntentRequest) (*PaymentIntent, error)
		ConfirmPayment(ctx context.Context, paymentID string, confirmationData map[string]interface{}) (*PaymentResult, error)
		GetPaymentHistory(ctx context.Context, userID string, params PaymentQueryParams) ([]models.Payment, error)
		ValidatePayment(ctx context.Context, paymentData map[string]interface{}) (*PaymentValidation, error)
	}

	UploadService interface {
		UploadFile(ctx context.Context, req UploadRequest) (*UploadResult, error)
		DeleteFile(ctx context.Context, fileID string) error
		GetFile(ctx context.Context, fileID string) (*models.File, error)
		GetUserFiles(ctx context.Context, userID string, params FileQueryParams) ([]models.File, error)
		GeneratePresignedURL(ctx context.Context, req PresignedURLRequest) (*PresignedURL, error)
		ValidateFile(ctx context.Context, fileInfo models.File) (*FileValidation, error)
		GetUploadQuota(ctx context.Context, userID string) (*UploadQuota, error)
	}

	NotificationService interface {
		CreateNotification(ctx context.Context, req NotificationCreateRequest) (*models.Notification, error)
		GetUserNotifications(ctx context.Context, userID string, params NotificationQueryParams) ([]models.Notification, error)
		MarkAsRead(ctx context.Context, notificationID string) error
		MarkAllAsRead(ctx context.Context, userID string) error
		DeleteNotification(ctx context.Context, notificationID string) error
		GetUnreadCount(ctx context.Context, userID string) (int64, error)
	}

	AdminService interface {
		GetDashboardStats(ctx context.Context) (*DashboardStats, error)
		GetUsers(ctx context.Context, params UserQueryParams) ([]models.User, error)
		GetSystemLogs(ctx context.Context, params SystemLogQuery) ([]models.SystemLog, error)
		UpdateSystemSettings(ctx context.Context, settings []models.Setting) error
		BanUser(ctx context.Context, userID string, reason string) error
		UnbanUser(ctx context.Context, userID string) error
	}

	CacheService interface {
		Get(key string) (interface{}, error)
		Set(key string, value interface{}, expiration time.Duration) error
		Delete(key string) error
		Exists(key string) (bool, error)
		Flush() error
	}
)

// ================================
// تطبيقات D1 SQL
// ================================

type (
	authServiceImpl struct {
		db *sql.DB
	}

	userServiceImpl struct {
		db *sql.DB
	}

	serviceServiceImpl struct {
		db *sql.DB
	}

	categoryServiceImpl struct {
		db *sql.DB
	}

	orderServiceImpl struct {
		db *sql.DB
	}

	paymentServiceImpl struct {
		db *sql.DB
	}

	uploadServiceImpl struct {
		db *sql.DB
	}

	notificationServiceImpl struct {
		db *sql.DB
	}

	adminServiceImpl struct {
		db *sql.DB
	}

	cacheServiceImpl struct {
		store map[string]interface{}
	}
)

// ================================
// أمثلة على تنفيذ AuthService باستخدام D1
// ================================

func (s *authServiceImpl) Register(ctx context.Context, req AuthRegisterRequest) (*AuthResponse, error) {
	userID := fmt.Sprintf("user_%d", time.Now().UnixNano())

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO users (id, email, username, password, first_name, last_name, phone, role, status, email_verified, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, req.Email, req.Username, "hashed_password", req.FirstName, req.LastName, req.Phone, "user", "active", false, time.Now(), time.Now(),
	)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:            userID,
		Email:         req.Email,
		Username:      req.Username,
		Role:          "user",
		Status:        "active",
		EmailVerified: false,
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  "access_" + userID,
		RefreshToken: "refresh_" + userID,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}, nil
}

func (s *authServiceImpl) Login(ctx context.Context, req AuthLoginRequest) (*AuthResponse, error) {
	var user models.User
	row := s.db.QueryRowContext(ctx, "SELECT id, email, username, role, status, email_verified FROM users WHERE email = ?", req.Email)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Role, &user.Status, &user.EmailVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("المستخدم غير موجود")
		}
		return nil, err
	}

	return &AuthResponse{
		User:         &user,
		AccessToken:  "access_" + user.ID,
		RefreshToken: "refresh_" + user.ID,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}, nil
}

// ================================
// دوال الإنشاء باستخدام D1
// ================================

func NewAuthService(db *sql.DB) AuthService {
	return &authServiceImpl{db: db}
}

func NewUserService(db *sql.DB) UserService {
	return &userServiceImpl{db: db}
}

func NewServiceService(db *sql.DB) ServiceService {
	return &serviceServiceImpl{db: db}
}

func NewCategoryService(db *sql.DB) CategoryService {
	return &categoryServiceImpl{db: db}
}

func NewOrderService(db *sql.DB) OrderService {
	return &orderServiceImpl{db: db}
}

func NewPaymentService(db *sql.DB) PaymentService {
	return &paymentServiceImpl{db: db}
}

func NewUploadService(db *sql.DB) UploadService {
	return &uploadServiceImpl{db: db}
}

func NewNotificationService(db *sql.DB) NotificationService {
	return &notificationServiceImpl{db: db}
}

func NewAdminService(db *sql.DB) AdminService {
	return &adminServiceImpl{db: db}
}

func NewCacheService() CacheService {
	return &cacheServiceImpl{store: make(map[string]interface{})}
}

// ================================
// Service Container مع D1
// ================================

type ServiceContainer struct {
	Auth         AuthService
	User         UserService
	Service      ServiceService
	Category     CategoryService
	Order        OrderService
	Payment      PaymentService
	Upload       UploadService
	Notification NotificationService
	Admin        AdminService
	Cache        CacheService
}

func NewServiceContainer(d1db *sql.DB) *ServiceContainer {
	return &ServiceContainer{
		Auth:         NewAuthService(d1db),
		User:         NewUserService(d1db),
		Service:      NewServiceService(d1db),
		Category:     NewCategoryService(d1db),
		Order:        NewOrderService(d1db),
		Payment:      NewPaymentService(d1db),
		Upload:       NewUploadService(d1db),
		Notification: NewNotificationService(d1db),
		Admin:        NewAdminService(d1db),
		Cache:        NewCacheService(),
	}
}