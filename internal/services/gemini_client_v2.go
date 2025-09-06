package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/genai"
)

type GeminiClientV2 struct {
	client *genai.Client
}

func NewGeminiClientV2(apiKey string) *GeminiClientV2 {
	ctx := context.Background()

	// Create client with API key
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatalf("Failed to create Genai client: %v", err)
	}

	return &GeminiClientV2{
		client: client,
	}
}

// GenerateEmbedding generates embedding for given text using official SDK
func (g *GeminiClientV2) GenerateEmbedding(text string) ([]float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Clean text
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	// Create content for embedding
	contents := []*genai.Content{
		genai.NewContentFromText(text, genai.RoleUser),
	}

	// Generate embedding using official SDK
	fmt.Printf("Calling Gemini API for embedding (text length: %d)...\n", len(text))
	result, err := g.client.Models.EmbedContent(ctx,
		"gemini-embedding-001",
		contents,
		nil, // Use default config
	)
	if err != nil {
		fmt.Printf("Gemini API error: %v\n", err)
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}
	fmt.Printf("Gemini API call successful\n")

	if len(result.Embeddings) == 0 || len(result.Embeddings[0].Values) == 0 {
		return nil, fmt.Errorf("no embedding values returned")
	}

	// Downsample to 768 dims to reduce memory/DB footprint
	values := result.Embeddings[0].Values
	targetDim := 768
	if len(values) > targetDim {
		pooled := make([]float32, targetDim)
		groupSize := len(values) / targetDim
		if groupSize < 1 {
			groupSize = 1
		}
		for i := 0; i < targetDim; i++ {
			start := i * groupSize
			end := start + groupSize
			if start >= len(values) {
				break
			}
			if end > len(values) {
				end = len(values)
			}
			var sum float64
			for j := start; j < end; j++ {
				sum += float64(values[j])
			}
			count := end - start
			if count > 0 {
				pooled[i] = float32(sum / float64(count))
			}
		}
		return pooled, nil
	}

	return values, nil
}

// Chat generates response using Gemini chat API
func (g *GeminiClientV2) Chat(conversation string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create content for chat
	contents := []*genai.Content{
		genai.NewContentFromText(conversation, genai.RoleUser),
	}

	// Generate response using Gemini 2.0 Flash
	result, err := g.client.Models.GenerateContent(ctx,
		"gemini-2.0-flash",
		contents,
		&genai.GenerateContentConfig{
			Temperature:     genai.Ptr(float32(0.7)),
			MaxOutputTokens: 2048,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate chat response: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	// Extract text from response
	response := ""
	for _, part := range result.Candidates[0].Content.Parts {
		if part.Text != "" {
			response += part.Text
		}
	}

	if response == "" {
		return "", fmt.Errorf("empty response generated")
	}

	return response, nil
}

// Close closes the client connection
func (g *GeminiClientV2) Close() error {
	// Note: genai.Client doesn't have Close method in current version
	return nil
}
