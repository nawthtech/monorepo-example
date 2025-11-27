package config

import (
	"os"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/logger"

	"github.com/caarlos0/env/v11"
)

type cors struct {
	AllowedOrigins []string `env:"ALLOWED_ORIGINS,required,notEmpty" envSeparator:","`
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

type Config struct {
	Environment    string
	DatabaseURL    string
	EncryptionKey  string
	JWTSecret      string
	Port           string
	Version        string
	Cors           *cors
	Redis          *redis
	Cache          *cache
}

var (
	Cors   = &cors{}
	Redis  = &redis{}
	Cache  = &cache{}
	AppConfig = &Config{}
)

func init() {
	// تحليل متغيرات البيئة
	toParse := []any{Cors, Redis, Cache}
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
	}

	// تعيين القيم الافتراضية للتخزين المؤقت
	if AppConfig.Cache.Prefix == "" {
		AppConfig.Cache.Prefix = "nawthtech:"
	}
	if AppConfig.Cache.DefaultTTL == 0 {
		AppConfig.Cache.DefaultTTL = 1 * time.Hour
	}
	if AppConfig.Cache.MaxRetries == 0 {
		AppConfig.Cache.MaxRetries = 3
	}

	// تعيين القيم الافتراضية لـ Redis
	if AppConfig.Redis.Host == "" {
		AppConfig.Redis.Host = "localhost"
	}
	if AppConfig.Redis.Port == "" {
		AppConfig.Redis.Port = "6379"
	}
	if AppConfig.Redis.DB == 0 {
		AppConfig.Redis.DB = 0
	}

	return AppConfig
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

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