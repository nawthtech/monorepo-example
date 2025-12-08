package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required")
	}

	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// اختبار توليد النص
	fmt.Println("=== توليد نص ===")
	model := client.GenerativeModel("gemini-pro")

	resp, err := model.GenerateContent(ctx,
		genai.Text("Explain how AI works in a few words"),
	)
	if err != nil {
		log.Fatal("Failed to generate text:", err)
	}

	if resp != nil && len(resp.Candidates) > 0 {
		for _, cand := range resp.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					fmt.Println(part)
				}
			}
		}
	}

	fmt.Println("\n✅ Test completed successfully")
}
