package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GeminiClient struct {
	apiKey string
	client *http.Client
}

// Gemini API request/response structures for embeddings
type GeminiEmbeddingRequest struct {
	Model    string                 `json:"model"`
	Content  GeminiEmbeddingContent `json:"content"`
	TaskType string                 `json:"taskType,omitempty"`
}

type GeminiEmbeddingContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiEmbeddingResponse struct {
	Embedding GeminiEmbeddingData `json:"embedding"`
}

type GeminiEmbeddingData struct {
	Values []float32 `json:"values"`
}

// Gemini API request/response structures for chat
type GeminiChatRequest struct {
	Contents         []GeminiContent         `json:"contents"`
	GenerationConfig *GeminiGenerationConfig `json:"generationConfig,omitempty"`
	SafetySettings   []GeminiSafetySetting   `json:"safetySettings,omitempty"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
	Role  string       `json:"role,omitempty"`
}

type GeminiGenerationConfig struct {
	Temperature     float32 `json:"temperature,omitempty"`
	TopK            int     `json:"topK,omitempty"`
	TopP            float32 `json:"topP,omitempty"`
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
}

type GeminiSafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

type GeminiChatResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

type GeminiCandidate struct {
	Content       GeminiContent        `json:"content"`
	FinishReason  string               `json:"finishReason,omitempty"`
	Index         int                  `json:"index,omitempty"`
	SafetyRatings []GeminiSafetyRating `json:"safetyRatings,omitempty"`
}

type GeminiSafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (c *GeminiClient) GenerateEmbedding(text string) ([]float32, error) {
	reqBody := GeminiEmbeddingRequest{
		Model: "models/embedding-001",
		Content: GeminiEmbeddingContent{
			Parts: []GeminiPart{
				{Text: text},
			},
		},
		TaskType: "RETRIEVAL_DOCUMENT",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/embedding-001:embedContent?key=%s", c.apiKey)

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gemini API error: %s", string(body))
	}

	var embeddingResp GeminiEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embeddingResp); err != nil {
		return nil, err
	}

	return embeddingResp.Embedding.Values, nil
}

func (c *GeminiClient) Chat(messages []Message) (string, error) {
	// Convert messages to Gemini format
	var contents []GeminiContent

	for _, msg := range messages {
		role := msg.Role
		if role == "assistant" {
			role = "model"
		}

		// Skip system messages for now as Gemini handles them differently
		if msg.Role == "system" {
			// Add system message as user message with instruction prefix
			contents = append(contents, GeminiContent{
				Parts: []GeminiPart{
					{Text: "System instructions: " + msg.Content},
				},
				Role: "user",
			})
			// Add empty model response to maintain conversation flow
			contents = append(contents, GeminiContent{
				Parts: []GeminiPart{
					{Text: "Understood. I will follow these instructions."},
				},
				Role: "model",
			})
		} else {
			contents = append(contents, GeminiContent{
				Parts: []GeminiPart{
					{Text: msg.Content},
				},
				Role: role,
			})
		}
	}

	reqBody := GeminiChatRequest{
		Contents: contents,
		GenerationConfig: &GeminiGenerationConfig{
			Temperature:     0.7,
			TopK:            40,
			TopP:            0.95,
			MaxOutputTokens: 2048,
		},
		SafetySettings: []GeminiSafetySetting{
			{
				Category:  "HARM_CATEGORY_HARASSMENT",
				Threshold: "BLOCK_MEDIUM_AND_ABOVE",
			},
			{
				Category:  "HARM_CATEGORY_HATE_SPEECH",
				Threshold: "BLOCK_MEDIUM_AND_ABOVE",
			},
			{
				Category:  "HARM_CATEGORY_SEXUALLY_EXPLICIT",
				Threshold: "BLOCK_MEDIUM_AND_ABOVE",
			},
			{
				Category:  "HARM_CATEGORY_DANGEROUS_CONTENT",
				Threshold: "BLOCK_MEDIUM_AND_ABOVE",
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", c.apiKey)

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("gemini API error: %s", string(body))
	}

	var chatResp GeminiChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Candidates) == 0 || len(chatResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini API")
	}

	return chatResp.Candidates[0].Content.Parts[0].Text, nil
}
