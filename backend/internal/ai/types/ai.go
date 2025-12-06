package types

import "time"  // ← أضف هذا

// AIRequest طلب AI عام
type AIRequest struct {
    Type        string                 `json:"type"` // text, image, video, analysis, translation
    Data        interface{}            `json:"data"`
    Options     map[string]interface{} `json:"options,omitempty"`
    Priority    int                    `json:"priority,omitempty"`
    UserContext *UserContext           `json:"user_context,omitempty"`
}

// AIResponse استجابة AI عامة
type AIResponse struct {
    Type        string                 `json:"type"`
    Success     bool                   `json:"success"`
    Data        interface{}            `json:"data,omitempty"`
    Error       string                 `json:"error,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
    ProcessingTime float64             `json:"processing_time"` // in seconds
    Provider    string                 `json:"provider,omitempty"`
}

// UserContext سياق المستخدم
type UserContext struct {
    UserID      string                 `json:"user_id"`
    UserTier    string                 `json:"user_tier"`
    SessionID   string                 `json:"session_id,omitempty"`
    Preferences map[string]interface{} `json:"preferences,omitempty"`
    History     []UsageRecord          `json:"history,omitempty"`
    Balance     float64                `json:"balance,omitempty"`
    Quota       map[string]int         `json:"quota,omitempty"`
}

// ModelInfo معلومات النموذج
type ModelInfo struct {
    Name        string                 `json:"name"`
    Provider    string                 `json:"provider"`
    Type        string                 `json:"type"`
    MaxTokens   int                    `json:"max_tokens"`
    Cost        float64                `json:"cost_per_token"`
    Languages   []string               `json:"languages"`
    Capabilities []string              `json:"capabilities"`
    Description string                 `json:"description,omitempty"`
    Version     string                 `json:"version,omitempty"`
    ReleasedAt  time.Time              `json:"released_at,omitempty"`
}

// ProviderCapabilities قدرات المزود
type ProviderCapabilities struct {
    Provider    string                 `json:"provider"`
    Text        bool                   `json:"text"`
    Image       bool                   `json:"image"`
    Video       bool                   `json:"video"`
    Analysis    bool                   `json:"analysis"`
    Translation bool                   `json:"translation"`
    Embedding   bool                   `json:"embedding"`
    Chat        bool                   `json:"chat"`
    Streaming   bool                   `json:"streaming"`
    Vision      bool                   `json:"vision"`
    MaxTokens   int                    `json:"max_tokens"`
    Models      []ModelInfo            `json:"models"`
}

// UsageSummary ملخص الاستخدام
type UsageSummary struct {
    Period      string                 `json:"period"` // daily, weekly, monthly
    StartDate   time.Time              `json:"start_date"`
    EndDate     time.Time              `json:"end_date"`
    TotalCost   float64                `json:"total_cost"`
    TotalRequests int64                `json:"total_requests"`
    Successful  int64                  `json:"successful"`
    Failed      int64                  `json:"failed"`
    ByProvider  map[string]UsageStats  `json:"by_provider"`
    ByType      map[string]UsageStats  `json:"by_type"`
}

// UsageStats إحصائيات الاستخدام
type UsageStats struct {
    Count       int64     `json:"count"`
    Cost        float64   `json:"cost"`
    SuccessRate float64   `json:"success_rate"`
    AvgLatency  float64   `json:"avg_latency"`
}

// FeatureFlag علم الميزة
type FeatureFlag struct {
    Name        string    `json:"name"`
    Enabled     bool      `json:"enabled"`
    Description string    `json:"description,omitempty"`
    EnabledFor  []string  `json:"enabled_for,omitempty"` // user tiers
    ExpiresAt   time.Time `json:"expires_at,omitempty"`
}

// SystemMetrics مقاييس النظام
type SystemMetrics struct {
    Timestamp   time.Time              `json:"timestamp"`
    CPUUsage    float64                `json:"cpu_usage"` // percentage
    MemoryUsage float64                `json:"memory_usage"` // MB
    DiskUsage   float64                `json:"disk_usage"` // percentage
    Goroutines  int                    `json:"goroutines"`
    ActiveConnections int              `json:"active_connections"`
    RequestRate float64                `json:"request_rate"` // requests per second
    ErrorRate   float64                `json:"error_rate"` // errors per second
    ProviderStatus map[string]bool     `json:"provider_status"`
}