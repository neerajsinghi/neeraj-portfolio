package blog

import (
	"errors"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestEditorCannotPublishOrDelete(t *testing.T) {
	principal := Principal{Roles: map[Role]bool{RoleEditor: true}}
	if !principal.CanEdit() {
		t.Fatal("editor should be able to edit")
	}
	if principal.CanPublish() || principal.CanDelete() {
		t.Fatal("editor should not be able to publish or delete")
	}
}

func TestScheduledPostRequiresFutureTime(t *testing.T) {
	valid := WriteInput{
		Slug: "scheduled-post", Title: "Scheduled post",
		Description:     "A sufficiently detailed description for a scheduled technical article.",
		ContentMarkdown: strings.Repeat("scheduled content ", 10), Status: StatusScheduled,
	}
	if err := valid.NormalizeAndValidate(); !errors.Is(err, ErrInvalid) {
		t.Fatalf("missing scheduled_at error = %v, want ErrInvalid", err)
	}

	past := time.Now().Add(-time.Minute)
	valid.ScheduledAt = &past
	if err := valid.NormalizeAndValidate(); !errors.Is(err, ErrInvalid) {
		t.Fatalf("past scheduled_at error = %v, want ErrInvalid", err)
	}

	future := time.Now().Add(time.Hour)
	valid.ScheduledAt = &future
	if err := valid.NormalizeAndValidate(); err != nil {
		t.Fatalf("future scheduled_at error = %v", err)
	}
}

func TestPublicFilterIncludesOnlyDueScheduledPosts(t *testing.T) {
	now := time.Now().UTC()
	filter := publicFilter(now)
	conditions := filter["$or"].(bson.A)
	scheduled := conditions[1].(bson.M)
	cutoff := scheduled["scheduled_at"].(bson.M)["$lte"].(time.Time)
	if scheduled["status"] != StatusScheduled || !cutoff.Equal(now) {
		t.Fatalf("scheduled public filter = %#v", scheduled)
	}
}

func TestWriteInputNormalizesTags(t *testing.T) {
	input := WriteInput{
		Slug:            "reliable-go-services",
		Title:           "Reliable Go Services",
		Description:     "A practical guide to building reliable Go services in production.",
		ContentMarkdown: strings.Repeat("content ", 20),
		Tags:            []string{"Go", " go ", "Reliability"},
	}
	if err := input.NormalizeAndValidate(); err != nil {
		t.Fatalf("NormalizeAndValidate() error = %v", err)
	}
	if len(input.Tags) != 2 || input.Tags[0] != "go" {
		t.Fatalf("normalized tags = %#v", input.Tags)
	}
	if input.Status != StatusDraft {
		t.Fatalf("default status = %q", input.Status)
	}
}
