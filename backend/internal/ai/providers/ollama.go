package providers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type OllamaProvider struct {
    baseURL string
    client  *http.Client
}

func NewOllamaProvider() *OllamaProvider {
    baseURL := "http://localhost:11434"
    return &OllamaProvider{
        baseURL: baseURL,
        client: &http.Client{Timeout: 300 * time.Second},
    }
}

func (p *OllamaProvider) GenerateText(model, prompt string) (string, error) {
    url := p.baseURL + "/api/generate"
    
    payload := map[string]interface{}{
        "model":  model,
        "prompt": prompt,
        "stream": false,
        "options": map[string]interface{}{
            "temperature": 0.7,
            "num_predict": 2000,
        },
    }
    
    jsonData, _ := json.Marshal(payload)
    
    resp, err := p.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("Ollama connection failed: %w", err)
    }
    defer resp.Body.Close()
    
    var result struct {
        Response string `json:"response"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", err
    }
    
    return result.Response, nil
}

// نماذج Ollama المجانية:
// - "llama3.2:3b"        // 3B parameters
// - "mistral:7b"         // 7B
// - "qwen2.5:7b"         // 7B, دعم عربي جيد
// - "phi3:mini"          // 3.8B
// - "gemma:7b"           // من Google