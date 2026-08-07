package main

import (
	"log"
	"net/http"
	"os"

	anthropicprovider "neeraj-portfolio/backend/internal/llm/anthropic"
	"neeraj-portfolio/backend/internal/server"
	// Swap the active provider by commenting/uncommenting one line below:
	// openaiprovider  "neeraj-portfolio/backend/internal/llm/openai"
	// grokprovider    "neeraj-portfolio/backend/internal/llm/grok"
	// geminiprovider  "neeraj-portfolio/backend/internal/llm/gemini"
)

func main() {
	server.LoadDotEnv(".env")

	// Wire the LLM provider here. To switch providers, replace anthropicprovider.New()
	// with any other implementation — nothing else in the codebase changes:
	//   openaiprovider.New()   → OpenAI (gpt-4o)
	//   grokprovider.New()     → xAI Grok (grok-3)
	//   geminiprovider.New()   → Google Gemini (gemini-2.0-flash)
	prov := anthropicprovider.New()
	h := server.NewHandler(prov)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("neeraj-agent backend listening on :%s (model=%s)", port, prov.ModelName())
	log.Fatal(http.ListenAndServe(":"+port, h))
}
