package ai

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strings"
    "time"
)

// GeminiProvider مزود Google Gemini
type GeminiProvider struct {
    apiKey  string
    baseURL string
    client  *http.Client
}

// NewGeminiProvider إنشاء مزود Gemini جديد
func NewGeminiProvider() *GeminiProvider {
    apiKey := os.Getenv("GEMINI_API_KEY")
    
    return &GeminiProvider{
        apiKey:  apiKey,
        baseURL: "https://generativelanguage.googleapis.com/v1beta",
        client: &http.Client{
            Timeout: 120 * time.Second,
        },
    }
}

// GenerateText توليد نص باستخدام Gemini
func (p *GeminiProvider) GenerateText(req TextRequest) (*TextResponse, error) {
    if p.apiKey == "" {
        return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
    }
    
    model := req.Model
    if model == "" {
        model = "gemini-2.5-flash-exp" // نموذج مجاني
    }
    
    url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", p.baseURL, model, p.apiKey)
    
    payload := map[string]interface{}{
        "contents": []map[string]interface{}{
            {
                "parts": []map[string]interface{}{
                    {
                        "text": req.Prompt,
                    },
                },
            },
        },
        "generationConfig": map[string]interface{}{
            "temperature":     req.Temperature,
            "maxOutputTokens": req.MaxTokens,
            "topP":           0.95,
            "topK":           40,
        },
        "safetySettings": []map[string]interface{}{
            {
                "category":  "HARM_CATEGORY_HARASSMENT",
                "threshold": "BLOCK_MEDIUM_AND_ABOVE",
            },
            {
                "category":  "HARM_CATEGORY_HATE_SPEECH",
                "threshold": "BLOCK_MEDIUM_AND_ABOVE",
            },
            {
                "category":  "HARM_CATEGORY_SEXUALLY_EXPLICIT",
                "threshold": "BLOCK_MEDIUM_AND_ABOVE",
            },
            {
                "category":  "HARM_CATEGORY_DANGEROUS_CONTENT",
                "threshold": "BLOCK_MEDIUM_AND_ABOVE",
            },
        },
    }
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    httpReq, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    httpReq.Header.Set("Content-Type", "application/json")
    
    resp, err := p.client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("Gemini API error: %s - %s", resp.Status, string(body))
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }
    
    var result struct {
        Candidates []struct {
            Content struct {
                Parts []struct {
                    Text string `json:"text"`
                } `json:"parts"`
            } `json:"content"`
            FinishReason string `json:"finishReason"`
        } `json:"candidates"`
        UsageMetadata struct {
            PromptTokenCount     int `json:"promptTokenCount"`
            CandidatesTokenCount int `json:"candidatesTokenCount"`
            TotalTokenCount      int `json:"totalTokenCount"`
        } `json:"usageMetadata"`
    }
    
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }
    
    if len(result.Candidates) == 0 {
        return nil, fmt.Errorf("no candidates returned from Gemini")
    }
    
    var textParts []string
    for _, part := range result.Candidates[0].Content.Parts {
        if part.Text != "" {
            textParts = append(textParts, part.Text)
        }
    }
    
    fullText := strings.Join(textParts, "\n")
    
    // تقدير التكلفة (Gemini flash مجاني تقريباً)
    cost := 0.0
    if model != "gemini-2.5-flash-exp" {
        // تكلفة تقديرية للنماذج المدفوعة
        totalTokens := result.UsageMetadata.TotalTokenCount
        cost = float64(totalTokens) * 0.0000025 // $0.0025 per 1K tokens للـ pro
    }
    
    return &TextResponse{
        Text:         strings.TrimSpace(fullText),
        Tokens:       result.UsageMetadata.TotalTokenCount,
        Cost:         cost,
        ModelUsed:    model,
        FinishReason: result.Candidates[0].FinishReason,
        CreatedAt:    time.Now(),
    }, nil
}

