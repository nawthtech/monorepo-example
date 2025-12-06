package types

import "time"

// ProviderInterface واجهة أساسية لجميع مزودي AI
type ProviderInterface interface {
    // العمليات الأساسية
    GenerateText(req TextRequest) (*TextResponse, error)
    GenerateImage(req ImageRequest) (*ImageResponse, error)
    GenerateVideo(req VideoRequest) (*VideoResponse, error)
    AnalyzeText(req AnalysisRequest) (*AnalysisResponse, error)
    AnalyzeImage(req AnalysisRequest) (*AnalysisResponse, error)
    TranslateText(req TranslationRequest) (*TranslationResponse, error)
    
    // معلومات المزود
    GetName() string
    GetType() string
    GetCost() float64
    IsAvailable() bool
    GetStats() *ProviderStats
    SupportsStreaming() bool
    SupportsEmbedding() bool
    GetMaxTokens() int
    GetSupportedLanguages() []string
}

// TextRequest طلب توليد نص
type TextRequest struct {
    Prompt        string                 `json:"prompt"`
    Model         string                 `json:"model,omitempty"`
    Temperature   float64                `json:"temperature,omitempty"`
    MaxTokens     int                    `json:"max_tokens,omitempty"`
    TopP          float64                `json:"top_p,omitempty"`
    FrequencyPenalty float64             `json:"frequency_penalty,omitempty"`
    PresencePenalty float64              `json:"presence_penalty,omitempty"`
    Stop          []string               `json:"stop,omitempty"`
    Stream        bool                   `json:"stream,omitempty"`
    User          string                 `json:"user,omitempty"`
    UserID        string                 `json:"user_id,omitempty"`
    UserTier      string                 `json:"user_tier,omitempty"`
    Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// TextResponse استجابة توليد نص
type TextResponse struct {
    Text         string    `json:"text"`
    Tokens       int       `json:"tokens"`
    Cost         float64   `json:"cost"`
    ModelUsed    string    `json:"model_used"`
    FinishReason string    `json:"finish_reason,omitempty"`
    CreatedAt    time.Time `json:"created_at"`
    RawResponse  string    `json:"raw_response,omitempty"`
}

// ImageRequest طلب توليد صورة
type ImageRequest struct {
    Prompt       string                 `json:"prompt"`
    Size         string                 `json:"size,omitempty"`
    Quality      string                 `json:"quality,omitempty"`
    Style        string                 `json:"style,omitempty"`
    N            int                    `json:"n,omitempty"`
    ResponseFormat string               `json:"response_format,omitempty"`
    User         string                 `json:"user,omitempty"`
    UserID       string                 `json:"user_id,omitempty"`
    UserTier     string                 `json:"user_tier,omitempty"`
    Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// ImageResponse استجابة توليد صورة
type ImageResponse struct {
    URL        string    `json:"url,omitempty"`
    ImageData  []byte    `json:"image_data,omitempty"`
    Size       string    `json:"size"`
    Format     string    `json:"format"`
    Cost       float64   `json:"cost"`
    ModelUsed  string    `json:"model_used"`
    CreatedAt  time.Time `json:"created_at"`
    Seed       int64     `json:"seed,omitempty"`
    Width      int       `json:"width,omitempty"`
    Height     int       `json:"height,omitempty"`
}

// VideoRequest طلب توليد فيديو
type VideoRequest struct {
    Prompt       string                 `json:"prompt"`
    Duration     int                    `json:"duration,omitempty"`
    Resolution   string                 `json:"resolution,omitempty"`
    FPS          int                    `json:"fps,omitempty"`
    Style        string                 `json:"style,omitempty"`
    User         string                 `json:"user,omitempty"`
    UserID       string                 `json:"user_id,omitempty"`
    UserTier     string                 `json:"user_tier,omitempty"`
    Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// VideoResponse استجابة توليد فيديو
type VideoResponse struct {
    URL         string                 `json:"url,omitempty"`
    VideoData   []byte                 `json:"video_data,omitempty"`
    Duration    int                    `json:"duration"`
    Resolution  string                 `json:"resolution"`
    FPS         int                    `json:"fps,omitempty"`
    Cost        float64                `json:"cost"`
    ModelUsed   string                 `json:"model_used"`
    CreatedAt   time.Time              `json:"created_at"`
    Status      string                 `json:"status,omitempty"`
    Progress    float64                `json:"progress,omitempty"`
    OperationID string                 `json:"operation_id,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// AnalysisRequest طلب تحليل
type AnalysisRequest struct {
    Text        string                 `json:"text,omitempty"`
    ImageData   []byte                 `json:"image_data,omitempty"`
    Prompt      string                 `json:"prompt,omitempty"`
    Model       string                 `json:"model,omitempty"`
    User        string                 `json:"user,omitempty"`
    UserID      string                 `json:"user_id,omitempty"`
    UserTier    string                 `json:"user_tier,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// AnalysisResponse استجابة تحليل
type AnalysisResponse struct {
    Result      string    `json:"result"`
    Confidence  float64   `json:"confidence,omitempty"`
    Cost        float64   `json:"cost"`
    Model       string    `json:"model"`
    CreatedAt   time.Time `json:"created_at"`
    Categories  []string  `json:"categories,omitempty"`
    Sentiment   string    `json:"sentiment,omitempty"`
    SentimentScore float64 `json:"sentiment_score,omitempty"`
}

// TranslationRequest طلب ترجمة
type TranslationRequest struct {
    Text        string                 `json:"text"`
    FromLang    string                 `json:"from_lang"`
    ToLang      string                 `json:"to_lang"`
    Model       string                 `json:"model,omitempty"`
    User        string                 `json:"user,omitempty"`
    UserID      string                 `json:"user_id,omitempty"`
    UserTier    string                 `json:"user_tier,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// TranslationResponse استجابة ترجمة
type TranslationResponse struct {
    TranslatedText string    `json:"translated_text"`
    Cost           float64   `json:"cost"`
    Model          string    `json:"model"`
    CreatedAt      time.Time `json:"created_at"`
    Accuracy       float64   `json:"accuracy,omitempty"`
    DetectedSourceLanguage string `json:"detected_source_language,omitempty"`
}

// ProviderStats إحصائيات المزود
type ProviderStats struct {
    Name        string    `json:"name"`
    Type        string    `json:"type"`
    IsAvailable bool      `json:"is_available"`
    Requests    int64     `json:"requests"`
    Successful  int64     `json:"successful"`
    Failed      int64     `json:"failed"`
    TotalCost   float64   `json:"total_cost"`
    AvgLatency  float64   `json:"avg_latency"`
    SuccessRate float64   `json:"success_rate"`
    LastUsed    time.Time `json:"last_used"`
    Features    []string  `json:"features,omitempty"`
    Models      []string  `json:"models,omitempty"`
}

// UsageRecord سجل استخدام
type UsageRecord struct {
    ID          string                 `json:"id"`
    UserID      string                 `json:"user_id"`
    UserTier    string                 `json:"user_tier"`
    Provider    string                 `json:"provider"`
    Type        string                 `json:"type"` // text, image, video, etc.
    Cost        float64                `json:"cost"`
    Quantity    int64                  `json:"quantity"`
    Latency     float64                `json:"latency"` // in milliseconds
    Success     bool                   `json:"success"`
    Timestamp   time.Time              `json:"timestamp"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
    Error       string                 `json:"error,omitempty"`
}

// CostConfig تكوين التكاليف
type CostConfig struct {
    Provider    string  `json:"provider"`
    Model       string  `json:"model"`
    Type        string  `json:"type"`
    CostPerToken float64 `json:"cost_per_token"`
    CostPerImage float64 `json:"cost_per_image"`
    CostPerSecond float64 `json:"cost_per_second"`
    MaxTokens   int     `json:"max_tokens"`
    Currency    string  `json:"currency"`
}

// ProviderConfig تكوين المزود
type ProviderConfig struct {
    Name        string                 `json:"name"`
    Type        string                 `json:"type"`
    APIKey      string                 `json:"api_key,omitempty"`
    BaseURL     string                 `json:"base_url,omitempty"`
    Timeout     int                    `json:"timeout"` // in seconds
    Enabled     bool                   `json:"enabled"`
    Priority    int                    `json:"priority"`
    Models      []string               `json:"models"`
    Capabilities []string              `json:"capabilities"` // text, image, video, etc.
    Config      map[string]interface{} `json:"config,omitempty"`
}

// ClientConfig تكوين العميل
type ClientConfig struct {
    DefaultProvider   string                 `json:"default_provider"`
    FallbackProviders []string               `json:"fallback_providers"`
    MaxRetries        int                    `json:"max_retries"`
    Timeout           int                    `json:"timeout"` // in seconds
    CacheEnabled      bool                   `json:"cache_enabled"`
    CacheTTL          int                    `json:"cache_ttl"` // in seconds
    RateLimit         int                    `json:"rate_limit"` // requests per minute
    Providers         []ProviderConfig       `json:"providers"`
    Costs             []CostConfig           `json:"costs"`
    Features          map[string]bool        `json:"features"`
}

// ErrorResponse استجابة خطأ
type ErrorResponse struct {
    Error       string                 `json:"error"`
    Code        string                 `json:"code,omitempty"`
    Message     string                 `json:"message"`
    Details     map[string]interface{} `json:"details,omitempty"`
    Timestamp   time.Time              `json:"timestamp"`
}

// StreamingChunk جزء من التدفق
type StreamingChunk struct {
    Content     string                 `json:"content"`
    Index       int                    `json:"index"`
    FinishReason string                `json:"finish_reason,omitempty"`
    Tokens      int                    `json:"tokens,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// EmbeddingRequest طلب تضمين
type EmbeddingRequest struct {
    Input       string                 `json:"input"`
    Model       string                 `json:"model,omitempty"`
    User        string                 `json:"user,omitempty"`
    UserID      string                 `json:"user_id,omitempty"`
    Dimensions  int                    `json:"dimensions,omitempty"`
}

// EmbeddingResponse استجابة تضمين
type EmbeddingResponse struct {
    Embedding   []float64              `json:"embedding"`
    Model       string                 `json:"model"`
    Cost        float64                `json:"cost"`
    Tokens      int                    `json:"tokens"`
    CreatedAt   time.Time              `json:"created_at"`
}

// ChatMessage رسالة محادثة
type ChatMessage struct {
    Role        string                 `json:"role"` // user, assistant, system
    Content     string                 `json:"content"`
    Name        string                 `json:"name,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ChatRequest طلب محادثة
type ChatRequest struct {
    Messages    []ChatMessage          `json:"messages"`
    Model       string                 `json:"model,omitempty"`
    Temperature float64                `json:"temperature,omitempty"`
    MaxTokens   int                    `json:"max_tokens,omitempty"`
    Stream      bool                   `json:"stream,omitempty"`
    User        string                 `json:"user,omitempty"`
    UserID      string                 `json:"user_id,omitempty"`
}

// ChatResponse استجابة محادثة
type ChatResponse struct {
    Message     ChatMessage            `json:"message"`
    Tokens      int                    `json:"tokens"`
    Cost        float64                `json:"cost"`
    ModelUsed   string                 `json:"model_used"`
    FinishReason string                `json:"finish_reason,omitempty"`
    CreatedAt   time.Time              `json:"created_at"`
}

// ProviderInfo معلومات المزود
type ProviderInfo struct {
    Name        string                 `json:"name"`
    Type        string                 `json:"type"`
    IsAvailable bool                   `json:"is_available"`
    Cost        float64                `json:"cost"`
    MaxTokens   int                    `json:"max_tokens"`
    Models      []string               `json:"models"`
    Languages   []string               `json:"languages"`
    Capabilities []string              `json:"capabilities"`
    Stats       *ProviderStats         `json:"stats,omitempty"`
}

// HealthStatus حالة الصحة
type HealthStatus struct {
    Status      string                 `json:"status"` // healthy, unhealthy, degraded
    Timestamp   time.Time              `json:"timestamp"`
    Providers   map[string]bool        `json:"providers"`
    Uptime      float64                `json:"uptime"` // in seconds
    MemoryUsage float64                `json:"memory_usage"` // in MB
    Goroutines  int                    `json:"goroutines"`
    Errors      []string               `json:"errors,omitempty"`
}