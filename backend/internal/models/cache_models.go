package models

import "time"

// CacheStats إحصائيات التخزين المؤقت
type CacheStats struct {
	Status          string `json:"status"`
	KeysCount       int64  `json:"keysCount"`
	UsedMemory      string `json:"usedMemory"`
	ConnectedClients int64  `json:"connectedClients"`
	Hits           int64  `json:"hits"`
	Misses         int64  `json:"misses"`
	HitRate        int    `json:"hitRate"`
	Uptime         int64  `json:"uptime"`
	Environment    string `json:"environment"`
	RetryCount     int    `json:"retryCount"`
}

// CacheHealthStatus حالة صحة التخزين المؤقت
type CacheHealthStatus struct {
	Status      string      `json:"status"`
	Message     string      `json:"message"`
	Error       string      `json:"error,omitempty"`
	Environment string      `json:"environment"`
	RetryCount  int         `json:"retryCount"`
	Stats       *CacheStats `json:"stats,omitempty"`
}

// CacheOperationResult نتيجة عملية التخزين المؤقت
type CacheOperationResult struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Key     string      `json:"key,omitempty"`
	Value   interface{} `json:"value,omitempty"`
	TTL     float64     `json:"ttl,omitempty"` // بالثواني
}

// CacheListResult نتيجة عملية القائمة
type CacheListResult struct {
	Key    string        `json:"key"`
	Values []interface{} `json:"values"`
	Count  int           `json:"count"`
	Start  int64         `json:"start"`
	Stop   int64         `json:"stop"`
}

// CacheHashResult نتيجة عملية الهاش
type CacheHashResult struct {
	Key   string                 `json:"key"`
	Hash  map[string]interface{} `json:"hash"`
	Size  int                    `json:"size"`
	Field string                 `json:"field,omitempty"`
}

// CacheKeysResult نتيجة عملية البحث عن المفاتيح
type CacheKeysResult struct {
	Pattern string   `json:"pattern"`
	Keys    []string `json:"keys"`
	Count   int      `json:"count"`
}

// CacheFlushResult نتيجة عملية المسح
type CacheFlushResult struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	Pattern      string `json:"pattern,omitempty"`
	DeletedCount int64  `json:"deletedCount,omitempty"`
}

// CacheIncrementResult نتيجة عملية الزيادة
type CacheIncrementResult struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

// CacheTTLResult نتيجة عملية وقت الصلاحية
type CacheTTLResult struct {
	Key string  `json:"key"`
	TTL float64 `json:"ttl"` // بالثواني
}

// CacheExistsResult نتيجة عملية التحقق من الوجود
type CacheExistsResult struct {
	Key    string `json:"key"`
	Exists bool   `json:"exists"`
}

// CacheConfig تكوين التخزين المؤقت
type CacheConfig struct {
	Enabled       bool          `json:"enabled"`
	RedisURL      string        `json:"redisUrl,omitempty"`
	RedisHost     string        `json:"redisHost,omitempty"`
	RedisPort     string        `json:"redisPort,omitempty"`
	Prefix        string        `json:"prefix"`
	DefaultTTL    time.Duration `json:"defaultTTL"`
	MaxRetries    int           `json:"maxRetries"`
	Environment   string        `json:"environment"`
}

// CacheMetrics مقاييس التخزين المؤقت
type CacheMetrics struct {
	Timestamp      time.Time `json:"timestamp"`
	OperationCount int64     `json:"operationCount"`
	HitRate        float64   `json:"hitRate"`
	MemoryUsage    string    `json:"memoryUsage"`
	ActiveKeys     int64     `json:"activeKeys"`
	ResponseTime   float64   `json:"responseTime"` // بالمللي ثانية
}

// CachePattern أنماط التخزين المؤقت
type CachePattern struct {
	Name        string `json:"name"`
	Pattern     string `json:"pattern"`
	Description string `json:"description"`
	TTL         time.Duration `json:"ttl"`
}

// CacheKeyTemplate قوالب مفاتيح التخزين المؤقت
type CacheKeyTemplate struct {
	UserProfile    string `json:"userProfile"`    // user:{id}
	ServiceList    string `json:"serviceList"`    // services:{category}
	Session        string `json:"session"`        // session:{token}
	RateLimit      string `json:"rateLimit"`      // rate:{ip}:{endpoint}
	SearchResults  string `json:"searchResults"`  // search:{query}:{page}
	APIResponse    string `json:"apiResponse"`    // api:{endpoint}:{params}
}