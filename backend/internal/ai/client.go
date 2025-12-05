package ai

import (
    "fmt"
    "log"
    "os"
    "sync"
    
    "github.com/nawthtech/nawthtech/backend/internal/ai/providers"
    "github.com/nawthtech/nawthtech/backend/internal/ai/services"
)

type Client struct {
    mu                sync.RWMutex
    textProviders     map[string]providers.TextProvider
    imageProviders    map[string]providers.ImageProvider
    videoProviders    map[string]providers.VideoProvider
    
    // Services
    ContentService    *services.ContentService
    AnalysisService   *services.AnalysisService
    StrategyService   *services.StrategyService
    MediaService      *services.MediaService
    TranslationService *services.TranslationService
}

func NewClient() (*Client, error) {
    c := &Client{
        textProviders:  make(map[string]providers.TextProvider),
        imageProviders: make(map[string]providers.ImageProvider),
        videoProviders: make(map[string]providers.VideoProvider),
    }
    
    // تهيئة مزودي النصوص
    if err := c.initTextProviders(); err != nil {
        log.Printf("Warning: Text providers init failed: %v", err)
    }
    
    // تهيئة مزودي الصور
    if err := c.initImageProviders(); err != nil {
        log.Printf("Warning: Image providers init failed: %v", err)
    }
    
    // تهيئة مزودي الفيديو
    if err := c.initVideoProviders(); err != nil {
        log.Printf("Warning: Video providers init failed: %v", err)
    }
    
    // إنشاء الخدمات
    c.initServices()
    
    log.Printf("🤖 AI Client initialized with %d text, %d image, %d video providers",
        len(c.textProviders), len(c.imageProviders), len(c.videoProviders))
    
    return c, nil
}

func (c *Client) initTextProviders() error {
    // 1. Gemini (مجاني - 60 request/دقيقة)
    if apiKey := os.Getenv("GEMINI_API_KEY"); apiKey != "" {
        gemini, err := providers.NewGeminiProvider()
        if err == nil {
            c.textProviders["gemini"] = gemini
            log.Println("✅ Gemini text provider initialized")
        }
    }
    
    // 2. Ollama (محلي - مجاني بالكامل)
    ollama := providers.NewOllamaProvider()
    c.textProviders["ollama"] = ollama
    log.Println("✅ Ollama text provider initialized")
    
    // 3. Hugging Face (مجاني - 30k tokens/شهر)
    if token := os.Getenv("HUGGINGFACE_TOKEN"); token != "" {
        hf := providers.NewHuggingFaceProvider()
        c.textProviders["huggingface"] = hf
        log.Println("✅ Hugging Face text provider initialized")
    }
    
    if len(c.textProviders) == 0 {
        return fmt.Errorf("no text providers available")
    }
    
    return nil
}

func (c *Client) GetTextProvider(name string) providers.TextProvider {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if name == "" || name == "auto" {
        // اختيار تلقائي: Gemini أولاً، ثم Ollama
        if provider, ok := c.textProviders["gemini"]; ok {
            return provider
        }
        return c.textProviders["ollama"]
    }
    
    return c.textProviders[name]
}

// Client عميل AI محدث
type Client struct {
    providers    map[string]Provider
    videoService *VideoService
    mu           sync.RWMutex
}

// NewClient إنشاء عميل AI جديد مع دعم الفيديوهات
func NewClient() (*Client, error) {
    c := &Client{
        providers: make(map[string]Provider),
    }
    
    // إضافة مزود Gemini للنصوص والصور
    gemini, err := NewGeminiProvider()
    if err == nil {
        c.providers["gemini"] = gemini
        fmt.Println("✅ Gemini provider initialized")
    }
    
    // إضافة مزود خاص للفيديوهات
    videoProvider, err := NewVideoProvider()
    if err == nil {
        c.providers["video"] = videoProvider
        fmt.Println("✅ Video provider initialized")
    } else {
        fmt.Printf("⚠️ Video provider unavailable: %v\n", err)
    }
    
    // إنشاء VideoService
    c.videoService = NewVideoService(videoProvider)
    
    if len(c.providers) == 0 {
        return nil, fmt.Errorf("no AI providers available")
    }
    
    return c, nil
}

// GenerateVideo توليد فيديو
func (c *Client) GenerateVideo(req VideoRequest) (*VideoResponse, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if videoProvider, ok := c.providers["video"]; ok && videoProvider.IsAvailable() {
        return videoProvider.GenerateVideo(req)
    }
    
    return nil, fmt.Errorf("video generation not available")
}

// GetVideoStatus الحصول على حالة فيديو
func (c *Client) GetVideoStatus(operationID string) (*VideoResponse, error) {
    if c.videoService != nil {
        return c.videoService.GetStatus(operationID)
    }
    return nil, fmt.Errorf("video service not available")
}