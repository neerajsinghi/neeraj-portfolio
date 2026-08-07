# neerajsinghi.com

Personal portfolio for **Neeraj Singhi** — a live site with an AI agent that answers questions about his work, grounded in his résumé, GitHub, and experience via **RAG** and a **tool-use loop**.

## Architecture

```
neeraj-portfolio/
├── backend/                   Go 1.24 API
│   ├── cmd/server/main.go       HTTP server — /api/chat (SSE), /api/github, /api/health
│   └── internal/
│       ├── agent/               Tool-use loop (provider-agnostic)
│       ├── kb/                  Knowledge base + TF-IDF retrieval (RAG)
│       ├── tools/               Tool registry: 8 tools for profile, projects, skills, etc.
│       ├── github/              Live GitHub repo fetcher
│       └── llm/                 Provider interface + Anthropic, OpenAI, Grok, Gemini impls
├── frontend/                  Next.js 14 (App Router, TypeScript)
│   ├── app/                     Layout + page
│   ├── components/
│   │   ├── sections/            Hero, About, Stack, Experience, Projects, Contact
│   │   ├── AgentConsole.tsx     SSE chat client
│   │   ├── GithubStrip.tsx      Live repo strip
│   │   └── Topology.tsx         Background animation
│   ├── lib/api.ts               Typed API client
│   └── types/                   Shared types (Repo, ChatItem)
├── infrastructure/terraform/  AWS infra (Amplify + Lambda default, optional EKS full mode)
│   └── modules/               lambda_api · amplify · vpc · ecr · eks
├── backend/k8s/               Kubernetes manifests (EKS deploy)
├── .github/workflows/         CI/CD pipelines
└── scripts/                   Bootstrap helpers
```

## How the agent works

1. The browser sends the chat history to `POST /api/chat`.
2. The Go backend runs a **tool-use loop** via the `llm.Provider` interface — swap one line in `cmd/server/main.go` to change the model provider.
3. `search_profile` performs **RAG**: TF-IDF retrieval over 17 knowledge-base documents in `internal/kb/`.
4. Each loop step streams back as **Server-Sent Events** (`tool`, `sources`, `text`), so the console shows tool calls and retrieved sources live before the final answer.
5. API keys never leave the server.

### Available tools

| Tool                     | Description                                 |
| ------------------------ | ------------------------------------------- |
| `search_profile`         | TF-IDF RAG over the KB (top-N configurable) |
| `list_projects`          | Projects with URLs and tech stack           |
| `get_links`              | Social / contact links                      |
| `get_github_repos`       | Live repos from GitHub API                  |
| `get_skills`             | Skills by category                          |
| `get_education`          | Education history                           |
| `get_certifications`     | Certifications                              |
| `get_experience_summary` | Work experience timeline                    |

### Supported LLM providers

All share the `llm.Provider` interface — zero changes needed in the agent loop.

| Provider      | Package                  | Default model       |
| ------------- | ------------------------ | ------------------- |
| Anthropic     | `internal/llm/anthropic` | `claude-sonnet-4-6` |
| OpenAI        | `internal/llm/openai`    | `gpt-4o`            |
| xAI Grok      | `internal/llm/grok`      | `grok-3`            |
| Google Gemini | `internal/llm/gemini`    | `gemini-2.0-flash`  |

## Run locally

**Prerequisites:** Go 1.24+, Node 20+, an API key for your chosen LLM provider.

### Backend

```bash
cd backend
cp .env.example .env          # set ANTHROPIC_API_KEY (and/or others)
go run ./cmd/server           # http://localhost:8080
```

### Frontend

```bash
cd frontend
cp .env.local.example .env.local   # NEXT_PUBLIC_API_BASE=http://localhost:8080
npm install
npm run dev                        # http://localhost:3000
```

### Docker Compose

```bash
cp backend/.env.example backend/.env   # add your key(s)
docker compose up --build              # backend only (frontend is on Amplify)

# optional: run frontend container locally too
docker compose --profile local-frontend up --build
```

## Configuration

### Backend env vars

| Variable            | Default             | Purpose                                   |
| ------------------- | ------------------- | ----------------------------------------- |
| `ANTHROPIC_API_KEY` | —                   | Required for Anthropic provider           |
| `ANTHROPIC_MODEL`   | `claude-sonnet-4-6` | Model override                            |
| `OPENAI_API_KEY`    | —                   | Required for OpenAI provider              |
| `XAI_API_KEY`       | —                   | Required for Grok provider                |
| `GEMINI_API_KEY`    | —                   | Required for Gemini provider              |
| `GITHUB_USER`       | `neerajsinghi`      | GitHub account for repo strip             |
| `GITHUB_TOKEN`      | —                   | Optional — raises GitHub rate limit       |
| `ALLOWED_ORIGIN`    | `*`                 | CORS origin (set to frontend URL in prod) |
| `PORT`              | `8080`              | HTTP listen port                          |

### Frontend env vars

