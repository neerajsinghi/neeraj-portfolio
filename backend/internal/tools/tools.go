package tools

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"neeraj-portfolio/backend/internal/blog"
	"neeraj-portfolio/backend/internal/github"
	"neeraj-portfolio/backend/internal/kb"
)

// Tool is the JSON-serializable tool definition sent to the model.
// This is the same shape an MCP server advertises — name, description, schema.
type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"input_schema"`
}

func obj(props map[string]any, required ...string) map[string]any {
	m := map[string]any{"type": "object", "properties": props}
	if len(required) > 0 {
		m["required"] = required
	}
	return m
}

// Tools is the registry. To expose these over a real MCP server, register the
// same names/schemas and forward calls to ExecuteTool.
var Tools = []Tool{
	{
		"search_profile",
		"Full-text search over Neeraj Singhi's profile knowledge base — résumé, LinkedIn, roles, skills, security work, AI/RAG work, performance wins, and leadership. Use this first for any factual question about his background, experience, or fit for a role. Returns the top matching passages.",
		obj(map[string]any{
			"query": map[string]any{"type": "string", "description": "natural-language search query, e.g. 'Dell license enforcement' or 'AI RAG experience'"},
			"top_n": map[string]any{"type": "integer", "description": "number of passages to return (1-5, default 3)", "default": 3},
		}, "query"),
	},
	{
		"list_projects",
		"List Neeraj's featured software projects with descriptions and full tech stacks. Use when asked about specific projects or built products.",
		obj(map[string]any{}),
	},
	{
		"get_skills",
		"Return Neeraj's full technical skill set, structured by category (languages, backend, cloud, AI/LLM, data, security, frontend). Use when asked about tech stack, skills, or whether he knows a specific technology.",
		obj(map[string]any{}),
	},
	{
		"get_education",
		"Return Neeraj's academic background — degrees, institutions, and years. Use when asked about his education or qualifications.",
		obj(map[string]any{}),
	},
	{
		"get_certifications",
		"Return Neeraj's professional certifications and recognition/awards. Use when asked about certifications, credentials, or accomplishments.",
		obj(map[string]any{}),
	},
	{
		"get_experience_summary",
		"Return a structured chronological list of all Neeraj's work roles (title, company, dates, one-line summary). Use for a high-level career overview or timeline questions.",
		obj(map[string]any{}),
	},
	{
		"get_links",
		"Return Neeraj's contact details: email, LinkedIn, GitHub, location, and open-to-work status. Use when asked how to reach him or what roles he's looking for.",
		obj(map[string]any{}),
	},
	{
		"get_github_repos",
		"Fetch Neeraj's public GitHub repositories live (name, description, stars, primary language). Use when asked about his open-source work or GitHub activity.",
		obj(map[string]any{}),
	},
	{
		"list_blogs",
		"List Neeraj's latest published blog posts with titles, summaries, tags, publication dates, and links. Use when asked what he has written or for a list of posts.",
		obj(map[string]any{
			"limit": map[string]any{"type": "integer", "description": "number of posts to return (1-10, default 5)", "default": 5},
		}),
	},
	{
		"search_blogs",
		"Search Neeraj's published blog content by topic and return grounded excerpts with links. Use for questions about subjects covered in his writing.",
		obj(map[string]any{
			"query": map[string]any{"type": "string", "description": "topic or natural-language blog search query"},
			"top_n": map[string]any{"type": "integer", "description": "number of matching posts (1-5, default 3)", "default": 3},
		}, "query"),
	},
	{
		"get_services",
		"Return the engineering services Neeraj can provide, grounded in his demonstrated experience. Use when asked how he can help, what he offers, consulting services, or engagement areas.",
		obj(map[string]any{}),
	},
}

// Executor runs tools with access to optional live data sources.
type Executor struct {
	blogStore blog.Store
}

func NewExecutor(blogStore blog.Store) *Executor {
	return &Executor{blogStore: blogStore}
}

// Execute runs a tool and returns (resultText, sourceIDs).
func (executor *Executor) Execute(name string, input map[string]any) (string, []string) {
	switch name {
	case "search_profile":
		q, _ := input["query"].(string)
		k := 3
		if n, ok := input["top_n"].(float64); ok && n >= 1 && n <= 5 {
			k = int(n)
		}
		docs := kb.Retrieve(q, k)
		if len(docs) == 0 {
			return "No matching passage found in the knowledge base.", nil
		}
		var b strings.Builder
		srcs := make([]string, 0, len(docs))
		for i, d := range docs {
			if i > 0 {
				b.WriteString("\n\n")
			}
			b.WriteString("## " + d.Title + "\n" + d.Text)
			srcs = append(srcs, d.ID)
		}
		return b.String(), srcs

	case "list_projects":
		lines := make([]string, 0, len(Projects))
		for _, p := range Projects {
			line := fmt.Sprintf("• %s\n  %s\n  Stack: %s", p.Name, p.Desc, p.Stack)
			if p.URL != "" {
				line += "\n  URL: " + p.URL
			}
			lines = append(lines, line)
		}
		return strings.Join(lines, "\n\n"), []string{"projects"}

	case "get_skills":
		return `Technical skills by category:

