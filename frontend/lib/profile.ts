// Display data for the portfolio (the agent's RAG corpus lives in the Go backend).

export const FACTS: [string, string][] = [
  ["experience", "10+ years"],
  ["core", "Go (Golang) · gRPC · microservices"],
  ["also", "Node.js · NestJS · TypeScript"],
  ["cloud", "AWS · Docker · K8s"],
  ["focus now", "AI · RAG · agents"],
  ["based in", "Delhi, India"],
  ["open to", "Senior / Staff / Lead · open to relocation"],
];

export const STACK: { title: string; tags: string[] }[] = [
  { title: "Languages", tags: ["Go", "TypeScript", "JavaScript", "Java", "C (cgo)"] },
  { title: "Backend & Distributed", tags: ["Microservices", "gRPC", "REST", "WebSockets", "SDK design", "Node.js", "NestJS"] },
  { title: "Cloud & Infra", tags: ["AWS (ECS/ECR/S3/SES/SNS)", "Docker", "Kubernetes", "Terraform"] },
  { title: "AI / LLM", tags: ["OpenAI", "RAG", "MCP", "Tool-use agents", "Semantic search"] },
  { title: "Observability & CI/CD", tags: ["Prometheus", "OpenTelemetry", "GitLab CI", "GitHub Actions"] },
  { title: "Data & Security", tags: ["MongoDB", "MySQL", "Redis", "PKI / mTLS", "Encryption & hashing"] },
];

export type Job = { role: string; co: string; when: string; dim?: boolean; pts: string[] };

export const EXPERIENCE: Job[] = [
  {
    role: "Lead Backend Engineer",
    co: "Independent / Freelance",
    when: "Feb 2026 — Present",
    pts: [
      "Leads a 5-engineer team building an AI-enabled consumer marketplace — owns architecture & delivery.",
      "Go + NestJS backend (40+ modules) on MongoDB with a Next.js frontend; OpenAI/RAG features.",
      "JWT/RBAC, rate limiting & HTTP hardening; Stripe, Firebase, Twilio; Redis + AWS, Dockerized.",
    ],
  },
  {
    role: "Senior Software Consultant",
    co: "Dell Technologies (via Objectwin)",
    when: "Oct 2022 — Jan 2026",
    pts: [
      "Built core parts of Dell's internal license-enforcement service: entitlement validation & feature gating.",
      "Usage metering & tamper detection using encryption and cryptographic hashing.",
      "Ported a performance-critical C component to a single pure-Go artifact; owned CI/CD on AWS ECS/ECR; mentored 5 engineers.",
    ],
  },
  {
    role: "Consulting Software Engineer",
    co: "Turing.com",
    when: "Mar 2022 — Aug 2022",
    pts: [
      "Built an AWS-Linux daemon capturing developer activity for customer visibility.",
      "Stabilized the Matching Kanban — resolved 10+ defects across Go, Node/NestJS, React.",
    ],
  },
  {
    role: "Software Development Consultant",
    co: "Truelancer · client AstraZeneca",
    when: "May 2021 — Apr 2022",
    pts: [
      "Delivered user-roster & session features for AstraZeneca's CRM (~30% perf, 20+ defects).",
      "Led code review on the shared repo.",
    ],
  },
  {
    role: "Full-Stack Engineer",
    co: "Freelance",
    when: "Apr 2019 — Aug 2021",
    pts: [
      "Medical News platform as Go microservices + Flutter on Docker/Kubernetes.",
      "Backend systems & APIs for social and stock-trading clients.",
    ],
  },
  {
    role: "R&D Engineer, Software 2",
    co: "Broadcom Inc.",
    when: "Jul 2017 — Apr 2019",
    dim: true,
    pts: [
      "Improved alarming-system accuracy ~20%; built SaaS monitoring probes (SQL, AD, system).",
      "Resolved critical P1 defects with customers.",
    ],
  },
  {
    role: "Software Intern · QA Engineer",
    co: "Cisco · Trignodev",
    when: "2013 — 2017",
    dim: true,
    pts: ["Cisco: IoT/CoAP protocol work on switches. Trignodev: GUI/functional/regression test suites."],
  },
];

