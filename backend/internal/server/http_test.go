package server

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"neeraj-portfolio/backend/internal/blog"
	"neeraj-portfolio/backend/internal/llm"
)

type fakeProvider struct{}

func (fakeProvider) Complete(context.Context, llm.Request) (llm.Response, error) {
	return llm.Response{StopReason: llm.StopReasonEndTurn}, nil
}
func (fakeProvider) ModelName() string { return "fake-model" }

func TestVersionedRoutesAliasLegacyRoutes(t *testing.T) {
	handler := NewHandlerWithDependencies(fakeProvider{}, &fakeBlogStore{}, nil)

	for _, path := range []string{"/api/health", "/api/v1/health"} {
		request := httptest.NewRequest("GET", path, nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != 200 {
			t.Fatalf("%s status = %d, want 200", path, response.Code)
		}
		if version := response.Header().Get("X-API-Version"); version != "v1" {
			t.Fatalf("%s X-API-Version = %q, want v1", path, version)
		}
	}
}

func TestVersionedPublicBlogDetailAliasesLegacyPath(t *testing.T) {
	store := &fakeBlogStore{posts: []blog.Post{{
		Slug: "hello", Title: "Hello world", Description: "A published post",
		ContentMarkdown: "Body", Tags: []string{"go"}, Status: blog.StatusPublished, UpdatedAt: time.Now().UTC(),
	}}}
	handler := NewHandlerWithDependencies(fakeProvider{}, store, nil)

	for _, path := range []string{"/api/blogs/hello", "/api/v1/blogs/hello"} {
		request := httptest.NewRequest("GET", path, nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != 200 {
			t.Fatalf("%s status = %d, want 200; body = %s", path, response.Code, response.Body.String())
		}
	}
}
