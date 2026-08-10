package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"neeraj-portfolio/backend/internal/auth"
	"neeraj-portfolio/backend/internal/blog"
	"neeraj-portfolio/backend/internal/content"
	"neeraj-portfolio/backend/internal/publisher"
)

type fakeAuthenticator struct {
	principal blog.Principal
	err       error
}

func (fake fakeAuthenticator) Authenticate(*http.Request) (blog.Principal, error) {
	return fake.principal, fake.err
}

type fakeBlogStore struct {
	posts           []blog.Post
	created         blog.WriteInput
	createErr       error
	updateErr       error
	publishDuePosts []blog.Post
	publishDueErr   error
}

type fakeExternalPublisher struct {
	platforms []string
}

func (fake *fakeExternalPublisher) PublishDEV(context.Context, content.Bundle, string, bool) (publisher.Result, error) {
	fake.platforms = append(fake.platforms, "devto")
	return publisher.Result{Platform: "devto", URL: "https://dev.to/example"}, nil
}

func (fake *fakeExternalPublisher) PublishLinkedIn(context.Context, content.Bundle, string, string, string) (publisher.Result, error) {
	fake.platforms = append(fake.platforms, "linkedin")
	return publisher.Result{Platform: "linkedin", URL: "https://linkedin.com/example"}, nil
}

func (fake *fakeBlogStore) ListPublished(context.Context, int) ([]blog.Post, error) {
	return fake.posts, nil
}
func (fake *fakeBlogStore) GetPublishedBySlug(context.Context, string) (blog.Post, error) {
	return fake.posts[0], nil
}
func (fake *fakeBlogStore) ListAll(context.Context, int) ([]blog.Post, error) {
	return fake.posts, nil
}
func (fake *fakeBlogStore) Create(_ context.Context, input blog.WriteInput, _ blog.Principal) (blog.Post, error) {
	fake.created = input
	return blog.Post{Slug: input.Slug, Status: input.Status}, fake.createErr
}
func (fake *fakeBlogStore) Update(context.Context, string, blog.WriteInput, blog.Principal) (blog.Post, error) {
	return blog.Post{}, fake.updateErr
}
func (fake *fakeBlogStore) Delete(context.Context, string, blog.Principal) error { return nil }
func (fake *fakeBlogStore) PublishDue(context.Context) ([]blog.Post, error) {
	return fake.publishDuePosts, fake.publishDueErr
}

func TestAdminListRequiresAuthenticationAndRole(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/admin/blogs", nil)

	unauthorized := httptest.NewRecorder()
	adminListHandler(&fakeBlogStore{}, fakeAuthenticator{err: auth.ErrUnauthorized})(unauthorized, request)
	if unauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated status = %d, want 401", unauthorized.Code)
	}

	forbidden := httptest.NewRecorder()
	adminListHandler(&fakeBlogStore{}, fakeAuthenticator{principal: blog.Principal{Roles: map[blog.Role]bool{}}})(forbidden, request)
	if forbidden.Code != http.StatusForbidden {
		t.Fatalf("unprivileged status = %d, want 403", forbidden.Code)
	}
}

func TestAdminUpdateMapsConflict(t *testing.T) {
	store := &fakeBlogStore{updateErr: blog.ErrConflict}
	authenticator := fakeAuthenticator{principal: blog.Principal{Roles: map[blog.Role]bool{blog.RoleEditor: true}}}
	body := `{"slug":"valid-slug","title":"Valid title","description":"A description long enough to satisfy the blog validation rules.","content_markdown":"This content is intentionally longer than one hundred characters so it passes validation before the fake store returns conflict.","tags":[],"status":"draft","version":1}`
	request := httptest.NewRequest(http.MethodPut, "/api/admin/blogs/abc", strings.NewReader(body))
	response := httptest.NewRecorder()

	adminUpdateHandler(store, authenticator)(response, request)
	if response.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body = %s", response.Code, response.Body.String())
	}
}

func TestAdminCreateHandlesUnavailableStore(t *testing.T) {
	authenticator := fakeAuthenticator{principal: blog.Principal{Roles: map[blog.Role]bool{blog.RoleEditor: true}}}
	request := httptest.NewRequest(http.MethodPost, "/api/admin/blogs", strings.NewReader(`{}`))
	response := httptest.NewRecorder()

	adminCreateHandler(nil, authenticator)(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", response.Code)
	}
}