export const PROJECTS: { tag: string; name: string; desc: string; stack: string[]; repo?: string }[] = [
  {
    tag: "AI · Marketplace",
    name: "Consumer Marketplace Platform",
    desc: "Mobile-first marketplace with AI-assisted listing generation, document parsing, and semantic search — built and led end to end.",
    stack: ["Go", "NestJS", "Next.js", "MongoDB", "OpenAI / RAG", "Redis", "AWS"],
  },
  {
    tag: "Microservices",
    name: "Medical News Platform",
    desc: "Curated medical-news platform built as Go microservices with a Flutter mobile client, deployed on Docker and Kubernetes.",
    stack: ["Go", "Flutter", "Docker", "Kubernetes"],
  },
  {
    tag: "Open Source · AI Security",
    name: "mcp-gateway",
    desc: "Security gateway for MCP servers: API-key and JWT auth, per-tool allowlists that also filter tools/list, rate limits, audit logs and OpenTelemetry.",
    stack: ["Go", "MCP", "JWT", "OpenTelemetry", "Docker"],
    repo: "https://github.com/neerajsinghi/mcp-gateway",
  },
  {
    tag: "Open Source · System Design",
    name: "system-design-patterns-go",
    desc: "48 system design patterns as small, tested, runnable Go examples: circuit breakers, sagas, consistent hashing, LSM trees, quorums, fencing tokens and more.",
    stack: ["Go", "Distributed systems", "Reliability", "Storage engines"],
    repo: "https://github.com/neerajsinghi/system-design-patterns-go",
  },
  {
    tag: "Open Source · DSA Course",
    name: "dsa-interview-go",
    desc: "A 21-module coding-interview course in Go: a lesson per pattern and 150+ interview problems, each with brute-force and optimal solutions checked by randomized tests.",
    stack: ["Go", "Algorithms", "Data structures", "Dynamic programming"],
    repo: "https://github.com/neerajsinghi/dsa-interview-go",
  },
  {
    tag: "Open Source · Security",
    name: "go-license",
    desc: "Offline Ed25519 license keys with feature gating, clock-rollback detection and a tamper-evident, hash-chained usage log, plus a CLI.",
    stack: ["Go", "Ed25519", "HMAC", "CLI"],
    repo: "https://github.com/neerajsinghi/go-license",
  },
  {
    tag: "Open Source · AI / RAG",
    name: "go-rag-kit",
    desc: "Production-minded RAG toolkit in pure Go: chunking, embeddings, pgvector and Atlas stores, BM25 + vector hybrid retrieval, and a recall/MRR eval harness.",
    stack: ["Go", "RAG", "BM25", "pgvector", "MongoDB Atlas"],
    repo: "https://github.com/neerajsinghi/go-rag-kit",
  },
  {
    tag: "Open Source · Observability",
    name: "llm-otel-go",
    desc: "OpenTelemetry tracing, token metrics and cost tracking for OpenAI and Anthropic calls in Go, including streaming, with a Grafana dashboard.",
    stack: ["Go", "OpenTelemetry", "Prometheus", "Grafana"],
    repo: "https://github.com/neerajsinghi/llm-otel-go",
  },
  {
    tag: "Open Source · Distributed",
    name: "ratelimit",
    desc: "Distributed rate limiting on Redis (GCRA and sliding window) for Go and NestJS, sharing the same Lua scripts so both behave identically.",
    stack: ["Go", "NestJS", "TypeScript", "Redis", "Lua"],
    repo: "https://github.com/neerajsinghi/ratelimit",
  },
  {
    tag: "Open Source · Microservices",
    name: "go-outbox",
    desc: "Transactional outbox relay and Idempotency-Key middleware for Go microservices, with leases, retries and dead-lettering on PostgreSQL, MySQL and SQLite.",
    stack: ["Go", "SQL", "Outbox pattern", "Idempotency"],
    repo: "https://github.com/neerajsinghi/go-outbox",
  },
  {
    tag: "Open Source · Performance",
    name: "cgo-to-pure-go",
    desc: "Case study on replacing a cgo dependency with pure Go: golden vectors, differential fuzzing, real benchmarks and a build-tag migration path.",
    stack: ["Go", "cgo", "C", "Fuzzing", "Benchmarks"],
    repo: "https://github.com/neerajsinghi/cgo-to-pure-go",
  },
  {
    tag: "Open Source · Template",
    name: "portfolio-rag-template",
    desc: "This site as a reusable GitHub template: a Next.js portfolio with a Go RAG agent that answers visitors' questions about you.",
    stack: ["Go", "Next.js", "RAG", "Docker"],
    repo: "https://github.com/neerajsinghi/portfolio-rag-template",
  },
];

// Rendered visibly on the homepage and as FAQPage JSON-LD; search engines and
// AI assistants quote these answers for "who is Neeraj Singhi"-style queries.
export const FAQ: { q: string; a: string }[] = [
  {
    q: "Who is Neeraj Singhi?",
    a: "Neeraj Singhi is a senior backend engineer and Go (Golang) developer based in Delhi, India, with 10+ years of experience. He currently works as a Lead Backend Engineer, leading a 5-engineer team building an AI-enabled consumer marketplace, and previously spent over three years as a Senior Software Consultant at Dell Technologies.",
  },
  {
    q: "What does Neeraj Singhi specialise in?",
    a: "Production Go services, microservices and gRPC/REST APIs, containerized systems on AWS, Docker and Kubernetes, security-focused backend work such as license enforcement and tamper detection, and AI integration with OpenAI, retrieval-augmented generation (RAG) and tool-based agent workflows.",
  },
  {
    q: "Is Neeraj Singhi a Golang developer?",
    a: "Yes. Go is his primary language. He has built Go microservices, SDKs and a pure-Go port of a performance-critical C component at Dell, and his current marketplace backend is written in Go and NestJS.",
  },
  {
    q: "Does Neeraj Singhi work with Node.js and TypeScript?",
    a: "Yes. Alongside Go he builds Node.js services with NestJS and TypeScript, including a 40+ module NestJS and Go backend in his current role and Node.js/NestJS work at Turing.com, plus Next.js and React on the frontend.",
  },
  {
    q: "What roles is Neeraj Singhi open to?",
    a: "Senior, Staff and Lead backend engineering roles, especially Go (Golang), distributed systems and AI/RAG work. He is based in Delhi, India and is open to remote work and relocation.",
  },
  {
    q: "How can I contact Neeraj Singhi?",
    a: "Email nsinghi2011@gmail.com, connect on LinkedIn at linkedin.com/in/neeraj-singhi-golang, or see his code at github.com/neerajsinghi. His résumé is available as a PDF on neerajsinghi.com.",
  },
];

export const CHIPS = [
  "Show me Neeraj's latest blogs",
  "What engineering services does he offer?",
  "What did Neeraj do at Dell?",
  "Is he a fit for a Staff backend role?",
  "Show his AI / RAG work",
];

export const LINKS = {
  email: "nsinghi2011@gmail.com",
  linkedin: "https://www.linkedin.com/in/neeraj-singhi-golang",
  github: "https://github.com/neerajsinghi",
};

export const API_BASE = process.env.NEXT_PUBLIC_API_BASE || "http://localhost:8080";
