package types

import "time"

// ============ الواجهات الأساسية ============

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

// ============ الواجهات المتخصصة ============

type TextProvider interface {
    GenerateText(req TextRequest) (*TextResponse, error)
    AnalyzeText(req AnalysisRequest) (*AnalysisResponse, error)
    TranslateText(req TranslationRequest) (*TranslationResponse, error)
    GetName() string
    SupportsStreaming() bool
    GetMaxTokens() int
}

type ImageProvider interface {
    GenerateImage(req ImageRequest) (*ImageResponse, error)
    AnalyzeImage(req AnalysisRequest) (*AnalysisResponse, error)
    GetName() string
}

type VideoProvider interface {
    GenerateVideo(req VideoRequest) (*VideoResponse, error)
    GetName() string
}

// MultiProviderInterface واجهة للمزود المتعدد
type MultiProviderInterface interface {
    ProviderInterface
    GetActiveProvider(providerType string) string
    SetActiveProvider(providerType, providerName string) error
    GetAvailableProviders(providerType string) []string
    GetProviderStats(providerType, providerName string) (*ProviderStats, error)
    RotateProvider(providerType string) error
    GetFallbackChain(providerType string) []string
}

// CostManagerInterface واجهة لإدارة التكاليف
type CostManagerInterface interface {
    RecordUsage(record *UsageRecord) error
    CanUseAI(userID, requestType string) (bool, string)
    GetUsageStatistics() map[string]interface{}
    GetUserQuotas(userID string) (map[string]*Quota, error)
    ResetUserQuotas(userID string) error
    SetLimits(monthly, daily float64)
    GetProviderStats(providerName string) (*ProviderStats, error)
}

// CacheManagerInterface واجهة لإدارة الذاكرة المؤقتة
type CacheManagerInterface interface {
    Get(key string) (interface{}, bool)
    Set(key string, value interface{}, ttl time.Duration)
    Delete(key string)
    Clear()
    Size() int
    GetStats() CacheStats
}

// AIClientInterface واجهة عميل AI
type AIClientInterface interface {
    // العمليات الأساسية
    GenerateText(prompt, provider string) (string, error)
    GenerateImage(prompt, provider string) (string, error)
    GenerateVideo(prompt, provider string) (string, error)
    
    // العمليات المتقدمة
    GenerateTextWithOptions(req TextRequest) (*TextResponse, error)
    GenerateImageWithOptions(req ImageRequest) (*ImageResponse, error)
    GenerateVideoWithOptions(req VideoRequest) (*VideoResponse, error)
    AnalyzeText(text, provider string) (*AnalysisResponse, error)
    AnalyzeTextWithOptions(req AnalysisRequest) (*AnalysisResponse, error)
    TranslateText(text, fromLang, toLang, provider string) (string, error)
    TranslateTextWithOptions(req TranslationRequest) (*TranslationResponse, error)
    AnalyzeImage(imageData []byte, prompt, provider string) (*AnalysisResponse, error)
    AnalyzeImageWithOptions(req AnalysisRequest) (*AnalysisResponse, error)
    
    // معلومات النظام
    GetVideoStatus(operationID string) (*VideoResponse, error)
    GetAvailableProviders() map[string][]string
    IsProviderAvailable(providerType, providerName string) bool
    GetProviderStats(providerName string) (*ProviderStats, error)
    GetUsageStatistics() map[string]interface{}
    
    // إدارة المزودين
    RegisterProvider(name string, provider ProviderInterface)
    RemoveProvider(name string)
    
    // إدارة الاتصال
    Close() error
}

// ============ هياكل الطلبات والاستجابات ============

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

// ============ الإحصائيات والتتبع ============

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

// UsageStats إحصائيات الاستخدام
type UsageStats struct {
    Count       int64     `json:"count"`
    Cost        float64   `json:"cost"`
    SuccessRate float64   `json:"success_rate"`
    AvgLatency  float64   `json:"avg_latency"`
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

// ============ أنواع إضافية ============

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

// ============ أنواع مساعدة ============

// Quota حد الاستخدام
type Quota struct {
    UserID      string                 `json:"user_id"`
    Type        string                 `json:"type"` // daily, monthly, etc.
    Limit       float64                `json:"limit"`
    Used        float64                `json:"used"`
    ResetAt     time.Time              `json:"reset_at"`
    Period      string                 `json:"period"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// CacheStats إحصائيات الذاكرة المؤقتة
type CacheStats struct {
    Hits        int64                  `json:"hits"`
    Misses      int64                  `json:"misses"`
    HitRate     float64                `json:"hit_rate"`
    Size        int                    `json:"size"`
    Items       int                    `json:"items"`
    MemoryUsage int64                  `json:"memory_usage"`
    Evictions   int64                  `json:"evictions,omitempty"`
}

// VideoOptions خيارات الفيديو
type VideoOptions struct {
    Duration     int                    `json:"duration"`
    Quality      string                 `json:"quality,omitempty"`
    AspectRatio  string                 `json:"aspect_ratio,omitempty"`
    Style        string                 `json:"style,omitempty"`
    OutputFormat string                 `json:"output_format,omitempty"`
    FPS          int                    `json:"fps,omitempty"`
    Resolution   string                 `json:"resolution,omitempty"`
    Audio        bool                   `json:"audio,omitempty"`
    Watermark    string                 `json:"watermark,omitempty"`
    Background   string                 `json:"background,omitempty"`
}

// TextOptions خيارات النصوص
type TextOptions struct {
    Model         string               `json:"model,omitempty"`
    Temperature   float64              `json:"temperature,omitempty"`
    MaxTokens     int                  `json:"max_tokens,omitempty"`
    TopP          float64              `json:"top_p,omitempty"`
    TopK          int                  `json:"top_k,omitempty"`
    FrequencyPenalty float64           `json:"frequency_penalty,omitempty"`
    PresencePenalty  float64           `json:"presence_penalty,omitempty"`
    StopSequences []string             `json:"stop_sequences,omitempty"`
    SystemPrompt  string               `json:"system_prompt,omitempty"`
    Language      string               `json:"language,omitempty"`
}

// ImageOptions خيارات الصور
type ImageOptions struct {
    Model         string               `json:"model,omitempty"`
    Size          string               `json:"size,omitempty"`
    Style         string               `json:"style,omitempty"`
    Quality       string               `json:"quality,omitempty"`
    AspectRatio   string               `json:"aspect_ratio,omitempty"`
    NumImages     int                  `json:"num_images,omitempty"`
    NegativePrompt string              `json:"negative_prompt,omitempty"`
    Seed          int64                `json:"seed,omitempty"`
}

// AnalysisOptions خيارات التحليل
type AnalysisOptions struct {
    Model         string               `json:"model,omitempty"`
    Task          string               `json:"task,omitempty"`
    Language      string               `json:"language,omitempty"`
    DetailLevel   string               `json:"detail_level,omitempty"`
    Format        string               `json:"format,omitempty"`
    MaxResults    int                  `json:"max_results,omitempty"`
}

// VideoInfo معلومات الفيديو
type VideoInfo struct {
    ID          string                 `json:"id"`
    Title       string                 `json:"title,omitempty"`
    Description string                 `json:"description,omitempty"`
    URL         string                 `json:"url,omitempty"`
    Duration    int                    `json:"duration"`
    Size        int64                  `json:"size"`
    Status      string                 `json:"status"`
    CreatedAt   time.Time              `json:"created_at"`
    UpdatedAt   time.Time              `json:"updated_at"`
    Cost        float64                `json:"cost"`
    Provider    string                 `json:"provider"`
}