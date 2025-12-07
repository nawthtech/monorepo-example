package video

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "time"
)

// SVDClient عميل Stable Video Diffusion
type SVDClient struct {
    apiKey  string
    baseURL string
    client  *http.Client
}

// NewSVDClient إنشاء عميل SVD جديد
func NewSVDClient() *SVDClient {
    apiKey := os.Getenv("STABILITY_API_KEY")
    baseURL := "https://api.stability.ai/v2alpha/generation"
    
    if apiKey == "" {
        // إذا لم يكن هناك مفتاح API، نستخدم URL بديل للاختبار
        baseURL = "https://api.stability.ai/v2alpha"
    }
    
    return &SVDClient{
        apiKey:  apiKey,
        baseURL: baseURL,
        client:  &http.Client{Timeout: 120 * time.Second},
    }
}

// GenerateVideo توليد فيديو باستخدام Stable Video Diffusion
func (c *SVDClient) GenerateVideo(req VideoRequest) (*VideoResponse, error) {
    if c.apiKey == "" {
        return nil, &VideoError{
            Code:    "api_key_missing",
            Message: "Stability API key is required for video generation",
        }
    }
    
    // التحقق من صحة الطلب
    if err := ValidateVideoRequest(req); err != nil {
        return nil, err
    }
    
    // إنشاء جسم الطلب
    reqBody := c.buildSVDRequest(req)
    jsonBody, err := json.Marshal(reqBody)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %v", err)
    }
    
    // إرسال الطلب
    httpReq, err := http.NewRequest("POST", c.baseURL+"/video", bytes.NewBuffer(jsonBody))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }
    
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
    httpReq.Header.Set("Accept", "application/json")
    
    // إرسال الطلب
    resp, err := c.client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("API request failed: %v", err)
    }
    defer resp.Body.Close()
    
    // قراءة الاستجابة
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %v", err)
    }
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
    }
    
    // تحليل الاستجابة
    var svdResponse SVDResponse
    if err := json.Unmarshal(body, &svdResponse); err != nil {
        return nil, fmt.Errorf("failed to parse response: %v", err)
    }
    
    // تحويل إلى VideoResponse
    return c.convertToVideoResponse(req, &svdResponse), nil
}

// buildSVDRequest بناء طلب SVD
func (c *SVDClient) buildSVDRequest(req VideoRequest) map[string]interface{} {
    // تحويل الدقة إلى أبعاد
    width, height, _ := ResolutionToDimensions(req.Resolution)
    
    requestBody := map[string]interface{}{
        "text_prompts": []map[string]interface{}{
            {
                "text":   req.Prompt,
                "weight": 1.0,
            },
        },
        "cfg_scale":     req.Options.CFGScale,
        "clip_guidance_preset": "FAST_BLUE",
        "sampler":       "K_DPMPP_2M",
        "frames":        c.calculateFrames(req.Duration),
        "fps":           10,
        "seed":          req.Options.Seed,
        "width":         width,
        "height":        height,
    }
    
    // إضافة النص السلبي إذا كان موجوداً
    if req.NegativePrompt != "" {
        requestBody["text_prompts"] = append(requestBody["text_prompts"].([]map[string]interface{}),
            map[string]interface{}{
                "text":   req.NegativePrompt,
                "weight": -1.0,
            })
    }
    
    // إضافة خيارات النمط
    if req.Style != "" {
        requestBody["style_preset"] = c.mapStyleToPreset(req.Style)
    }
    
    return requestBody
}

// convertToVideoResponse تحويل استجابة SVD إلى VideoResponse
func (c *SVDClient) convertToVideoResponse(req VideoRequest, svdResp *SVDResponse) *VideoResponse {
    width, height, _ := ResolutionToDimensions(req.Resolution)
    
    return &VideoResponse{
        Success:    svdResp.Artifacts != nil && len(svdResp.Artifacts) > 0,
        VideoURL:   c.extractVideoURL(svdResp),
        Duration:   req.Duration,
        Width:      width,
        Height:     height,
        Resolution: req.Resolution,
        Format:     string(FormatMP4),
        Provider:   "stability_svd",
        Cost:       c.calculateCost(req),
        Status:     string(VideoJobCompleted),
        CreatedAt:  time.Now(),
        Timestamp:  time.Now().Unix(),
    }
}

// calculateCost حساب تكلفة التوليد
func (c *SVDClient) calculateCost(req VideoRequest) float64 {
    // Stability AI له خطة مجانية: 25 توليد/شهر
    // بعد ذلك: $0.02 لكل فيديو
    // هذا حساب تقديري
    frames := c.calculateFrames(req.Duration)
    return float64(frames) * 0.0001 // $0.0001 لكل إطار
}

// calculateFrames حساب عدد الإطارات
func (c *SVDClient) calculateFrames(duration int) int {
    fps := 10 // SVD يدعم 10 FPS
    frames := duration * fps
    
    // الحدود المسموح بها من Stability AI
    if frames < 14 {
        return 14
    }
    if frames > 100 {
        return 100
    }
    return frames
}

