package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/models"
)

// CacheKeyGenerator مولد مفاتيح التخزين المؤقت
type CacheKeyGenerator struct {
	prefix string
}

// NewCacheKeyGenerator إنشاء مولد مفاتيح جديد
func NewCacheKeyGenerator(prefix string) *CacheKeyGenerator {
	return &CacheKeyGenerator{
		prefix: prefix,
	}
}

// Generate إنشاء مفتاح تخزين مؤقت
func (g *CacheKeyGenerator) Generate(pattern string, params ...interface{}) string {
	key := pattern
	if len(params) > 0 {
		key = fmt.Sprintf(pattern, params...)
	}
	return g.prefix + key
}

// CacheSerializer مُسلسل التخزين المؤقت
type CacheSerializer struct{}

// Serialize تسلسل القيمة للتخزين
func (s *CacheSerializer) Serialize(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		jsonData, err := json.Marshal(v)
		if err != nil {
			return "", fmt.Errorf("فشل في تسلسل القيمة: %v", err)
		}
		return string(jsonData), nil
	}
}

// Deserialize إعادة تسلسل القيمة من التخزين
func (s *CacheSerializer) Deserialize(value string, target interface{}) error {
	if target == nil {
		return nil
	}

	switch t := target.(type) {
	case *string:
		*t = value
		return nil
	default:
		return json.Unmarshal([]byte(value), target)
	}
}

// CacheValidator مدقق التخزين المؤقت
type CacheValidator struct{}

// ValidateKey التحقق من صحة المفتاح
func (v *CacheValidator) ValidateKey(key string) error {
	if key == "" {
		return fmt.Errorf("المفتاح لا يمكن أن يكون فارغاً")
	}
	
	if len(key) > 256 {
		return fmt.Errorf("طول المفتاح يتجاوز الحد المسموح (256 حرف)")
	}
	
	// منع أحرف خاصة قد تسبب مشاكل في Redis
	forbiddenChars := []string{" ", "\n", "\r", "\t", "\x00"}
	for _, char := range forbiddenChars {
		if strings.Contains(key, char) {
			return fmt.Errorf("المفتاح يحتوي على أحرف غير مسموحة")
		}
	}
	
	return nil
}

// ValidateTTL التحقق من صحة وقت الصلاحية
func (v *CacheValidator) ValidateTTL(ttl time.Duration) error {
	if ttl < 0 {
		return fmt.Errorf("وقت الصلاحية لا يمكن أن يكون سالباً")
	}
	
	if ttl > 365*24*time.Hour {
		return fmt.Errorf("وقت الصلاحية يتجاوز الحد المسموح (سنة واحدة)")
	}
	
	return nil
}

// CachePatternManager مدير أنماط التخزين المؤقت
type CachePatternManager struct {
	patterns map[string]models.CachePattern
}

// NewCachePatternManager إنشاء مدير أنماط جديد
func NewCachePatternManager() *CachePatternManager {
	return &CachePatternManager{
		patterns: make(map[string]models.CachePattern),
	}
}

// RegisterPattern تسجيل نمط تخزين مؤقت
func (m *CachePatternManager) RegisterPattern(name, pattern, description string, ttl time.Duration) {
	m.patterns[name] = models.CachePattern{
		Name:        name,
		Pattern:     pattern,
		Description: description,
		TTL:         ttl,
	}
}

// GetPattern الحصول على نمط محدد
func (m *CachePatternManager) GetPattern(name string) (models.CachePattern, bool) {
	pattern, exists := m.patterns[name]
	return pattern, exists
}

// GetAllPatterns الحصول على جميع الأنماط
func (m *CachePatternManager) GetAllPatterns() []models.CachePattern {
	patterns := make([]models.CachePattern, 0, len(m.patterns))
	for _, pattern := range m.patterns {
		patterns = append(patterns, pattern)
	}
	return patterns
}

// CacheMetricsCollector جامع مقاييس التخزين المؤقت
type CacheMetricsCollector struct {
	metrics *models.CacheMetrics
}

// NewCacheMetricsCollector إنشاء جامع مقاييس جديد
func NewCacheMetricsCollector() *CacheMetricsCollector {
	return &CacheMetricsCollector{
		metrics: &models.CacheMetrics{
			Timestamp: time.Now(),
		},
	}
}

