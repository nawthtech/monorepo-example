package models

// ImageModel نموذج صور
type ImageModel struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Provider     string   `json:"provider"`
	Resolution   string   `json:"resolution"` // 512x512, 1024x1024
	CostPerImage float64  `json:"cost_per_image"`
	Styles       []string `json:"styles"`
	IsLocal      bool     `json:"is_local"`
}

// النماذج المجانية لتوليد الصور
var FreeImageModels = []ImageModel{
	{
		ID:           "stable-diffusion-xl",
		Name:         "Stable Diffusion XL",
		Provider:     "Stability AI",
		Resolution:   "1024x1024",
		CostPerImage: 0.0,
		Styles: []string{
			"realistic",
			"anime",
			"digital-art",
			"photographic",
		},
		IsLocal: true,
	},
	{
		ID:           "playground-v2",
		Name:         "Playground v2",
		Provider:     "Playground AI",
		Resolution:   "1024x1024",
		CostPerImage: 0.0,
		Styles: []string{
			"realistic",
			"cinematic",
			"anime",
			"3d-render",
		},
		IsLocal: false, // API مجاني محدود
	},
	{
		ID:           "flux-dev",
		Name:         "FLUX.1 Dev",
		Provider:     "Black Forest Labs",
		Resolution:   "1024x1024",
		CostPerImage: 0.0,
		Styles: []string{
			"realistic",
			"illustration",
			"photography",
		},
		IsLocal: true,
	},
	{
		ID:           "dall-e-2",
		Name:         "DALL-E 2",
		Provider:     "OpenAI Compatible",
		Resolution:   "1024x1024",
		CostPerImage: 0.0,
		Styles: []string{
			"digital-art",
			"photo",
			"painting",
		},
		IsLocal: false,
	},
}
