package services

import (
    "context"
    "github.com/nawthtech/nawthtech/backend/internal/ai/types"
)

type TranslationService struct {
    textProvider types.TextProvider
}

func NewTranslationService(provider types.TextProvider) *TranslationService {
    return &TranslationService{
        textProvider: provider,
    }
}

func (s *TranslationService) Translate(ctx context.Context, text string, targetLang string) (string, error) {
    prompt := "Translate to " + targetLang + ": " + text
    return s.textProvider.GenerateText(ctx, prompt, nil)
}