| Variable               | Default                 | Purpose                       |
| ---------------------- | ----------------------- | ----------------------------- |
| `NEXT_PUBLIC_API_BASE` | `http://localhost:8080` | Backend URL the browser calls |

## Customise

- **Agent knowledge** → `backend/internal/kb/kb.go`
- **Tools / capabilities** → `backend/internal/tools/tools.go`
- **Site content** → `frontend/lib/profile.ts`
- **Switch LLM provider** → change one line in `backend/cmd/server/main.go`
- **Resume PDF** → replace `frontend/public/Neeraj_Singhi_Resume.pdf`

## Infrastructure & CI/CD

Deployed on AWS using Terraform. No long-lived AWS credentials — GitHub Actions authenticates via **OIDC**.

```
Default  → AWS Lambda + API Gateway backend (serverless)
Optional → EKS (Kubernetes) backend via ECR image (full mode)
Frontend → AWS Amplify       (Next.js SSR)
Secrets  → AWS SSM Parameter Store (SecureString)
```

### Cost modes

- `infrastructure_mode="budget"` (default): provisions Lambda + API Gateway + Amplify. This is the mode to target monthly spend under ~$5 at low traffic.
- `infrastructure_mode="full"`: provisions VPC + EKS + ECR + Amplify (significantly higher monthly baseline).

### Custom domain (Route53)

Set these in `terraform.tfvars`:

- `enable_custom_domains = true`
- `root_domain = "neerajsinghi.com"`
- `frontend_subdomain = ""` (apex domain) or e.g. `"app"`
- `api_subdomain = "api"` (backend becomes `api.neerajsinghi.com`)

When enabled:

- Frontend host is auto-mapped in Amplify.
- API Gateway is mapped to your API custom domain with ACM cert + Route53 alias records.
- Backend CORS `ALLOWED_ORIGIN` is auto-derived from the frontend host.

### First-time setup

```bash
# 1. Store API keys in SSM (one-time)
export ANTHROPIC_API_KEY=...
export OPENAI_API_KEY=...
export XAI_API_KEY=...
export GEMINI_API_KEY=...
export ALLOWED_ORIGIN=https://your-amplify-url
./scripts/bootstrap-secrets.sh

# 2. Provision AWS infrastructure
cd infrastructure/terraform
cp terraform.tfvars.example terraform.tfvars   # fill in values
terraform init
cd ../..
./scripts/build-lambda.sh
cd infrastructure/terraform
terraform apply

# Note: keep infrastructure_mode="budget" for low-cost operation.
# Switch to "full" only when you explicitly want EKS.

# — or — run the GitHub Actions workflow:
# Actions → "Bootstrap — Store secrets in SSM" → Run workflow
```

### GitHub repository configuration

Add these under **Settings → Secrets and variables → Actions**:

| Key                 | Type     | Value                                               |
| ------------------- | -------- | --------------------------------------------------- |
| `AWS_ROLE_ARN`      | Secret   | ARN from `terraform output github_actions_role_arn` |
| `AMPLIFY_APP_ID`    | Secret   | From `terraform output amplify_app_id`              |
| `ANTHROPIC_API_KEY` | Secret   | Your Anthropic key                                  |
| `OPENAI_API_KEY`    | Secret   | Your OpenAI key                                     |
| `XAI_API_KEY`       | Secret   | Your xAI key                                        |
| `GEMINI_API_KEY`    | Secret   | Your Gemini key                                     |
| `ALLOWED_ORIGIN`    | Secret   | Your Amplify frontend URL                           |
| `AWS_REGION`        | Variable | e.g. `ap-south-1`                                   |

### Environment matrix

SSM Parameter Store (path prefix: `/neeraj-portfolio`):

- `anthropic-api-key`
- `openai-api-key`
- `xai-api-key`
- `gemini-api-key`
- `allowed-origin` (used when `enable_custom_domains=false`; auto-derived in terraform otherwise)

Terraform variables (`terraform.tfvars`):

- `aws_region`, `project`, `environment`
- `infrastructure_mode`, `enable_lambda_backend`
- `github_repo`, `github_org_repo`, `github_access_token`
- `enable_custom_domains`, `root_domain`, `frontend_subdomain`, `frontend_enable_www`, `api_subdomain`
- `lambda_zip_path`, `lambda_memory_size`, `lambda_timeout`
- `amplify_branch`, `amplify_frontend_app_root`, `amplify_enable_auto_build`

### GitHub Actions workflows

| Workflow                | Trigger                        | Jobs                                         |
| ----------------------- | ------------------------------ | -------------------------------------------- |
| `backend.yml`           | Push to `main` → `backend/**`  | Build+vet → build lambda zip → lambda deploy |
| `frontend.yml`          | Push to `main` → `frontend/**` | Type-check+build → Amplify deploy            |
| `bootstrap-secrets.yml` | Manual (`workflow_dispatch`)   | GitHub secrets → SSM Parameter Store         |

In `budget` mode, backend deploy targets Lambda and frontend deploy targets Amplify.

