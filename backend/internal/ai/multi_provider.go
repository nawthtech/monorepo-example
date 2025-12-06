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
    ProviderClaude      ProviderType = "claude"
    ProviderCohere      ProviderType = "cohere"
    ProviderLuma        ProviderType = "luma"
    ProviderRunway      ProviderType = "runway"
    ProviderPika        ProviderType = "pika"
)

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
    SelectProvider(userTier, promptType, providerType string, providers map[ProviderType]ProviderInterface) ProviderType
    GetFallbackChain(primary ProviderType, providerType string) []ProviderType
}

// TieredStrategy إستراتيجية التوجيه حسب خطة المستخدم
type TieredStrategy struct {
    providerConfigs map[ProviderType]ProviderConfig
}

// ProviderConfig تكوين المزود
type ProviderConfig struct {
    Priority    int
    CostPerToken float64
    MaxTokens   int
    Speed       float64 // 0-1
    Quality     float64 // 0-1
    Availability float64 // 0-1
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
    
    // 1. Gemini Provider
    if apiKey := getEnvWithFallback("GEMINI_API_KEY", ""); apiKey != "" {
        gemini, err := NewGeminiProvider()
        if err == nil {
            mp.providers[ProviderGemini] = gemini
            mp.textProviders["gemini"] = gemini
            mp.imageProviders["gemini"] = gemini
            log.Println("✅ Gemini provider initialized")
        } else {
            log.Printf("⚠️ Gemini provider failed: %v", err)
        }
    }
    
    // 2. Ollama Provider (دائمًا متاح محليًا)
    ollama := NewOllamaProvider()
    mp.providers[ProviderOllama] = ollama
    mp.textProviders["ollama"] = ollama
    log.Println("✅ Ollama provider initialized")
    
    // 3. Hugging Face Provider
    if token := getEnvWithFallback("HUGGINGFACE_TOKEN", ""); token != "" {
        hf := NewHuggingFaceProvider()
        mp.providers[ProviderHuggingFace] = hf
        mp.textProviders["huggingface"] = hf
        mp.imageProviders["huggingface"] = hf
        log.Println("✅ Hugging Face provider initialized")
    }
    
    // 4. OpenAI Provider (Claude يعتبر كبديل)
    if apiKey := getEnvWithFallback("OPENAI_API_KEY", ""); apiKey != "" {
        openai, err := NewOpenAIProvider()
        if err == nil {
            mp.providers[ProviderOpenAI] = openai
            mp.textProviders["openai"] = openai
            mp.imageProviders["openai"] = openai
            log.Println("✅ OpenAI provider initialized")
        }
    }
    
    // 5. Luma Video Provider
    if apiKey := getEnvWithFallback("LUMA_API_KEY", ""); apiKey != "" {
        luma, err := NewVideoProvider("luma")
        if err == nil {
            mp.providers[ProviderLuma] = luma
            mp.videoProviders["luma"] = luma
            log.Println("✅ Luma video provider initialized")
        }
    }
    
    // 6. Runway Video Provider
    if apiKey := getEnvWithFallback("RUNWAY_API_KEY", ""); apiKey != "" {
        runway, err := NewVideoProvider("runway")
        if err == nil {
            mp.providers[ProviderRunway] = runway
            mp.videoProviders["runway"] = runway
            log.Println("✅ Runway video provider initialized")
        }
    }
    
    // 7. Pika Video Provider
    if apiKey := getEnvWithFallback("PIKA_API_KEY", ""); apiKey != "" {
        pika, err := NewVideoProvider("pika")
        if err == nil {
            mp.providers[ProviderPika] = pika
            mp.videoProviders["pika"] = pika
            log.Println("✅ Pika video provider initialized")
        }
    }
    
    if len(mp.providers) == 0 {
        return fmt.Errorf("no AI providers available")
    }
    
    return nil
}