Languages: Go (primary), TypeScript, JavaScript, Java, C (cgo)

Backend & Distributed Systems:
  Microservices, gRPC, REST, WebSockets, SDK design, Node.js, NestJS

Cloud & Infrastructure:
  AWS (ECS, ECR, S3, SES, SNS), Docker, Kubernetes, Terraform

AI / LLM:
  OpenAI, Anthropic Claude, Google Gemini, xAI Grok
  RAG pipelines, MCP / tool-use agents, semantic search, TF-IDF, embeddings
  Provider-agnostic LLM abstraction (llm.Provider interface)

Observability & CI/CD:
  Prometheus, OpenTelemetry, GitLab CI, GitHub Actions

Data Stores:
  MongoDB, MySQL, Redis

Security:
  PKI / mTLS, AES encryption, cryptographic hashing (SHA)
  JWT / RBAC, rate limiting, IP blacklisting, HTTP hardening
  License enforcement, entitlement validation, tamper detection

Frontend (supporting):
  React, Next.js, Redux`, []string{"skills"}

	case "get_education":
		return `Education:

• M.Eng. in Computer Science — BITS Pilani, Pilani Campus (Aug 2015 – Jul 2017)
• B.Tech in Computer Science — CIITM, Jaipur (Sep 2008 – Jul 2012)`, []string{"education"}

	case "get_certifications":
		return `Certifications:
• AWS Certified Solutions Architect – Associate
• AWS Certified Cloud Practitioner
• Go Design Patterns

Recognition & Awards:
• Hackathon Winner
• Above & Beyond Award`, []string{"certs"}

	case "get_experience_summary":
		lines := make([]string, 0, len(ExperienceSummary))
		for _, e := range ExperienceSummary {
			lines = append(lines, fmt.Sprintf("• %s | %s | %s\n  %s", e.Role, e.Company, e.Period, e.Summary))
		}
		return strings.Join(lines, "\n\n"), []string{"role_current", "role_dell", "role_turing", "role_truelancer", "role_freelance", "role_broadcom", "role_early"}

	case "get_links":
		return `Contact & availability:
