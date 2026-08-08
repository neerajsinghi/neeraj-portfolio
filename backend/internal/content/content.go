package content

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"neeraj-portfolio/backend/internal/llm"
)

type Bundle struct {
	Title        string   `json:"title"`
	Slug         string   `json:"slug"`
	Description  string   `json:"description"`
	CanonicalURL string   `json:"canonical_url,omitempty"`
	Tags         []string `json:"tags"`
	Article      string   `json:"article_markdown"`
	LinkedIn     string   `json:"linkedin_post"`
	Social       string   `json:"social_post"`
	CreatedAt    string   `json:"created_at"`
}

type Generator struct {
	provider llm.Provider
}

func NewGenerator(provider llm.Provider) *Generator {
	return &Generator{provider: provider}
}

func (g *Generator) Generate(ctx context.Context, topic, audience, canonicalBase string) (Bundle, error) {
	if strings.TrimSpace(topic) == "" {
		return Bundle{}, errors.New("topic is required")
	}
	if audience == "" {
		audience = "senior backend, cloud, and AI engineers"
	}

	prompt := fmt.Sprintf(`Create an original, technically rigorous article about %q for %s.
Neeraj's editorial authority comes from 10+ years building production Go backends, microservices, SDKs, distributed systems, security controls, data-intensive services, and AWS delivery systems. Keep the article aligned with that background without inventing personal anecdotes or claiming he implemented the specific example.

Editorial requirements:
- Make backend engineering the center of the article. Prefer Go examples and internals where they clarify the design.
- Explain production mechanics, failure modes, operational consequences, and design tradeoffs instead of giving a generic overview.
- Favor system design, Go internals and language concepts, concurrency, networking, storage, caching, reliability, security, observability, APIs, distributed systems, and modern backend infrastructure.
- Cover AI only when the main problem is backend architecture, retrieval, tool execution, data flow, reliability, evaluation, security, or cost control.
- Avoid frontend tutorials, beginner syntax walkthroughs, career advice, vague trend summaries, and listicles about tools.
- Include at least one concrete architecture, request flow, algorithm, or coherent Go example and end with a practical decision framework.

Return only one JSON object matching this shape:
{"title":"...","slug":"lowercase-hyphenated","description":"120-160 characters","tags":["up-to-4"],"article_markdown":"1200-1800 word markdown article","linkedin_post":"professional post under 2500 characters with the article URL placeholder {{CANONICAL_URL}}","social_post":"post under 280 characters with {{CANONICAL_URL}}"}
Use concrete examples and practical tradeoffs. Do not invent personal experiences, employers, benchmarks, quotations, or sources. Do not wrap the JSON in Markdown.`, topic, audience)

	response, err := g.provider.Complete(ctx, llm.Request{
		System:    "You are Neeraj Singhi's senior backend technical editor. Write precise, production-oriented engineering prose without hype, filler, beginner-level padding, or fabricated claims.",
		Messages:  []llm.Message{{Role: "user", Content: prompt}},
		MaxTokens: 6000,
	})
	if err != nil {
		return Bundle{}, err
	}

	raw := responseText(response)
	bundle, err := ParseBundle(raw)
	if err != nil {
		return Bundle{}, err
	}
	bundle.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	if canonicalBase != "" {
		bundle.CanonicalURL = strings.TrimRight(canonicalBase, "/") + "/blog/" + bundle.Slug
		bundle.LinkedIn = strings.ReplaceAll(bundle.LinkedIn, "{{CANONICAL_URL}}", bundle.CanonicalURL)
		bundle.Social = strings.ReplaceAll(bundle.Social, "{{CANONICAL_URL}}", bundle.CanonicalURL)
	}
	return bundle, nil
}

func ParseBundle(raw string) (Bundle, error) {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "```") {
		firstNewline := strings.IndexByte(raw, '\n')
		lastFence := strings.LastIndex(raw, "```")
		if firstNewline >= 0 && lastFence > firstNewline {
			raw = strings.TrimSpace(raw[firstNewline+1 : lastFence])
		}
	}

	var bundle Bundle
	if err := json.Unmarshal([]byte(raw), &bundle); err != nil {
		return Bundle{}, fmt.Errorf("decode generated bundle: %w", err)
	}
	if bundle.Title == "" || bundle.Slug == "" || bundle.Description == "" || bundle.Article == "" || bundle.LinkedIn == "" || bundle.Social == "" {
		return Bundle{}, errors.New("generated bundle is missing required fields")
	}
	if len(bundle.Tags) == 0 || len(bundle.Tags) > 4 {
		return Bundle{}, errors.New("generated bundle must contain 1 to 4 tags")
	}
	return bundle, nil
}

func responseText(response llm.Response) string {
	var text strings.Builder
	for _, block := range response.Content {
		if block.Type == "text" {
			text.WriteString(block.Text)
		}
	}
	return text.String()
}