// GenerateText توليد نص
func (mp *MultiProvider) GenerateText(req TextRequest) (*TextResponse, error) {
    startTime := time.Now()
    
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    // تحديد المزود المناسب
    providerType := mp.strategy.SelectProvider(req.UserTier, "text", "text", mp.providers)
    provider, exists := mp.textProviders[string(providerType)]
    
    // إذا لم يكن المزود متوفراً، استخدم التسلسل الاحتياطي
    if !exists || !provider.IsAvailable() {
        fallbackChain := mp.strategy.GetFallbackChain(providerType, "text")
        for _, fbType := range fallbackChain {
            if fbProvider, fbExists := mp.textProviders[string(fbType)]; fbExists && fbProvider.IsAvailable() {
                provider = fbProvider
                mp.stats.FallbackCount[fbType]++
                log.Printf("🔄 Fallback from %s to %s", providerType, fbType)
                break
            }
        }
    }
    
    if provider == nil {
        return nil, fmt.Errorf("no available text provider")
    }
    
    // توليد النص
    response, err := provider.GenerateText(req)
    
    // تسجيل الاستخدام
    if mp.costManager != nil {
        record := &UsageRecord{
            UserID:     req.UserID,
            UserTier:   req.UserTier,
            Provider:   provider.GetName(),
            Type:       "text",
            Cost:       provider.GetCost() * float64(len(req.Prompt)/4), // تقدير تقريبي
            Quantity:   int64(len(req.Prompt)),
            Latency:    float64(time.Since(startTime).Milliseconds()),
            Success:    err == nil,
            Timestamp:  time.Now(),
            Metadata: map[string]interface{}{
                "model": req.Model,
                "tokens": len(req.Prompt),
            },
        }
        mp.costManager.RecordUsage(record)
    }
    
    // تحديث الإحصائيات
    mp.updateRequestStats(providerType, err == nil, provider.GetCost())
    
    return response, err
}

// GenerateImage توليد صورة
func (mp *MultiProvider) GenerateImage(req ImageRequest) (*ImageResponse, error) {
    startTime := time.Now()
    
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    // تحديد المزود المناسب
    providerType := mp.strategy.SelectProvider(req.UserTier, "image", "image", mp.providers)
    provider, exists := mp.imageProviders[string(providerType)]
    
    // التسلسل الاحتياطي
    if !exists || !provider.IsAvailable() {
        fallbackChain := mp.strategy.GetFallbackChain(providerType, "image")
        for _, fbType := range fallbackChain {
            if fbProvider, fbExists := mp.imageProviders[string(fbType)]; fbExists && fbProvider.IsAvailable() {
                provider = fbProvider
                mp.stats.FallbackCount[fbType]++
                break
            }
        }
    }
    
    if provider == nil {
        return nil, fmt.Errorf("no available image provider")
    }
    
    // توليد الصورة
    response, err := provider.GenerateImage(req)
    
    // تسجيل الاستخدام
    if mp.costManager != nil {
        record := &UsageRecord{
            UserID:     req.UserID,
            UserTier:   req.UserTier,
            Provider:   provider.GetName(),
            Type:       "image",
            Cost:       provider.GetCost(),
            Quantity:   1,
            Latency:    float64(time.Since(startTime).Milliseconds()),
            Success:    err == nil,
            Timestamp:  time.Now(),
        }
        mp.costManager.RecordUsage(record)
    }
    
    // تحديث الإحصائيات
    mp.updateRequestStats(providerType, err == nil, provider.GetCost())
    
    return response, err
}

// GenerateVideo توليد فيديو
func (mp *MultiProvider) GenerateVideo(req VideoRequest) (*VideoResponse, error) {
    startTime := time.Now()
    
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    // تحديد المزود المناسب
    providerType := mp.strategy.SelectProvider(req.UserTier, "video", "video", mp.providers)
    provider, exists := mp.videoProviders[string(providerType)]
    
    // التسلسل الاحتياطي
    if !exists || !provider.IsAvailable() {
        fallbackChain := mp.strategy.GetFallbackChain(providerType, "video")
        for _, fbType := range fallbackChain {
            if fbProvider, fbExists := mp.videoProviders[string(fbType)]; fbExists && fbProvider.IsAvailable() {
                provider = fbProvider
                mp.stats.FallbackCount[fbType]++
                break
            }
        }
    }
    
    if provider == nil {
        return nil, fmt.Errorf("no available video provider")
    }
    
    // توليد الفيديو
    response, err := provider.GenerateVideo(req)
    
    // تسجيل الاستخدام
    if mp.costManager != nil {
        record := &UsageRecord{
            UserID:     req.UserID,
            UserTier:   req.UserTier,
            Provider:   provider.GetName(),
            Type:       "video",
            Cost:       provider.GetCost(),
            Quantity:   1,
            Latency:    float64(time.Since(startTime).Milliseconds()),
            Success:    err == nil,
            Timestamp:  time.Now(),
            Metadata: map[string]interface{}{
                "duration": req.Duration,
                "aspect_ratio": req.AspectRatio,
            },
        }
        mp.costManager.RecordUsage(record)
    }
    
    // تحديث الإحصائيات
    mp.updateRequestStats(providerType, err == nil, provider.GetCost())
    
    return response, err
}

