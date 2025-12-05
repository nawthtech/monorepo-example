package ai

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
    
    "google.golang.org/genai"
)

// VideoProvider مزود خاص لتوليد الفيديوهات
type VideoProvider struct {
    client *genai.Client
    apiKey string
}

// NewVideoProvider إنشاء مزود فيديوهات جديد
func NewVideoProvider() (*VideoProvider, error) {
    apiKey := os.Getenv("GEMINI_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
    }
    
    ctx := context.Background()
    client, err := genai.NewClient(ctx, &genai.ClientConfig{
        APIKey: apiKey,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create video client: %w", err)
    }
    
    return &VideoProvider{
        client: client,
        apiKey: apiKey,
    }, nil
}

// VideoRequest طلب توليد فيديو
type VideoRequest struct {
    Prompt       string `json:"prompt" binding:"required"`
    Duration     int    `json:"duration"`     // بالثواني
    AspectRatio  string `json:"aspect_ratio"` // 16:9, 1:1, 9:16
    Style        string `json:"style"`        // realistic, animated, cinematic
    OutputFormat string `json:"output_format"` // mp4, gif
}

// VideoResponse استجابة توليد فيديو
type VideoResponse struct {
    Success      bool     `json:"success"`
    VideoURL     string   `json:"video_url,omitempty"`
    VideoData    []byte   `json:"video_data,omitempty"`
    Duration     float64  `json:"duration"`
    Size         int64    `json:"size"`
    ModelUsed    string   `json:"model_used"`
    GenerationID string   `json:"generation_id"`
    Status       string   `json:"status"`
    Error        string   `json:"error,omitempty"`
}

// GenerateVideo توليد فيديو باستخدام Veo
func (p *VideoProvider) GenerateVideo(req VideoRequest) (*VideoResponse, error) {
    ctx := context.Background()
    
    // بناء prompt محسن
    prompt := p.buildVideoPrompt(req.Prompt, req.Style, req.Duration)
    
    // إعداد config
    var config genai.GenerateVideosConfig
    
    // تحديد النموذج (Veo 2.0)
    modelName := "veo-2.0-generate-001"
    
    // Call the GenerateVideo method.
    operation, err := p.client.Models.GenerateVideos(ctx, modelName, prompt, nil, nil, &config)
    if err != nil {
        return nil, fmt.Errorf("failed to start video generation: %w", err)
    }
    
    log.Printf("🎬 Video generation started. Operation: %s", operation.Name)
    
    // انتظار اكتمال العملية
    videoURL, videoData, err := p.waitForVideoCompletion(ctx, operation)
    if err != nil {
        return nil, fmt.Errorf("video generation failed: %w", err)
    }
    
    return &VideoResponse{
        Success:      true,
        VideoURL:     videoURL,
        VideoData:    videoData,
        Duration:     float64(req.Duration),
        Size:         int64(len(videoData)),
        ModelUsed:    modelName,
        GenerationID: operation.Name,
        Status:       "completed",
    }, nil
}

// waitForVideoCompletion انتظار اكتمال توليد الفيديو
func (p *VideoProvider) waitForVideoCompletion(ctx context.Context, operation *genai.VideosOperation) (string, []byte, error) {
    maxAttempts := 30 // 10 دقائق كحد أقصى
    attempt := 0
    
    for !operation.Done && attempt < maxAttempts {
        attempt++
        log.Printf("⏳ Waiting for video generation... Attempt %d/%d", attempt, maxAttempts)
        
        time.Sleep(20 * time.Second)
        
        var err error
        operation, err = p.client.Operations.GetVideosOperation(ctx, operation, nil)
        if err != nil {
            return "", nil, fmt.Errorf("failed to check operation status: %w", err)
        }
    }
    
    if !operation.Done {
        return "", nil, fmt.Errorf("video generation timed out after %d attempts", maxAttempts)
    }
    
    // تحقق من وجود أخطاء
    if operation.Error != nil {
        return "", nil, fmt.Errorf("video generation error: %v", operation.Error)
    }
    
    log.Printf("✅ Video generation completed successfully")
    
    // تنزيل الفيديو إذا كان متاحاً
    if p.client.ClientConfig().Backend != genai.BackendVertexAI {
        for _, v := range operation.Response.GeneratedVideos {
            data, err := p.client.Files.Download(ctx, genai.NewDownloadURIFromGeneratedVideo(v), nil)
            if err != nil {
                log.Printf("⚠️ Failed to download video: %v", err)
                continue
            }
            
            log.Printf("📥 Video downloaded. Size: %d bytes", len(data))
            return v.Video.URI, data, nil
        }
    }
    
    // إذا كان VertexAI، إرجاع URI فقط
    if len(operation.Response.GeneratedVideos) > 0 {
        return operation.Response.GeneratedVideos[0].Video.URI, nil, nil
    }
    
    return "", nil, fmt.Errorf("no video generated")
}

// buildVideoPrompt بناء prompt فيديو محسن
func (p *VideoProvider) buildVideoPrompt(prompt, style string, duration int) string {
    styleMap := map[string]string{
        "realistic":  "realistic, cinematic, high-quality video",
        "animated":   "animated, cartoon style, vibrant colors",
        "cinematic":  "cinematic, movie-like, dramatic lighting",
        "corporate":  "corporate, professional, clean animation",
        "social":     "social media optimized, eye-catching, vertical format",
    }
    
    styleDesc := styleMap[style]
    if styleDesc == "" {
        styleDesc = "high-quality, professional"
    }
    
    return fmt.Sprintf(`Create a %s video: "%s"
    
    Requirements:
    - Style: %s
    - Duration: %d seconds
    - Professional quality
    - Smooth animation
    - Clear visual storytelling
    - Optimized for digital platforms`, 
    styleDesc, prompt, styleDesc, duration)
}

// GetVideoStatus الحصول على حالة فيديو
func (p *VideoProvider) GetVideoStatus(operationID string) (*genai.VideosOperation, error) {
    ctx := context.Background()
    
    operation := &genai.VideosOperation{Name: operationID}
    
    op, err := p.client.Operations.GetVideosOperation(ctx, operation, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get video status: %w", err)
    }
    
    return op, nil
}

// SaveVideoToFile حفظ الفيديو في ملف
func (p *VideoProvider) SaveVideoToFile(videoData []byte, filename string) error {
    if len(videoData) == 0 {
        return fmt.Errorf("no video data to save")
    }
    
    // التأكد من امتداد الملف
    if len(filename) < 4 || filename[len(filename)-4:] != ".mp4" {
        filename = filename + ".mp4"
    }
    
    // إنشاء مجلد videos إذا لم يكن موجوداً
    if err := os.MkdirAll("videos", 0755); err != nil {
        return fmt.Errorf("failed to create videos directory: %w", err)
    }
    
    filepath := "videos/" + filename
    
    if err := os.WriteFile(filepath, videoData, 0644); err != nil {
        return fmt.Errorf("failed to save video file: %w", err)
    }
    
    log.Printf("💾 Video saved to: %s", filepath)
    return nil
}

// GenerateNawthTechVideo توليد فيديو مخصص لـ NawthTech
func (p *VideoProvider) GenerateNawthTechVideo(videoType, topic string) (*VideoResponse, error) {
    prompts := map[string]string{
        "explainer":   "An animated explainer video about %s for digital marketing and business growth",
        "promotional": "A promotional video showcasing %s for NawthTech platform with futuristic UI elements",
        "tutorial":    "A step-by-step tutorial video showing how to use %s on NawthTech platform",
        "testimonial": "A video testimonial animation for %s with customer success stories",
        "social":      "A short, engaging social media video about %s optimized for Instagram and TikTok",
    }
    
    promptTemplate, exists := prompts[videoType]
    if !exists {
        promptTemplate = "A professional video about %s for digital growth"
    }
    
    prompt := fmt.Sprintf(promptTemplate, topic)
    
    req := VideoRequest{
        Prompt:      prompt,
        Duration:    30, // ثواني
        AspectRatio: "16:9",
        Style:       "animated",
    }
    
    return p.GenerateVideo(req)
}