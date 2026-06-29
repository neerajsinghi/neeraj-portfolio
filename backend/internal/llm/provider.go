// Package llm defines the provider-agnostic interface and types for LLM
// interactions. Swap providers by injecting a different implementation —
// no agent or handler code needs to change.
package llm

import (
	"context"
	"encoding/json"
)

// Provider is the interface any LLM backend must satisfy.
type Provider interface {
	// Complete sends a request to the model and returns the full response.
	Complete(ctx context.Context, req Request) (Response, error)
	// ModelName returns the identifier of the model in use.
	ModelName() string
}

// Request is the provider-agnostic input to Complete.
type Request struct {
	System    string
	Messages  []Message
	Tools     []ToolDef
	MaxTokens int
}

// Response is the provider-agnostic output from Complete.
type Response struct {
	Content    []Block
	StopReason StopReason
}

// StopReason indicates why the model stopped generating.
type StopReason string

const (
	StopReasonEndTurn StopReason = "end_turn"
	StopReasonToolUse StopReason = "tool_use"
)

// Message is a single turn in the conversation.
type Message struct {
	Role    string // "user" | "assistant"
	Content any    // string | []Block
}

// Block is a content unit within a message: plain text, a tool call, or a tool result.
type Block struct {
	Type string // "text" | "tool_use" | "tool_result"

	// Type == "text"
	Text string

	// Type == "tool_use"
	ID    string
	Name  string
	Input json.RawMessage

	// Type == "tool_result"
	ToolUseID string
	Content   any // string result returned to the model
}

// ToolDef is the provider-agnostic tool definition sent with each request.
type ToolDef struct {
	Name        string
	Description string
	InputSchema map[string]any
}
