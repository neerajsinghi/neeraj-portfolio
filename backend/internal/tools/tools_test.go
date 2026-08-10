package tools

import (
	"context"
	"strings"
	"testing"
	"time"

	"neeraj-portfolio/backend/internal/blog"
)

type fakeBlogStore struct {
	posts []blog.Post
}

func (fake fakeBlogStore) ListPublished(context.Context, int) ([]blog.Post, error) {
	return fake.posts, nil
}
func (fake fakeBlogStore) GetPublishedBySlug(context.Context, string) (blog.Post, error) {
	return blog.Post{}, blog.ErrNotFound
}
func (fake fakeBlogStore) ListAll(context.Context, int) ([]blog.Post, error) {
	return fake.posts, nil
}
func (fake fakeBlogStore) Create(context.Context, blog.WriteInput, blog.Principal) (blog.Post, error) {
	return blog.Post{}, nil
}
func (fake fakeBlogStore) Update(context.Context, string, blog.WriteInput, blog.Principal) (blog.Post, error) {
	return blog.Post{}, nil
}
func (fake fakeBlogStore) Delete(context.Context, string, blog.Principal) error { return nil }
func (fake fakeBlogStore) PublishDue(context.Context) ([]blog.Post, error)      { return nil, nil }
func (fake fakeBlogStore) PendingExternalPublish(context.Context) ([]blog.Post, error) {
	return nil, nil
}
func (fake fakeBlogStore) RecordExternalPublish(context.Context, string, string, string) error {
	return nil
}

func TestBlogToolsReturnGroundedPostsAndSources(t *testing.T) {
	now := time.Now().UTC()
	executor := NewExecutor(fakeBlogStore{posts: []blog.Post{
		{Slug: "go-reliability", Title: "Reliable Go services", Description: "Practical reliability patterns.", ContentMarkdown: "Retries, idempotency, and observability in Go.", Tags: []string{"go", "reliability"}, PublishedAt: &now},
		{Slug: "rag-agents", Title: "Grounded RAG agents", Description: "Tool-use agent design.", ContentMarkdown: "Retrieval and model context for production agents.", Tags: []string{"ai", "rag"}, PublishedAt: &now},
	}})

	listed, sources := executor.Execute("list_blogs", map[string]any{"limit": float64(2)})
	if !strings.Contains(listed, "https://neerajsinghi.com/blogs/go-reliability") || len(sources) != 2 {
		t.Fatalf("list_blogs result = %q, sources = %#v", listed, sources)
	}

	found, sources := executor.Execute("search_blogs", map[string]any{"query": "RAG", "top_n": float64(1)})
	if !strings.Contains(found, "Grounded RAG agents") || len(sources) != 1 || sources[0] != "blog:rag-agents" {
		t.Fatalf("search_blogs result = %q, sources = %#v", found, sources)
	}
}

func TestBlogToolsHandleUnavailableStorage(t *testing.T) {
	result, sources := NewExecutor(nil).Execute("list_blogs", map[string]any{})
	if !strings.Contains(result, "temporarily unavailable") || len(sources) != 0 {
		t.Fatalf("result = %q, sources = %#v", result, sources)
	}
}
