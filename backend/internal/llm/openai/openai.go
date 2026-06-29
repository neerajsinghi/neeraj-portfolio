// Package openai implements llm.Provider for any OpenAI-compatible Chat
// Completions API (OpenAI, xAI/Grok, Google Gemini OpenAI-compat, etc.).
// Use New() for OpenAI, or NewWithConfig() to point at a different endpoint.
package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"neeraj-portfolio/backend/internal/llm"
)

// Config holds the settings for an OpenAI-compatible provider.
type Config struct {
	// BaseURL is the API root, without a trailing slash.
	// e.g. "https://api.openai.com/v1"
	BaseURL string
	// APIKey is sent as "Authorization: Bearer <key>".
	APIKey string
	// Model is the model identifier (e.g. "gpt-4o", "grok-3").
	Model string
}

// Provider implements llm.Provider for an OpenAI-compatible endpoint.
type Provider struct {
	cfg        Config
	httpClient *http.Client
}

// New returns a Provider using the standard OpenAI endpoint.
// Model defaults to gpt-4o; override with OPENAI_MODEL.
func New() *Provider {
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4o"
	}
	return NewWithConfig(Config{
		BaseURL: "https://api.openai.com/v1",
		APIKey:  os.Getenv("OPENAI_API_KEY"),
		Model:   model,
	})
}

// NewWithConfig returns a Provider with the given configuration.
// Use this to target Grok, Gemini, local models, or any other
// OpenAI-compatible endpoint.
func NewWithConfig(cfg Config) *Provider {
	return &Provider{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 45 * time.Second},
	}
}

func (p *Provider) ModelName() string { return p.cfg.Model }

// ── OpenAI wire types (internal to this package) ──────────────────────────

type apiFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type apiToolCall struct {
	ID       string      `json:"id"`
	Type     string      `json:"type"` // "function"
	Function apiFunction `json:"function"`
}

type apiMessage struct {
	Role       string        `json:"role"`
	Content    string        `json:"content,omitempty"`
	ToolCalls  []apiToolCall `json:"tool_calls,omitempty"`
	ToolCallID string        `json:"tool_call_id,omitempty"`
}

type apiFunctionDef struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type apiTool struct {
	Type     string         `json:"type"` // "function"
	Function apiFunctionDef `json:"function"`
}

type apiRequest struct {
	Model     string       `json:"model"`
	MaxTokens int          `json:"max_tokens"`
	Messages  []apiMessage `json:"messages"`
	Tools     []apiTool    `json:"tools,omitempty"`
}

type apiChoice struct {
	Message      apiMessage `json:"message"`
	FinishReason string     `json:"finish_reason"`
}

type apiError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type apiErrorWrapper struct {
	Error apiError `json:"error"`
}

type apiResponse struct {
	Choices []apiChoice     `json:"choices"`
	Error   *apiErrorWrapper `json:"error"`
}

// ── Provider.Complete ─────────────────────────────────────────────────────

func (p *Provider) Complete(ctx context.Context, req llm.Request) (llm.Response, error) {
	if p.cfg.APIKey == "" {
		return llm.Response{}, errors.New("API key is not set for provider " + p.cfg.BaseURL)
	}
	if len(req.Messages) == 0 {
		return llm.Response{}, errors.New("messages must not be empty")
	}

	// System prompt is prepended as a system-role message.
	msgs := make([]apiMessage, 0, len(req.Messages)+1)
	if req.System != "" {
		msgs = append(msgs, apiMessage{Role: "system", Content: req.System})
	}
	msgs = append(msgs, toAPIMessages(req.Messages)...)

	body, _ := json.Marshal(apiRequest{
		Model:     p.cfg.Model,
		MaxTokens: req.MaxTokens,
		Messages:  msgs,
		Tools:     toAPITools(req.Tools),
	})

	url := strings.TrimRight(p.cfg.BaseURL, "/") + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return llm.Response{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.cfg.APIKey)

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
			return llm.Response{}, errors.New(r.Error.Error.Message)
		}
		return llm.Response{}, fmt.Errorf("%s %d: %s", p.cfg.BaseURL, resp.StatusCode, truncate(string(raw), 200))
	}
	if r.Error != nil {
		return llm.Response{}, errors.New(r.Error.Error.Message)
	}
	if len(r.Choices) == 0 {
		return llm.Response{}, errors.New("provider returned no choices")
	}

	return fromAPIResponse(r.Choices[0]), nil
}

// ── Conversion: llm → OpenAI ──────────────────────────────────────────────

func toAPITools(defs []llm.ToolDef) []apiTool {
	out := make([]apiTool, len(defs))
	for i, d := range defs {
		out[i] = apiTool{
			Type: "function",
			Function: apiFunctionDef{
				Name:        d.Name,
				Description: d.Description,
				Parameters:  d.InputSchema,
			},
		}
	}
	return out
}

// toAPIMessages converts canonical messages to OpenAI wire format.
// A single llm.Message with tool_result blocks becomes multiple tool messages.
func toAPIMessages(msgs []llm.Message) []apiMessage {
	var out []apiMessage
	for _, m := range msgs {
		switch v := m.Content.(type) {
		case string:
			out = append(out, apiMessage{Role: m.Role, Content: v})
		case []llm.Block:
			out = append(out, blocksToAPIMessages(m.Role, v)...)
		default:
			// Fallback: marshal to string
			b, _ := json.Marshal(v)
			out = append(out, apiMessage{Role: m.Role, Content: string(b)})
		}
	}
	return out
}

func blocksToAPIMessages(role string, blocks []llm.Block) []apiMessage {
	// Tool result blocks → individual "tool" role messages.
	if len(blocks) > 0 && blocks[0].Type == "tool_result" {
		msgs := make([]apiMessage, 0, len(blocks))
		for _, b := range blocks {
			content, _ := b.Content.(string)
			msgs = append(msgs, apiMessage{
				Role:       "tool",
				Content:    content,
				ToolCallID: b.ToolUseID,
			})
		}
		return msgs
	}

	// Assistant message: collect text and tool_use blocks.
	msg := apiMessage{Role: role}
	var textParts []string
	for _, b := range blocks {
		switch b.Type {
		case "text":
			textParts = append(textParts, b.Text)
		case "tool_use":
			msg.ToolCalls = append(msg.ToolCalls, apiToolCall{
				ID:   b.ID,
				Type: "function",
				Function: apiFunction{
					Name:      b.Name,
					Arguments: string(b.Input),
				},
			})
		}
	}
	if len(textParts) > 0 {
		msg.Content = strings.Join(textParts, "\n")
	}
	return []apiMessage{msg}
}

// ── Conversion: OpenAI → llm ──────────────────────────────────────────────

func fromAPIResponse(choice apiChoice) llm.Response {
	var blocks []llm.Block

	if choice.Message.Content != "" {
		blocks = append(blocks, llm.Block{Type: "text", Text: choice.Message.Content})
	}
	for _, tc := range choice.Message.ToolCalls {
		blocks = append(blocks, llm.Block{
			Type:  "tool_use",
			ID:    tc.ID,
			Name:  tc.Function.Name,
			Input: json.RawMessage(tc.Function.Arguments),
		})
	}

	stopReason := llm.StopReasonEndTurn
	if choice.FinishReason == "tool_calls" {
		stopReason = llm.StopReasonToolUse
	}
	return llm.Response{Content: blocks, StopReason: stopReason}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
