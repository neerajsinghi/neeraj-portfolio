package agent

import (
	"context"
	"encoding/json"
	"strings"

	"neeraj-portfolio/backend/internal/llm"
	"neeraj-portfolio/backend/internal/tools"
)

const systemPrompt = "You are the portfolio agent for Neeraj Singhi, a senior backend & AI engineer. " +
	"Answer a visitor's questions about Neeraj — his experience, skills, projects, and fit for roles. " +
	"Rules: (1) ALWAYS ground answers in the tools; call search_profile (or the other tools) before answering " +
	"factual questions, and never invent facts. (2) Be concise and specific — 2 to 4 short paragraphs max, in a " +
	"warm, professional voice. (3) Refer to him as 'Neeraj'. (4) If asked whether he fits a role, search first, then " +
	"give an honest, evidence-based take. (5) If a question is unrelated to Neeraj or his career, briefly say so and " +
	"steer back. You are speaking to recruiters and engineers, so be substantive, not salesy."

// Turn is a simple text message in the visitor-facing conversation history.
type Turn struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Agent runs the tool-use loop against an LLM provider.
type Agent struct {
	provider llm.Provider
}

// New returns an Agent backed by the given provider.
func New(p llm.Provider) *Agent {
	return &Agent{provider: p}
}

// toolDefs converts the tools registry to provider-agnostic ToolDef slice once.
var toolDefs = func() []llm.ToolDef {
	defs := make([]llm.ToolDef, len(tools.Tools))
	for i, t := range tools.Tools {
		defs[i] = llm.ToolDef{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: t.InputSchema,
		}
	}
	return defs
}()

// Run executes the tool-use loop, emitting events as it goes:
//
//	tool    {name, input}   – the model decided to call a tool (MCP-style)
//	sources {name, sources} – which KB chunks the tool returned
//	text    {text}          – an answer paragraph from the model
//	done    {}              – finished
func (a *Agent) Run(history []Turn, emit func(event string, data any)) error {
	msgs := make([]llm.Message, 0, len(history)+6)
	for _, t := range history {
		role := t.Role
		if role != "assistant" {
			role = "user"
		}
		msgs = append(msgs, llm.Message{Role: role, Content: t.Content})
	}

	for step := 0; step < 5; step++ {
		resp, err := a.provider.Complete(context.Background(), llm.Request{
			System:    systemPrompt,
			Messages:  msgs,
			Tools:     toolDefs,
			MaxTokens: 1000,
		})
		if err != nil {
			return err
		}

		msgs = append(msgs, llm.Message{Role: "assistant", Content: resp.Content})

		var toolUses []llm.Block
		for _, b := range resp.Content {
			switch b.Type {
			case "text":
				if strings.TrimSpace(b.Text) != "" {
					emit("text", map[string]string{"text": b.Text})
				}
			case "tool_use":
				toolUses = append(toolUses, b)
			}
		}

		if resp.StopReason != llm.StopReasonToolUse || len(toolUses) == 0 {
			break
		}

		results := make([]llm.Block, 0, len(toolUses))
		for _, tu := range toolUses {
			var input map[string]any
			_ = json.Unmarshal(tu.Input, &input)
			emit("tool", map[string]any{"name": tu.Name, "input": input})
			text, srcs := tools.ExecuteTool(tu.Name, input)
			emit("sources", map[string]any{"name": tu.Name, "sources": srcs})
			results = append(results, llm.Block{Type: "tool_result", ToolUseID: tu.ID, Content: text})
		}
		msgs = append(msgs, llm.Message{Role: "user", Content: results})
	}

	emit("done", map[string]any{})
	return nil
}
