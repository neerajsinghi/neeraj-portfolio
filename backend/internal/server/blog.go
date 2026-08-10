package server

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"neeraj-portfolio/backend/internal/auth"
	"neeraj-portfolio/backend/internal/blog"
	"neeraj-portfolio/backend/internal/content"
	"neeraj-portfolio/backend/internal/publisher"
)

type publicPost struct {
	Slug            string     `json:"slug"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	ContentMarkdown string     `json:"content_markdown,omitempty"`
	Tags            []string   `json:"tags"`
	PublishedAt     *time.Time `json:"published_at,omitempty"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func publicListHandler(store blog.Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if store == nil {
			writeAPIError(writer, http.StatusServiceUnavailable, "blog storage is not configured")
			return
		}
		limit, _ := strconv.Atoi(request.URL.Query().Get("limit"))
		ctx, cancel := context.WithTimeout(request.Context(), 5*time.Second)
		defer cancel()
		posts, err := store.ListPublished(ctx, limit)
		if err != nil {
			writeAPIError(writer, http.StatusInternalServerError, "could not load blog posts")
			return
		}
		summaries := make([]publicPost, 0, len(posts))
		for _, post := range posts {
			summaries = append(summaries, toPublicPost(post, false))
		}
		writeJSON(writer, http.StatusOK, map[string]any{"posts": summaries})
	}
}

func publicDetailHandler(store blog.Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if store == nil {
			writeAPIError(writer, http.StatusServiceUnavailable, "blog storage is not configured")
			return
		}
		ctx, cancel := context.WithTimeout(request.Context(), 5*time.Second)
		defer cancel()
		post, err := store.GetPublishedBySlug(ctx, request.PathValue("slug"))
		if errors.Is(err, blog.ErrNotFound) {
			writeAPIError(writer, http.StatusNotFound, "blog post not found")
			return
		}
		if err != nil {
			writeAPIError(writer, http.StatusInternalServerError, "could not load blog post")
			return
		}
		writeJSON(writer, http.StatusOK, toPublicPost(post, true))
	}
}

func adminListHandler(store blog.Store, authenticator auth.Authenticator) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		principal, ok := requireEditor(writer, request, authenticator)
		if !ok {
			return
		}
		_ = principal
		if store == nil {
			writeAPIError(writer, http.StatusServiceUnavailable, "blog storage is not configured")
			return
		}
		limit, _ := strconv.Atoi(request.URL.Query().Get("limit"))
		ctx, cancel := context.WithTimeout(request.Context(), 5*time.Second)
		defer cancel()
		posts, err := store.ListAll(ctx, limit)
		if err != nil {
			writeAPIError(writer, http.StatusInternalServerError, "could not load blog posts")
			return
		}
		writeJSON(writer, http.StatusOK, map[string]any{"posts": posts})
	}
}

func adminCreateHandler(store blog.Store, authenticator auth.Authenticator) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		principal, ok := requireEditor(writer, request, authenticator)
		if !ok {
			return
		}
		if store == nil {
			writeAPIError(writer, http.StatusServiceUnavailable, "blog storage is not configured")
			return
		}
		var input blog.WriteInput
		if !decodeJSON(writer, request, &input) {
			return
		}
		ctx, cancel := context.WithTimeout(request.Context(), 5*time.Second)
		defer cancel()
		post, err := store.Create(ctx, input, principal)
		if writeBlogError(writer, err) {
			return
		}
		writeJSON(writer, http.StatusCreated, post)
	}
}

