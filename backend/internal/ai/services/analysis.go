package services

import (
    "context"
    "github.com/nawthtech/nawthtech/backend/internal/ai/types"
)

type AnalysisService struct {
    textProvider types.TextProvider
}

func NewAnalysisService(provider types.TextProvider) *AnalysisService {
    return &AnalysisService{
        textProvider: provider,
    }
}

func (s *AnalysisService) AnalyzeMarketTrends(ctx context.Context, industry string, timeframe string) (string, error) {
    prompt := "Analyze market trends for " + industry + " for " + timeframe
    return s.textProvider.GenerateText(ctx, prompt, nil)
}
