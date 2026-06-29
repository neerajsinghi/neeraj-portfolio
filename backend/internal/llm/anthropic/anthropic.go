// Package anthropic implements llm.Provider for the Anthropic Messages API.
// To switch to a different provider, create another package that satisfies
// llm.Provider and inject it in cmd/server/main.go.
package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"neeraj-portfolio/backend/internal/llm"
)

const apiURL = "https://api.anthropic.com/v1/messages"

// Provider implements llm.Provider for Anthropic.
type Provider struct {
	model      string
	httpClient *http.Client
}

// New returns an Anthropic Provider. The model is read from ANTHROPIC_MODEL
// (defaults to claude-sonnet-4-6).
func New() *Provider {
	model := os.Getenv("ANTHROPIC_MODEL")
	if model == "" {
		model = "claude-sonnet-4-6"
	}
	return &Provider{
		model:      model,
		httpClient: &http.Client{Timeout: 45 * time.Second},
	}
}

func (p *Provider) ModelName() string { return p.model }

// ── Anthropic wire types (internal to this package) ───────────────────────

type apiBlock struct {
	Type      string          `json:"type"`
	Text      string          `json:"text,omitempty"`
	ID        string          `json:"id,omitempty"`
	Name      string          `json:"name,omitempty"`
	Input     json.RawMessage `json:"input,omitempty"`
	ToolUseID string          `json:"tool_use_id,omitempty"`
	Content   any             `json:"content,omitempty"`
}

type apiMessage struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

type apiTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"input_schema"`
}

type apiRequest struct {
	Model     string       `json:"model"`
	MaxTokens int          `json:"max_tokens"`
	System    string       `json:"system"`
	Tools     []apiTool    `json:"tools"`
	Messages  []apiMessage `json:"messages"`
}

type apiError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type apiResponse struct {
	Content    []apiBlock `json:"content"`
	StopReason string     `json:"stop_reason"`
	Error      *apiError  `json:"error"`
}

// ── Provider.Complete ─────────────────────────────────────────────────────

func (p *Provider) Complete(ctx context.Context, req llm.Request) (llm.Response, error) {
	key := os.Getenv("ANTHROPIC_API_KEY")
	if key == "" {
		return llm.Response{}, errors.New("ANTHROPIC_API_KEY is not set on the server")
	}

	body, _ := json.Marshal(apiRequest{
		Model:     p.model,
		MaxTokens: req.MaxTokens,
		System:    req.System,
		Tools:     toAPITools(req.Tools),
		Messages:  toAPIMessages(req.Messages),
	})

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		return llm.Response{}, err
	}
	httpReq.Header.Set("content-type", "application/json")
	httpReq.Header.Set("x-api-key", key)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return llm.Response{}, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	var r apiResponse
	_ = json.Unmarshal(raw, &r)
	if resp.StatusCode != http.StatusOK {
		if r.Error != nil {
			return llm.Response{}, errors.New(r.Error.Message)
		}
		return llm.Response{}, fmt.Errorf("anthropic %d: %s", resp.StatusCode, truncate(string(raw), 200))
	}
	if r.Error != nil {
		return llm.Response{}, errors.New(r.Error.Message)
	}

	stopReason := llm.StopReasonEndTurn
	if r.StopReason == "tool_use" {
		stopReason = llm.StopReasonToolUse
	}
	return llm.Response{
		Content:    fromAPIBlocks(r.Content),
		StopReason: stopReason,
	}, nil
}

// ── Conversion helpers ────────────────────────────────────────────────────

func toAPITools(defs []llm.ToolDef) []apiTool {
	out := make([]apiTool, len(defs))
	for i, d := range defs {
		out[i] = apiTool{Name: d.Name, Description: d.Description, InputSchema: d.InputSchema}
	}
	return out
}

func toAPIMessages(msgs []llm.Message) []apiMessage {
	out := make([]apiMessage, len(msgs))
	for i, m := range msgs {
		switch v := m.Content.(type) {
		case []llm.Block:
			blocks := make([]apiBlock, len(v))
			for j, b := range v {
				blocks[j] = toAPIBlock(b)
			}
			out[i] = apiMessage{Role: m.Role, Content: blocks}
		default:
			out[i] = apiMessage{Role: m.Role, Content: v}
		}
	}
	return out
}

func toAPIBlock(b llm.Block) apiBlock {
	return apiBlock{
		Type:      b.Type,
		Text:      b.Text,
		ID:        b.ID,
		Name:      b.Name,
		Input:     b.Input,
		ToolUseID: b.ToolUseID,
		Content:   b.Content,
	}
}

func fromAPIBlocks(blocks []apiBlock) []llm.Block {
	out := make([]llm.Block, len(blocks))
	for i, b := range blocks {
		out[i] = llm.Block{
			Type:      b.Type,
			Text:      b.Text,
			ID:        b.ID,
			Name:      b.Name,
			Input:     b.Input,
			ToolUseID: b.ToolUseID,
			Content:   b.Content,
		}
	}
	return out
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
