package ai

import (
    "fmt"
    "log"
    "strings"
    "sync"
    "time"
)

// ProviderType نوع المزود
type ProviderType string

// ثوابت أنواع المزودين
const (
    ProviderGemini      ProviderType = "gemini"
    ProviderOpenAI      ProviderType = "openai"
    ProviderOllama      ProviderType = "ollama"
    ProviderHuggingFace ProviderType = "huggingface"
    ProviderLuma        ProviderType = "luma"
    ProviderRunway      ProviderType = "runway"
    ProviderPika        ProviderType = "pika"
)

// MultiProviderStats إحصائيات المزود المتعدد
type MultiProviderStats struct {
    TotalRequests     int64
    Successful        int64
    Failed            int64
    TotalCost         float64
    ProviderStats     map[ProviderType]*ProviderStats
    LastRotation      map[string]time.Time
    FallbackCount     map[ProviderType]int64
}

// RoutingStrategy واجهة إستراتيجية التوجيه
type RoutingStrategy interface {
    SelectProvider(userTier, promptType, providerType string) ProviderType
    GetFallbackChain(primary ProviderType, providerType string) []ProviderType
}

// MPProviderConfig تكوين المزود للمزود المتعدد
type MPProviderConfig struct {
    Priority    int
    CostPerToken float64
    MaxTokens   int
    Speed       float64 // 0-1
    Quality     float64 // 0-1
    Availability float64 // 0-1
}

// MultiProvider مزود متعدد يدعم عدة مزودين AI
type MultiProvider struct {
    mu          sync.RWMutex
    providers   map[ProviderType]ProviderInterface
    textProviders map[string]ProviderInterface
    imageProviders map[string]ProviderInterface
    videoProviders map[string]ProviderInterface
    strategy    RoutingStrategy
    costManager *CostManager
    failover    *FailoverManager
    stats       *MultiProviderStats
}

// NewMultiProvider إنشاء مزود متعدد جديد
func NewMultiProvider() (*MultiProvider, error) {
    mp := &MultiProvider{
        providers: make(map[ProviderType]ProviderInterface),
        textProviders: make(map[string]ProviderInterface),
        imageProviders: make(map[string]ProviderInterface),
        videoProviders: make(map[string]ProviderInterface),
        strategy: &DefaultStrategy{},
        stats: &MultiProviderStats{
            ProviderStats: make(map[ProviderType]*ProviderStats),
            LastRotation:  make(map[string]time.Time),
            FallbackCount: make(map[ProviderType]int64),
        },
    }
    
    // تهيئة مدير التكاليف
    cm, err := NewCostManager()
    if err != nil {
        log.Printf("Warning: Failed to initialize cost manager: %v", err)
    }
    mp.costManager = cm
    
    // تهيئة مدير الفشل
    mp.failover = NewFailoverManager(mp)
    
    // تهيئة المزودين
    if err := mp.initProviders(); err != nil {
        return nil, fmt.Errorf("failed to initialize providers: %w", err)
    }
    
    // تهيئة الإحصائيات
    mp.updateProviderStats()
    
    log.Printf("🤖 MultiProvider initialized with %d total providers", len(mp.providers))
    
    return mp, nil
}

// initProviders تهيئة جميع المزودين
func (mp *MultiProvider) initProviders() error {
    mp.mu.Lock()
    defer mp.mu.Unlock()
    
    // 1. Ollama Provider (دائمًا متاح محليًا)
    ollama := NewOllamaProvider()
    mp.providers[ProviderOllama] = ollama
    mp.textProviders["ollama"] = ollama
    log.Println("✅ Ollama provider initialized")
    
    // المزودين الآخرين يحتاجون إلى API keys
    // يمكن إضافتهم لاحقًا
    
    if len(mp.providers) == 0 {
        return fmt.Errorf("no AI providers available")
    }
    
    return nil
}

// GenerateText توليد نص
func (mp *MultiProvider) GenerateText(req TextRequest) (*TextResponse, error) {
    // البحث عن مزود نصوص
    for _, provider := range mp.textProviders {
        if provider.IsAvailable() {
            return provider.GenerateText(req)
        }
    }
    return nil, fmt.Errorf("no available text provider")
}

// GenerateImage توليد صورة
func (mp *MultiProvider) GenerateImage(req ImageRequest) (*ImageResponse, error) {
    // البحث عن مزود صور
    for _, provider := range mp.imageProviders {
        if provider.IsAvailable() {
            return provider.GenerateImage(req)
        }
    }
    return nil, fmt.Errorf("no available image provider")
}

// GenerateVideo توليد فيديو
func (mp *MultiProvider) GenerateVideo(req VideoRequest) (*VideoResponse, error) {
    // البحث عن مزود فيديو
    for _, provider := range mp.videoProviders {
        if provider.IsAvailable() {
            return provider.GenerateVideo(req)
        }
    }
    return nil, fmt.Errorf("no available video provider")
}