func contentImportHandler(store blog.Store, importToken string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		providedToken := strings.TrimSpace(strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer "))
		if importToken == "" || subtle.ConstantTimeCompare([]byte(providedToken), []byte(importToken)) != 1 {
			writeAPIError(writer, http.StatusUnauthorized, "valid import token required")
			return
		}
		if store == nil {
			writeAPIError(writer, http.StatusServiceUnavailable, "blog storage is not configured")
			return
		}

		var bundle content.Bundle
		if !decodeJSON(writer, request, &bundle) {
			return
		}
		input := blog.WriteInput{
			Slug: bundle.Slug, Title: bundle.Title, Description: bundle.Description,
			ContentMarkdown: bundle.Article, LinkedInPost: bundle.LinkedIn, SocialPost: bundle.Social,
			Tags: bundle.Tags, Status: blog.StatusDraft,
		}
		principal := blog.Principal{
			Subject: "github-content-import", Email: "github-actions@neerajsinghi.com",
			Roles: map[blog.Role]bool{blog.RoleEditor: true},
		}
		ctx, cancel := context.WithTimeout(request.Context(), 5*time.Second)
		defer cancel()
		post, err := store.Create(ctx, input, principal)
		if errors.Is(err, blog.ErrConflict) {
			writeJSON(writer, http.StatusOK, map[string]any{"imported": false, "slug": bundle.Slug})
			return
		}
		if writeBlogError(writer, err) {
			return
		}
		writeJSON(writer, http.StatusCreated, map[string]any{"imported": true, "post": post})
	}
}

type externalPublisher interface {
	PublishDEV(context.Context, content.Bundle, string, bool) (publisher.Result, error)
	PublishLinkedIn(context.Context, content.Bundle, string, string, string) (publisher.Result, error)
}

// scheduledPublishHandler is called by the Atlas Scheduled Trigger instead of
// letting it mutate MongoDB directly, so due posts go through the same
// version/revision/audit path as a human-triggered publish.
func scheduledPublishHandler(store blog.Store, token string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		providedToken := strings.TrimSpace(strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer "))
		if token == "" || subtle.ConstantTimeCompare([]byte(providedToken), []byte(token)) != 1 {
			writeAPIError(writer, http.StatusUnauthorized, "valid scheduled-publish token required")
			return
		}
		if store == nil {
			writeAPIError(writer, http.StatusServiceUnavailable, "blog storage is not configured")
			return
		}
		ctx, cancel := context.WithTimeout(request.Context(), 15*time.Second)
		defer cancel()
		posts, err := store.PublishDue(ctx)
		if err != nil {
			writeAPIError(writer, http.StatusInternalServerError, "could not publish due posts")
			return
		}
		slugs := make([]string, 0, len(posts))
		for _, post := range posts {
			slugs = append(slugs, post.Slug)
		}
		writeJSON(writer, http.StatusOK, map[string]any{"published": len(posts), "slugs": slugs})
	}
}

type externalPublishRequest struct {
	content.Bundle
	Platforms []string `json:"platforms"`
}

func adminExternalPublishHandler(client externalPublisher, authenticator auth.Authenticator) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		principal, ok := requireEditor(writer, request, authenticator)
		if !ok {
			return
		}
		if !principal.CanPublish() {
			writeAPIError(writer, http.StatusForbidden, "admin role required for publishing")
			return
		}
		var input externalPublishRequest
		if !decodeJSON(writer, request, &input) {
			return
		}
		if len(input.Platforms) == 0 {
			writeAPIError(writer, http.StatusBadRequest, "select at least one external platform")
			return
		}
		if input.CanonicalURL == "" {
			input.CanonicalURL = "https://neerajsinghi.com/blogs/" + input.Slug
		}
		input.LinkedIn = strings.ReplaceAll(input.LinkedIn, "{{CANONICAL_URL}}", input.CanonicalURL)
		input.Social = strings.ReplaceAll(input.Social, "{{CANONICAL_URL}}", input.CanonicalURL)

		ctx, cancel := context.WithTimeout(request.Context(), 45*time.Second)
		defer cancel()
		results := make([]publisher.Result, 0, len(input.Platforms))
		seen := make(map[string]bool, len(input.Platforms))
		for _, platform := range input.Platforms {
			platform = strings.TrimSpace(strings.ToLower(platform))
			if seen[platform] {
				continue
			}
			seen[platform] = true
			var result publisher.Result
			var err error
			switch platform {
			case "devto":
				result, err = client.PublishDEV(ctx, input.Bundle, os.Getenv("DEVTO_API_KEY"), true)
			case "linkedin":
				result, err = client.PublishLinkedIn(ctx, input.Bundle, os.Getenv("LINKEDIN_ACCESS_TOKEN"), os.Getenv("LINKEDIN_AUTHOR_URN"), os.Getenv("LINKEDIN_VERSION"))
			default:
				writeAPIError(writer, http.StatusBadRequest, "unsupported publishing platform: "+platform)
				return
			}
			if err != nil {
				writeAPIError(writer, http.StatusBadGateway, "could not publish to "+platform+": "+err.Error())
				return
			}
			results = append(results, result)
		}
		writeJSON(writer, http.StatusOK, map[string]any{"results": results})
	}
}

