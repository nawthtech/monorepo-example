package providers

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "image"
    "image/jpeg"
    "io"
    "net/http"
    "os"
)

type StabilityProvider struct {
    apiKey  string
    client  *http.Client
}

func NewStabilityProvider() *StabilityProvider {
    return &StabilityProvider{
        apiKey: os.Getenv("STABILITY_API_KEY"),
        client: &http.Client{},
    }
}

func (p *StabilityProvider) GenerateImage(prompt string, width, height int) (image.Image, error) {
    url := "https://api.stability.ai/v2beta/stable-image/generate/core"
    
    payload := map[string]interface{}{
        "prompt": prompt,
        "output_format": "jpeg",
        "width": width,
        "height": height,
    }
    
    jsonData, _ := json.Marshal(payload)
    
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    req.Header.Set("Authorization", "Bearer "+p.apiKey)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := p.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    
    var result struct {
        Image string `json:"image"`
    }
    
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }
    
    // تحويل base64 إلى image
    imgData, err := base64.StdEncoding.DecodeString(result.Image)
    if err != nil {
        return nil, err
    }
    
    img, _, err := image.Decode(bytes.NewReader(imgData))
    return img, err
}

// مجاني: 25 صورة/شهر مع API key