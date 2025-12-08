package providers

import (
	"context"
)

type GeminiProvider struct {
	apiKey string
}

func NewGeminiProvider(apiKey string) *GeminiProvider {
	return &GeminiProvider{
		apiKey: apiKey,
	}
}

func (g *GeminiProvider) GenerateText(ctx context.Context, prompt string, options map[string]interface{}) (string, error) {
	return "Gemini generated text for: " + prompt, nil
}

func (g *GeminiProvider) Name() string {
	return "gemini"
}
