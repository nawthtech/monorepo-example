package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/nawthtech/nawthtech/backend/internal/logger"
)

// ========== هياكل التكوين ==========

// Cors تكوين CORS
type Cors struct {
	AllowedOrigins   []string `env:"ALLOWED_ORIGINS" envSeparator:","`
	AllowedMethods   []string `env:"ALLOWED_METHODS" envSeparator:","`
	AllowedHeaders   []string `env:"ALLOWED_HEADERS" envSeparator:","`
	ExposedHeaders   []string `env:"EXPOSED_HEADERS" envSeparator:","`
	AllowCredentials bool     `env:"ALLOW_CREDENTIALS"`
	MaxAge           int      `env:"MAX_AGE"`
}

// Redis تكوين Redis
type Redis struct {
	URL      string `env:"REDIS_URL"`
	Host     string `env:"REDIS_HOST"`
	Port     string `env:"REDIS_PORT"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB"`
}

// Cache تكوين التخزين المؤقت
type Cache struct {
	Enabled    bool          `env:"CACHE_ENABLED"`
	Prefix     string        `env:"CACHE_PREFIX"`
	DefaultTTL time.Duration `env:"CACHE_DEFAULT_TTL"`
	MaxRetries int           `env:"CACHE_MAX_RETRIES"`
}

// ServicesConfig تكوين الخدمات
type ServicesConfig struct {
	MaxServicesPerUser      int           `env:"SERVICES_MAX_PER_USER"`
	MaxActiveServices       int           `env:"SERVICES_MAX_ACTIVE"`
	DefaultPaginationLimit  int           `env:"SERVICES_PAGINATION_LIMIT"`
	MaxPaginationLimit      int           `env:"SERVICES_MAX_PAGINATION_LIMIT"`
	SearchCacheTTL          time.Duration `env:"SERVICES_SEARCH_CACHE_TTL"`
	FeaturedCacheTTL        time.Duration `env:"SERVICES_FEATURED_CACHE_TTL"`
	MaxImagesPerService     int           `env:"SERVICES_MAX_IMAGES"`
	MaxFeaturesPerService   int           `env:"SERVICES_MAX_FEATURES"`
	MaxTagsPerService       int           `env:"SERVICES_MAX_TAGS"`
	MinTitleLength          int           `env:"SERVICES_MIN_TITLE_LENGTH"`
	MaxTitleLength          int           `env:"SERVICES_MAX_TITLE_LENGTH"`
	MinDescriptionLength    int           `env:"SERVICES_MIN_DESCRIPTION_LENGTH"`
	MaxDescriptionLength    int           `env:"SERVICES_MAX_DESCRIPTION_LENGTH"`
	MinPrice                float64       `env:"SERVICES_MIN_PRICE"`
	MaxPrice                float64       `env:"SERVICES_MAX_PRICE"`
	MinDuration             int           `env:"SERVICES_MIN_DURATION"`
	MaxDuration             int           `env:"SERVICES_MAX_DURATION"`
	AutoApproveServices     bool          `env:"SERVICES_AUTO_APPROVE"`
	AllowServiceEditing     bool          `env:"SERVICES_ALLOW_EDITING"`
	EnableServiceReviews    bool          `env:"SERVICES_ENABLE_REVIEWS"`
	EnableServiceRatings    bool          `env:"SERVICES_ENABLE_RATINGS"`
	EnableServiceBookings   bool          `env:"SERVICES_ENABLE_BOOKINGS"`
	EnableServicePromotions bool          `env:"SERVICES_ENABLE_PROMOTIONS"`
	RateLimitCreate         int           `env:"SERVICES_RATE_LIMIT_CREATE"`
	RateLimitUpdate         int           `env:"SERVICES_RATE_LIMIT_UPDATE"`
	RateLimitSearch         int           `env:"SERVICES_RATE_LIMIT_SEARCH"`
}

// Upload تكوين الرفع
type Upload struct {
	MaxFileSize    int64    `env:"UPLOAD_MAX_FILE_SIZE"`
	AllowedTypes   []string `env:"UPLOAD_ALLOWED_TYPES" envSeparator:","`
	ImageMaxWidth  int      `env:"UPLOAD_IMAGE_MAX_WIDTH"`
	ImageMaxHeight int      `env:"UPLOAD_IMAGE_MAX_HEIGHT"`
	StoragePath    string   `env:"UPLOAD_STORAGE_PATH"`
}

