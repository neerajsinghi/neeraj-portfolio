package kb

import (
	"math"
	"regexp"
	"sort"
	"strings"
)

// Doc is one chunk of Neeraj's profile (résumé + LinkedIn derived).
type Doc struct {
	ID    string
	Title string
	Tags  string
	Text  string
}

// KB is the corpus the search_profile tool retrieves over.
var KB = []Doc{
	{"summary", "Summary", "senior backend engineer go aws sdk security ai rag leader 10 years backend ai distributed systems",
		"Senior backend engineer with 10+ years building production Go services, SDKs, and containerized systems on AWS. Currently leads a 5-engineer team delivering an AI-enabled consumer marketplace platform (Go, NestJS, Next.js, MongoDB), owning architecture and technical direction. Track record spans security-conscious backend work — license enforcement, usage metering, tamper detection with encryption and hashing — and modern AI integration with OpenAI, RAG, and tool-based workflows. Proven technical lead, mentor, and architect who operates as the technical anchor on a team."},
	{"role_current", "Lead Backend Engineer (current)", "lead freelance 2026 present marketplace go nestjs nextjs mongodb openai rag stripe redis aws team mentoring ai agents architecture",
		"Lead Backend Engineer (Feb 2026–present). Leads a 5-engineer team (3 backend, 2 frontend) building a consumer marketplace platform, owning architecture and delivery. Architected a modular Go + NestJS backend (40+ modules) on MongoDB with a Next.js frontend and versioned REST + Swagger. Built AI features with OpenAI — listing/job-posting generation, document/resume parsing, semantic search — via RAG and tool-based integrations. Implemented JWT/RBAC, rate limiting, IP blacklisting, HTTP hardening; integrated Stripe, Firebase, Twilio. Improved performance with MongoDB query refactors, Redis caching, compression and pagination on AWS (S3/SES/SNS), Dockerized."},
	{"role_dell", "Senior Software Consultant — Dell Technologies", "dell objectwin 2022 2026 license enforcement entitlement tamper detection encryption hashing cgo c to go ci cd gitlab github actions ecs ecr mentoring security pki",
		"Senior Software Consultant at Dell Technologies (via Objectwin Technology), Oct 2022–Jan 2026. Developed part of Dell's license-enforcement service used internally across teams. Built core enforcement components: entitlement validation and feature/capability gating. Implemented usage metering and tamper detection using encryption and cryptographic hashing. Ported a performance-critical C component to Go (reconciling cgo and memory-model differences) shipping a single pure-Go artifact. Established testing/release standards with automated CI gates, set code-review standards, mentored 5 engineers (Go/Java), and owned the release pipeline across GitLab CI and GitHub Actions on AWS ECS/ECR."},
	{"role_turing", "Consulting Software Engineer — Turing.com", "turing 2022 daemon aws linux vm developer activity matching kanban go node nestjs react mysql",
		"Consulting Software Engineer at Turing.com, Mar–Aug 2022. Built and deployed a background daemon on an AWS Linux VM that captured developer activity (keystrokes, commits) to give Turing's customers visibility into engineers' work. Stabilized the Matching Kanban system, resolving 10+ critical defects across Go, Node.js/NestJS and React on MySQL."},
	{"role_truelancer", "Software Development Consultant — Truelancer (AstraZeneca)", "truelancer astrazeneca crm 2021 2022 user roster sessions performance code review",
		"Software Development Consultant via Truelancer (client: AstraZeneca), May 2021–Apr 2022. Delivered user-roster and session functionality for AstraZeneca's CRM, improving performance ~30% and resolving 20+ defects. Led code review on the shared repo."},
	{"role_freelance", "Full-Stack Engineer — Freelance", "freelance 2019 2021 medical news platform microservices go flutter docker kubernetes social stock trading apis",
		"Full-Stack Engineer (freelance), Apr 2019–Aug 2021. Architected a Medical News platform as containerized microservices (Go backend, Flutter client) on Docker/Kubernetes. Designed backend systems and APIs for a social platform and stock-trading clients, owning delivery from requirements to deployment."},
	{"role_broadcom", "R&D Engineer Software 2 — Broadcom", "broadcom 2017 2019 alarming system saas monitoring probes sql active directory p1 defects observability",
		"R&D Engineer, Software 2 at Broadcom Inc., Jul 2017–Apr 2019, Hyderabad. Improved alarming-system accuracy by ~20% and built SaaS monitoring probes (SQL, Active Directory, system) to strengthen customer observability. Partnered with customers to resolve critical P1 defects."},
	{"role_early", "Earlier roles — Cisco, QA", "cisco intern 2017 coap iot switches qa trignodev 2013 2014 testing kareermatrix",
		"Software Intern at Cisco (Jan–Jun 2017): IoT protocol work (CoAP) on Cisco switches. QA Engineer at Trignodev Software (Oct 2013–Dec 2014): built and executed GUI, functional and regression test suites for Kareermatrix.com."},
	{"skills", "Skills & tech stack", "go typescript javascript java c microservices grpc rest websockets sdk node nestjs aws ecs ecr s3 ses sns docker kubernetes terraform prometheus opentelemetry openai anthropic rag mcp tool-use agents semantic search mongodb mysql redis pki mtls encryption react nextjs redux gitlab github actions",
		"Languages: Go (primary), TypeScript, JavaScript, Java, C (cgo). Backend & distributed systems: microservices, gRPC, REST, WebSockets, SDK design, Node.js, NestJS. Cloud & infra: AWS (ECS, ECR, S3, SES, SNS), Docker, Kubernetes, Terraform. Observability: Prometheus, OpenTelemetry. AI/LLM: OpenAI, Anthropic Claude, RAG, MCP, tool-use agents, semantic search, embeddings. CI/CD: GitLab CI, GitHub Actions. Data: MongoDB, MySQL, Redis. Security: PKI/mTLS, AES encryption, cryptographic hashing, JWT/RBAC, rate limiting. Frontend (supporting): React, Next.js, Redux."},
	{"leadership", "Team leadership & mentoring", "team lead mentor architecture technical direction 5 engineers code review standards junior senior staff ownership",
		"Experienced technical lead with a hands-on mentoring track record. Currently leads a 5-engineer team (3 backend, 2 frontend) owning architecture and delivery of a consumer marketplace platform. At Dell: set code-review and testing standards, mentored 5 Go/Java engineers, owned the CI/CD release pipeline end-to-end. Comfortable running full technical direction — from vague requirements to production, including API design, documentation, cross-team coordination, and sprint planning. Believes in setting clear standards and then getting out of the team's way."},
	{"security", "Security engineering", "security license enforcement entitlement tamper detection encryption hashing pki mtls jwt rbac rate limiting ip blacklist http hardening aes sha cryptographic pure-go",
		"Strong security engineering background across multiple roles. At Dell: built core license-enforcement components — entitlement validation, feature gating, usage metering, and tamper detection backed by AES encryption and SHA cryptographic hashing; ported a C component to a dependency-free pure-Go binary eliminating all cgo risk. Current role: JWT/RBAC, rate limiting, IP blacklisting, HTTP hardening, Stripe PCI-compliant payment integration. Infrastructure-level: PKI/mTLS for service-to-service communication. Operates with a security-first mindset at every system boundary."},
	{"devops_infra", "DevOps & infrastructure", "docker kubernetes aws ecs ecr s3 ses sns terraform ci cd gitlab github actions containerized multi-stage dockerfile prometheus opentelemetry observability",
		"Solid DevOps and infrastructure experience: Dockerized production services with multi-stage builds, Kubernetes deployments, and AWS (ECS, ECR, S3, SES, SNS). Owns release pipelines across GitLab CI and GitHub Actions with automated quality gates and environment promotion. Terraform for infrastructure-as-code. Prometheus + OpenTelemetry for metrics and distributed tracing. Hands-on with container optimisation, environment-based configuration, and zero-downtime deployments."},
	{"performance", "Performance engineering", "performance optimization mongodb query refactor redis caching pagination compression latency throughput 30 percent improvement p1 defects",
		"Consistent performance impact across roles: MongoDB query refactors cutting latency, Redis caching layer reducing DB load, compression and pagination shrinking payload sizes on AWS. At AstraZeneca: ~30% CRM performance improvement, resolved 20+ defects. At Broadcom: improved alarming-system accuracy ~20%. At Turing: resolved 10+ critical defects in a distributed Matching Kanban. Always profiles and measures before optimising; documents the before/after delta."},
	{"education", "Education", "bits pilani meng masters computer science ciitm jaipur btech",
		"M.Eng. in Computer Science from BITS Pilani (Aug 2015–Jul 2017). B.Tech in Computer Science from CIITM, Jaipur (Sep 2008–Jul 2012)."},
	{"certs", "Certifications & recognition", "aws solutions architect associate cloud practitioner go design patterns hackathon winner above and beyond award",
		"Certifications: AWS Certified Solutions Architect – Associate, AWS Certified Cloud Practitioner, Go Design Patterns. Recognition: Hackathon Winner, Above & Beyond Award."},
	{"ai_focus", "AI / RAG / agents focus", "ai rag retrieval augmented generation openai anthropic claude gemini grok mcp tool-use agents semantic search embeddings llm provider abstraction multi-provider tfidf",
		"Deep focus on AI engineering: integrating OpenAI and Anthropic APIs, building retrieval-augmented generation (RAG) pipelines with TF-IDF and embedding-based search, and MCP-style tool-use agent loops. This portfolio itself is a live demonstration: a Go backend with a provider-agnostic LLM abstraction (llm.Provider interface) supporting Anthropic Claude, OpenAI GPT, xAI Grok, and Google Gemini with a single config swap; a TF-IDF retrieval corpus over a curated knowledge base; and a structured tool-use loop that grounds every model response in real data. At the marketplace: OpenAI-powered listing generation, document/résumé parsing, and semantic search via RAG, deployed on AWS."},
	{"contact", "Contact & availability", "email linkedin github open senior staff role hire location delhi relocation relocate",
		"Based in Delhi, India; open to relocation. Open to Senior / Staff backend (and backend+AI) roles. Email nsinghi2011@gmail.com, LinkedIn in/neeraj-singhi-golang, GitHub github.com/neerajsinghi."},
}

