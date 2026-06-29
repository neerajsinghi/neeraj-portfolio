// Package grok implements llm.Provider for xAI's Grok using the
// OpenAI-compatible endpoint. Only the base URL, API key, and model differ
// from the openai package — all protocol logic is shared.
package grok

import (
	"os"

	"neeraj-portfolio/backend/internal/llm/openai"
)

// New returns a Provider pointing at the xAI Grok endpoint.
// Model defaults to grok-3; override with GROK_MODEL.
// API key is read from XAI_API_KEY.
func New() *openai.Provider {
	model := os.Getenv("GROK_MODEL")
	if model == "" {
		model = "grok-3"
	}
	return openai.NewWithConfig(openai.Config{
		BaseURL: "https://api.x.ai/v1",
		APIKey:  os.Getenv("XAI_API_KEY"),
		Model:   model,
	})
}