// Email تكوين البريد
type Email struct {
	Enabled   bool   `env:"EMAIL_ENABLED"`
	Host      string `env:"EMAIL_HOST"`
	Port      int    `env:"EMAIL_PORT"`
	Username  string `env:"EMAIL_USERNAME"`
	Password  string `env:"EMAIL_PASSWORD"`
	FromEmail string `env:"EMAIL_FROM_EMAIL"`
	FromName  string `env:"EMAIL_FROM_NAME"`
}

// AuthConfig تكوين المصادقة
type AuthConfig struct {
	JWTSecret         string        `env:"JWT_SECRET"`
	JWTExpiration     time.Duration `env:"JWT_EXPIRATION"`
	RefreshExpiration time.Duration `env:"REFRESH_EXPIRATION"`
	BCryptCost        int           `env:"BCRYPT_COST"`
}

// DatabaseConfig تكوين قاعدة البيانات
type DatabaseConfig struct {
	URL          string `env:"DATABASE_URL"`
	MaxIdleConns int    `env:"DB_MAX_IDLE_CONNS"`
	MaxOpenConns int    `env:"DB_MAX_OPEN_CONNS"`
	MaxLifetime  int    `env:"DB_MAX_LIFETIME"`
}

// Config التكوين الرئيسي
type Config struct {
	Environment   string         `env:"ENVIRONMENT"`
	Port          string         `env:"PORT"`
	Version       string         `env:"APP_VERSION"`
	EncryptionKey string         `env:"ENCRYPTION_KEY"`
	Database      DatabaseConfig `envPrefix:"DB_"`
	Auth          AuthConfig     `envPrefix:"AUTH_"`
	Cors          Cors           `envPrefix:"CORS_"`
	Redis         Redis          `envPrefix:"REDIS_"`
	Cache         Cache          `envPrefix:"CACHE_"`
	Services      ServicesConfig `envPrefix:"SERVICES_"`
	Upload        Upload         `envPrefix:"UPLOAD_"`
	Email         Email          `envPrefix:"EMAIL_"`
}

// ========== متغيرات عامة ==========

var (
	appConfig *Config
)

// ========== التهيئة ==========

