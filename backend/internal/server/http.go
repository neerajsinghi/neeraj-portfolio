package server

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"neeraj-portfolio/backend/internal/agent"
	"neeraj-portfolio/backend/internal/github"
	"neeraj-portfolio/backend/internal/llm"
)

// NewHandler builds the backend HTTP handler with CORS middleware.
func NewHandler(prov llm.Provider) http.Handler {
	ag := agent.New(prov)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "model": prov.ModelName()})
	})
	mux.HandleFunc("/api/chat", chatHandler(ag))
	mux.HandleFunc("/api/github", reposHandler)
	return withCORS(mux)
}

// withCORS allows the Next.js dev/prod origin. Set ALLOWED_ORIGIN in prod.
func withCORS(next http.Handler) http.Handler {
	origin := os.Getenv("ALLOWED_ORIGIN")
	if origin == "" {
		origin = "*"
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type chatRequest struct {
	Messages []agent.Turn `json:"messages"`
}

// chatHandler streams the agent run as Server-Sent Events.
func chatHandler(ag *agent.Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req chatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Messages) == 0 {
			http.Error(w, "expected JSON {messages:[{role,content}]}", http.StatusBadRequest)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")

		emit := func(event string, data any) {
			payload, _ := json.Marshal(data)
			_, _ = w.Write([]byte("event: " + event + "\ndata: "))
			_, _ = w.Write(payload)
			_, _ = w.Write([]byte("\n\n"))
			flusher.Flush()
		}

		if err := ag.Run(req.Messages, emit); err != nil {
			emit("error", map[string]string{"message": err.Error()})
		}
	}
}

// reposHandler returns the cached public repo list for the live "from GitHub" strip.
func reposHandler(w http.ResponseWriter, _ *http.Request) {
	repos, err := github.FetchRepos()
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{"user": github.User(), "repos": []github.Repo{}, "error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": github.User(), "repos": repos})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// LoadDotEnv loads KEY=VALUE lines from a .env file if present, without
// overriding variables already set in the environment.
func LoadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k, v = strings.TrimSpace(k), strings.Trim(strings.TrimSpace(v), `"'`)
		if _, exists := os.LookupEnv(k); !exists {
			_ = os.Setenv(k, v)
		}
	}
}
