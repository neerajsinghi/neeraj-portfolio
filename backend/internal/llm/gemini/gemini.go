// Package gemini implements llm.Provider for Google Gemini using its
// OpenAI-compatible endpoint. Only the base URL, API key, and model differ
// from the openai package — all protocol logic is shared.
package gemini

import (
	"os"

	"neeraj-portfolio/backend/internal/llm/openai"
)

// New returns a Provider pointing at Google's Gemini OpenAI-compatible endpoint.
// Model defaults to gemini-2.0-flash; override with GEMINI_MODEL.
// API key is read from GEMINI_API_KEY.
func New() *openai.Provider {
	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = "gemini-2.0-flash"
	}
	return openai.NewWithConfig(openai.Config{
		BaseURL: "https://generativelanguage.googleapis.com/v1beta/openai",
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Model:   model,
	})
}
