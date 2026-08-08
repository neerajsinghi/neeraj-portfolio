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
│       ├── tools/               Tool registry: profile, projects, live blogs, and services
│       ├── github/              Live GitHub repo fetcher
│       └── llm/                 Provider interface + Anthropic, OpenAI, Grok, Gemini impls
├── frontend/                  Next.js 16 (App Router, TypeScript)
│   ├── app/                     Layout + page
│   ├── components/
│   │   ├── sections/            Hero, About, Stack, Experience, Projects, Contact
│   │   ├── AgentConsole.tsx     SSE chat client
│   │   ├── GithubStrip.tsx      Live repo strip
│   │   └── Topology.tsx         Background animation
│   ├── lib/api.ts               Typed API client
│   └── types/                   Shared types (Repo, ChatItem)
├── admin/                     Separate Next.js editorial console
│   ├── app/                     Cognito callback + noindex metadata
│   ├── components/              Draft editor, preview, publish/archive controls
│   └── lib/                     PKCE auth + typed admin API
├── infrastructure/terraform/  AWS infra (Amplify + Lambda default, optional EKS full mode)
│   └── modules/               lambda_api · amplify · vpc · ecr · eks
├── backend/k8s/               Kubernetes manifests (EKS deploy)
├── .github/workflows/         CI/CD pipelines
└── scripts/                   Bootstrap helpers
```

## How the agent works

1. The browser sends the chat history to `POST /api/chat`.
2. The Go backend runs a **tool-use loop** via the `llm.Provider` interface — swap one line in `cmd/server/main.go` to change the model provider.
3. `search_profile` performs **RAG** over the curated profile corpus, while `list_blogs` and `search_blogs` retrieve currently published MongoDB posts.
4. Each loop step streams back as **Server-Sent Events** (`tool`, `sources`, `text`), so the console shows tool calls and retrieved sources live before the final answer.
5. API keys never leave the server.

### Available tools

| Tool                     | Description                                 |
| ------------------------ | ------------------------------------------- |
| `search_profile`         | TF-IDF RAG over the KB (top-N configurable) |
| `list_projects`          | Projects with URLs and tech stack           |
| `get_links`              | Social / contact links                      |
| `get_github_repos`       | Live repos from GitHub API                  |
| `list_blogs`             | Latest published posts with canonical links |
| `search_blogs`           | Topic search over published blog content    |
| `get_services`           | Grounded engineering service areas          |
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

# MongoDB + backend + both frontends (Cognito values are required for admin login)
docker compose --profile local-frontend --profile local-admin up --build
```

## Configuration

### Backend env vars

| Variable               | Default             | Purpose                                   |
| ---------------------- | ------------------- | ----------------------------------------- |
| `ANTHROPIC_API_KEY`    | —                   | Required for Anthropic provider           |
| `ANTHROPIC_MODEL`      | `claude-sonnet-4-6` | Model override                            |
| `OPENAI_API_KEY`       | —                   | Required for OpenAI provider              |
| `XAI_API_KEY`          | —                   | Required for Grok provider                |
| `GEMINI_API_KEY`       | —                   | Required for Gemini provider              |
| `GITHUB_USER`          | `neerajsinghi`      | GitHub account for repo strip             |
| `GITHUB_TOKEN`         | —                   | Optional — raises GitHub rate limit       |
| `ALLOWED_ORIGIN`       | `*`                 | CORS origin (set to frontend URL in prod) |
| `PORT`                 | `8080`              | HTTP listen port                          |
| `MONGODB_URI`          | —                   | MongoDB Atlas or local connection string  |
| `MONGODB_DATABASE`     | `neeraj_portfolio`  | Blog database name                        |
| `COGNITO_REGION`       | —                   | Cognito user pool region                  |
| `COGNITO_USER_POOL_ID` | —                   | Admin user pool ID                        |
| `COGNITO_CLIENT_ID`    | —                   | Admin web app client ID                   |

### Frontend env vars

| Variable               | Default                 | Purpose                       |
| ---------------------- | ----------------------- | ----------------------------- |
| `NEXT_PUBLIC_API_BASE` | `http://localhost:8080` | Backend URL the browser calls |

The `admin/` app additionally receives `NEXT_PUBLIC_COGNITO_DOMAIN`, `NEXT_PUBLIC_COGNITO_CLIENT_ID`, and `NEXT_PUBLIC_COGNITO_REDIRECT_URI` from Terraform/Amplify.

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
- Admin is deployed as a separate Amplify app at `admin.neerajsinghi.com`.
- Cognito uses authorization code + PKCE, mandatory TOTP MFA, token revocation, and groups named `admin` and `editor`.

### First-time setup

```bash
# 1. Store API keys in SSM (one-time)
export ANTHROPIC_API_KEY=...
export OPENAI_API_KEY=...
export XAI_API_KEY=...
export GEMINI_API_KEY=...
export ALLOWED_ORIGIN=https://your-amplify-url
export MONGODB_URI='mongodb+srv://...'
./scripts/bootstrap-secrets.sh

# 2. Provision AWS infrastructure
cd infrastructure/terraform
cp terraform.tfvars.example terraform.tfvars   # fill in values
terraform init
cd ../..
./scripts/build-lambda.sh
cd infrastructure/terraform
terraform apply

# Create the first operator without putting a password in Terraform or source.
POOL_ID=$(terraform output -raw cognito_user_pool_id)
aws cognito-idp admin-create-user --user-pool-id "$POOL_ID" --username you@example.com
aws cognito-idp admin-add-user-to-group --user-pool-id "$POOL_ID" --username you@example.com --group-name admin

# Note: keep infrastructure_mode="budget" for low-cost operation.
# Switch to "full" only when you explicitly want EKS.

# — or — run the GitHub Actions workflow:
# Actions → "Bootstrap — Store secrets in SSM" → Run workflow
```