// GetTextProvider الحصول على مزود نصوص محدد
func (mp *MultiProvider) GetTextProvider(name string) ProviderInterface {
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    return mp.textProviders[name]
}

// GetImageProvider الحصول على مزود صور محدد
func (mp *MultiProvider) GetImageProvider(name string) ProviderInterface {
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    return mp.imageProviders[name]
}

// GetVideoProvider الحصول على مزود فيديوهات محدد
func (mp *MultiProvider) GetVideoProvider(name string) ProviderInterface {
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    return mp.videoProviders[name]
}

// GetAvailableProviders الحصول على المزودين المتاحين
func (mp *MultiProvider) GetAvailableProviders() map[string][]string {
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    result := make(map[string][]string)
    
    // مزودي النصوص
    textProviders := make([]string, 0, len(mp.textProviders))
    for name := range mp.textProviders {
        textProviders = append(textProviders, name)
    }
    result["text"] = textProviders
    
    // مزودي الصور
    imageProviders := make([]string, 0, len(mp.imageProviders))
    for name := range mp.imageProviders {
        imageProviders = append(imageProviders, name)
    }
    result["image"] = imageProviders
    
    // مزودي الفيديو
    videoProviders := make([]string, 0, len(mp.videoProviders))
    for name := range mp.videoProviders {
        videoProviders = append(videoProviders, name)
    }
    result["video"] = videoProviders
    
    return result
}

// SetRoutingStrategy تعيين إستراتيجية التوجيه
func (mp *MultiProvider) SetRoutingStrategy(strategy RoutingStrategy) {
    mp.mu.Lock()
    defer mp.mu.Unlock()
    
    mp.strategy = strategy
}

// GetStats الحصول على إحصائيات المزود المتعدد
func (mp *MultiProvider) GetStats() *MultiProviderStats {
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    return mp.stats
}

// GetProviderStats الحصول على إحصائيات مزود محدد
func (mp *MultiProvider) GetProviderStats(providerType ProviderType) (*ProviderStats, error) {
    mp.mu.RLock()
    defer mp.mu.RUnlock()
    
    if stats, exists := mp.stats.ProviderStats[providerType]; exists {
        return stats, nil
    }
    
    return nil, fmt.Errorf("provider stats not found: %s", providerType)
}

// updateRequestStats تحديث إحصائيات الطلب
func (mp *MultiProvider) updateRequestStats(providerType ProviderType, success bool, cost float64) {
    mp.mu.Lock()
    defer mp.mu.Unlock()
    
    mp.stats.TotalRequests++
    if success {
        mp.stats.Successful++
    } else {
        mp.stats.Failed++
    }
    mp.stats.TotalCost += cost
    
    // تحديث إحصائيات المزود المحدد
    if _, exists := mp.stats.ProviderStats[providerType]; !exists {
        mp.stats.ProviderStats[providerType] = &ProviderStats{
            Name: string(providerType),
            Type: getProviderType(providerType),
        }
    }
    
    stats := mp.stats.ProviderStats[providerType]
    stats.Requests++
    if success {
        stats.Successful++
    } else {
        stats.Failed++
    }
    stats.TotalCost += cost
    stats.LastUsed = time.Now()
    stats.SuccessRate = float64(stats.Successful) / float64(stats.Requests) * 100
}