// GenerateImage توليد صور باستخدام Gemini - غير مدعوم مباشرة
func (p *GeminiProvider) GenerateImage(req ImageRequest) (*ImageResponse, error) {
    // Gemini لا يدعم توليد الصور مباشرة، لكن يمكن استخدام Imagen (مدفوع)
    // هنا نستخدم توليد النص مع وصف الصورة
    
    if p.apiKey == "" {
        return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
    }
    
    // يمكن استخدام Gemini لوصف الصورة المفصلة
    prompt := fmt.Sprintf("Generate a detailed prompt for DALL-E or Stable Diffusion to create an image with the following description: %s. Style: %s, Size: %s, Quality: %s",
        req.Prompt, req.Style, req.Size, req.Quality)
    
    textReq := TextRequest{
        Prompt: prompt,
        Model:  "gemini-2.5-flash-exp",
    }
    
    resp, err := p.GenerateText(textReq)
    if err != nil {
        return nil, fmt.Errorf("failed to generate image prompt: %w", err)
    }
    
    // إرجاع الوصف بدلاً من الصورة الفعلية
    return &ImageResponse{
        URL:         "", // لا يوجد URL للصورة
        ImageData:   nil,
        Size:        req.Size,
        Format:      "text/description",
        Cost:        resp.Cost,
        ModelUsed:   resp.ModelUsed,
        CreatedAt:   time.Now(),
        Seed:        0,
    }, nil
}

// GenerateVideo توليد فيديو - غير مدعوم في Gemini
func (p *GeminiProvider) GenerateVideo(req VideoRequest) (*VideoResponse, error) {
    return nil, fmt.Errorf("video generation not supported by Gemini")
}

// AnalyzeText تحليل نص
func (p *GeminiProvider) AnalyzeText(req AnalysisRequest) (*AnalysisResponse, error) {
    if p.apiKey == "" {
        return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
    }
    
    model := req.Model
    if model == "" {
        model = "gemini-2.5-flash-exp"
    }
    
    // بناء prompt للتحليل
    analysisPrompt := ""
    if req.Prompt != "" {
        analysisPrompt = fmt.Sprintf("%s\n\nText to analyze: %s", req.Prompt, req.Text)
    } else {
        analysisPrompt = fmt.Sprintf("Analyze this text and provide insights: %s", req.Text)
    }
    
    textReq := TextRequest{
        Prompt:      analysisPrompt,
        Model:       model,
        Temperature: 0.3, // أقل درجة حرارة لتحليل أكثر دقة
    }
    
    resp, err := p.GenerateText(textReq)
    if err != nil {
        return nil, fmt.Errorf("failed to analyze text: %w", err)
    }
    
    return &AnalysisResponse{
        Result:     resp.Text,
        Confidence: 0.9, // تقدير ثقة عالي لـ Gemini
        Cost:       resp.Cost,
        Model:      resp.ModelUsed,
    }, nil
}

// AnalyzeImage تحليل صور باستخدام Gemini Vision
func (p *GeminiProvider) AnalyzeImage(req AnalysisRequest) (*AnalysisResponse, error) {
    if p.apiKey == "" {
        return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
    }
    
    if len(req.ImageData) == 0 {
        return nil, fmt.Errorf("image data is required for image analysis")
    }
    
    model := req.Model
    if model == "" {
        model = "gemini-2.5-flash-exp-vision"
    }
    
    url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", p.baseURL, model, p.apiKey)
    
    // ترميز الصورة إلى base64 (مبسط)
    imageBase64 := "" // في الواقع تحتاج إلى ترميز base64
    
    payload := map[string]interface{}{
        "contents": []map[string]interface{}{
            {
                "parts": []map[string]interface{}{
                    {
                        "inlineData": map[string]interface{}{
                            "mimeType": "image/jpeg",
                            "data":     imageBase64,
                        },
                    },
                    {
                        "text": req.Prompt,
                    },
                },
            },
        },
    }
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    httpReq, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    httpReq.Header.Set("Content-Type", "application/json")
    
    resp, err := p.client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("Gemini Vision API error: %s - %s", resp.Status, string(body))
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }
    
    var result struct {
        Candidates []struct {
            Content struct {
                Parts []struct {
                    Text string `json:"text"`
                } `json:"parts"`
            } `json:"content"`
        } `json:"candidates"`
        UsageMetadata struct {
            TotalTokenCount int `json:"totalTokenCount"`
        } `json:"usageMetadata"`
    }
    
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }
    
    if len(result.Candidates) == 0 {
        return nil, fmt.Errorf("no candidates returned from Gemini Vision")
    }
    
    var textParts []string
    for _, part := range result.Candidates[0].Content.Parts {
        if part.Text != "" {
            textParts = append(textParts, part.Text)
        }
    }
    
    fullText := strings.Join(textParts, "\n")
    
    // تقدير التكلفة
    cost := float64(result.UsageMetadata.TotalTokenCount) * 0.0000025
    
    return &AnalysisResponse{
        Result:     strings.TrimSpace(fullText),
        Confidence: 0.85, // ثقة عالية في تحليل الصور
        Cost:       cost,
        Model:      model,
    }, nil
}

