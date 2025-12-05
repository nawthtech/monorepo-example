package models

// TextModel نموذج نصي
type TextModel struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Provider    string   `json:"provider"`
    MaxTokens   int      `json:"max_tokens"`
    CostPer1K   float64  `json:"cost_per_1k"` // 0 = مجاني
    Languages   []string `json:"languages"`
    Capabilities []string `json:"capabilities"`
    IsLocal     bool     `json:"is_local"`
    IsAvailable bool     `json:"is_available"`
}

// النماذج النصية المجانية
var FreeTextModels = []TextModel{
    {
        ID:        "gemini-2.0-flash",
        Name:      "Gemini 2.0 Flash",
        Provider:  "Google",
        MaxTokens: 8192,
        CostPer1K: 0.0,
        Languages: []string{"en", "ar", "fr", "es", "de"},
        Capabilities: []string{
            "text_generation",
            "translation",
            "summarization",
            "question_answering",
        },
        IsLocal:     false,
        IsAvailable: true,
    },
    {
        ID:        "llama-3.2-3b",
        Name:      "Llama 3.2 3B",
        Provider:  "Meta (via Ollama)",
        MaxTokens: 4096,
        CostPer1K: 0.0,
        Languages: []string{"en", "ar", "es"},
        Capabilities: []string{
            "text_generation",
            "summarization",
            "code_generation",
        },
        IsLocal:     true,
        IsAvailable: true,
    },
    {
        ID:        "mistral-7b",
        Name:      "Mistral 7B",
        Provider:  "Mistral AI",
        MaxTokens: 32768,
        CostPer1K: 0.0,
        Languages: []string{"en", "fr", "es", "de", "it"},
        Capabilities: []string{
            "text_generation",
            "translation",
            "summarization",
        },
        IsLocal:     true,
        IsAvailable: true,
    },
    {
        ID:        "qwen2.5-7b",
        Name:      "Qwen 2.5 7B",
        Provider:  "Alibaba",
        MaxTokens: 32768,
        CostPer1K: 0.0,
        Languages: []string{"en", "zh", "ar", "fr", "es"},
        Capabilities: []string{
            "text_generation",
            "translation",
            "code_generation",
            "reasoning",
        },
        IsLocal:     true,
        IsAvailable: true,
    },
    {
        ID:        "phi-3-mini",
        Name:      "Phi-3 Mini",
        Provider:  "Microsoft",
        MaxTokens: 4096,
        CostPer1K: 0.0,
        Languages: []string{"en", "es", "fr", "de"},
        Capabilities: []string{
            "text_generation",
            "summarization",
            "instruction_following",
        },
        IsLocal:     true,
        IsAvailable: true,
    },
}