// AnalyzeText تحليل نص
func (mp *MultiProvider) AnalyzeText(req AnalysisRequest) (*AnalysisResponse, error) {
    // البحث عن مزود يدعم تحليل النصوص
    for _, provider := range mp.textProviders {
        if provider.IsAvailable() {
            return provider.AnalyzeText(req)
        }
    }
    return nil, fmt.Errorf("no available text analysis provider")
}

// AnalyzeImage تحليل صورة
func (mp *MultiProvider) AnalyzeImage(req AnalysisRequest) (*AnalysisResponse, error) {
    // البحث عن مزود يدعم تحليل الصور
    for _, provider := range mp.imageProviders {
        if provider.IsAvailable() {
            return provider.AnalyzeImage(req)
        }
    }
    return nil, fmt.Errorf("no available image analysis provider")
}

// TranslateText ترجمة نص
func (mp *MultiProvider) TranslateText(req TranslationRequest) (*TranslationResponse, error) {
    // البحث عن مزود يدعم الترجمة
    for _, provider := range mp.textProviders {
        if provider.IsAvailable() {
            return provider.TranslateText(req)
        }
    }
    return nil, fmt.Errorf("no available translation provider")
}

// GetName اسم المزود
func (mp *MultiProvider) GetName() string {
    return "MultiProvider"
}

// GetType نوع المزود
func (mp *MultiProvider) GetType() string {
    return "multi"
}

// IsAvailable التحقق من التوفر
func (mp *MultiProvider) IsAvailable() bool {
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    return len(mp.providers) > 0
}

// GetCost التكلفة
func (mp *MultiProvider) GetCost() float64 {
    return 0.0 // سيتم حسابها بناءً على الاستخدام الفعلي
}

// GetStats الحصول على إحصائيات
func (mp *MultiProvider) GetStats() *ProviderStats {
    return &ProviderStats{
        Name:        mp.GetName(),
        Type:        mp.GetType(),
        IsAvailable: mp.IsAvailable(),
        Requests:    mp.stats.TotalRequests,
        Successful:  mp.stats.Successful,
        Failed:      mp.stats.Failed,
        TotalCost:   mp.stats.TotalCost,
        SuccessRate: 0,
    }
}

// SupportsStreaming يدعم التدفق
func (mp *MultiProvider) SupportsStreaming() bool {
    // Ollama يدعم التدفق
    if provider, ok := mp.providers[ProviderOllama]; ok {
        // تحقق إذا كان المزود يدعم التدفق
        if streamingProvider, ok := provider.(interface{ SupportsStreaming() bool }); ok {
            return streamingProvider.SupportsStreaming()
        }
    }
    return false
}

// SupportsEmbedding يدعم التضمين
func (mp *MultiProvider) SupportsEmbedding() bool {
    // Ollama يدعم التضمين
    if provider, ok := mp.providers[ProviderOllama]; ok {
        if embeddingProvider, ok := provider.(interface{ SupportsEmbedding() bool }); ok {
            return embeddingProvider.SupportsEmbedding()
        }
    }
    return false
}

// GetMaxTokens الحد الأقصى للرموز
func (mp *MultiProvider) GetMaxTokens() int {
    // العودة إلى القيمة الافتراضية
    return 2048
}

// GetSupportedLanguages اللغات المدعومة
func (mp *MultiProvider) GetSupportedLanguages() []string {
    return []string{"ar", "en", "es", "fr", "de"}
}

// DefaultStrategy إستراتيجية افتراضية
type DefaultStrategy struct{}

func (s *DefaultStrategy) SelectProvider(userTier, promptType, providerType string) ProviderType {
    // إستراتيجية بسيطة: استخدام Ollama للمستخدمين المجانيين
    if providerType == "text" || providerType == "" {
        return ProviderOllama
    }
    
    // للأنواع الأخرى، استخدام أول مزود متاح
    switch providerType {
    case "image":
        return ProviderOllama // Ollama قد يدعم الصور في المستقبل
    case "video":
        return ProviderLuma
    default:
        return ProviderOllama
    }
}

func (s *DefaultStrategy) GetFallbackChain(primary ProviderType, providerType string) []ProviderType {
    // سلسلة احتياطية بسيطة
    if providerType == "text" {
        return []ProviderType{ProviderOllama}
    }
    
    return []ProviderType{ProviderOllama}
}

// FailoverManager مدير التسلسل الاحتياطي
type FailoverManager struct {
    multiProvider *MultiProvider
    failoverCache map[string]ProviderType
}

func NewFailoverManager(mp *MultiProvider) *FailoverManager {
    return &FailoverManager{
        multiProvider: mp,
        failoverCache: make(map[string]ProviderType),
    }
}