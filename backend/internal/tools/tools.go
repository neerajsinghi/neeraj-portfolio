package tools

import (
	"fmt"
	"strings"

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
}

// ExecuteTool runs a tool and returns (resultText, sourceIDs).
func ExecuteTool(name string, input map[string]any) (string, []string) {
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
	}
	return "Unknown tool.", nil
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