Email:    nsinghi2011@gmail.com
LinkedIn: https://www.linkedin.com/in/neeraj-singhi-golang
GitHub:   https://github.com/neerajsinghi
Location: Delhi, India — open to relocation
Seeking:  Senior / Staff backend or backend+AI roles`, []string{"contact"}

	case "get_github_repos":
		repos, err := github.FetchRepos()
		if err != nil || len(repos) == 0 {
			return "No public GitHub repositories were returned (the account may be private, empty, or rate-limited).", []string{"github"}
		}
		lines := make([]string, 0, 8)
		for i, r := range repos {
			if i >= 8 {
				break
			}
			line := "• " + r.Name
			if r.Language != "" {
				line += " (" + r.Language + ")"
			}
			if r.Stars > 0 {
				line += fmt.Sprintf(" ★%d", r.Stars)
			}
			if r.Description != "" {
				line += " — " + r.Description
			}
			lines = append(lines, line)
		}
		return strings.Join(lines, "\n"), []string{"github"}

	case "list_blogs":
		limit := boundedInt(input["limit"], 5, 1, 10)
		posts, err := executor.publishedBlogs(limit)
		if err != nil {
			return "Published blogs are temporarily unavailable.", nil
		}
		if len(posts) == 0 {
			return "No published blog posts are available yet.", nil
		}
		return formatBlogs(posts, false)

	case "search_blogs":
		query, _ := input["query"].(string)
		limit := boundedInt(input["top_n"], 3, 1, 5)
		posts, err := executor.publishedBlogs(100)
		if err != nil {
			return "Published blogs are temporarily unavailable.", nil
		}
		matches := rankBlogs(posts, query, limit)
		if len(matches) == 0 {
			return "No published blog posts matched that topic.", nil
		}
		return formatBlogs(matches, true)

	case "get_services":
		return `Engineering services grounded in Neeraj's production experience:

