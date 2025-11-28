package models

import "time"

// Content نموذج المحتوى
type Content struct {
	ID           string                 `json:"id"`
	Topic        string                 `json:"topic"`
	Content      string                 `json:"content"`
	Platform     string                 `json:"platform"`
	ContentType  string                 `json:"content_type"`
	Tone         string                 `json:"tone"`
	Keywords     []string               `json:"keywords"`
	Language     string                 `json:"language"`
	Status       string                 `json:"status"`
	Analysis     *ContentAnalysis       `json:"analysis"`
	Optimization *PlatformOptimization  `json:"optimization"`
	Metadata     map[string]interface{} `json:"metadata"`
	Performance  *ContentPerformance    `json:"performance"`
	CreatedBy    string                 `json:"created_by"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// BatchContent محتوى جماعي
type BatchContent struct {
	ID        string        `json:"id"`
	Content   []Content     `json:"content"`
	Summary   *BatchSummary `json:"summary"`
	CreatedBy string        `json:"created_by"`
	CreatedAt time.Time     `json:"created_at"`
}

// BatchSummary ملخص المحتوى الجماعي
type BatchSummary struct {
	TotalContent      int      `json:"total_content"`
	Platforms         []string `json:"platforms"`
	AverageQuality    int      `json:"average_quality"`
	ConsistencyScore  int      `json:"consistency_score"`
}

// PlatformOptimization تحسين المحتوى للمنصة
type PlatformOptimization struct {
	Platform        string   `json:"platform"`
	ContentType     string   `json:"content_type"`
	OptimizedLength string   `json:"optimized_length"`
	Formatting      string   `json:"formatting"`
	Hashtags        []string `json:"hashtags"`
}

// ContentAnalysis تحليل المحتوى
type ContentAnalysis struct {
	ContentID        string   `json:"content_id"`
	AnalysisType     string   `json:"analysis_type"`
	SentimentScore   int      `json:"sentiment_score"`
	SEOScore         int      `json:"seo_score"`
	EngagementScore  int      `json:"engagement_score"`
	ReadabilityScore int      `json:"readability_score"`
	OverallScore     int      `json:"overall_score"`
	Recommendations  []string `json:"recommendations"`
	GeneratedAt      time.Time `json:"generated_at"`
}

// ContentOptimization تحسين المحتوى
type ContentOptimization struct {
	OriginalContent   string                 `json:"original_content"`
	OptimizedContent  string                 `json:"optimized_content"`
	Improvements      []string               `json:"improvements"`
	Metrics           map[string]interface{} `json:"metrics"`
	GeneratedAt       time.Time              `json:"generated_at"`
}

// ContentPerformance أداء المحتوى
type ContentPerformance struct {
	ContentID      string                 `json:"content_id"`
	Timeframe      string                 `json:"timeframe"`
	Metrics        map[string]interface{} `json:"metrics"`
	Analysis       *PerformanceAnalysis   `json:"analysis"`
	Recommendations []string              `json:"recommendations"`
	GeneratedAt    time.Time              `json:"generated_at"`
}

// PerformanceAnalysis تحليل الأداء
type PerformanceAnalysis struct {
	Patterns       []string `json:"patterns"`
	SuccessFactors []string `json:"success_factors"`
	Insights       []string `json:"insights"`
}