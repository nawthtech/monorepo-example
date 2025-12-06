package types

import "time"

// ProviderInterface واجهة أساسية للمزودين
type ProviderInterface interface {
    // الخصائص الأساسية
    GetName() string
    GetType() string
    IsAvailable() bool
    GetCost() float64
    GetStats() *ProviderStats
    
    // الوظائف الأساسية
    GenerateText(req TextRequest) (*TextResponse, error)
    GenerateImage(req ImageRequest) (*ImageResponse, error)
    GenerateVideo(req VideoRequest) (*VideoResponse, error)
    AnalyzeText(req AnalysisRequest) (*AnalysisResponse, error)
    AnalyzeImage(req AnalysisRequest) (*AnalysisResponse, error)
    TranslateText(req TranslationRequest) (*TranslationResponse, error)
    
    // وظائف إضافية
    SupportsStreaming() bool
    SupportsEmbedding() bool
    GetMaxTokens() int
    GetSupportedLanguages() []string
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

// VideoServiceInterface واجهة لفصل التبعيات
type VideoServiceInterface interface {
    GenerateVideo(prompt string, options VideoOptions) (*VideoResponse, error)
    GetVideoStatus(operationID string) (*VideoResponse, error)
    CancelVideoGeneration(operationID string) error
    ListVideos() ([]VideoInfo, error)
    DownloadVideo(operationID string) ([]byte, error)
}

// TextServiceInterface واجهة لخدمة النصوص
type TextServiceInterface interface {
    GenerateText(prompt string, options TextOptions) (*TextResponse, error)
    AnalyzeText(text string, options AnalysisOptions) (*AnalysisResult, error)
    TranslateText(text, sourceLang, targetLang string, options TranslationOptions) (*TranslationResult, error)
    SummarizeText(text string, options SummaryOptions) (*SummaryResult, error)
}

// ImageServiceInterface واجهة لخدمة الصور
type ImageServiceInterface interface {
    GenerateImage(prompt string, options ImageOptions) (*ImageResponse, error)
    AnalyzeImage(imageData []byte, options AnalysisOptions) (*AnalysisResult, error)
    EditImage(imageData []byte, prompt string, options EditOptions) (*ImageResponse, error)
    UpscaleImage(imageData []byte, options UpscaleOptions) (*ImageResponse, error)
}

// VoiceServiceInterface واجهة لخدمة الصوت
type VoiceServiceInterface interface {
    GenerateSpeech(text string, options VoiceOptions) (*VoiceResponse, error)
    RecognizeSpeech(audioData []byte, options RecognitionOptions) (*RecognitionResult, error)
    ListVoices() ([]VoiceInfo, error)
}

// EmbeddingServiceInterface واجهة لخدمة التضمين
type EmbeddingServiceInterface interface {
    GenerateEmbedding(req EmbeddingRequest) (*EmbeddingResponse, error)
    CompareEmbeddings(embeddings1, embeddings2 []float64) (float64, error)
}

// ConversationServiceInterface واجهة لخدمة المحادثة
type ConversationServiceInterface interface {
    StartConversation(context ConversationContext) (string, error)
    AddMessage(conversationID string, message ConversationMessage) (*ConversationMessage, error)
    GetConversation(conversationID string) (*ConversationContext, error)
    EndConversation(conversationID string) error
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

// MonitoringInterface واجهة للمراقبة
type MonitoringInterface interface {
    RecordMetric(metric MonitorMetric) error
    GetHealthStatus() *HealthStatus
    GetMetrics(name string, start, end time.Time) ([]MonitorMetric, error)
    Alert(alertType, message string, severity int) error
}

// AuditInterface واجهة للتدقيق
type AuditInterface interface {
    Log(action AuditLog) error
    GetLogs(userID, action string, start, end time.Time) ([]AuditLog, error)
    ExportLogs(format string, start, end time.Time) ([]byte, error)
}

// ConfigManagerInterface واجهة لإدارة التكوين
type ConfigManagerInterface interface {
    LoadConfig(path string) (*AIConfig, error)
    SaveConfig(config *AIConfig, path string) error
    GetConfig() *AIConfig
    UpdateConfig(updates map[string]interface{}) error
    ReloadConfig() error
}

// RateLimiterInterface واجهة لمحدد المعدل
type RateLimiterInterface interface {
    Allow(userID string) bool
    Wait(userID string) error
    Reset(userID string) error
    GetStats(userID string) *RateLimitStats
}

// StorageInterface واجهة للتخزين
type StorageInterface interface {
    Save(key string, data []byte) error
    Load(key string) ([]byte, error)
    Delete(key string) error
    List(prefix string) ([]string, error)
    Exists(key string) bool
}

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

// RateLimitStats إحصائيات محدد المعدل
type RateLimitStats struct {
    UserID      string                 `json:"user_id"`
    Allowed     int64                  `json:"allowed"`
    Denied      int64                  `json:"denied"`
    Remaining   int                    `json:"remaining"`
    ResetAt     time.Time              `json:"reset_at"`
    Window      string                 `json:"window"` // second, minute, hour, day
}

// VoiceInfo معلومات الصوت
type VoiceInfo struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Language    string                 `json:"language"`
    Gender      string                 `json:"gender,omitempty"`
    Age         string                 `json:"age,omitempty"`
    Style       string                 `json:"style,omitempty"`
    Provider    string                 `json:"provider"`
    Available   bool                   `json:"available"`
    SampleURL   string                 `json:"sample_url,omitempty"`
}

// RecognitionOptions خيارات التعرف على الصوت
type RecognitionOptions struct {
    Language    string                 `json:"language"`
    Model       string                 `json:"model,omitempty"`
    Format      string                 `json:"format"` // wav, mp3, etc.
    SampleRate  int                    `json:"sample_rate,omitempty"`
    Channels    int                    `json:"channels,omitempty"`
    SpeakerDiarization bool            `json:"speaker_diarization,omitempty"`
    Punctuation bool                   `json:"punctuation,omitempty"`
    ProfanityFilter bool               `json:"profanity_filter,omitempty"`
}

// RecognitionResult نتيجة التعرف على الصوت
type RecognitionResult struct {
    Text        string                 `json:"text"`
    Confidence  float64                `json:"confidence"`
    Segments    []RecognitionSegment   `json:"segments,omitempty"`
    Language    string                 `json:"language,omitempty"`
    Cost        float64                `json:"cost"`
    ModelUsed   string                 `json:"model_used"`
    Duration    int                    `json:"duration"` // in milliseconds
}

// RecognitionSegment جزء التعرف
type RecognitionSegment struct {
    Text        string                 `json:"text"`
    Start       int                    `json:"start"` // in milliseconds
    End         int                    `json:"end"`   // in milliseconds
    Confidence  float64                `json:"confidence"`
    Speaker     string                 `json:"speaker,omitempty"`
}

// StreamingHandler واجهة لمعالج التدفق
type StreamingHandler interface {
    OnChunk(chunk StreamingResponse) error
    OnComplete(response *TextResponse) error
    OnError(err error) error
}