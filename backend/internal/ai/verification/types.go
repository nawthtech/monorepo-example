// backend/internal/ai/verification/types.go
package verification

import "time"

// VerificationCriteria معايير التحقق
type VerificationCriteria struct {
	Toxicity   bool `json:"toxicity"`
	Factuality bool `json:"factuality"`
	Coherence  bool `json:"coherence"`
	Relevance  bool `json:"relevance"`
	Safety     bool `json:"safety"`
	Moderation bool `json:"moderation"`
	Bias       bool `json:"bias"`
}

// VerificationResult نتيجة التحقق
type VerificationResult struct {
	IsValid     bool                      `json:"isValid"`
	Confidence  float64                   `json:"confidence"`
	Reason      string                    `json:"reason"`
	Issues      []string                  `json:"issues"`
	Suggestions []string                  `json:"suggestions"`
	Categories  map[string]CategoryResult `json:"categories"`
	Metrics     VerificationMetrics       `json:"metrics"`
	Error       string                    `json:"error,omitempty"`
	Metadata    map[string]interface{}    `json:"metadata,omitempty"`
}

// CategoryResult نتيجة الفئة
type CategoryResult struct {
	Passed      bool    `json:"passed"`
	Score       float64 `json:"score"`
	Explanation string  `json:"explanation"`
	Details     []string `json:"details,omitempty"`
}

// VerificationMetrics مقاييس التحقق
type VerificationMetrics struct {
	Latency      int64     `json:"latency_ms"`
	InputTokens  int       `json:"input_tokens,omitempty"`
	OutputTokens int       `json:"output_tokens,omitempty"`
	TotalTokens  int       `json:"total_tokens,omitempty"`
	Cost         float64   `json:"cost,omitempty"`
	Model        string    `json:"model"`
	Provider     string    `json:"provider"`
	Timestamp    string    `json:"timestamp"`
}

// BatchVerificationResult نتيجة التحقق الدفعي
type BatchVerificationResult struct {
	Total            int                  `json:"total"`
	Valid            int                  `json:"valid"`
	Invalid          int                  `json:"invalid"`
	AverageConfidence float64             `json:"average_confidence"`
	TotalCost        float64             `json:"total_cost"`
	TotalTokens      int                 `json:"total_tokens"`
	Results          []*VerificationResult `json:"results"`
	Summary          map[string]int      `json:"summary"`
	ProcessingTime   time.Duration       `json:"processing_time"`
}

// VerificationOption خيارات التحقق
type VerificationOption func(map[string]interface{})

// WithModel تحديد النموذج
func WithModel(model string) VerificationOption {
	return func(opts map[string]interface{}) {
		opts["model"] = model
	}
}

// WithTemperature تحديد درجة الحرارة
func WithTemperature(temp float32) VerificationOption {
	return func(opts map[string]interface{}) {
		opts["temperature"] = temp
	}
}

// WithMaxTokens تحديد الحد الأقصى للـ Tokens
func WithMaxTokens(tokens int) VerificationOption {
	return func(opts map[string]interface{}) {
		opts["maxTokens"] = tokens
	}
}

// WithType تحديد نوع التحقق
func WithType(typ string) VerificationOption {
	return func(opts map[string]interface{}) {
		opts["type"] = typ
	}
}

// WithContext تحديد السياق
func WithContext(context string) VerificationOption {
	return func(opts map[string]interface{}) {
		opts["context"] = context
	}
}

// WithCustomCriteria تحديد معايير مخصصة
func WithCustomCriteria(criteria VerificationCriteria) VerificationOption {
	return func(opts map[string]interface{}) {
		opts["customCriteria"] = criteria
	}
}

// VerificationRequest طلب التحقق
type VerificationRequest struct {
	Content  string                 `json:"content" binding:"required"`
	Type     string                 `json:"type,omitempty"`
	Context  string                 `json:"context,omitempty"`
	Model    string                 `json:"model,omitempty"`
	Options  map[string]interface{} `json:"options,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// BatchVerificationRequest طلب التحقق الدفعي
type BatchVerificationRequest struct {
	Contents []string               `json:"contents" binding:"required"`
	Type     string                 `json:"type,omitempty"`
	Model    string                 `json:"model,omitempty"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

// VerificationResponse استجابة التحقق
type VerificationResponse struct {
	Success bool                `json:"success"`
	Data    *VerificationResult `json:"data,omitempty"`
	Error   string              `json:"error,omitempty"`
	Message string              `json:"message,omitempty"`
}

// BatchVerificationResponse استجابة التحقق الدفعي
type BatchVerificationResponse struct {
	Success bool                      `json:"success"`
	Data    *BatchVerificationResult  `json:"data,omitempty"`
	Error   string                    `json:"error,omitempty"`
	Message string                    `json:"message,omitempty"`
}

// ProviderConfig تكوين المزود
type ProviderConfig struct {
	Name     string            `json:"name"`
	APIKey   string            `json:"api_key"`
	BaseURL  string            `json:"base_url"`
	Models   map[string]string `json:"models"`
	Cost     CostConfig        `json:"cost"`
	Timeout  time.Duration     `json:"timeout"`
	MaxRetries int             `json:"max_retries"`
}

// CostConfig تكوين التكاليف
type CostConfig struct {
	InputTokens  float64 `json:"input_tokens_per_1k"`
	OutputTokens float64 `json:"output_tokens_per_1k"`
}

// Stats إحصائيات التحقق
type Stats struct {
	TotalVerifications  int64     `json:"total_verifications"`
	Successful          int64     `json:"successful"`
	Failed              int64     `json:"failed"`
	AverageConfidence   float64   `json:"average_confidence"`
	TotalCost           float64   `json:"total_cost"`
	TotalTokens         int64     `json:"total_tokens"`
	MostCommonIssues    map[string]int `json:"most_common_issues"`
	VerificationByType  map[string]int `json:"verification_by_type"`
	LastUpdated         time.Time `json:"last_updated"`
}

// HealthCheckResult نتيجة فحص الصحة
type HealthCheckResult struct {
	Status      string    `json:"status"`
	LLMProvider string    `json:"llm_provider"`
	Model       string    `json:"model"`
	Connected   bool      `json:"connected"`
	Latency     int64     `json:"latency_ms"`
	Timestamp   time.Time `json:"timestamp"`
	Message     string    `json:"message"`
}