# neerajsinghi.com

Personal portfolio for **Neeraj Singhi** with a live AI agent that answers questions
about his work — grounded in his résumé, GitHub and LinkedIn via **RAG** and a
**tool-use (MCP-style) loop**.

```
neeraj-portfolio/
├── backend/         Go 1.26 API — agent loop, RAG, tools, live GitHub
│   ├── main.go        HTTP server: /api/chat (SSE), /api/github, /api/health, CORS
│   ├── agent.go       Anthropic client + tool-use loop (streams events)
│   ├── kb.go          knowledge base + TF-IDF retrieval (the RAG step)
│   └── tools.go       tool registry: search_profile, list_projects, get_links, get_github_repos
└── frontend/        Next.js (App Router, TS) — the site + agent console
    ├── app/           layout, page, global styles
    ├── components/    AgentConsole (SSE client), Topology, GithubStrip
    └── lib/profile.ts display data + API base
```

## How the agent works

1. The browser sends the chat history to the Go backend (`POST /api/chat`).
2. The backend runs an **Anthropic tool-use loop**. The model decides which tool to
   call — the same contract an **MCP** server exposes (`name`, `description`, `input_schema`).
3. `search_profile` performs **RAG**: TF-IDF retrieval over the knowledge base in `kb.go`.
4. Each step is streamed back as **Server-Sent Events** (`tool`, `sources`, `text`), so the
   console shows the agent's tool calls and retrieved sources live before the answer.
5. The Anthropic API key never leaves the server.

To expose the tools over a **real MCP server**, register the same names/schemas from
`tools.go` and forward calls to `ExecuteTool`.

## Run it locally

**Prerequisites:** Go 1.26+, Node 18+, an Anthropic API key.

### 1. Backend

```bash
cd backend
cp .env.example .env          # then put your ANTHROPIC_API_KEY in .env
go run .                      # serves http://localhost:8080
```

### 2. Frontend

```bash
cd frontend
cp .env.local.example .env.local   # NEXT_PUBLIC_API_BASE=http://localhost:8080
npm install
npm run dev                        # http://localhost:3000
```

Open http://localhost:3000 and ask the agent something.

## Configuration

| Variable             | Where      | Default              | Purpose                                   |
| -------------------- | ---------- | -------------------- | ----------------------------------------- |
| `ANTHROPIC_API_KEY`  | backend    | —                    | required; kept server-side                |
| `ANTHROPIC_MODEL`    | backend    | `claude-sonnet-4-6`  | model used by the agent                   |
| `GITHUB_USER`        | backend    | `neerajsinghi`       | account for the live repo strip + tool    |
| `GITHUB_TOKEN`       | backend    | —                    | optional; raises GitHub rate limit        |
| `ALLOWED_ORIGIN`     | backend    | `*`                  | CORS origin for the frontend in prod      |
| `NEXT_PUBLIC_API_BASE` | frontend | `http://localhost:8080` | backend URL the browser calls          |

## Customize

- **Profile facts the agent uses** → `backend/kb.go` (`KB`).
- **What the site renders** → `frontend/lib/profile.ts`.
- **Résumé download** → replace `frontend/public/Neeraj_Singhi_Resume.pdf`.

## Deploy

- **Backend:** any container host (Fly.io, Render, Cloud Run, ECS). `go build -o server .`
- **Frontend:** Vercel/Netlify (set `NEXT_PUBLIC_API_BASE` to the deployed backend, and
  set the backend's `ALLOWED_ORIGIN` to your site's URL).