### GitHub repository configuration

Add these under **Settings → Secrets and variables → Actions**:

| Key                    | Type     | Value                                               |
| ---------------------- | -------- | --------------------------------------------------- |
| `AWS_ROLE_ARN`         | Secret   | ARN from `terraform output github_actions_role_arn` |
| `AMPLIFY_APP_ID`       | Secret   | From `terraform output amplify_app_id`              |
| `ADMIN_AMPLIFY_APP_ID` | Secret   | From `terraform output admin_amplify_app_id`        |
| `MONGODB_URI`          | Secret   | Atlas SRV URI for the Lambda runtime                |
| `ANTHROPIC_API_KEY`    | Secret   | Your Anthropic key                                  |
| `OPENAI_API_KEY`       | Secret   | Your OpenAI key                                     |
| `XAI_API_KEY`          | Secret   | Your xAI key                                        |
| `GEMINI_API_KEY`       | Secret   | Your Gemini key                                     |
| `ALLOWED_ORIGIN`       | Secret   | Your Amplify frontend URL                           |
| `AWS_REGION`           | Variable | e.g. `ap-south-1`                                   |

### Environment matrix

SSM Parameter Store (path prefix: `/neeraj-portfolio`):

- `anthropic-api-key`
- `openai-api-key`
- `xai-api-key`
- `gemini-api-key`
- `anthropic-model`
- `openai-model`
- `grok-model`
- `gemini-model`
- `github-user`
- `github-token` (optional)
- `port`
- `allowed-origin` (used when `enable_custom_domains=false`; auto-derived in terraform otherwise)
- `mongodb-uri`
- `mongodb-database`

Terraform variables (`terraform.tfvars`):

- `aws_region`, `project`, `environment`
- `infrastructure_mode`, `enable_lambda_backend`
- `github_repo`, `github_org_repo`, `github_access_token`
- `enable_custom_domains`, `root_domain`, `frontend_subdomain`, `frontend_enable_www`, `api_subdomain`
- `admin_subdomain`, `admin_frontend_app_root`
- `lambda_zip_path`, `lambda_memory_size`, `lambda_timeout`
- `amplify_branch`, `amplify_frontend_app_root`, `amplify_enable_auto_build`

### GitHub Actions workflows

| Workflow                | Trigger                        | Jobs                                              |
| ----------------------- | ------------------------------ | ------------------------------------------------- |
| `backend.yml`           | Push to `main` → `backend/**`  | Build+vet → build lambda zip → lambda deploy      |
| `frontend.yml`          | Push to `main` → `frontend/**` | Type-check+build → Amplify deploy                 |
| `admin.yml`             | Push to `main` → `admin/**`    | Type-check+build → admin Amplify deploy           |
| `bootstrap-secrets.yml` | Manual (`workflow_dispatch`)   | GitHub secrets → SSM Parameter Store              |
| `content-draft.yml`     | Daily at 08:00 IST or manual   | Generate an article bundle and open a review PR   |
| `content-publish.yml`   | Manual approval only           | Publish an approved bundle to DEV and/or LinkedIn |

In `budget` mode, backend deploy targets Lambda and frontend deploy targets Amplify.

### MongoDB Atlas networking

Use an Atlas M0 free cluster for the blog collections and store its SRV connection string only in SSM/GitHub Secrets. The Lambda intentionally stays outside a VPC, which gives it outbound internet access without a NAT gateway (a NAT gateway alone would exceed the budget). Atlas must therefore allow the Lambda's changing public egress addresses; for this low-risk portfolio dataset, use an Atlas database user limited to the blog database, TLS, a strong generated password, and Atlas network access appropriate to serverless egress. Do not place private or authentication data in blog documents. If fixed IP allowlisting becomes mandatory, use Atlas Data API/HTTPS or accept the cost of a VPC egress design instead of adding an unplanned NAT gateway.

## Daily content publishing

The content workflow is approval-first. The daily action generates a canonical article bundle under `content/drafts/` and opens a pull request; it does not publish. After review and merge, run **Content - Publish Approved** with the bundle path and approval checkbox.

Generate or preview locally from `backend`:

```bash
go run ./cmd/content generate \
  --topic "Idempotency patterns for Go services" \
  --output ../content/drafts/idempotency-patterns.json

# Creates a private DEV draft; LinkedIn is skipped without --approve.
go run ./cmd/content publish \
  --file ../content/drafts/idempotency-patterns.json \
  --platforms devto
```

Public publishing requires `--approve`. Configure `DEVTO_API_KEY` and, for LinkedIn, `LINKEDIN_ACCESS_TOKEN`, `LINKEDIN_AUTHOR_URN`, and `LINKEDIN_VERSION`. Medium remains a manual import because Medium no longer issues new API integration tokens. See the `daily-content-publisher` skill for the editorial and platform rules.

For GitHub Actions, store `DEVTO_API_KEY` and `LINKEDIN_ACCESS_TOKEN` as repository secrets, and `LINKEDIN_AUTHOR_URN`, `LINKEDIN_VERSION`, and `ANTHROPIC_MODEL` as repository variables. The daily workflow also needs the repository setting **Allow GitHub Actions to create and approve pull requests** enabled. Tokens must never be committed to a draft bundle.