func adminUpdateHandler(store blog.Store, authenticator auth.Authenticator) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		principal, ok := requireEditor(writer, request, authenticator)
		if !ok {
			return
		}
		if store == nil {
			writeAPIError(writer, http.StatusServiceUnavailable, "blog storage is not configured")
			return
		}
		var input blog.WriteInput
		if !decodeJSON(writer, request, &input) {
			return
		}
		ctx, cancel := context.WithTimeout(request.Context(), 5*time.Second)
		defer cancel()
		post, err := store.Update(ctx, request.PathValue("id"), input, principal)
		if writeBlogError(writer, err) {
			return
		}
		writeJSON(writer, http.StatusOK, post)
	}
}

func adminDeleteHandler(store blog.Store, authenticator auth.Authenticator) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		principal, ok := requireEditor(writer, request, authenticator)
		if !ok {
			return
		}
		if store == nil {
			writeAPIError(writer, http.StatusServiceUnavailable, "blog storage is not configured")
			return
		}
		ctx, cancel := context.WithTimeout(request.Context(), 5*time.Second)
		defer cancel()
		if writeBlogError(writer, store.Delete(ctx, request.PathValue("id"), principal)) {
			return
		}
		writer.WriteHeader(http.StatusNoContent)
	}
}

func requireEditor(writer http.ResponseWriter, request *http.Request, authenticator auth.Authenticator) (blog.Principal, bool) {
	if authenticator == nil {
		writeAPIError(writer, http.StatusServiceUnavailable, "admin authentication is not configured")
		return blog.Principal{}, false
	}
	principal, err := authenticator.Authenticate(request)
	if err != nil {
		writeAPIError(writer, http.StatusUnauthorized, "valid bearer token required")
		return blog.Principal{}, false
	}
	if !principal.CanEdit() {
		writeAPIError(writer, http.StatusForbidden, "editor or admin role required")
		return blog.Principal{}, false
	}
	return principal, true
}

func decodeJSON(writer http.ResponseWriter, request *http.Request, destination any) bool {
	request.Body = http.MaxBytesReader(writer, request.Body, 1<<20)
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		writeAPIError(writer, http.StatusBadRequest, "invalid request body")
		return false
	}
	return true
}

func writeBlogError(writer http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, blog.ErrNotFound):
		writeAPIError(writer, http.StatusNotFound, err.Error())
	case errors.Is(err, blog.ErrConflict):
		writeAPIError(writer, http.StatusConflict, err.Error())
	case errors.Is(err, blog.ErrForbidden):
		writeAPIError(writer, http.StatusForbidden, err.Error())
	case errors.Is(err, blog.ErrInvalid):
		writeAPIError(writer, http.StatusBadRequest, err.Error())
	default:
		writeAPIError(writer, http.StatusInternalServerError, "blog storage operation failed")
	}
	return true
}

func writeAPIError(writer http.ResponseWriter, status int, message string) {
	writeJSON(writer, status, map[string]string{"error": message})
}

func toPublicPost(post blog.Post, includeContent bool) publicPost {
	result := publicPost{
		Slug: post.Slug, Title: post.Title, Description: post.Description,
		Tags: post.Tags, PublishedAt: post.PublishedAt, UpdatedAt: post.UpdatedAt,
	}
	if includeContent {
		result.ContentMarkdown = post.ContentMarkdown
	}
	return result
}