• Go backend architecture — APIs, microservices, gRPC, SDKs, concurrency, and distributed-system design.
• AI and RAG integration — retrieval pipelines, tool-use agents, semantic search, document processing, and provider integrations.
• Cloud and delivery — AWS serverless/container deployments, Docker, Kubernetes, Terraform, CI/CD, and observability.
• Security engineering — JWT/RBAC, PKI/mTLS, entitlement systems, tamper detection, rate limiting, and HTTP hardening.
• Performance and reliability — MongoDB query design, Redis caching, profiling, pagination, compression, and production stabilization.
• Technical leadership — architecture reviews, delivery planning, code-review standards, and mentoring backend teams.`, []string{"services", "skills", "leadership", "security", "performance"}
	}
	return "Unknown tool.", nil
}

func (executor *Executor) publishedBlogs(limit int) ([]blog.Post, error) {
	if executor.blogStore == nil {
		return nil, errors.New("blog storage is not configured")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return executor.blogStore.ListPublished(ctx, limit)
}

func boundedInt(value any, fallback, minimum, maximum int) int {
	if number, ok := value.(float64); ok && int(number) >= minimum && int(number) <= maximum {
		return int(number)
	}
	return fallback
}

func rankBlogs(posts []blog.Post, query string, limit int) []blog.Post {
	tokens := strings.Fields(strings.ToLower(query))
	type scoredPost struct {
		post  blog.Post
		score int
	}
	scored := make([]scoredPost, 0, len(posts))
	for _, post := range posts {
		title := strings.ToLower(post.Title)
		body := strings.ToLower(post.Description + " " + strings.Join(post.Tags, " ") + " " + post.ContentMarkdown)
		score := 0
		for _, token := range tokens {
			if strings.Contains(title, token) {
				score += 3
			}
			if strings.Contains(body, token) {
				score++
			}
		}
		if score > 0 {
			scored = append(scored, scoredPost{post: post, score: score})
		}
	}
	sort.SliceStable(scored, func(i, j int) bool { return scored[i].score > scored[j].score })
	if len(scored) > limit {
		scored = scored[:limit]
	}
	result := make([]blog.Post, len(scored))
	for index, match := range scored {
		result[index] = match.post
	}
	return result
}

func formatBlogs(posts []blog.Post, includeExcerpt bool) (string, []string) {
	lines := make([]string, 0, len(posts))
	sources := make([]string, 0, len(posts))
	for _, post := range posts {
		line := fmt.Sprintf("• %s\n  %s", post.Title, post.Description)
		if len(post.Tags) > 0 {
			line += "\n  Topics: " + strings.Join(post.Tags, ", ")
		}
		if post.PublishedAt != nil {
			line += "\n  Published: " + post.PublishedAt.Format("2 Jan 2006")
		}
		if includeExcerpt {
			excerptRunes := []rune(strings.TrimSpace(post.ContentMarkdown))
			if len(excerptRunes) > 700 {
				excerptRunes = append(excerptRunes[:700], '…')
			}
			excerpt := string(excerptRunes)
			line += "\n  Excerpt: " + excerpt
		}
		line += "\n  URL: https://neerajsinghi.com/blogs/" + post.Slug
		lines = append(lines, line)
		sources = append(sources, "blog:"+post.Slug)
	}
	return strings.Join(lines, "\n\n"), sources
}

// Project is a featured project shown on the portfolio site.
type Project struct {
	Name  string
	Desc  string
	Stack string
	URL   string
}

var Projects = []Project{
	{
		Name:  "Consumer Marketplace Platform",
		Desc:  "AI-enabled mobile-first marketplace: listing generation, document/résumé parsing, semantic search. Led a 5-engineer team end-to-end.",
		Stack: "Go, NestJS, Next.js, MongoDB, OpenAI/RAG, Redis, Stripe, Firebase, AWS",
	},
	{
		Name:  "Medical News Platform",
		Desc:  "Curated medical-news platform as Go microservices with a Flutter mobile client, deployed on Docker and Kubernetes.",
		Stack: "Go, Flutter, Docker, Kubernetes",
	},
	{
		Name:  "Portfolio Agent (this site)",
		Desc:  "Provider-agnostic LLM agent backed by TF-IDF RAG over a curated knowledge base. Supports Anthropic, OpenAI, Grok, and Gemini via a single llm.Provider interface.",
		Stack: "Go, Anthropic/OpenAI/Grok/Gemini, RAG, TF-IDF, Next.js, SSE, Docker",
		URL:   "https://github.com/neerajsinghi",
	},
	{
		Name:  "mcp-gateway (open source)",
		Desc:  "Security gateway for MCP servers: API-key and JWT auth, per-tool allow/deny policies that also filter tools/list, per-identity and per-tool rate limits, structured audit logs and OpenTelemetry tracing.",
		Stack: "Go, MCP, JSON-RPC, JWT, OpenTelemetry, Docker",
		URL:   "https://github.com/neerajsinghi/mcp-gateway",
	},
	{
		Name:  "system-design-patterns-go (open source)",
		Desc:  "Catalog of 48 system design patterns, each a small tested Go implementation with a runnable simulation and a write-up of trade-offs, across reliability, data and storage, distributed coordination, and architecture.",
		Stack: "Go (standard library only), distributed systems, storage engines",
		URL:   "https://github.com/neerajsinghi/system-design-patterns-go",
	},
	{
		Name:  "dsa-interview-go (open source)",
		Desc:  "A 21-module data structures and algorithms course for coding interviews: a lesson per pattern plus 150+ interview problems, each with brute-force and optimal Go solutions verified by table-driven and randomized tests.",
		Stack: "Go (standard library only), algorithms, data structures",
		URL:   "https://github.com/neerajsinghi/dsa-interview-go",
	},
	{
		Name:  "go-license (open source)",
		Desc:  "Offline Ed25519-signed license keys with feature gating, seat limits, clock-rollback detection and a tamper-evident hash-chained usage log, plus a golicense CLI.",
		Stack: "Go, Ed25519, HMAC-SHA256, CLI",
		URL:   "https://github.com/neerajsinghi/go-license",
	},
	{
		Name:  "go-rag-kit (open source)",
		Desc:  "Production-minded RAG toolkit in pure Go: markdown-aware chunking, OpenAI embeddings, in-memory/pgvector/MongoDB Atlas stores, BM25 + vector hybrid retrieval with RRF, and a recall@k/MRR eval CLI.",
		Stack: "Go, RAG, BM25, pgvector, MongoDB Atlas Vector Search",
		URL:   "https://github.com/neerajsinghi/go-rag-kit",
	},
	{
		Name:  "llm-otel-go (open source)",
		Desc:  "OpenTelemetry tracing, token metrics and cost tracking for OpenAI and Anthropic API calls in Go (including streaming), following the GenAI semantic conventions, with a Grafana dashboard.",
		Stack: "Go, OpenTelemetry, Prometheus, Jaeger, Grafana",
		URL:   "https://github.com/neerajsinghi/llm-otel-go",
	},
	{
		Name:  "ratelimit (open source)",
		Desc:  "Distributed rate limiting on Redis (GCRA and sliding window) for Go and NestJS, sharing the same atomic Lua scripts so both behave identically.",
		Stack: "Go, NestJS, TypeScript, Redis, Lua",
		URL:   "https://github.com/neerajsinghi/ratelimit",
	},
	{
		Name:  "go-outbox (open source)",
		Desc:  "Transactional outbox relay and Idempotency-Key HTTP middleware for Go microservices, with leases, backoff retries and dead-lettering on PostgreSQL, MySQL and SQLite.",
		Stack: "Go, database/sql, PostgreSQL, MySQL, SQLite",
		URL:   "https://github.com/neerajsinghi/go-outbox",
	},
	{
		Name:  "cgo-to-pure-go (open source)",
		Desc:  "Case study on replacing a cgo dependency with pure Go: golden vectors, differential fuzzing, measured benchmarks and a build-tag migration path.",
		Stack: "Go, cgo, C, fuzzing, benchmarks",
		URL:   "https://github.com/neerajsinghi/cgo-to-pure-go",
	},
	{
		Name:  "portfolio-rag-template (open source)",
		Desc:  "This site as a reusable GitHub template: a Next.js portfolio with a Go RAG agent that answers visitors' questions about its owner.",
		Stack: "Go, Next.js, RAG, Docker",
		URL:   "https://github.com/neerajsinghi/portfolio-rag-template",
	},
}

// ExperienceSummary is the structured work history used by get_experience_summary.
type Experience struct {
	Role    string
	Company string
	Period  string
	Summary string
}

var ExperienceSummary = []Experience{
	{"Lead Backend Engineer", "Independent / Freelance", "Feb 2026 – Present",
		"Leads a 5-engineer team building an AI-enabled consumer marketplace (Go, NestJS, MongoDB, OpenAI/RAG)."},
	{"Senior Software Consultant", "Dell Technologies (via Objectwin)", "Oct 2022 – Jan 2026",
		"Built license-enforcement service: entitlement validation, tamper detection, C→Go port; owned CI/CD on AWS ECS/ECR."},
	{"Consulting Software Engineer", "Turing.com", "Mar 2022 – Aug 2022",
		"Built AWS developer-activity daemon; stabilised Matching Kanban (Go, Node.js, React, MySQL)."},
	{"Software Development Consultant", "Truelancer · client AstraZeneca", "May 2021 – Apr 2022",
		"Delivered CRM user-roster & session features, ~30% perf improvement, 20+ defect fixes."},
	{"Full-Stack Engineer", "Freelance", "Apr 2019 – Aug 2021",
		"Medical News microservices (Go + Flutter on K8s); social platform and stock-trading backend APIs."},
	{"R&D Engineer, Software 2", "Broadcom Inc.", "Jul 2017 – Apr 2019",
		"Improved alarming-system accuracy ~20%; built SaaS monitoring probes (SQL, Active Directory)."},
	{"Software Intern", "Cisco Systems", "Jan 2017 – Jun 2017",
		"IoT/CoAP protocol work on Cisco switches."},
	{"QA Engineer", "Trignodev Software", "Oct 2013 – Dec 2014",
		"GUI, functional, and regression test suites for Kareermatrix.com."},
}
