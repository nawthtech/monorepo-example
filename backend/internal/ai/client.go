package ai

import (
    "fmt"
    "log"
    "os"
    "sync"
)

// Client عميل AI مبسط
type Client struct {
    mu                sync.RWMutex
    providers         map[string]ProviderInterface
    multiProvider     *MultiProvider
}

// NewClient إنشاء عميل AI جديد
func NewClient() (*Client, error) {
    c := &Client{
        providers: make(map[string]ProviderInterface),
    }
    
    // إنشاء مزود متعدد
    mp, err := NewMultiProvider()
    if err != nil {
        return nil, fmt.Errorf("failed to create multi-provider: %w", err)
    }
    c.multiProvider = mp
    
    // تهيئة مزود Ollama (دائمًا متاح محليًا)
    ollama := NewOllamaProvider()
    c.providers["ollama"] = ollama
    log.Println("✅ Ollama provider initialized")
    
    // محاولة تهيئة مزود Gemini إذا كان هناك API key
    if apiKey := os.Getenv("GEMINI_API_KEY"); apiKey != "" {
        // Gemini سيتم إضافته لاحقاً
        log.Println("⚠️ Gemini API key found but provider not implemented yet")
    }
    
    log.Printf("🤖 AI Client initialized with %d providers", len(c.providers))
    
    return c, nil
}

// GenerateText توليد نص
func (c *Client) GenerateText(prompt, provider string) (string, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if provider == "" || provider == "auto" {
        // استخدام MultiProvider للاختيار التلقائي
        req := TextRequest{
            Prompt: prompt,
            Model:  "llama3.2:3b",
        }
        
        resp, err := c.multiProvider.GenerateText(req)
        if err != nil {
            return "", err
        }
        return resp.Text, nil
    }
    
    // استخدام مزود محدد
    p, exists := c.providers[provider]
    if !exists {
        return "", fmt.Errorf("provider %s not found", provider)
    }
    
    req := TextRequest{
        Prompt: prompt,
    }
    
    resp, err := p.GenerateText(req)
    if err != nil {
        return "", err
    }
    
    return resp.Text, nil
}

// GenerateImage توليد صورة
func (c *Client) GenerateImage(prompt, provider string) (string, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    req := ImageRequest{
        Prompt: prompt,
    }
    
    resp, err := c.multiProvider.GenerateImage(req)
    if err != nil {
        return "", err
    }
    
    return resp.URL, nil
}

// GenerateVideo توليد فيديو
func (c *Client) GenerateVideo(prompt, provider string) (string, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    req := VideoRequest{
        Prompt: prompt,
        Duration: 30, // 30 ثانية افتراضياً
    }
    
    resp, err := c.multiProvider.GenerateVideo(req)
    if err != nil {
        return "", err
    }
    
    return resp.URL, nil
}

// AnalyzeText تحليل نص
func (c *Client) AnalyzeText(text, provider string) (*AnalysisResponse, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    req := AnalysisRequest{
        Text: text,
    }
    
    return c.multiProvider.AnalyzeText(req)
}

// TranslateText ترجمة نص
func (c *Client) TranslateText(text, fromLang, toLang, provider string) (string, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    req := TranslationRequest{
        Text:     text,
        FromLang: fromLang,
        ToLang:   toLang,
    }
    
    resp, err := c.multiProvider.TranslateText(req)
    if err != nil {
        return "", err
    }
    
    return resp.TranslatedText, nil
}

// GetVideoStatus الحصول على حالة فيديو
func (c *Client) GetVideoStatus(operationID string) (*VideoResponse, error) {
    // هذه وظيفة تحتاج إلى VideoService
    // سنعود إليها لاحقاً
    return nil, fmt.Errorf("video service not available yet")
}

// GetAvailableProviders الحصول على المزودين المتاحين
func (c *Client) GetAvailableProviders() map[string][]string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    providers := make(map[string][]string)
    
    // إضافة Ollama كمزود نص
    providers["text"] = []string{"ollama", "auto"}
    
    // MultiProvider قد يكون لديه مزودين آخرين
    // سنضيفهم لاحقاً
    
    return providers
}

// IsProviderAvailable التحقق من توفر مزود
func (c *Client) IsProviderAvailable(providerType, providerName string) bool {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if providerName == "auto" {
        return c.multiProvider.IsAvailable()
    }
    
    if p, exists := c.providers[providerName]; exists {
        return p.IsAvailable()
    }
    
    return false
}

// Close إغلاق العميل وتحرير الموارد
func (c *Client) Close() error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    log.Println("Closing AI client...")
    
    // إغلاق جميع المزودين
    for name, provider := range c.providers {
        if closer, ok := provider.(interface{ Close() error }); ok {
            if err := closer.Close(); err != nil {
                log.Printf("Error closing provider %s: %v", name, err)
            }
        }
    }
    
    return nil
}