// updateProviderStats تحديث إحصائيات جميع المزودين
func (mp *MultiProvider) updateProviderStats() {
    mp.mu.Lock()
    defer mp.mu.Unlock()
    
    for providerType, provider := range mp.providers {
        if _, exists := mp.stats.ProviderStats[providerType]; !exists {
            mp.stats.ProviderStats[providerType] = &ProviderStats{
                Name: string(providerType),
                Type: getProviderType(providerType),
            }
        }
        
        stats := mp.stats.ProviderStats[providerType]
        stats.IsAvailable = provider.IsAvailable()
        
        // الحصول على الإحصائيات من المزود نفسه إذا كان يدعمها
        if providerStats := provider.GetStats(); providerStats != nil {
            stats.Requests = providerStats.Requests
            stats.Successful = providerStats.Successful
            stats.Failed = providerStats.Failed
            stats.TotalCost = providerStats.TotalCost
            stats.AvgLatency = providerStats.AvgLatency
            stats.SuccessRate = providerStats.SuccessRate
        }
    }
}

// Helper Functions

func getEnvWithFallback(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}

func getProviderType(providerType ProviderType) string {
    switch providerType {
    case ProviderGemini, ProviderOpenAI, ProviderOllama, ProviderHuggingFace:
        return "text"
    case ProviderLuma, ProviderRunway, ProviderPika:
        return "video"
    default:
        return "mixed"
    }
}

func getUserTier(userID string) string {
    // هذه دالة مساعدة - يجب تنفيذها حسب نظام المستخدمين
    return "free" // مؤقت
}

func classifyPrompt(prompt string) string {
    prompt = strings.ToLower(prompt)
    
    keywords := map[string][]string{
        "analysis": {"analyze", "analysis", "compare", "evaluate", "assess"},
        "strategy": {"strategy", "plan", "marketing", "business", "growth"},
        "creative": {"creative", "story", "poem", "song", "script"},
        "technical": {"code", "algorithm", "technical", "explain", "how to"},
    }
    
    for category, words := range keywords {
        for _, word := range words {
            if strings.Contains(prompt, word) {
                return category
            }
        }
    }
    
    return "general"
}

// DefaultStrategy إستراتيجية افتراضية
type DefaultStrategy struct{}

func (s *DefaultStrategy) SelectProvider(userTier, promptType, providerType string, providers map[ProviderType]ProviderInterface) ProviderType {
    // منطق بسيط: استخدام Ollama للمستخدمين المجانيين، Gemini للمستخدمين المميزين
    switch userTier {
    case "free":
        if providerType == "video" {
            return ProviderLuma // Luma مجاني للمستخدمين المجانيين
        }
        return ProviderOllama // Ollama مجاني بالكامل
    case "premium":
        if providerType == "video" {
            return ProviderRunway // Runway أفضل جودة
        }
        return ProviderGemini // Gemini للمستخدمين المميزين
    case "enterprise":
        if providerType == "video" {
            return ProviderPika // Pika للشركات
        }
        return ProviderOpenAI // OpenAI للأمور المتقدمة
    default:
        return ProviderOllama
    }
}

func (s *DefaultStrategy) GetFallbackChain(primary ProviderType, providerType string) []ProviderType {
    chains := map[ProviderType][]ProviderType{
        ProviderGemini: {ProviderHuggingFace, ProviderOllama},
        ProviderOpenAI: {ProviderGemini, ProviderHuggingFace, ProviderOllama},
        ProviderOllama: {ProviderGemini, ProviderHuggingFace},
        ProviderLuma: {ProviderRunway, ProviderPika},
        ProviderRunway: {ProviderLuma, ProviderPika},
        ProviderPika: {ProviderLuma, ProviderRunway},
    }
    
    if chain, exists := chains[primary]; exists {
        return chain
    }
    
    // سلسلة احتياطية افتراضية
    if providerType == "text" {
        return []ProviderType{ProviderGemini, ProviderHuggingFace, ProviderOllama}
    } else if providerType == "video" {
        return []ProviderType{ProviderLuma, ProviderRunway, ProviderPika}
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