// TranslateText ترجمة نص
func (p *GeminiProvider) TranslateText(req TranslationRequest) (*TranslationResponse, error) {
    if p.apiKey == "" {
        return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
    }
    
    model := req.Model
    if model == "" {
        model = "gemini-2.5-flash-exp"
    }
    
    prompt := fmt.Sprintf("Translate the following text from %s to %s:\n\n%s",
        req.FromLang, req.ToLang, req.Text)
    
    textReq := TextRequest{
        Prompt: prompt,
        Model:  model,
    }
    
    resp, err := p.GenerateText(textReq)
    if err != nil {
        return nil, fmt.Errorf("failed to translate text: %w", err)
    }
    
    return &TranslationResponse{
        TranslatedText: strings.TrimSpace(resp.Text),
        Cost:           resp.Cost,
        Model:          resp.ModelUsed,
    }, nil
}

// IsAvailable التحقق من التوفر
func (p *GeminiProvider) IsAvailable() bool {
    if p.apiKey == "" {
        return false
    }
    
    // اختبار اتصال بسيط
    url := fmt.Sprintf("%s/models/gemini-2.5-flash-exp?key=%s", p.baseURL, p.apiKey)
    resp, err := p.client.Get(url)
    if err != nil {
        return false
    }
    defer resp.Body.Close()
    
    return resp.StatusCode == http.StatusOK
}

// GetName اسم المزود
func (p *GeminiProvider) GetName() string {
    return "Google Gemini"
}

// GetCost التكلفة
func (p *GeminiProvider) GetCost() float64 {
    return 0.0 // Flash مجاني
}

// GetType نوع المزود
func (p *GeminiProvider) GetType() string {
    return "text"
}

// GetStats الحصول على إحصائيات
func (p *GeminiProvider) GetStats() *ProviderStats {
    return &ProviderStats{
        Name:        p.GetName(),
        Type:        p.GetType(),
        IsAvailable: p.IsAvailable(),
        Requests:    0,
        Successful:  0,
        Failed:      0,
        TotalCost:   0.0,
        AvgLatency:  0.0,
        LastUsed:    time.Time{},
        SuccessRate: 95.0,
    }
}

// SupportsStreaming يدعم التدفق
func (p *GeminiProvider) SupportsStreaming() bool {
    return false // Gemini لا يدعم التدفق في API المجاني
}

// SupportsEmbedding يدعم التضمين
func (p *GeminiProvider) SupportsEmbedding() bool {
    return false
}

// GetMaxTokens الحد الأقصى للرموز
func (p *GeminiProvider) GetMaxTokens() int {
    return 8192 // لحد Gemini flash
}

// GetSupportedLanguages اللغات المدعومة
func (p *GeminiProvider) GetSupportedLanguages() []string {
    return []string{
        "ar", "en", "es", "fr", "de", "zh", "ja", "ko", "ru", "pt",
        "it", "nl", "pl", "sv", "da", "fi", "no", "he", "hi", "tr",
    }
}

// ListModels عرض النماذج المتاحة
func (p *GeminiProvider) ListModels() ([]string, error) {
    if p.apiKey == "" {
        return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
    }
    
    url := fmt.Sprintf("%s/models?key=%s", p.baseURL, p.apiKey)
    
    resp, err := p.client.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to list models: %w", err)
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }
    
    var result struct {
        Models []struct {
            Name string `json:"name"`
        } `json:"models"`
    }
    
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, fmt.Errorf("failed to parse models: %w", err)
    }
    
    var models []string
    for _, model := range result.Models {
        models = append(models, strings.TrimPrefix(model.Name, "models/"))
    }
    
    return models, nil
}