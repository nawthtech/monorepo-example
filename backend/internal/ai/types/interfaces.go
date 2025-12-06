package types

import "context"

// TextProvider واجهة لمزودي النصوص
type TextProvider interface {
    GenerateText(ctx context.Context, prompt string, options map[string]interface{}) (string, error)
    Name() string
}

// ImageProvider واجهة لمزودي الصور
type ImageProvider interface {
    GenerateImage(ctx context.Context, prompt string, options map[string]interface{}) ([]byte, error)
    Name() string
}

// VideoProvider واجهة لمزودي الفيديو
type VideoProvider interface {
    GenerateVideo(ctx context.Context, prompt string, options map[string]interface{}) ([]byte, error)
    Name() string
}

// AIRequest طلب AI عام
type AIRequest struct {
    Prompt  string                 `json:"prompt"`
    Type    string                 `json:"type"`
    Options map[string]interface{} `json:"options"`
}

// AIResponse استجابة AI عامة
type AIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data"`
    Error   string      `json:"error,omitempty"`
}
