package blog

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var ErrInvalid = errors.New("invalid blog input")

type Status string

const (
	StatusDraft     Status = "draft"
	StatusScheduled Status = "scheduled"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleEditor Role = "editor"
)

type Author struct {
	Subject string `json:"subject" bson:"subject"`
	Email   string `json:"email" bson:"email"`
}

type Post struct {
	ID              string     `json:"id" bson:"-"`
	Slug            string     `json:"slug" bson:"slug"`
	Title           string     `json:"title" bson:"title"`
	Description     string     `json:"description" bson:"description"`
	ContentMarkdown string     `json:"content_markdown" bson:"content_markdown"`
	Tags            []string   `json:"tags" bson:"tags"`
	Status          Status     `json:"status" bson:"status"`
	Author          Author     `json:"author" bson:"author"`
	CreatedAt       time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" bson:"updated_at"`
	ScheduledAt     *time.Time `json:"scheduled_at,omitempty" bson:"scheduled_at,omitempty"`
	PublishedAt     *time.Time `json:"published_at,omitempty" bson:"published_at,omitempty"`
	SchemaVersion   int        `json:"schema_version" bson:"schema_version"`
	Version         int        `json:"version" bson:"version"`
}

type WriteInput struct {
	Slug            string     `json:"slug"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	ContentMarkdown string     `json:"content_markdown"`
	Tags            []string   `json:"tags"`
	Status          Status     `json:"status"`
	Version         int        `json:"version,omitempty"`
	ScheduledAt     *time.Time `json:"scheduled_at,omitempty"`
}

type Principal struct {
	Subject string
	Email   string
	Roles   map[Role]bool
}

func (principal Principal) CanEdit() bool {
	return principal.Roles[RoleAdmin] || principal.Roles[RoleEditor]
}

func (principal Principal) CanPublish() bool {
	return principal.Roles[RoleAdmin]
}

func (principal Principal) CanDelete() bool {
	return principal.Roles[RoleAdmin]
}

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func (input *WriteInput) NormalizeAndValidate() error {
	input.Slug = strings.TrimSpace(strings.ToLower(input.Slug))
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)
	input.ContentMarkdown = strings.TrimSpace(input.ContentMarkdown)
	if input.Status == "" {
		input.Status = StatusDraft
	}
	if !slugPattern.MatchString(input.Slug) {
		return errors.Join(ErrInvalid, errors.New("slug must contain lowercase letters, numbers, and single hyphens"))
	}
	if len(input.Title) < 5 || len(input.Title) > 120 {
		return errors.Join(ErrInvalid, errors.New("title must contain 5 to 120 characters"))
	}
	if len(input.Description) < 40 || len(input.Description) > 180 {
		return errors.Join(ErrInvalid, errors.New("description must contain 40 to 180 characters"))
	}
	if len(input.ContentMarkdown) < 100 {
		return errors.Join(ErrInvalid, errors.New("content_markdown must contain at least 100 characters"))
	}
	if len(input.Tags) > 6 {
		return errors.Join(ErrInvalid, errors.New("tags must contain no more than 6 values"))
	}
	if input.Status != StatusDraft && input.Status != StatusScheduled && input.Status != StatusPublished && input.Status != StatusArchived {
		return errors.Join(ErrInvalid, errors.New("status must be draft, scheduled, published, or archived"))
	}
	if input.Status == StatusScheduled {
		if input.ScheduledAt == nil {
			return errors.Join(ErrInvalid, errors.New("scheduled_at is required for scheduled posts"))
		}
		if !input.ScheduledAt.After(time.Now().UTC()) {
			return errors.Join(ErrInvalid, errors.New("scheduled_at must be in the future"))
		}
		value := input.ScheduledAt.UTC()
		input.ScheduledAt = &value
	} else {
		input.ScheduledAt = nil
	}
	seen := make(map[string]bool, len(input.Tags))
	normalizedTags := make([]string, 0, len(input.Tags))
	for _, tag := range input.Tags {
		tag = strings.TrimSpace(strings.ToLower(tag))
		if tag != "" && !seen[tag] {
			seen[tag] = true
			normalizedTags = append(normalizedTags, tag)
		}
	}
	input.Tags = normalizedTags
	return nil
}
