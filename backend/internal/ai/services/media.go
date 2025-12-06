package services

import (
    "context"
    "github.com/nawthtech/nawthtech/backend/internal/ai/types"
)

type MediaService struct {
    imageProvider types.ImageProvider
    videoProvider types.VideoProvider
}

func NewMediaService(imageProvider types.ImageProvider, videoProvider types.VideoProvider) *MediaService {
    return &MediaService{
        imageProvider: imageProvider,
        videoProvider: videoProvider,
    }
}

func (s *MediaService) GenerateSocialMediaImage(ctx context.Context, platform string, prompt string, style string) ([]byte, error) {
    options := map[string]interface{}{
        "platform": platform,
        "style":    style,
    }
    return s.imageProvider.GenerateImage(ctx, prompt, options)
}