func TestContentImportCreatesEditableDraft(t *testing.T) {
	store := &fakeBlogStore{}
	body := `{"title":"Production Go Timeouts","slug":"production-go-timeouts","description":"A detailed guide to timeout budgets in production Go services.","tags":["Go","Reliability"],"article_markdown":"This article body is intentionally longer than one hundred characters so the imported blog input satisfies validation and remains editable in the admin console.","linkedin_post":"LinkedIn copy","social_post":"Social copy","created_at":"2026-08-09T00:00:00Z"}`
	request := httptest.NewRequest(http.MethodPost, "/api/internal/content/import", strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer import-secret")
	response := httptest.NewRecorder()

	contentImportHandler(store, "import-secret")(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", response.Code, response.Body.String())
	}
	if store.created.Status != blog.StatusDraft || store.created.Slug != "production-go-timeouts" || store.created.ContentMarkdown == "" {
		t.Fatalf("unexpected imported input: %#v", store.created)
	}
}

func TestContentImportIsAuthenticatedAndIdempotent(t *testing.T) {
	body := `{"title":"Production Go Timeouts","slug":"production-go-timeouts","description":"A detailed guide to timeout budgets in production Go services.","tags":["go"],"article_markdown":"This article body is intentionally longer than one hundred characters so the imported blog input satisfies validation and remains editable in the admin console.","linkedin_post":"LinkedIn copy","social_post":"Social copy"}`

	unauthorized := httptest.NewRecorder()
	contentImportHandler(&fakeBlogStore{}, "import-secret")(unauthorized, httptest.NewRequest(http.MethodPost, "/api/internal/content/import", strings.NewReader(body)))
	if unauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("unauthorized status = %d, want 401", unauthorized.Code)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/internal/content/import", strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer import-secret")
	response := httptest.NewRecorder()
	contentImportHandler(&fakeBlogStore{createErr: blog.ErrConflict}, "import-secret")(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"imported":false`) {
		t.Fatalf("duplicate status = %d, want 200; body = %s", response.Code, response.Body.String())
	}
}

func TestExternalPublishRequiresAdminAndUsesSelectedPlatforms(t *testing.T) {
	body := `{"title":"Production Go Timeouts","slug":"production-go-timeouts","description":"A detailed guide to timeout budgets in production Go services.","tags":["go"],"article_markdown":"Article body","linkedin_post":"Read this: {{CANONICAL_URL}}","social_post":"Read this: {{CANONICAL_URL}}","platforms":["devto","linkedin"]}`

	editorRequest := httptest.NewRequest(http.MethodPost, "/api/admin/publish", strings.NewReader(body))
	editorResponse := httptest.NewRecorder()
	adminExternalPublishHandler(&fakeExternalPublisher{}, fakeAuthenticator{principal: blog.Principal{Roles: map[blog.Role]bool{blog.RoleEditor: true}}})(editorResponse, editorRequest)
	if editorResponse.Code != http.StatusForbidden {
		t.Fatalf("editor status = %d, want 403", editorResponse.Code)
	}

	client := &fakeExternalPublisher{}
	adminRequest := httptest.NewRequest(http.MethodPost, "/api/admin/publish", strings.NewReader(body))
	adminResponse := httptest.NewRecorder()
	adminExternalPublishHandler(client, fakeAuthenticator{principal: blog.Principal{Roles: map[blog.Role]bool{blog.RoleAdmin: true}}})(adminResponse, adminRequest)
	if adminResponse.Code != http.StatusOK || len(client.platforms) != 2 {
		t.Fatalf("admin status = %d, platforms = %v; body = %s", adminResponse.Code, client.platforms, adminResponse.Body.String())
	}
}

func TestPublicBlogResponseRedactsAuthorIdentity(t *testing.T) {
	now := time.Now().UTC()
	store := &fakeBlogStore{posts: []blog.Post{{
		Slug: "public-post", Title: "Public post", Description: "Public description",
		ContentMarkdown: "Public body", Tags: []string{"go"}, Status: blog.StatusPublished,
		Author: blog.Author{Subject: "cognito-subject", Email: "private@example.com"}, UpdatedAt: now,
	}}}
	request := httptest.NewRequest(http.MethodGet, "/api/blogs/public-post", nil)
	response := httptest.NewRecorder()

	publicDetailHandler(store)(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	if strings.Contains(response.Body.String(), "private@example.com") || strings.Contains(response.Body.String(), "cognito-subject") {
		t.Fatalf("public response exposed private author identity: %s", response.Body.String())
	}
}

func TestInternalStoreErrorIsNotExposed(t *testing.T) {
	response := httptest.NewRecorder()
	writeBlogError(response, errors.New("mongodb credential detail"))
	if response.Code != http.StatusInternalServerError || strings.Contains(response.Body.String(), "credential") {
		t.Fatalf("unexpected internal error response: status=%d body=%s", response.Code, response.Body.String())
	}
}

func TestScheduledPublishRequiresToken(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/internal/scheduled-publish", nil)
	response := httptest.NewRecorder()

	scheduledPublishHandler(&fakeBlogStore{}, "publish-secret")(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", response.Code)
	}
}

func TestScheduledPublishReturnsPublishedSlugs(t *testing.T) {
	store := &fakeBlogStore{publishDuePosts: []blog.Post{{Slug: "due-post", Status: blog.StatusPublished}}}
	request := httptest.NewRequest(http.MethodPost, "/api/internal/scheduled-publish", nil)
	request.Header.Set("Authorization", "Bearer publish-secret")
	response := httptest.NewRecorder()

	scheduledPublishHandler(store, "publish-secret")(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"due-post"`) {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}
