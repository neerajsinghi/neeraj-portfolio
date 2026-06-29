// Display data for the portfolio (the agent's RAG corpus lives in the Go backend).

export const FACTS: [string, string][] = [
  ["experience", "10+ years"],
  ["core", "Go · gRPC · microservices"],
  ["cloud", "AWS · Docker · K8s"],
  ["focus now", "AI · RAG · agents"],
  ["based in", "Delhi, India"],
  ["open to", "Senior / Staff · open to relocation"],
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

export const PROJECTS: { tag: string; name: string; desc: string; stack: string[] }[] = [
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
];

export const CHIPS = [
  "What did Neeraj do at Dell?",
  "Is he a fit for a Staff backend role?",
  "Show his AI / RAG work",
  "What's his security / PKI experience?",
];

export const LINKS = {
  email: "nsinghi2011@gmail.com",
  linkedin: "https://www.linkedin.com/in/neeraj-singhi-golang",
  github: "https://github.com/neerajsinghi",
};

export const API_BASE = process.env.NEXT_PUBLIC_API_BASE || "http://localhost:8080";
