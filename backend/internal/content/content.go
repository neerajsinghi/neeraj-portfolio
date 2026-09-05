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
	// Type distinguishes a full blog bundle ("blog", the default) from a
	// "linkedin" bundle that has no article and must never be published to DEV.
	Type string `json:"type,omitempty"`
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

// GenerateLinkedInOnly creates a short LinkedIn-only bundle with no article,
// for alternate-day posts that are never published as a blog.
func (g *Generator) GenerateLinkedInOnly(ctx context.Context, topic, audience string) (Bundle, error) {
	if strings.TrimSpace(topic) == "" {
		return Bundle{}, errors.New("topic is required")
	}
	if audience == "" {
		audience = "senior backend, cloud, and AI engineers"
	}

	prompt := fmt.Sprintf(`Write an original, professional LinkedIn post about %q for %s.
Neeraj's voice comes from 10+ years building production Go backends, microservices, SDKs, distributed systems, security controls, data-intensive services, and AWS delivery systems. Keep opinions and framing consistent with that background without inventing specific personal anecdotes, employers, or metrics.

Requirements:
- This is not a blog article. Do not write article prose, headings, or a long-form structure.
- Write a single professional LinkedIn post, under 2500 characters, with a strong hook line, 2-4 short paragraphs or a brief list, and a closing thought or question that invites engagement.
- Topics: career growth, engineering leadership, industry commentary, hiring/interviewing, team practices, or personal engineering philosophy. Avoid technical tutorials or code samples.
- Do not fabricate quotations, statistics, employers, or events.

Return only one JSON object matching this shape:
{"title":"short internal label","slug":"lowercase-hyphenated","description":"one-line summary for internal review","tags":["up-to-4"],"linkedin_post":"the post, under 2500 characters","social_post":"a shorter under-280-character variant of the same post"}
Do not wrap the JSON in Markdown.`, topic, audience)

	response, err := g.provider.Complete(ctx, llm.Request{
		System:    "You are Neeraj Singhi's personal LinkedIn ghostwriter. Write in a direct, credible, non-hype professional voice without fabricated claims.",
		Messages:  []llm.Message{{Role: "user", Content: prompt}},
		MaxTokens: 2000,
	})
	if err != nil {
		return Bundle{}, err
	}

	bundle, err := ParseLinkedInBundle(responseText(response))
	if err != nil {
		return Bundle{}, err
	}
	bundle.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	return bundle, nil
}

func ParseBundle(raw string) (Bundle, error) {
	bundle, err := decodeBundle(raw)
	if err != nil {
		return Bundle{}, err
	}
	bundle.Title = strings.TrimSpace(bundle.Title)
	bundle.Slug = strings.TrimSpace(bundle.Slug)
	bundle.Description = strings.TrimSpace(bundle.Description)
	bundle.Article = strings.TrimSpace(bundle.Article)
	bundle.LinkedIn = strings.TrimSpace(bundle.LinkedIn)
	bundle.Social = strings.TrimSpace(bundle.Social)
	if bundle.Title == "" || bundle.Slug == "" || bundle.Description == "" || bundle.Article == "" || bundle.LinkedIn == "" || bundle.Social == "" {
		return Bundle{}, errors.New("generated bundle is missing required fields")
	}
	if len(bundle.Title) < 5 || len(bundle.Title) > 120 {
		return Bundle{}, errors.New("generated bundle title must contain 5 to 120 characters")
	}
	if len(bundle.Description) < 40 || len(bundle.Description) > 180 {
		return Bundle{}, errors.New("generated bundle description must contain 40 to 180 characters")
	}
	if len(bundle.Article) < 100 {
		return Bundle{}, errors.New("generated bundle article must contain at least 100 characters")
	}
	if len(bundle.Tags) == 0 || len(bundle.Tags) > 4 {
		return Bundle{}, errors.New("generated bundle must contain 1 to 4 tags")
	}
	bundle.Type = "blog"
	return bundle, nil
}

// ParseLinkedInBundle parses a LinkedIn-only bundle, which has no article and
// is never eligible for DEV publishing.
func ParseLinkedInBundle(raw string) (Bundle, error) {
	bundle, err := decodeBundle(raw)
	if err != nil {
		return Bundle{}, err
	}
	if bundle.Title == "" || bundle.Slug == "" || bundle.LinkedIn == "" || bundle.Social == "" {
		return Bundle{}, errors.New("generated bundle is missing required fields")
	}
	if len(bundle.Tags) == 0 || len(bundle.Tags) > 4 {
		return Bundle{}, errors.New("generated bundle must contain 1 to 4 tags")
	}
	bundle.Type = "linkedin"
	return bundle, nil
}

func decodeBundle(raw string) (Bundle, error) {
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