// Load تحميل الإعدادات
func Load() *Config {
	if appConfig != nil {
		return appConfig
	}

	appConfig = &Config{
		Environment:   getEnv("ENVIRONMENT", "development"),
		Port:          getEnv("PORT", "3000"),
		Version:       getEnv("APP_VERSION", "1.0.0"),
		EncryptionKey: getEnv("ENCRYPTION_KEY", "default-encryption-key-change-in-production"),
		Database: DatabaseConfig{
			URL:          getEnv("DATABASE_URL", ""),
			MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 100),
			MaxLifetime:  getEnvInt("DB_MAX_LIFETIME", 3600),
		},
		Auth: AuthConfig{
			JWTSecret:         getEnv("JWT_SECRET", "default-jwt-secret-change-in-production"),
			JWTExpiration:     getEnvDuration("JWT_EXPIRATION", 24*time.Hour),
			RefreshExpiration: getEnvDuration("REFRESH_EXPIRATION", 7*24*time.Hour),
			BCryptCost:        getEnvInt("BCRYPT_COST", 12),
		},
		Cors: Cors{
			AllowedOrigins:   getEnvSlice("ALLOWED_ORIGINS", []string{}, ","),
			AllowedMethods:   getEnvSlice("ALLOWED_METHODS", []string{}, ","),
			AllowedHeaders:   getEnvSlice("ALLOWED_HEADERS", []string{}, ","),
			ExposedHeaders:   getEnvSlice("EXPOSED_HEADERS", []string{}, ","),
			AllowCredentials: getEnvBool("ALLOW_CREDENTIALS", true),
			MaxAge:           getEnvInt("MAX_AGE", 86400),
		},
		Redis: Redis{
			URL:      getEnv("REDIS_URL", ""),
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Cache: Cache{
			Enabled:    getEnvBool("CACHE_ENABLED", true),
			Prefix:     getEnv("CACHE_PREFIX", "nawthtech:"),
			DefaultTTL: getEnvDuration("CACHE_DEFAULT_TTL", 1*time.Hour),
			MaxRetries: getEnvInt("CACHE_MAX_RETRIES", 3),
		},
		Services: ServicesConfig{
			MaxServicesPerUser:      getEnvInt("SERVICES_MAX_PER_USER", 50),
			MaxActiveServices:       getEnvInt("SERVICES_MAX_ACTIVE", 20),
			DefaultPaginationLimit:  getEnvInt("SERVICES_PAGINATION_LIMIT", 20),
			MaxPaginationLimit:      getEnvInt("SERVICES_MAX_PAGINATION_LIMIT", 100),
			SearchCacheTTL:          getEnvDuration("SERVICES_SEARCH_CACHE_TTL", 5*time.Minute),
			FeaturedCacheTTL:        getEnvDuration("SERVICES_FEATURED_CACHE_TTL", 30*time.Minute),
			MaxImagesPerService:     getEnvInt("SERVICES_MAX_IMAGES", 10),
			MaxFeaturesPerService:   getEnvInt("SERVICES_MAX_FEATURES", 20),
			MaxTagsPerService:       getEnvInt("SERVICES_MAX_TAGS", 15),
			MinTitleLength:          getEnvInt("SERVICES_MIN_TITLE_LENGTH", 3),
			MaxTitleLength:          getEnvInt("SERVICES_MAX_TITLE_LENGTH", 200),
			MinDescriptionLength:    getEnvInt("SERVICES_MIN_DESCRIPTION_LENGTH", 10),
			MaxDescriptionLength:    getEnvInt("SERVICES_MAX_DESCRIPTION_LENGTH", 2000),
			MinPrice:                getEnvFloat("SERVICES_MIN_PRICE", 0),
			MaxPrice:                getEnvFloat("SERVICES_MAX_PRICE", 1000000),
			MinDuration:             getEnvInt("SERVICES_MIN_DURATION", 1),
			MaxDuration:             getEnvInt("SERVICES_MAX_DURATION", 365),
			AutoApproveServices:     getEnvBool("SERVICES_AUTO_APPROVE", true),
			AllowServiceEditing:     getEnvBool("SERVICES_ALLOW_EDITING", true),
			EnableServiceReviews:    getEnvBool("SERVICES_ENABLE_REVIEWS", true),
			EnableServiceRatings:    getEnvBool("SERVICES_ENABLE_RATINGS", true),
			EnableServiceBookings:   getEnvBool("SERVICES_ENABLE_BOOKINGS", true),
			EnableServicePromotions: getEnvBool("SERVICES_ENABLE_PROMOTIONS", true),
			RateLimitCreate:         getEnvInt("SERVICES_RATE_LIMIT_CREATE", 10),
			RateLimitUpdate:         getEnvInt("SERVICES_RATE_LIMIT_UPDATE", 30),
			RateLimitSearch:         getEnvInt("SERVICES_RATE_LIMIT_SEARCH", 60),
		},
		Upload: Upload{
			MaxFileSize:    getEnvInt64("UPLOAD_MAX_FILE_SIZE", 10*1024*1024),
			AllowedTypes:   getEnvSlice("UPLOAD_ALLOWED_TYPES", []string{"image/jpeg", "image/png", "image/gif", "image/webp", "application/pdf"}, ","),
			ImageMaxWidth:  getEnvInt("UPLOAD_IMAGE_MAX_WIDTH", 1920),
			ImageMaxHeight: getEnvInt("UPLOAD_IMAGE_MAX_HEIGHT", 1080),
			StoragePath:    getEnv("UPLOAD_STORAGE_PATH", "./uploads"),
		},
		Email: Email{
			Enabled:   getEnvBool("EMAIL_ENABLED", false),
			Host:      getEnv("EMAIL_HOST", ""),
			Port:      getEnvInt("EMAIL_PORT", 587),
			Username:  getEnv("EMAIL_USERNAME", ""),
			Password:  getEnv("EMAIL_PASSWORD", ""),
			FromEmail: getEnv("EMAIL_FROM_EMAIL", "noreply@nawthtech.com"),
			FromName:  getEnv("EMAIL_FROM_NAME", "نوذ تك"),
		},
	}

	// تعيين القيم الافتراضية من ملف cors.go
	setCorsDefaults()

	// التحقق من صحة الإعدادات
	if err := validateConfig(); err != nil {
		logger.Stderr.Error("invalid configuration", logger.ErrAttr(err))
		os.Exit(1)
	}

	// تحليل متغيرات البيئة باستخدام env package
	if err := env.Parse(appConfig); err != nil {
		logger.Stderr.Error("failed to parse environment variables", logger.ErrAttr(err))
		os.Exit(1)
	}

	logger.Stdout.Info("تم تحميل الإعدادات بنجاح",
		"environment", appConfig.Environment,
		"port", appConfig.Port,
		"version", appConfig.Version,
	)

	return appConfig
}

