package types

import "fmt"

// أخطاء مخصصة
var (
    ErrNoProviderAvailable   = &ProviderError{Code: "NO_PROVIDER_AVAILABLE", Message: "No AI provider is currently available"}
    ErrProviderNotFound      = &ProviderError{Code: "PROVIDER_NOT_FOUND", Message: "AI provider not found"}
    ErrProviderUnavailable   = &ProviderError{Code: "PROVIDER_UNAVAILABLE", Message: "AI provider is currently unavailable"}
    ErrInvalidAPIKey         = &ProviderError{Code: "INVALID_API_KEY", Message: "Invalid API key"}
    ErrRateLimitExceeded     = &ProviderError{Code: "RATE_LIMIT_EXCEEDED", Message: "Rate limit exceeded"}
    ErrInsufficientQuota     = &ProviderError{Code: "INSUFFICIENT_QUOTA", Message: "Insufficient quota"}
    ErrRequestTimeout        = &ProviderError{Code: "REQUEST_TIMEOUT", Message: "Request timeout"}
    ErrInvalidRequest        = &ProviderError{Code: "INVALID_REQUEST", Message: "Invalid request parameters"}
    ErrServiceUnavailable    = &ProviderError{Code: "SERVICE_UNAVAILABLE", Message: "Service temporarily unavailable"}
    ErrInternalServerError   = &ProviderError{Code: "INTERNAL_SERVER_ERROR", Message: "Internal server error"}
    ErrVideoNotSupported     = &ProviderError{Code: "VIDEO_NOT_SUPPORTED", Message: "Video generation not supported"}
    ErrImageNotSupported     = &ProviderError{Code: "IMAGE_NOT_SUPPORTED", Message: "Image generation not supported"}
    ErrTextNotSupported      = &ProviderError{Code: "TEXT_NOT_SUPPORTED", Message: "Text generation not supported"}
    ErrModelNotSupported     = &ProviderError{Code: "MODEL_NOT_SUPPORTED", Message: "Model not supported by this provider"}
    ErrContentFiltered       = &ProviderError{Code: "CONTENT_FILTERED", Message: "Content was filtered"}
)

// ProviderError خطأ خاص بالمزود
type ProviderError struct {
    Code        string `json:"code"`
    Message     string `json:"message"`
    Provider    string `json:"provider,omitempty"`
    Details     string `json:"details,omitempty"`
    StatusCode  int    `json:"status_code,omitempty"`
}

func (e *ProviderError) Error() string {
    if e.Provider != "" {
        return fmt.Sprintf("[%s] %s: %s", e.Provider, e.Code, e.Message)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// ValidationError خطأ في التحقق
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// CostError خطأ في التكلفة
type CostError struct {
    Message     string  `json:"message"`
    Required    float64 `json:"required"`
    Available   float64 `json:"available"`
    UserID      string  `json:"user_id,omitempty"`
}

func (e *CostError) Error() string {
    return fmt.Sprintf("cost error: %s (required: %.4f, available: %.4f)", e.Message, e.Required, e.Available)
}

// ErrorMap خريطة الأخطاء
var ErrorMap = map[string]*ProviderError{
    "NO_PROVIDER_AVAILABLE":   ErrNoProviderAvailable,
    "PROVIDER_NOT_FOUND":      ErrProviderNotFound,
    "PROVIDER_UNAVAILABLE":    ErrProviderUnavailable,
    "INVALID_API_KEY":         ErrInvalidAPIKey,
    "RATE_LIMIT_EXCEEDED":     ErrRateLimitExceeded,
    "INSUFFICIENT_QUOTA":      ErrInsufficientQuota,
    "REQUEST_TIMEOUT":         ErrRequestTimeout,
    "INVALID_REQUEST":         ErrInvalidRequest,
    "SERVICE_UNAVAILABLE":     ErrServiceUnavailable,
    "INTERNAL_SERVER_ERROR":   ErrInternalServerError,
    "VIDEO_NOT_SUPPORTED":     ErrVideoNotSupported,
    "IMAGE_NOT_SUPPORTED":     ErrImageNotSupported,
    "TEXT_NOT_SUPPORTED":      ErrTextNotSupported,
    "MODEL_NOT_SUPPORTED":     ErrModelNotSupported,
    "CONTENT_FILTERED":        ErrContentFiltered,
}