package models

// VideoModel نموذج فيديو
type VideoModel struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Provider    string   `json:"provider"`
    MaxDuration int      `json:"max_duration"` // بالثواني
    CostPerSecond float64 `json:"cost_per_second"`
    IsLocal     bool     `json:"is_local"`
    Limitations string   `json:"limitations"`
}

// النماذج المجانية لتوليد الفيديو
var FreeVideoModels = []VideoModel{
    {
        ID:          "stable-video-diffusion",
        Name:        "Stable Video Diffusion",
        Provider:    "Stability AI",
        MaxDuration: 4,
        CostPerSecond: 0.0,
        IsLocal:     true,
        Limitations: "Image-to-video only, 4 seconds max",
    },
    {
        ID:          "modelscope-t2v",
        Name:        "ModelScope T2V",
        Provider:    "Alibaba",
        MaxDuration: 3,
        CostPerSecond: 0.0,
        IsLocal:     true,
        Limitations: "3 seconds max, lower quality",
    },
    {
        ID:          "zeroscope-v2",
        Name:        "Zeroscope v2",
        Provider:    "Community",
        MaxDuration: 3,
        CostPerSecond: 0.0,
        IsLocal:     true,
        Limitations: "576x320 resolution, 3 seconds",
    },
    {
        ID:          "pika-labs-free",
        Name:        "Pika Labs Free",
        Provider:    "Pika Labs",
        MaxDuration: 3,
        CostPerSecond: 0.0,
        IsLocal:     false,
        Limitations: "30 generations/week, 3 seconds max",
    },
}