var wordRe = regexp.MustCompile(`[a-z0-9+#.]+`)

func tokenize(s string) []string { return wordRe.FindAllString(strings.ToLower(s), -1) }

// docFreq is document frequency, computed once at startup.
var docFreq = func() map[string]int {
	df := map[string]int{}
	for _, d := range KB {
		seen := map[string]bool{}
		for _, t := range tokenize(d.Title + " " + d.Tags + " " + d.Text) {
			if !seen[t] {
				seen[t] = true
				df[t]++
			}
		}
	}
	return df
}()

func scoreDoc(qTokens []string, d Doc) float64 {
	tf := map[string]int{}
	// Index tags and body text normally.
	for _, t := range tokenize(d.Tags + " " + d.Text) {
		tf[t]++
	}
	// Boost title matches — a query matching the title is strongly relevant.
	for _, t := range tokenize(d.Title) {
		tf[t] += 3
	}
	var s float64
	for _, q := range qTokens {
		if c, ok := tf[q]; ok {
			idf := math.Log(1 + float64(len(KB))/float64(maxInt(docFreq[q], 1)))
			s += (1 + math.Log(float64(c))) * idf
		}
	}
	return s
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Retrieve returns the top-k most relevant docs for a query (the RAG step).
func Retrieve(query string, k int) []Doc {
	q := tokenize(query)
	type sd struct {
		d Doc
		s float64
	}
	scored := make([]sd, 0, len(KB))
	for _, d := range KB {
		if s := scoreDoc(q, d); s > 0 {
			scored = append(scored, sd{d, s})
		}
	}
	sort.SliceStable(scored, func(i, j int) bool { return scored[i].s > scored[j].s })
	out := []Doc{}
	for i := 0; i < len(scored) && i < k; i++ {
		out = append(out, scored[i].d)
	}
	return out
}
