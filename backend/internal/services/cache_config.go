package services

import "time"

// DefaultCacheConfig التكوين الافتراضي لخدمة التخزين المؤقت
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		RedisURL:      os.Getenv("REDIS_URL"),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       0,
		Prefix:        "nawthtech:",
		DefaultTTL:    1 * time.Hour,
		MaxRetries:    3,
	}
}

// getEnv دالة مساعدة للحصول على متغيرات البيئة مع قيمة افتراضية
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}