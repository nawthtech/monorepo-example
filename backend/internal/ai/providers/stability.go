package providers

import (
	"context"
)

type StabilityProvider struct {
	apiKey string
}

func NewStabilityProvider(apiKey string) *StabilityProvider {
	return &StabilityProvider{
		apiKey: apiKey,
	}
}

func (s *StabilityProvider) GenerateImage(ctx context.Context, prompt string, options map[string]interface{}) ([]byte, error) {
	return []byte("stability image data"), nil
}

func (s *StabilityProvider) Name() string {
	return "stability"
}
