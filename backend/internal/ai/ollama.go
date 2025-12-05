package ai

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
    "time"
)

// OllamaProvider مزود Ollama المحلي
type OllamaProvider struct {
    baseURL    string
    httpClient *http.Client
    models     []string
}

// NewOllamaProvider إنشاء مزود Ollama جديد
func NewOllamaProvider() (*OllamaProvider, error) {
    baseURL := os.Getenv("OLLAMA_HOST")
    if baseURL == "" {
        baseURL = "http://localhost:11434"
    }
    
    provider := &OllamaProvider{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 300 * time.Second,
        },
    }
    
    // تحميل قائمة النماذج المتاحة
    if err := provider.loadModels(); err != nil {
        return nil, fmt.Errorf("failed to load Ollama models: %w", err)
    }
    
    return provider, nil
}

// loadModels تحميل النماذج المتاحة من Ollama
func (p *OllamaProvider) loadModels() error {
    url := p.baseURL + "/api/tags"
    
    resp, err := p.httpClient.Get(url)
    if err != nil {
        return fmt.Errorf("failed to connect to Ollama: %w", err)
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    
    var result struct {
        Models []struct {
            Name string `json:"name"`
        } `json:"models"`
    }
    
    if err := json.Unmarshal(body, &result); err != nil {
        return err
    }
    
    for _, model := range result.Models {
        p.models = append(p.models, model.Name)
    }
    
    return nil
}

// Generate توليد نص
func (p *OllamaProvider) Generate(prompt string, options ...Option) (string, error) {
    opts := &Options{
        Model:       "llama3.2:3b",
        Temperature: 0.7,
        MaxTokens:   2000,
    }
    
    for _, opt := range options {
        opt(opts)
    }
    
    // التأكد من وجود النموذج
    if !p.hasModel(opts.Model) {
        // استخدام النموذج الافتراضي
        opts.Model = "llama3.2:3b"
    }
    
    return p.generateText(prompt, opts)
}

// generateText توليد نص باستخدام Ollama
func (p *OllamaProvider) generateText(prompt string, opts *Options) (string, error) {
    url := p.baseURL + "/api/generate"
    
    request := map[string]interface{}{
        "model":  opts.Model,
        "prompt": prompt,
        "stream": false,
        "options": map[string]interface{}{
            "temperature": opts.Temperature,
            "num_predict": opts.MaxTokens,
            "top_p":       opts.TopP,
            "top_k":       opts.TopK,
            "repeat_penalty": opts.RepetitionPenalty,
        },
    }
    
    if opts.SystemPrompt != "" {
        request["system"] = opts.SystemPrompt
    }
    
    jsonData, err := json.Marshal(request)
    if err != nil {
        return "", fmt.Errorf("failed to marshal request: %w", err)
    }
    
    resp, err := p.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("Ollama request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("Ollama API error: %s - %s", resp.Status, string(body))
    }
    
    var result struct {
        Response string `json:"response"`
        Done     bool   `json:"done"`
        Model    string `json:"model"`
        CreatedAt string `json:"created_at"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", fmt.Errorf("failed to decode response: %w", err)
    }
    
    return strings.TrimSpace(result.Response), nil
}

// GenerateStream توليد نص بشكل متدفق
func (p *OllamaProvider) GenerateStream(prompt string, options ...Option) (<-chan string, <-chan error, context.CancelFunc) {
    opts := &Options{
        Model:       "llama3.2:3b",
        Temperature: 0.7,
        MaxTokens:   2000,
    }
    
    for _, opt := range options {
        opt(opts)
    }
    
    ctx, cancel := context.WithCancel(context.Background())
    textChan := make(chan string)
    errChan := make(chan error, 1)
    
    go func() {
        defer close(textChan)
        defer close(errChan)
        
        url := p.baseURL + "/api/generate"
        
        request := map[string]interface{}{
            "model":  opts.Model,
            "prompt": prompt,
            "stream": true,
            "options": map[string]interface{}{
                "temperature": opts.Temperature,
                "num_predict": opts.MaxTokens,
            },
        }
        
        jsonData, err := json.Marshal(request)
        if err != nil {
            errChan <- err
            return
        }
        
        req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
        if err != nil {
            errChan <- err
            return
        }
        req.Header.Set("Content-Type", "application/json")
        
        resp, err := p.httpClient.Do(req)
        if err != nil {
            errChan <- err
            return
        }
        defer resp.Body.Close()
        
        if resp.StatusCode != http.StatusOK {
            body, _ := io.ReadAll(resp.Body)
            errChan <- fmt.Errorf("Ollama API error: %s - %s", resp.Status, string(body))
            return
        }
        
        decoder := json.NewDecoder(resp.Body)
        var fullText strings.Builder
        
        for {
            select {
            case <-ctx.Done():
                return
            default:
                var chunk struct {
                    Response string `json:"response"`
                    Done     bool   `json:"done"`
                }
                
                if err := decoder.Decode(&chunk); err != nil {
                    if err == io.EOF {
                        return
                    }
                    errChan <- err
                    return
                }
                
                if chunk.Response != "" {
                    fullText.WriteString(chunk.Response)
                    textChan <- chunk.Response
                }
                
                if chunk.Done {
                    return
                }
            }
        }
    }()
    
    return textChan, errChan, cancel
}

// Embed توليد embeddings
func (p *OllamaProvider) Embed(text string, model string) ([]float64, error) {
    if model == "" {
        model = "nomic-embed-text"
    }
    
    url := p.baseURL + "/api/embed"
    
    request := map[string]interface{}{
        "model": model,
        "input": text,
    }
    
    jsonData, err := json.Marshal(request)
    if err != nil {
        return nil, err
    }
    
    resp, err := p.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result struct {
        Embedding []float64 `json:"embedding"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return result.Embedding, nil
}

// PullModel سحب نموذج جديد
func (p *OllamaProvider) PullModel(model string) error {
    url := p.baseURL + "/api/pull"
    
    request := map[string]interface{}{
        "name": model,
        "stream": false,
    }
    
    jsonData, err := json.Marshal(request)
    if err != nil {
        return err
    }
    
    resp, err := p.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("failed to pull model: %s - %s", resp.Status, string(body))
    }
    
    // إضافة النموذج إلى القائمة
    p.models = append(p.models, model)
    
    return nil
}

// ListModels عرض النماذج المتاحة
func (p *OllamaProvider) ListModels() ([]ModelInfo, error) {
    url := p.baseURL + "/api/tags"
    
    resp, err := p.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result struct {
        Models []struct {
            Name       string `json:"name"`
            ModifiedAt string `json:"modified_at"`
            Size       int64  `json:"size"`
            Digest     string `json:"digest"`
        } `json:"models"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    var models []ModelInfo
    for _, m := range result.Models {
        models = append(models, ModelInfo{
            ID:   m.Name,
            Name: m.Name,
            Size: m.Size,
        })
    }
    
    return models, nil
}

// ModelInfo معلومات النموذج
type ModelInfo struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Size int64  `json:"size"`
}

// hasModel التحقق من وجود النموذج
func (p *OllamaProvider) hasModel(modelName string) bool {
    for _, m := range p.models {
        if m == modelName {
            return true
        }
    }
    return false
}

// GetStats الحصول على إحصائيات
func (p *OllamaProvider) GetStats() (map[string]interface{}, error) {
    // محاولة الحصول على إحصائيات من Ollama
    url := p.baseURL + "/api/version"
    
    resp, err := p.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var versionInfo struct {
        Version string `json:"version"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&versionInfo); err != nil {
        return nil, err
    }
    
    stats := map[string]interface{}{
        "provider":        "ollama",
        "version":         versionInfo.Version,
        "models_count":    len(p.models),
        "models":          p.models,
        "status":          "online",
        "supports_stream": true,
        "supports_embed":  true,
    }
    
    return stats, nil
}

// IsAvailable التحقق من التوفر
func (p *OllamaProvider) IsAvailable() bool {
    resp, err := p.httpClient.Get(p.baseURL + "/api/tags")
    return err == nil && resp.StatusCode == http.StatusOK
}

// GetName اسم المزود
func (p *OllamaProvider) GetName() string {
    return "Ollama"
}

// GetCost التكلفة (مجاني بالكامل)
func (p *OllamaProvider) GetCost() float64 {
    return 0.0
}

// Options خيارات التوليد
type Options struct {
    Model             string
    Temperature       float64
    MaxTokens         int
    TopP              float64
    TopK              int
    RepetitionPenalty float64
    SystemPrompt      string
}

// Option دالة لتعديل الخيارات
type Option func(*Options)

// WithModel تحديد النموذج
func WithModel(model string) Option {
    return func(o *Options) {
        o.Model = model
    }
}

// WithTemperature تحديد درجة الحرارة
func WithTemperature(temp float64) Option {
    return func(o *Options) {
        o.Temperature = temp
    }
}

// WithMaxTokens تحديد الحد الأقصى للرموز
func WithMaxTokens(tokens int) Option {
    return func(o *Options) {
        o.MaxTokens = tokens
    }
}

// WithSystemPrompt إضافة نظام prompt
func WithSystemPrompt(prompt string) Option {
    return func(o *Options) {
        o.SystemPrompt = prompt
    }
}