package services

import (
	"context"
	"github.com/nawthtech/nawthtech/backend/internal/ai/types"
)
package services

import (
    "bytes"
    "encoding/base64"
    "fmt"
    "image"
    "image/jpeg"
    "io"
)

type MediaService struct {
    imageProvider types.ImageProvider
    videoProvider types.VideoProvider
}

func NewMediaService(imageProv types.ImageProvider, videoProv types.VideoProvider) *MediaService {
    return &MediaService{
        imageProvider: imageProv,
        videoProvider: videoProv,
    }
}

// GenerateSocialMediaImage توليد صورة لوسائل التواصل
func (s *MediaService) GenerateSocialMediaImage(platform, theme, text string) ([]byte, error) {
    prompt := fmt.Sprintf(`
    Create a social media image for %s platform.
    
    Theme: %s
    Text to include: "%s"
    
    Style requirements:
    - Brand colors: purple (#7A3EF0) and neon cyan (#00F6FF)
    - Modern, futuristic design
    - Clean, professional layout
    - Optimized for %s dimensions
    - Include subtle AI/tech elements
    
    Make it eye-catching and shareable.
    `, platform, theme, text, platform)
    
    var width, height int
    switch platform {
    case "instagram":
        width, height = 1080, 1080
    case "twitter":
        width, height = 1200, 675
    case "linkedin":
        width, height = 1200, 627
    default:
        width, height = 1024, 1024
    }
    
    img, err := s.imageProvider.GenerateImage(prompt, width, height)
    if err != nil {
        return nil, err
    }
    
    // تحويل image إلى bytes
    var buf bytes.Buffer
    if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
        return nil, err
    }
    
    return buf.Bytes(), nil
}

// GenerateProductImage توليد صورة منتج
func (s *MediaService) GenerateProductImage(productName, description string) ([]byte, error) {
    prompt := fmt.Sprintf(`
    Create a professional product image for: %s
    
    Product description: %s
    
    Style:
    - Clean white background
    - Professional product photography
    - 3D render style
    - Modern, minimalist
    - Good lighting and shadows
    - Product should be centered
    
    Make it look premium and trustworthy.
    `, productName, description)
    
    img, err := s.imageProvider.GenerateImage(prompt, 1024, 1024)
    if err != nil {
        return nil, err
    }
    
    var buf bytes.Buffer
    if err := jpeg.Encode(&buf, img, nil); err != nil {
        return nil, err
    }
    
    return buf.Bytes(), nil
}

// GenerateShortVideo توليد فيديو قصير
func (s *MediaService) GenerateShortVideo(topic, style string, duration int) ([]byte, error) {
    prompt := fmt.Sprintf(`
    Create a %d-second %s style video about: %s
    
    Requirements:
    - Smooth animation
    - Professional quality
    - Engaging visuals
    - Clear messaging
    - Optimized for social media
    
    Include text overlays if needed.
    `, duration, style, topic)
    
    return s.videoProvider.GenerateVideo(prompt, duration)
}