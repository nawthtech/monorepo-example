package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"github.com/nawthtech/nawthtech/backend/internal/services"
)

// ========== هياكل التكوين ==========

type cors struct {
	AllowedOrigins []string `env:"ALLOWED_ORIGINS" envSeparator:","`
}

type redis struct {
	URL      string `env:"REDIS_URL"`
	Host     string `env:"REDIS_HOST"`
	Port     string `env:"REDIS_PORT"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB"`
}

type cache struct {
	Enabled    bool          `env:"CACHE_ENABLED"`
	Prefix     string        `env:"CACHE_PREFIX"`
	DefaultTTL time.Duration `env:"CACHE_DEFAULT_TTL"`
	MaxRetries int           `env:"CACHE_MAX_RETRIES"`
}

type servicesConfig struct {
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

type upload struct {
	MaxFileSize    int64    `env:"UPLOAD_MAX_FILE_SIZE"`
	AllowedTypes   []string `env:"UPLOAD_ALLOWED_TYPES" envSeparator:","`
	ImageMaxWidth  int      `env:"UPLOAD_IMAGE_MAX_WIDTH"`
	ImageMaxHeight int      `env:"UPLOAD_IMAGE_MAX_HEIGHT"`
	StoragePath    string   `env:"UPLOAD_STORAGE_PATH"`
}

type email struct {
	Enabled   bool   `env:"EMAIL_ENABLED"`
	Host      string `env:"EMAIL_HOST"`
	Port      int    `env:"EMAIL_PORT"`
	Username  string `env:"EMAIL_USERNAME"`
	Password  string `env:"EMAIL_PASSWORD"`
	FromEmail string `env:"EMAIL_FROM_EMAIL"`
	FromName  string `env:"EMAIL_FROM_NAME"`
}

type Config struct {
	Environment   string         `env:"ENVIRONMENT"`
	DatabaseURL   string         `env:"DATABASE_URL"`
	EncryptionKey string         `env:"ENCRYPTION_KEY"`
	JWTSecret     string         `env:"JWT_SECRET"`
	Port          string         `env:"PORT"`
	Version       string         `env:"APP_VERSION"`
	Cors          *cors          `envPrefix:"CORS_"`
	Redis         *redis         `envPrefix:"REDIS_"`
	Cache         *cache         `envPrefix:"CACHE_"`
	Services      *servicesConfig `envPrefix:"SERVICES_"`
	Upload        *upload        `envPrefix:"UPLOAD_"`
	Email         *email         `envPrefix:"EMAIL_"`
}

// ========== متغيرات عامة ==========

var (
	Cors     = &cors{}
	Redis    = &redis{}
	Cache    = &cache{}
	Services = &servicesConfig{}
	Upload   = &upload{}
	Email    = &email{}
	AppConfig = &Config{}
)

// ========== التهيئة ==========

func init() {
	// تحليل متغيرات البيئة
	toParse := []any{Cors, Redis, Cache, Services, Upload, Email}
	errors := []error{}

	for _, v := range toParse {
		if err := env.Parse(v); err != nil {
			if er, ok := err.(env.AggregateError); ok {
				errors = append(errors, er.Errors...)
				continue
			}
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		logger.Stderr.Error("errors found while parsing environment variables", logger.ErrorsAttr(errors...))
		os.Exit(1)
	}
}

// Load تحميل الإعدادات
func Load() *Config {
	AppConfig = &Config{
		Environment:   getEnv("ENVIRONMENT", "development"),
		DatabaseURL:   getEnv("DATABASE_URL", ""),
		EncryptionKey: getEnv("ENCRYPTION_KEY", ""),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		Port:          getEnv("PORT", "3000"),
		Version:       getEnv("APP_VERSION", "1.0.0"),
		Cors:          Cors,
		Redis:         Redis,
		Cache:         Cache,
		Services:      Services,
		Upload:        Upload,
		Email:         Email,
	}

	// تعيين القيم الافتراضية
	setDefaults()

	// التحقق من صحة الإعدادات
	if err := validateConfig(); err != nil {
		logger.Stderr.Error("invalid configuration", logger.ErrAttr(err))
		os.Exit(1)
	}

	return AppConfig
}

// ========== دوال تعيين القيم الافتراضية ==========

func setDefaults() {
	setCacheDefaults()
	setServicesDefaults()
	setUploadDefaults()
	setEmailDefaults()
	setRedisDefaults()
	setCorsDefaults()
}

func setCacheDefaults() {
	if AppConfig.Cache.Prefix == "" {
		AppConfig.Cache.Prefix = "nawthtech:"
	}
	if AppConfig.Cache.DefaultTTL == 0 {
		AppConfig.Cache.DefaultTTL = 1 * time.Hour
	}
	if AppConfig.Cache.MaxRetries == 0 {
		AppConfig.Cache.MaxRetries = 3
	}
}

func setServicesDefaults() {
	if AppConfig.Services.MaxServicesPerUser == 0 {
		AppConfig.Services.MaxServicesPerUser = 50
	}
	if AppConfig.Services.MaxActiveServices == 0 {
		AppConfig.Services.MaxActiveServices = 20
	}
	if AppConfig.Services.DefaultPaginationLimit == 0 {
		AppConfig.Services.DefaultPaginationLimit = 20
	}
	if AppConfig.Services.MaxPaginationLimit == 0 {
		AppConfig.Services.MaxPaginationLimit = 100
	}
	if AppConfig.Services.SearchCacheTTL == 0 {
		AppConfig.Services.SearchCacheTTL = 5 * time.Minute
	}
	if AppConfig.Services.FeaturedCacheTTL == 0 {
		AppConfig.Services.FeaturedCacheTTL = 30 * time.Minute
	}
	if AppConfig.Services.MaxImagesPerService == 0 {
		AppConfig.Services.MaxImagesPerService = 10
	}
	if AppConfig.Services.MaxFeaturesPerService == 0 {
		AppConfig.Services.MaxFeaturesPerService = 20
	}
	if AppConfig.Services.MaxTagsPerService == 0 {
		AppConfig.Services.MaxTagsPerService = 15
	}
	if AppConfig.Services.MinTitleLength == 0 {
		AppConfig.Services.MinTitleLength = 3
	}
	if AppConfig.Services.MaxTitleLength == 0 {
		AppConfig.Services.MaxTitleLength = 200
	}
	if AppConfig.Services.MinDescriptionLength == 0 {
		AppConfig.Services.MinDescriptionLength = 10
	}
	if AppConfig.Services.MaxDescriptionLength == 0 {
		AppConfig.Services.MaxDescriptionLength = 2000
	}
	if AppConfig.Services.MinPrice == 0 {
		AppConfig.Services.MinPrice = 0
	}
	if AppConfig.Services.MaxPrice == 0 {
		AppConfig.Services.MaxPrice = 1000000
	}
	if AppConfig.Services.MinDuration == 0 {
		AppConfig.Services.MinDuration = 1
	}
	if AppConfig.Services.MaxDuration == 0 {
		AppConfig.Services.MaxDuration = 365
	}
	if AppConfig.Services.RateLimitCreate == 0 {
		AppConfig.Services.RateLimitCreate = 10 // 10 خدمات في الدقيقة
	}
	if AppConfig.Services.RateLimitUpdate == 0 {
		AppConfig.Services.RateLimitUpdate = 30 // 30 تحديث في الدقيقة
	}
	if AppConfig.Services.RateLimitSearch == 0 {
		AppConfig.Services.RateLimitSearch = 60 // 60 بحث في الدقيقة
	}
}

func setUploadDefaults() {
	if AppConfig.Upload.MaxFileSize == 0 {
		AppConfig.Upload.MaxFileSize = 10 * 1024 * 1024 // 10MB
	}
	if len(AppConfig.Upload.AllowedTypes) == 0 {
		AppConfig.Upload.AllowedTypes = []string{
			"image/jpeg", "image/png", "image/gif", "image/webp",
			"application/pdf", "application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		}
	}
	if AppConfig.Upload.ImageMaxWidth == 0 {
		AppConfig.Upload.ImageMaxWidth = 1920
	}
	if AppConfig.Upload.ImageMaxHeight == 0 {
		AppConfig.Upload.ImageMaxHeight = 1080
	}
	if AppConfig.Upload.StoragePath == "" {
		AppConfig.Upload.StoragePath = "./uploads"
	}
}

func setEmailDefaults() {
	if AppConfig.Email.FromEmail == "" {
		AppConfig.Email.FromEmail = "noreply@nawthtech.com"
	}
	if AppConfig.Email.FromName == "" {
		AppConfig.Email.FromName = "نوذ تك"
	}
	if AppConfig.Email.Port == 0 {
		AppConfig.Email.Port = 587
	}
}

func setRedisDefaults() {
	if AppConfig.Redis.Host == "" {
		AppConfig.Redis.Host = "localhost"
	}
	if AppConfig.Redis.Port == "" {
		AppConfig.Redis.Port = "6379"
	}
	if AppConfig.Redis.DB == 0 {
		AppConfig.Redis.DB = 0
	}
}

func setCorsDefaults() {
	if len(AppConfig.Cors.AllowedOrigins) == 0 {
		AppConfig.Cors.AllowedOrigins = []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"https://nawthtech.com",
			"https://www.nawthtech.com",
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
	return nil
}

func validateRequiredFields() error {
	if AppConfig.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if AppConfig.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if AppConfig.EncryptionKey == "" {
		return fmt.Errorf("ENCRYPTION_KEY is required")
	}
	return nil
}

func validateServicesConfig() error {
	if AppConfig.Services.MinPrice < 0 {
		return fmt.Errorf("SERVICES_MIN_PRICE يجب أن يكون أكبر من أو يساوي الصفر")
	}

	if AppConfig.Services.MaxPrice <= AppConfig.Services.MinPrice {
		return fmt.Errorf("SERVICES_MAX_PRICE يجب أن يكون أكبر من SERVICES_MIN_PRICE")
	}

	if AppConfig.Services.MinDuration < 1 {
		return fmt.Errorf("SERVICES_MIN_DURATION يجب أن يكون على الأقل 1")
	}

	if AppConfig.Services.MaxDuration < AppConfig.Services.MinDuration {
		return fmt.Errorf("SERVICES_MAX_DURATION يجب أن يكون أكبر من أو يساوي SERVICES_MIN_DURATION")
	}

	if AppConfig.Services.MinTitleLength < 1 {
		return fmt.Errorf("SERVICES_MIN_TITLE_LENGTH يجب أن يكون على الأقل 1")
	}

	if AppConfig.Services.MaxTitleLength < AppConfig.Services.MinTitleLength {
		return fmt.Errorf("SERVICES_MAX_TITLE_LENGTH يجب أن يكون أكبر من أو يساوي SERVICES_MIN_TITLE_LENGTH")
	}

	return nil
}

func validateUploadConfig() error {
	if AppConfig.Upload.MaxFileSize <= 0 {
		return fmt.Errorf("UPLOAD_MAX_FILE_SIZE يجب أن يكون أكبر من الصفر")
	}

	if AppConfig.Upload.ImageMaxWidth <= 0 {
		return fmt.Errorf("UPLOAD_IMAGE_MAX_WIDTH يجب أن يكون أكبر من الصفر")
	}

	if AppConfig.Upload.ImageMaxHeight <= 0 {
		return fmt.Errorf("UPLOAD_IMAGE_MAX_HEIGHT يجب أن يكون أكبر من الصفر")
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

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// ========== دوال الوصول العامة ==========

// GetCacheConfig الحصول على تكوين التخزين المؤقت
func (c *Config) GetCacheConfig() services.CacheConfig {
	return services.CacheConfig{
		RedisURL:      c.Redis.URL,
		RedisHost:     c.Redis.Host,
		RedisPort:     c.Redis.Port,
		RedisPassword: c.Redis.Password,
		RedisDB:       c.Redis.DB,
		Prefix:        c.Cache.Prefix,
		DefaultTTL:    c.Cache.DefaultTTL,
		MaxRetries:    c.Cache.MaxRetries,
	}
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

// GetCORSAllowedOrigins الحصول على النطاقات المسموح بها في CORS
func (c *Config) GetCORSAllowedOrigins() []string {
	return c.Cors.AllowedOrigins
}

// GetVersion الحصول على نسخة التطبيق
func (c *Config) GetVersion() string {
	return c.Version
}

// GetEnvironment الحصول على البيئة
func (c *Config) GetEnvironment() string {
	return c.Environment
}

// GetPort الحصول على المنفذ
func (c *Config) GetPort() string {
	return c.Port
}

// GetJWTSecret الحصول على مفتاح JWT
func (c *Config) GetJWTSecret() string {
	return c.JWTSecret
}

// GetDatabaseURL الحصول على رابط قاعدة البيانات
func (c *Config) GetDatabaseURL() string {
	return c.DatabaseURL
}

// GetEncryptionKey الحصول على مفتاح التشفير
func (c *Config) GetEncryptionKey() string {
	return c.EncryptionKey
}