package providers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "time"
)

type HuggingFaceProvider struct {
    apiToken string
    client   *http.Client
}

func NewHuggingFaceProvider() *HuggingFaceProvider {
    return &HuggingFaceProvider{
        apiToken: os.Getenv("HUGGINGFACE_TOKEN"),
        client: &http.Client{Timeout: 120 * time.Second},
    }
}

func (p *HuggingFaceProvider) GenerateText(model, prompt string) (string, error) {
    url := fmt.Sprintf("https://api-inference.huggingface.co/models/%s", model)
    
    payload := map[string]interface{}{
        "inputs": prompt,
        "parameters": map[string]interface{}{
            "max_new_tokens": 500,
            "temperature":    0.7,
        },
    }
    
    jsonData, _ := json.Marshal(payload)
    
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    req.Header.Set("Authorization", "Bearer "+p.apiToken)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := p.client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    
    var result []map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return "", err
    }
    
    if len(result) > 0 {
        if text, ok := result[0]["generated_text"].(string); ok {
            return text, nil
        }
    }
    
    return "", fmt.Errorf("no text generated")
}

// نماذج مجانية على HuggingFace:
// - "google/flan-t5-xl"
// - "microsoft/phi-2"
// - "mistralai/Mistral-7B-Instruct-v0.2"
// - "Qwen/Qwen2.5-7B-Instruct"