package server

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"neeraj-portfolio/backend/internal/agent"
	"neeraj-portfolio/backend/internal/auth"
	"neeraj-portfolio/backend/internal/blog"
	"neeraj-portfolio/backend/internal/github"
	"neeraj-portfolio/backend/internal/llm"
	"neeraj-portfolio/backend/internal/publisher"
)

// NewHandler builds the backend HTTP handler with CORS middleware.
func NewHandler(prov llm.Provider) http.Handler {
	var blogStore blog.Store
	if os.Getenv("MONGODB_URI") != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		store, err := blog.NewMongoStore(ctx, os.Getenv("MONGODB_URI"), os.Getenv("MONGODB_DATABASE"))
		if err != nil {
			log.Printf("blog storage unavailable: %v", err)
		} else {
			blogStore = store
		}
	}
	var authenticator auth.Authenticator
	if os.Getenv("COGNITO_USER_POOL_ID") != "" {
		cognito, err := auth.NewCognito(os.Getenv("COGNITO_REGION"), os.Getenv("COGNITO_USER_POOL_ID"), os.Getenv("COGNITO_CLIENT_ID"))
		if err != nil {
			log.Printf("admin authentication unavailable: %v", err)
		} else {
			authenticator = cognito
		}
	}
	return NewHandlerWithDependencies(prov, blogStore, authenticator)
}

// currentAPIVersion is the canonical, supported API version. Routes are
// mounted under /api/{currentAPIVersion}/... and also aliased at the
// legacy unversioned /api/... paths for backward compatibility.
const currentAPIVersion = "v1"

// NewHandlerWithDependencies builds a handler with injectable blog dependencies for tests.
func NewHandlerWithDependencies(prov llm.Provider, blogStore blog.Store, authenticator auth.Authenticator) http.Handler {
	ag := agent.NewWithBlogStore(prov, blogStore)
	publishClient := publisher.NewClient(&http.Client{Timeout: 30 * time.Second})

	healthHandler := func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"ok": true, "model": prov.ModelName(), "api_version": currentAPIVersion,
			"blog_configured": blogStore != nil, "admin_auth_configured": authenticator != nil,
		})
	}

	mux := http.NewServeMux()
	registerVersioned(mux, "", "/api/health", healthHandler)
	registerVersioned(mux, "", "/api/chat", chatHandler(ag))
	registerVersioned(mux, "", "/api/github", reposHandler)
	registerVersioned(mux, "GET", "/api/blogs", publicListHandler(blogStore))
	registerVersioned(mux, "GET", "/api/blogs/{slug}", publicDetailHandler(blogStore))
	registerVersioned(mux, "GET", "/api/admin/blogs", adminListHandler(blogStore, authenticator))
	registerVersioned(mux, "POST", "/api/admin/blogs", adminCreateHandler(blogStore, authenticator))
	registerVersioned(mux, "PUT", "/api/admin/blogs/{id}", adminUpdateHandler(blogStore, authenticator))
	registerVersioned(mux, "DELETE", "/api/admin/blogs/{id}", adminDeleteHandler(blogStore, authenticator))
	registerVersioned(mux, "POST", "/api/admin/publish", adminExternalPublishHandler(publishClient, authenticator))
	registerVersioned(mux, "POST", "/api/internal/content/import", contentImportHandler(blogStore, os.Getenv("CONTENT_IMPORT_TOKEN")))
	registerVersioned(mux, "POST", "/api/internal/scheduled-publish", scheduledPublishHandler(blogStore, publishClient, os.Getenv("SCHEDULED_PUBLISH_TOKEN")))
	return withCORS(withAPIVersionHeader(mux))
}

// registerVersioned mounts handler at the canonical /api/{currentAPIVersion}/...
// path and, for backward compatibility, at the equivalent legacy /api/... path.
// method may be empty to register the pattern for all HTTP methods.
func registerVersioned(mux *http.ServeMux, method, path string, handler http.HandlerFunc) {
	versionedPath := "/api/" + currentAPIVersion + strings.TrimPrefix(path, "/api")
	if method == "" {
		mux.HandleFunc(path, handler)
		mux.HandleFunc(versionedPath, handler)
		return
	}
	mux.HandleFunc(method+" "+path, handler)
	mux.HandleFunc(method+" "+versionedPath, handler)
}

// withAPIVersionHeader advertises the canonical API version on every response.
func withAPIVersionHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-API-Version", currentAPIVersion)
		next.ServeHTTP(w, r)
	})
}

// withCORS allows configured comma-separated frontend origins.
func withCORS(next http.Handler) http.Handler {
	configured := strings.Split(os.Getenv("ALLOWED_ORIGIN"), ",")
	allowed := make(map[string]bool, len(configured))
	for _, origin := range configured {
		if origin = strings.TrimSpace(origin); origin != "" {
			allowed[origin] = true
		}
	}
	if len(allowed) == 0 {
		allowed["*"] = true
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowed["*"] {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else if allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Max-Age", "600")
		if r.Method == http.MethodOptions {
			if !allowed["*"] && !allowed[origin] {
				w.WriteHeader(http.StatusForbidden)
				return
			}
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
