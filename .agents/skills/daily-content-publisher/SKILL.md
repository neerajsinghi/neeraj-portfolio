---
name: daily-content-publisher
description: Draft, review, adapt, schedule, or publish original technical blogs and social posts for Neeraj Singhi. Use for daily posts, blog ideas, Medium, LinkedIn, DEV, cross-posting, content calendars, or publishing automation.
metadata:
  version: 1.0.0
---

# Daily Content Publisher

## Workflow

1. Pick one narrow, advanced problem for senior backend and Go engineers from the editorial pillars below.
2. Check existing drafts before selecting a topic. Avoid repeating the same thesis within 90 days.
3. Research claims that may have changed. Prefer official documentation and primary sources.
4. Generate a canonical bundle with `go run ./cmd/content generate` from `backend`.
5. Review the JSON bundle against the quality gate below. Edit the bundle directly when corrections are needed.
6. Present the draft and adaptations for approval.
7. Publish only after explicit approval. Use `--approve` for public writes.
8. Record returned URLs. Export the Markdown and social copy for manual-only platforms.

## Editorial Positioning

Neeraj's authority comes from 10+ years building production Go services, microservices, SDKs, distributed systems, security controls, data-intensive backends, and AWS delivery systems. Content should reinforce that positioning.

- Keep at least 70% of each article focused on backend mechanics, architecture, failure modes, and tradeoffs.
- Prefer production-quality Go examples over language-neutral pseudocode.
- Write for senior engineers. Skip basic syntax, generic introductions, career advice, frontend tutorials, and technology listicles.
- Cover new technology only through a concrete backend problem, operational constraint, or migration decision.
- Cover AI when the core subject is backend architecture: RAG, retrieval, tool execution, evaluation, security, reliability, observability, or cost.
- Never imply Neeraj used a specific design at an employer unless that fact exists in the profile corpus.

## Topic Pillars

Rotate across these areas while selecting a different thesis within each pillar:

- Go runtime and internals: scheduler, garbage collector, escape analysis, allocations, memory model, compiler, reflection.
- Go concurrency: cancellation, backpressure, worker coordination, race prevention, goroutine leaks, bounded parallelism.
- Go concepts at scale: interfaces, generics, method sets, errors, package boundaries, API design, testing seams.
- Networking and APIs: `net/http`, connection pools, gRPC, WebSockets, streaming, protocol evolution, graceful shutdown.
- Distributed systems: consistency, idempotency, retries, ordering, leases, sagas, outbox, recovery.
- Data systems: MongoDB modeling and indexes, MySQL transactions, Redis caching, pagination, query planning.
- Performance: profiling, latency budgets, contention, batching, compression, capacity, load shedding.
- Security: JWT/RBAC, PKI/mTLS, rate limiting, entitlements, tamper evidence, encryption, secret boundaries.
- Reliability and observability: OpenTelemetry, Prometheus, SLOs, tracing, degradation, incident diagnosis.
- AWS and delivery: Lambda, API Gateway, ECS/ECR, S3, events, Docker, Kubernetes, Terraform, CI/CD.
- Systems engineering: pure Go versus cgo, binary distribution, daemons, metering, cryptographic integrity.
- Modern backend technology: eBPF, WebAssembly components, event streaming, platform engineering, new Go capabilities.
- AI backend engineering: RAG, semantic retrieval, tool execution, evaluation, model gateways, prompt security, token economics.

## Quality Gate

- The title makes a specific promise and avoids clickbait.
- The article has one defensible thesis, concrete examples, and meaningful tradeoffs.
- Claims about Neeraj are supported by repository profile or knowledge-base facts.
- Time-sensitive factual claims have primary sources.
- Code is coherent and does not expose credentials or unsafe defaults.
- The conclusion adds a decision framework or next action instead of repeating the introduction.
- The description is useful in search results; tags are relevant and no more than four.
- LinkedIn copy is useful without requiring the click and stays below 2,500 characters.
- Social copy stays below 280 characters.
- The article URL is canonical on syndication targets that support canonical metadata.

## Publishing Rules

- No public network write without current, explicit approval.
- DEV: default to an unpublished draft. `--approve` changes `published` to true.
- LinkedIn: `--approve` is mandatory because API creation publishes immediately.
- Medium: export/import manually; Medium does not issue new API integration tokens.
- Substack: export manually unless an official publishing API becomes available.
- Do not publish identical full articles everywhere without canonical handling.
- Do not send secrets, unpublished client information, prompts, or private résumé details to publishing APIs.

## Required Environment

- Generation: `ANTHROPIC_API_KEY`; optional `ANTHROPIC_MODEL`.
- DEV: `DEVTO_API_KEY`.
- LinkedIn: `LINKEDIN_ACCESS_TOKEN`, `LINKEDIN_AUTHOR_URN`, `LINKEDIN_VERSION` in `YYYYMM` format, and the `w_member_social` or appropriate organization permission.

See [platforms.md](references/platforms.md) before adding another connector.