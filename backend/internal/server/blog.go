package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"neeraj-portfolio/backend/internal/auth"
	"neeraj-portfolio/backend/internal/blog"
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