// ========== دوال تعيين القيم الافتراضية ==========

func setCorsDefaults() {
	if len(appConfig.Cors.AllowedOrigins) == 0 {
		appConfig.Cors.AllowedOrigins = getAllowedOrigins()
	}
	if len(appConfig.Cors.AllowedMethods) == 0 {
		appConfig.Cors.AllowedMethods = []string{
			"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD",
		}
	}
	if len(appConfig.Cors.AllowedHeaders) == 0 {
		appConfig.Cors.AllowedHeaders = []string{
			"Content-Type", "Authorization", "X-Requested-With", "X-API-Key",
			"Accept", "Origin", "X-Client-Version", "X-Device-ID", "X-Platform",
			"X-Plausible-Token", "X-Matomo-Token", "X-Fathom-Key", "X-Request-ID",
			"Cache-Control", "X-CSRF-Token",
		}
	}
	if len(appConfig.Cors.ExposedHeaders) == 0 {
		appConfig.Cors.ExposedHeaders = []string{
			"X-Request-ID", "X-Response-Time", "X-API-Version",
			"X-RateLimit-Limit", "X-RateLimit-Remaining", "X-Total-Count",
			"Content-Length",
		}
	}
}

// ========== دوال التحقق من الصحة ==========

func validateConfig() error {
	if err := validateRequiredFields(); err != nil {
		return err
	}
	if err := validateServicesConfig(); err != nil {
		return err
	}
	if err := validateUploadConfig(); err != nil {
		return err
	}
	if err := validateAuthConfig(); err != nil {
		return err
	}
	return nil
}

func validateRequiredFields() error {
	if appConfig.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if appConfig.Auth.JWTSecret == "" || appConfig.Auth.JWTSecret == "default-jwt-secret-change-in-production" {
		return fmt.Errorf("JWT_SECRET is required and must be changed in production")
	}
	if appConfig.EncryptionKey == "" || appConfig.EncryptionKey == "default-encryption-key-change-in-production" {
		return fmt.Errorf("ENCRYPTION_KEY is required and must be changed in production")
	}
	return nil
}

func validateServicesConfig() error {
	if appConfig.Services.MinPrice < 0 {
		return fmt.Errorf("SERVICES_MIN_PRICE يجب أن يكون أكبر من أو يساوي الصفر")
	}

	if appConfig.Services.MaxPrice <= appConfig.Services.MinPrice {
		return fmt.Errorf("SERVICES_MAX_PRICE يجب أن يكون أكبر من SERVICES_MIN_PRICE")
	}

	if appConfig.Services.MinDuration < 1 {
		return fmt.Errorf("SERVICES_MIN_DURATION يجب أن يكون على الأقل 1")
	}

	if appConfig.Services.MaxDuration < appConfig.Services.MinDuration {
		return fmt.Errorf("SERVICES_MAX_DURATION يجب أن يكون أكبر من أو يساوي SERVICES_MIN_DURATION")
	}

	if appConfig.Services.MinTitleLength < 1 {
		return fmt.Errorf("SERVICES_MIN_TITLE_LENGTH يجب أن يكون على الأقل 1")
	}

	if appConfig.Services.MaxTitleLength < appConfig.Services.MinTitleLength {
		return fmt.Errorf("SERVICES_MAX_TITLE_LENGTH يجب أن يكون أكبر من أو يساوي SERVICES_MIN_TITLE_LENGTH")
	}

	return nil
}

