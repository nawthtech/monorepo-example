package services

import (
    "context"
    "github.com/nawthtech/nawthtech/backend/internal/ai/types"
)

type ContentService struct {
    textProvider types.TextProvider
}

func NewContentService(provider types.TextProvider) *ContentService {
    return &ContentService{
        textProvider: provider,
    }
}

func (s *ContentService) GenerateBlogPost(ctx context.Context, topic string, options map[string]interface{}) (string, error) {
    prompt := "Write a blog post about: " + topic
    return s.textProvider.GenerateText(ctx, prompt, options)
}
