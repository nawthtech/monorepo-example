// backend/internal/ai/verification/verifier.go
package verification

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"github.com/sashabaranov/go-openai"
	"github.com/nawthtech/nawthtech/backend/internal/config"
)

// Verifier محقق المحتوى باستخدام LLM
type Verifier struct {
	client      *openai.Client
	logger      *logrus.Logger
	config      *config.Config
	model       string
	maxRetries  int
	timeout     time.Duration
	temperature float32
	maxTokens   int
	criteria    VerificationCriteria
	sentryHub   *sentry.Hub
}

// NewVerifier إنشاء محقق جديد
func NewVerifier(cfg *config.Config, logger *logrus.Logger) (*Verifier, error) {
	if cfg.AI.OpenAIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	client := openai.NewClient(cfg.AI.OpenAIKey)
	
	verifier := &Verifier{
		client:      client,
		logger:      logger,
		config:      cfg,
		model:       cfg.AI.DefaultModel,
		maxRetries:  cfg.AI.MaxRetries,
		timeout:     time.Duration(cfg.AI.TimeoutMS) * time.Millisecond,
		temperature: cfg.AI.Temperature,
		maxTokens:   cfg.AI.MaxTokens,
		criteria: VerificationCriteria{
			Toxicity:   cfg.Verification.CheckToxicity,
			Factuality: cfg.Verification.CheckFactuality,
			Coherence:  cfg.Verification.CheckCoherence,
			Relevance:  cfg.Verification.CheckRelevance,
			Safety:     cfg.Verification.CheckSafety,
			Moderation: cfg.Verification.CheckModeration,
			Bias:       cfg.Verification.CheckBias,
		},
	}

	// Initialize Sentry hub
	if cfg.Sentry.DSN != "" {
		verifier.sentryHub = sentry.CurrentHub().Clone()
	}

	return verifier, nil
}

// Verify يحقق من محتوى باستخدام LLM
func (v *Verifier) Verify(ctx context.Context, input string, opts ...VerificationOption) (*VerificationResult, error) {
	// بدء معاملة Sentry
	var span *sentry.Span
	if v.sentryHub != nil {
		span = sentry.StartSpan(ctx, "llm.verification")
		defer span.Finish()
	}

	startTime := time.Now()

	// إعداد الخيارات
	options := v.getVerificationOptions(opts...)

	// تسجيل بدء التحقق
	v.logger.WithFields(logrus.Fields{
		"input_length": len(input),
		"model":        v.model,
		"criteria":     v.criteria,
	}).Info("Starting LLM verification")

	// بناء الـ prompt
	prompt := v.buildVerificationPrompt(input, options)

	// استدعاء LLM مع إعادة المحاولة
	response, err := v.callLLMWithRetry(ctx, prompt, options)
	if err != nil {
		v.logger.WithError(err).Error("LLM verification failed")
		
		// تسجيل الخطأ في Sentry
		if v.sentryHub != nil {
			v.sentryHub.CaptureException(err)
		}
		
		return v.handleVerificationError(err, input, options, startTime), nil
	}

	// تحليل الاستجابة
	result := v.parseVerificationResult(response, input, startTime)

	// تسجيل النتيجة
	v.logger.WithFields(logrus.Fields{
		"isValid":    result.IsValid,
		"confidence": result.Confidence,
		"latency":    result.Metrics.Latency,
	}).Info("LLM verification completed")

	// إرسال إلى Sentry
	if v.sentryHub != nil {
		v.sendToSentry(result, input, options)
	}

	return result, nil
}