func validateUploadConfig() error {
	if appConfig.Upload.MaxFileSize <= 0 {
		return fmt.Errorf("UPLOAD_MAX_FILE_SIZE يجب أن يكون أكبر من الصفر")
	}

	if appConfig.Upload.ImageMaxWidth <= 0 {
		return fmt.Errorf("UPLOAD_IMAGE_MAX_WIDTH يجب أن يكون أكبر من الصفر")
	}

	if appConfig.Upload.ImageMaxHeight <= 0 {
		return fmt.Errorf("UPLOAD_IMAGE_MAX_HEIGHT يجب أن يكون أكبر من الصفر")
	}

	return nil
}

func validateAuthConfig() error {
	if appConfig.Auth.JWTExpiration <= 0 {
		return fmt.Errorf("JWT_EXPIRATION يجب أن يكون أكبر من الصفر")
	}

	if appConfig.Auth.RefreshExpiration <= 0 {
		return fmt.Errorf("REFRESH_EXPIRATION يجب أن يكون أكبر من الصفر")
	}

	if appConfig.Auth.BCryptCost < 4 || appConfig.Auth.BCryptCost > 31 {
		return fmt.Errorf("BCRYPT_COST يجب أن يكون بين 4 و 31")
	}

	return nil
}

// ========== دوال مساعدة ==========

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string, separator string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, separator)
	}
	return defaultValue
}

// ========== دوال الوصول العامة ==========

// GetConfig الحصول على التكوين الحالي
func GetConfig() *Config {
	if appConfig == nil {
		return Load()
	}
	return appConfig
}

// IsDevelopment التحقق إذا كانت البيئة تطوير
func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.Environment) == "development"
}

// IsProduction التحقق إذا كانت البيئة إنتاج
func (c *Config) IsProduction() bool {
	return strings.ToLower(c.Environment) == "production"
}

// IsStaging التحقق إذا كانت البيئة تجريبية
func (c *Config) IsStaging() bool {
	return strings.ToLower(c.Environment) == "staging"
}

// IsCacheEnabled التحقق من تفعيل التخزين المؤقت
func (c *Config) IsCacheEnabled() bool {
	return c.Cache.Enabled
}

// GetRedisAddress الحصول على عنوان Redis
func (c *Config) GetRedisAddress() string {
	if c.Redis.URL != "" {
		return c.Redis.URL
	}
	return c.Redis.Host + ":" + c.Redis.Port
}

// GetDatabaseURL الحصول على رابط قاعدة البيانات
func (c *Config) GetDatabaseURL() string {
	return c.Database.URL
}

// GetJWTSecret الحصول على مفتاح JWT
func (c *Config) GetJWTSecret() string {
	return c.Auth.JWTSecret
}

// GetEncryptionKey الحصول على مفتاح التشفير
func (c *Config) GetEncryptionKey() string {
	return c.EncryptionKey
}

// GetPort الحصول على المنفذ
func (c *Config) GetPort() string {
	return c.Port
}

// GetVersion الحصول على نسخة التطبيق
func (c *Config) GetVersion() string {
	return c.Version
}

// GetEnvironment الحصول على البيئة
func (c *Config) GetEnvironment() string {
	return c.Environment
}

// GetCORSAllowedOrigins الحصول على النطاقات المسموح بها في CORS
func (c *Config) GetCORSAllowedOrigins() []string {
	return c.Cors.AllowedOrigins
}

// GetCacheConfig الحصول على تكوين التخزين المؤقت
func (c *Config) GetCacheConfig() map[string]interface{} {
	return map[string]interface{}{
		"enabled":     c.Cache.Enabled,
		"prefix":      c.Cache.Prefix,
		"default_ttl": c.Cache.DefaultTTL,
		"max_retries": c.Cache.MaxRetries,
		"redis_url":   c.GetRedisAddress(),
		"redis_db":    c.Redis.DB,
	}
}