// mapStyleToPreset تحويل النمط إلى preset
func (c *SVDClient) mapStyleToPreset(style string) string {
    styleMap := map[string]string{
        "realistic":   "realistic",
        "anime":       "anime",
        "cartoon":     "cartoon",
        "artistic":    "enhance",
        "cinematic":   "cinematic",
        "minimal":     "low-poly",
    }
    
    if preset, ok := styleMap[style]; ok {
        return preset
    }
    return "realistic" // افتراضي
}

// extractVideoURL استخراج رابط الفيديو من الاستجابة
func (c *SVDClient) extractVideoURL(svdResp *SVDResponse) string {
    if svdResp.Artifacts == nil || len(svdResp.Artifacts) == 0 {
        return ""
    }
    
    // في Stability AI، الفيديوهات تُرجع كـ base64 أو رابط
    // هذا تنفيذ مبسط
    for _, artifact := range svdResp.Artifacts {
        if artifact.Base64 != "" {
            // في حالة البيانات المشفرة base64
            return fmt.Sprintf("data:video/mp4;base64,%s", artifact.Base64)
        }
    }
    
    return ""
}

// IsAvailable التحقق من توفر الخدمة
func (c *SVDClient) IsAvailable() bool {
    if c.apiKey == "" {
        return false
    }
    
    // اختبار الاتصال بـ API
    req, err := http.NewRequest("GET", c.baseURL+"/user/account", nil)
    if err != nil {
        return false
    }
    
    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    
    resp, err := c.client.Do(req)
    if err != nil {
        return false
    }
    defer resp.Body.Close()
    
    return resp.StatusCode == http.StatusOK
}

// GetUsageInfo الحصول على معلومات الاستخدام
func (c *SVDClient) GetUsageInfo() (map[string]interface{}, error) {
    if c.apiKey == "" {
        return nil, fmt.Errorf("API key not set")
    }
    
    req, err := http.NewRequest("GET", c.baseURL+"/user/balance", nil)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var balanceInfo map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&balanceInfo); err != nil {
        return nil, err
    }
    
    return balanceInfo, nil
}

// SVDResponse استجابة Stable Video Diffusion
type SVDResponse struct {
    ID         string            `json:"id"`
    Model      string            `json:"model"`
    Created    int64             `json:"created"`
    Artifacts  []SVDArtifact     `json:"artifacts"`
    Usage      SVDUsage          `json:"usage"`
}

// SVDArtifact قطعة فيديو
type SVDArtifact struct {
    Base64    string  `json:"base64"`
    Seed      int64   `json:"seed"`
    MIMEType  string  `json:"mime_type"`
    Classifer float64 `json:"classifier"`
}

// SVDUsage استخدام API
type SVDUsage struct {
    Images    int     `json:"images"`
    Credits   float64 `json:"credits"`
    Remaining float64 `json:"remaining"`
}

// SVDProvider مزود SVD (للتكامل مع نظام الفيديو)
type SVDProvider struct {
    client *SVDClient
}

// NewSVDProvider إنشاء مزود SVD جديد
func NewSVDProvider() *SVDProvider {
    return &SVDProvider{
        client: NewSVDClient(),
    }
}

// GenerateVideo توليد فيديو
func (p *SVDProvider) GenerateVideo(req VideoRequest) (*VideoResponse, error) {
    return p.client.GenerateVideo(req)
}

// Name اسم المزود
func (p *SVDProvider) Name() string {
    return "stability_svd"
}

// IsAvailable التحقق من التوفر
func (p *SVDProvider) IsAvailable() bool {
    return p.client.IsAvailable()
}

// IsLocal هل المزود محلي؟
func (p *SVDProvider) IsLocal() bool {
    return false
}

// IsFree هل المزود مجاني؟
func (p *SVDProvider) IsFree() bool {
    // Stability AI له حد مجاني 25 توليد/شهر
    return true
}

// SupportsResolution دعم الدقة
func (p *SVDProvider) SupportsResolution(resolution string) bool {
    // SVD يدعم دقات محددة
    supported := []string{
        "512x512", "576x1024", "1024x576",
        "768x768", "1024x1024",
    }
    
    for _, res := range supported {
        if res == resolution {
            return true
        }
    }
    return false
}

// GetCapabilities قدرات المزود
func (p *SVDProvider) GetCapabilities() map[string]interface{} {
    return map[string]interface{}{
        "provider":          "stability_svd",
        "local":             false,
        "free_tier":         true,
        "free_limit":        25, // توليد/شهر
        "max_duration":      10, // ثواني
        "max_frames":        100,
        "fps":               10,
        "supported_resolutions": []string{
            "512x512", "576x1024", "1024x576",
            "768x768", "1024x1024",
        },
        "price_per_video":   0.02, // دولار بعد الحد المجاني
    }
}