// buildVerificationPrompt بناء prompt التحقق
func (v *Verifier) buildVerificationPrompt(input string, options map[string]interface{}) string {
	var prompt strings.Builder

	prompt.WriteString("Please verify the following content and provide a structured JSON response.\n\n")
	
	if context, ok := options["context"].(string); ok && context != "" {
		prompt.WriteString(fmt.Sprintf("Context: %s\n\n", context))
	}
	
	prompt.WriteString(fmt.Sprintf("Content to verify: \"%s\"\n\n", input))
	prompt.WriteString("Verification Criteria (check all that apply):\n")
	
	if v.criteria.Toxicity {
		prompt.WriteString("- Toxicity: Check for toxic, hateful, harmful, or abusive content\n")
	}
	if v.criteria.Factuality {
		prompt.WriteString("- Factuality: Verify factual accuracy, check for misinformation\n")
	}
	if v.criteria.Coherence {
		prompt.WriteString("- Coherence: Check logical flow, consistency, and clarity\n")
	}
	if v.criteria.Relevance {
		prompt.WriteString("- Relevance: Check if content is relevant to the intended topic\n")
	}
	if v.criteria.Safety {
		prompt.WriteString("- Safety: Check for dangerous, illegal, or policy-violating content\n")
	}
	if v.criteria.Moderation {
		prompt.WriteString("- Moderation: Check for inappropriate, explicit, or offensive content\n")
	}
	if v.criteria.Bias {
		prompt.WriteString("- Bias: Check for political, racial, gender, or other biases\n")
	}
	
	prompt.WriteString("\nRespond with a JSON object in this exact format:\n")
	prompt.WriteString(`{
  "isValid": boolean,
  "confidence": number between 0 and 1,
  "reason": "brief explanation",
  "issues": ["specific issue 1", "specific issue 2"],
  "suggestions": ["suggestion 1", "suggestion 2"],
  "categories": {
    "toxicity": {"passed": boolean, "score": number, "explanation": "details"},
    "factuality": {"passed": boolean, "score": number, "explanation": "details"},
    "coherence": {"passed": boolean, "score": number, "explanation": "details"},
    "relevance": {"passed": boolean, "score": number, "explanation": "details"},
    "safety": {"passed": boolean, "score": number, "explanation": "details"},
    "moderation": {"passed": boolean, "score": number, "explanation": "details"},
    "bias": {"passed": boolean, "score": number, "explanation": "details"}
  }
}`)
	
	prompt.WriteString("\n\nImportant: Only include categories that were requested. Return \"passed\": false and score 0 for categories not checked.")
	
	return prompt.String()
}