// GetServicesConfig الحصول على تكوين الخدمات
func (c *Config) GetServicesConfig() map[string]interface{} {
	return map[string]interface{}{
		"max_services_per_user":      c.Services.MaxServicesPerUser,
		"max_active_services":        c.Services.MaxActiveServices,
		"default_pagination_limit":   c.Services.DefaultPaginationLimit,
		"max_pagination_limit":       c.Services.MaxPaginationLimit,
		"search_cache_ttl":           c.Services.SearchCacheTTL,
		"featured_cache_ttl":         c.Services.FeaturedCacheTTL,
		"max_images_per_service":     c.Services.MaxImagesPerService,
		"max_features_per_service":   c.Services.MaxFeaturesPerService,
		"max_tags_per_service":       c.Services.MaxTagsPerService,
		"min_title_length":           c.Services.MinTitleLength,
		"max_title_length":           c.Services.MaxTitleLength,
		"min_description_length":     c.Services.MinDescriptionLength,
		"max_description_length":     c.Services.MaxDescriptionLength,
		"min_price":                  c.Services.MinPrice,
		"max_price":                  c.Services.MaxPrice,
		"min_duration":               c.Services.MinDuration,
		"max_duration":               c.Services.MaxDuration,
		"auto_approve_services":      c.Services.AutoApproveServices,
		"allow_service_editing":      c.Services.AllowServiceEditing,
		"enable_service_reviews":     c.Services.EnableServiceReviews,
		"enable_service_ratings":     c.Services.EnableServiceRatings,
		"enable_service_bookings":    c.Services.EnableServiceBookings,
		"enable_service_promotions":  c.Services.EnableServicePromotions,
		"rate_limit_create":          c.Services.RateLimitCreate,
		"rate_limit_update":          c.Services.RateLimitUpdate,
		"rate_limit_search":          c.Services.RateLimitSearch,
	}
}

// GetUploadConfig الحصول على تكوين الرفع
func (c *Config) GetUploadConfig() map[string]interface{} {
	return map[string]interface{}{
		"max_file_size":    c.Upload.MaxFileSize,
		"allowed_types":    c.Upload.AllowedTypes,
		"image_max_width":  c.Upload.ImageMaxWidth,
		"image_max_height": c.Upload.ImageMaxHeight,
		"storage_path":     c.Upload.StoragePath,
	}
}

// GetEmailConfig الحصول على تكوين البريد
func (c *Config) GetEmailConfig() map[string]interface{} {
	return map[string]interface{}{
		"enabled":    c.Email.Enabled,
		"host":       c.Email.Host,
		"port":       c.Email.Port,
		"username":   c.Email.Username,
		"from_email": c.Email.FromEmail,
		"from_name":  c.Email.FromName,
	}
}

// GetAuthConfig الحصول على تكوين المصادقة
func (c *Config) GetAuthConfig() map[string]interface{} {
	return map[string]interface{}{
		"jwt_secret":          c.Auth.JWTSecret,
		"jwt_expiration":      c.Auth.JWTExpiration,
		"refresh_expiration":  c.Auth.RefreshExpiration,
		"bcrypt_cost":         c.Auth.BCryptCost,
	}
}

// GetDatabaseConfig الحصول على تكوين قاعدة البيانات
func (c *Config) GetDatabaseConfig() map[string]interface{} {
	return map[string]interface{}{
		"url":           c.Database.URL,
		"max_idle_conns": c.Database.MaxIdleConns,
		"max_open_conns": c.Database.MaxOpenConns,
		"max_lifetime":   c.Database.MaxLifetime,
	}
}

// GetCORSConfig الحصول على تكوين CORS
func (c *Config) GetCORSConfig() map[string]interface{} {
	return map[string]interface{}{
		"allowed_origins":   c.Cors.AllowedOrigins,
		"allowed_methods":   c.Cors.AllowedMethods,
		"allowed_headers":   c.Cors.AllowedHeaders,
		"exposed_headers":   c.Cors.ExposedHeaders,
		"allow_credentials": c.Cors.AllowCredentials,
		"max_age":           c.Cors.MaxAge,
	}
}