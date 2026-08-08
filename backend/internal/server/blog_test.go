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
)

type fakeAuthenticator struct {
	principal blog.Principal
	err       error
}

func (fake fakeAuthenticator) Authenticate(*http.Request) (blog.Principal, error) {
	return fake.principal, fake.err
}

type fakeBlogStore struct {
	posts     []blog.Post
	updateErr error
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
func (fake *fakeBlogStore) Create(context.Context, blog.WriteInput, blog.Principal) (blog.Post, error) {
	return blog.Post{}, nil
}
func (fake *fakeBlogStore) Update(context.Context, string, blog.WriteInput, blog.Principal) (blog.Post, error) {
	return blog.Post{}, fake.updateErr
}
func (fake *fakeBlogStore) Delete(context.Context, string, blog.Principal) error { return nil }

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