// RecordOperation تسجيل عملية تخزين مؤقت
func (c *CacheMetricsCollector) RecordOperation(duration time.Duration) {
	c.metrics.OperationCount++
	c.metrics.ResponseTime = (c.metrics.ResponseTime*float64(c.metrics.OperationCount-1) + duration.Seconds()*1000) / float64(c.metrics.OperationCount)
}

// GetMetrics الحصول على المقاييس الحالية
func (c *CacheMetricsCollector) GetMetrics() *models.CacheMetrics {
	return c.metrics
}

// CacheKeyTemplates قوالب مفاتيح التخزين المؤقت
var CacheKeyTemplates = models.CacheKeyTemplate{
	UserProfile:   "user:%s",
	ServiceList:   "services:%s:%d", // services:category:page
	Session:       "session:%s",
	RateLimit:     "rate:%s:%s",     // rate:ip:endpoint
	SearchResults: "search:%s:%d",   // search:query:page
	APIResponse:   "api:%s:%s",      // api:endpoint:params_hash
}

// CacheHelper مساعد التخزين المؤقت
type CacheHelper struct {
	keyGenerator *CacheKeyGenerator
	serializer   *CacheSerializer
	validator    *CacheValidator
}

// NewCacheHelper إنشاء مساعد تخزين مؤقت جديد
func NewCacheHelper(prefix string) *CacheHelper {
	return &CacheHelper{
		keyGenerator: NewCacheKeyGenerator(prefix),
		serializer:   &CacheSerializer{},
		validator:    &CacheValidator{},
	}
}

// GenerateUserKey إنشاء مفتاح ملف تعريف المستخدم
func (h *CacheHelper) GenerateUserKey(userID string) string {
	return h.keyGenerator.Generate(CacheKeyTemplates.UserProfile, userID)
}

// GenerateServiceKey إنشاء مفتاح قائمة الخدمات
func (h *CacheHelper) GenerateServiceKey(category string, page int) string {
	return h.keyGenerator.Generate(CacheKeyTemplates.ServiceList, category, page)
}

// GenerateSessionKey إنشاء مفتاح الجلسة
func (h *CacheHelper) GenerateSessionKey(token string) string {
	return h.keyGenerator.Generate(CacheKeyTemplates.Session, token)
}

// GenerateRateLimitKey إنشاء مفتاح تحديد المعدل
func (h *CacheHelper) GenerateRateLimitKey(ip, endpoint string) string {
	return h.keyGenerator.Generate(CacheKeyTemplates.RateLimit, ip, endpoint)
}

// GenerateSearchKey إنشاء مفتاح نتائج البحث
func (h *CacheHelper) GenerateSearchKey(query string, page int) string {
	return h.keyGenerator.Generate(CacheKeyTemplates.SearchResults, query, page)
}

// GenerateAPIKey إنشاء مفتاح استجابة API
func (h *CacheHelper) GenerateAPIKey(endpoint, paramsHash string) string {
	return h.keyGenerator.Generate(CacheKeyTemplates.APIResponse, endpoint, paramsHash)
}

// ValidateAndSerialize التحقق من الصحة وتسلسل القيمة
func (h *CacheHelper) ValidateAndSerialize(key string, value interface{}, ttl time.Duration) (string, error) {
	if err := h.validator.ValidateKey(key); err != nil {
		return "", err
	}
	
	if err := h.validator.ValidateTTL(ttl); err != nil {
		return "", err
	}
	
	serialized, err := h.serializer.Serialize(value)
	if err != nil {
		return "", err
	}
	
	return serialized, nil
}

// CacheContext سياق التخزين المؤقت
type CacheContext struct {
	context.Context
	BypassCache bool
	CustomTTL   time.Duration
	Tags        []string
}

// WithCacheOptions إضافة خيارات التخزين المؤقت إلى السياق
func WithCacheOptions(ctx context.Context, bypassCache bool, customTTL time.Duration, tags ...string) context.Context {
	return &CacheContext{
		Context:     ctx,
		BypassCache: bypassCache,
		CustomTTL:   customTTL,
		Tags:        tags,
	}
}

// GetCacheOptionsFromContext الحصول على خيارات التخزين المؤقت من السياق
func GetCacheOptionsFromContext(ctx context.Context) (bypassCache bool, customTTL time.Duration, tags []string) {
	if cacheCtx, ok := ctx.(*CacheContext); ok {
		return cacheCtx.BypassCache, cacheCtx.CustomTTL, cacheCtx.Tags
	}
	return false, 0, nil
}