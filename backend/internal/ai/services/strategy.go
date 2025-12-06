package services

import (
    "context"
    "github.com/nawthtech/nawthtech/backend/internal/ai/types"
)

type StrategyService struct {
    textProvider types.TextProvider
}

func NewStrategyService(provider types.TextProvider) *StrategyService {
    return &StrategyService{
        textProvider: provider,
    }
}

func (s *StrategyService) GenerateMarketingStrategy(ctx context.Context, product string) (string, error) {
    prompt := "Generate marketing strategy for: " + product
    return s.textProvider.GenerateText(ctx, prompt, nil)
}