// callLLMWithRetry استدعاء LLM مع إعادة المحاولة
func (v *Verifier) callLLMWithRetry(ctx context.Context, prompt string, options map[string]interface{}, retryCount ...int) (string, error) {
	currentRetry := 0
	if len(retryCount) > 0 {
		currentRetry = retryCount[0]
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, v.timeout)
	defer cancel()

	req := openai.ChatCompletionRequest{
		Model:       v.getModelFromOptions(options),
		Messages:    []openai.ChatCompletionMessage{{Role: "user", Content: prompt}},
		Temperature: v.getTemperatureFromOptions(options),
		MaxTokens:   v.getMaxTokensFromOptions(options),
	}

	resp, err := v.client.CreateChatCompletion(ctxWithTimeout, req)
	if err != nil {
		if currentRetry < v.maxRetries && v.isRetryableError(err) {
			v.logger.WithError(err).Warnf("LLM call failed, retry %d/%d", currentRetry+1, v.maxRetries)
			
			// انتظار أسي
			backoff := time.Duration(1<<uint(currentRetry)) * time.Second
			time.Sleep(backoff)
			
			return v.callLLMWithRetry(ctx, prompt, options, currentRetry+1)
		}
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	return resp.Choices[0].Message.Content, nil
}

// parseVerificationResult تحليل استجابة LLM
func (v *Verifier) parseVerificationResult(response, input string, startTime time.Time) *VerificationResult {
	// محاولة تحليل JSON
	var parsedResult map[string]interface{}
	if err := json.Unmarshal([]byte(response), &parsedResult); err == nil {
		return v.parseStructuredResult(parsedResult, input, startTime)
	}

	// التحليل اليدوي إذا فشل تحليل JSON
	return v.parseUnstructuredResult(response, input, startTime)
}

// parseStructuredResult تحليل النتيجة المهيكلة
func (v *Verifier) parseStructuredResult(data map[string]interface{}, input string, startTime time.Time) *VerificationResult {
	result := &VerificationResult{
		IsValid:    false,
		Confidence: 0,
		Reason:     "Verification inconclusive",
		Issues:     []string{},
		Suggestions: []string{},
		Categories: make(map[string]CategoryResult),
		Metrics: VerificationMetrics{
			Latency: time.Since(startTime).Milliseconds(),
			Model:   v.model,
			Provider: "openai",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// استخراج isValid
	if val, ok := data["isValid"].(bool); ok {
		result.IsValid = val
	}

	// استخراج confidence
	if val, ok := data["confidence"].(float64); ok {
		result.Confidence = v.normalizeScore(val)
	}

	// استخراج reason
	if val, ok := data["reason"].(string); ok {
		result.Reason = val
	}

	// استخراج issues
	if issues, ok := data["issues"].([]interface{}); ok {
		for _, issue := range issues {
			if str, ok := issue.(string); ok {
				result.Issues = append(result.Issues, str)
			}
		}
	}

	// استخراج suggestions
	if suggestions, ok := data["suggestions"].([]interface{}); ok {
		for _, suggestion := range suggestions {
			if str, ok := suggestion.(string); ok {
				result.Suggestions = append(result.Suggestions, str)
			}
		}
	}

	// استخراج categories
	if categories, ok := data["categories"].(map[string]interface{}); ok {
		for key, val := range categories {
			if categoryData, ok := val.(map[string]interface{}); ok {
				category := CategoryResult{
					Passed:      false,
					Score:       0,
					Explanation: "",
				}

				if passed, ok := categoryData["passed"].(bool); ok {
					category.Passed = passed
				}
				if score, ok := categoryData["score"].(float64); ok {
					category.Score = v.normalizeScore(score)
				}
				if explanation, ok := categoryData["explanation"].(string); ok {
					category.Explanation = explanation
				}

				result.Categories[key] = category
			}
		}
	}

	return result
}

// parseUnstructuredResult تحليل النتيجة غير المهيكلة
func (v *Verifier) parseUnstructuredResult(text, input string, startTime time.Time) *VerificationResult {
	result := &VerificationResult{
		IsValid:    false,
		Confidence: 0.5,
		Reason:     "Automated analysis",
		Issues:     []string{},
		Suggestions: []string{},
		Categories: make(map[string]CategoryResult),
		Metrics: VerificationMetrics{
			Latency: time.Since(startTime).Milliseconds(),
			Model:   v.model,
			Provider: "openai",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// تحديد isValid
	positiveWords := []string{"valid", "passed", "ok", "good", "safe", "appropriate"}
	negativeWords := []string{"invalid", "failed", "bad", "unsafe", "inappropriate", "toxic"}

	textLower := strings.ToLower(text)
	
	isPositive := false
	isNegative := false

	for _, word := range positiveWords {
		if strings.Contains(textLower, word) {
			isPositive = true
			break
		}
	}

	for _, word := range negativeWords {
		if strings.Contains(textLower, word) {
			isNegative = true
			break
		}
	}

	result.IsValid = isPositive && !isNegative

	// استخراج confidence باستخدام regex
	re := regexp.MustCompile(`confidence.*?(\d+\.?\d*)`)
	if matches := re.FindStringSubmatch(textLower); len(matches) > 1 {
		if confidence, err := fmt.Sscanf(matches[1], "%f", &result.Confidence); err == nil && confidence > 0 {
			result.Confidence = v.normalizeScore(result.Confidence)
		}
	}

	// استخراج issues
	issueRe := regexp.MustCompile(`(?i)(issues?|problems?|concerns?)[:\s]+([^.\n]+)`)
	if matches := issueRe.FindAllStringSubmatch(text, -1); matches != nil {
		for _, match := range matches {
			if len(match) > 2 {
				result.Issues = append(result.Issues, strings.TrimSpace(match[2]))
			}
		}
	}

	// استخراج suggestions
	suggestionRe := regexp.MustCompile(`(?i)(suggestions?|recommendations?|improvements?)[:\s]+([^.\n]+)`)
	if matches := suggestionRe.FindAllStringSubmatch(text, -1); matches != nil {
		for _, match := range matches {
			if len(match) > 2 {
				result.Suggestions = append(result.Suggestions, strings.TrimSpace(match[2]))
			}
		}
	}

	// إذا لم توجد issues
	if len(result.Issues) == 0 {
		result.Issues = []string{"No specific issues detected"}
	}

	// إذا لم توجد suggestions
	if len(result.Suggestions) == 0 {
		result.Suggestions = []string{"No suggestions provided"}
	}

	return result
}

// handleVerificationError معالجة أخطاء التحقق
func (v *Verifier) handleVerificationError(err error, input string, options map[string]interface{}, startTime time.Time) *VerificationResult {
	return &VerificationResult{
		IsValid:    false,
		Confidence: 0,
		Reason:     fmt.Sprintf("Verification failed: %v", err),
		Issues:     []string{"LLM API error", err.Error()},
		Suggestions: []string{"Please try again later", "Check API key configuration"},
		Categories: make(map[string]CategoryResult),
		Metrics: VerificationMetrics{
			Latency: time.Since(startTime).Milliseconds(),
			Model:   v.getModelFromOptions(options),
			Provider: "openai",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		Error: err.Error(),
		Metadata: map[string]interface{}{
			"errorType":  fmt.Sprintf("%T", err),
			"retryable":  v.isRetryableError(err),
			"inputLength": len(input),
		},
	}
}

// sendToSentry إرسال البيانات إلى Sentry
func (v *Verifier) sendToSentry(result *VerificationResult, input string, options map[string]interface{}) {
	if v.sentryHub == nil {
		return
	}

	v.sentryHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetExtra("verification_result", result)
		scope.SetTag("verification.success", fmt.Sprintf("%v", result.IsValid))
		scope.SetTag("verification.type", v.getTypeFromOptions(options))
		
		if result.Error != "" {
			scope.SetTag("verification.error", "true")
		}
	})

	// تسجيل حدث مخصص
	v.sentryHub.CaptureMessage("LLM Verification Completed", &sentry.Event{
		Level:   sentry.LevelInfo,
		Message: "LLM verification completed",
		Extra: map[string]interface{}{
			"input_length":       len(input),
			"verification_result": result,
			"model":              v.getModelFromOptions(options),
		},
	})
}

// Helper methods
func (v *Verifier) getVerificationOptions(opts ...VerificationOption) map[string]interface{} {
	options := map[string]interface{}{
		"model":       v.model,
		"temperature": v.temperature,
		"maxTokens":   v.maxTokens,
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

func (v *Verifier) getModelFromOptions(options map[string]interface{}) string {
	if model, ok := options["model"].(string); ok && model != "" {
		return model
	}
	return v.model
}

func (v *Verifier) getTemperatureFromOptions(options map[string]interface{}) float32 {
	if temp, ok := options["temperature"].(float32); ok {
		return temp
	}
	return v.temperature
}

func (v *Verifier) getMaxTokensFromOptions(options map[string]interface{}) int {
	if tokens, ok := options["maxTokens"].(int); ok && tokens > 0 {
		return tokens
	}
	return v.maxTokens
}

func (v *Verifier) getTypeFromOptions(options map[string]interface{}) string {
	if typ, ok := options["type"].(string); ok && typ != "" {
		return typ
	}
	return "general"
}

func (v *Verifier) normalizeScore(score float64) float64 {
	if score <= 0 {
		return 0
	}
	if score >= 1 {
		return 1
	}
	if score > 100 {
		return score / 100
	}
	if score > 10 {
		return score / 10
	}
	if score > 5 {
		return score / 5
	}
	return score
}

func (v *Verifier) isRetryableError(err error) bool {
	// تحقق من أخطاء الشبكة
	if strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "connection") ||
		strings.Contains(err.Error(), "network") {
		return true
	}

	// تحقق من أخطاء rate limit
	if strings.Contains(err.Error(), "rate limit") ||
		strings.Contains(err.Error(), "too many requests") {
		return true
	}

	// تحقق من أخطاء الخادم
	if strings.Contains(err.Error(), "500") ||
		strings.Contains(err.Error(), "502") ||
		strings.Contains(err.Error(), "503") ||
		strings.Contains(err.Error(), "504") {
		return true
	}

	return false
}

// BatchVerify التحقق الدفعي
func (v *Verifier) BatchVerify(ctx context.Context, inputs []string, opts ...VerificationOption) ([]*VerificationResult, error) {
	results := make([]*VerificationResult, 0, len(inputs))
	batchSize := 3

	for i := 0; i < len(inputs); i += batchSize {
		end := i + batchSize
		if end > len(inputs) {
			end = len(inputs)
		}

		batch := inputs[i:end]
		batchResults := make([]*VerificationResult, len(batch))

		for j, input := range batch {
			result, err := v.Verify(ctx, input, opts...)
			if err != nil {
				v.logger.WithError(err).Error("Batch verification failed for item")
				result = v.handleVerificationError(err, input, v.getVerificationOptions(opts...), time.Now())
			}
			batchResults[j] = result
		}

		results = append(results, batchResults...)

		// تأخير بين الدفعات لمنع rate limiting
		if i+batchSize < len(inputs) {
			time.Sleep(1 * time.Second)
		}
	}

	return results, nil
}

// TestConnection اختبار اتصال LLM
func (v *Verifier) TestConnection(ctx context.Context) (bool, error) {
	testPrompt := "Hello, please respond with 'OK'"
	
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := openai.ChatCompletionRequest{
		Model:    v.model,
		Messages: []openai.ChatCompletionMessage{{Role: "user", Content: testPrompt}},
		MaxTokens: 10,
	}

	resp, err := v.client.CreateChatCompletion(ctxWithTimeout, req)
	if err != nil {
		return false, err
	}

	if len(resp.Choices) == 0 {
		return false, fmt.Errorf("no response from LLM")
	}

	responseText := strings.ToLower(resp.Choices[0].Message.Content)
	return strings.Contains(responseText, "ok"), nil
}