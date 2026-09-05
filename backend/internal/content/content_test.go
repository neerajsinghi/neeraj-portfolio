package content

import (
	"strings"
	"testing"
)

func TestParseBundleAcceptsFencedJSON(t *testing.T) {
	raw := "```json\n{\"title\":\"Reliable Go Services\",\"slug\":\"reliable-go-services\",\"description\":\"A practical guide to reliability patterns for production Go services and the tradeoffs behind them.\",\"tags\":[\"go\",\"backend\"],\"article_markdown\":\"# Reliable Go Services\\n\\nThis article body is intentionally longer than one hundred characters so ParseBundle accepts it during validation tests.\",\"linkedin_post\":\"A practical reliability guide {{CANONICAL_URL}}\",\"social_post\":\"Reliable Go services: {{CANONICAL_URL}}\"}\n```"

	bundle, err := ParseBundle(raw)
	if err != nil {
		t.Fatalf("ParseBundle() error = %v", err)
	}
	if bundle.Slug != "reliable-go-services" {
		t.Fatalf("ParseBundle() slug = %q", bundle.Slug)
	}
}

func TestParseBundleRejectsIncompleteBundle(t *testing.T) {
	if _, err := ParseBundle(`{"title":"Incomplete"}`); err == nil {
		t.Fatal("ParseBundle() expected an error")
	}
}

func TestParseBundleRejectsTooLongTitle(t *testing.T) {
	raw := `{"title":"Kubernetes Readiness Gates as Release Safety Primitives: Traffic Admission, Rollout Sequencing, and the Premature-Ready Failure Mode","slug":"readiness-gates-release-safety","description":"A practical deep dive into rollout sequencing and readiness gates for production Kubernetes workloads.","tags":["kubernetes"],"article_markdown":"This article body is intentionally longer than one hundred characters so ParseBundle catches title-length violations before import jobs fail.","linkedin_post":"LinkedIn copy {{CANONICAL_URL}}","social_post":"Social copy {{CANONICAL_URL}}"}`
	if _, err := ParseBundle(raw); err == nil {
		t.Fatal("ParseBundle() expected an error for title longer than 120 characters")
	}
}

func TestParseBundleAcceptsUnicodeTitleAtCharacterLimit(t *testing.T) {
	title := strings.Repeat("界", 120)
	raw := `{"title":"` + title + `","slug":"unicode-title-limit","description":"A practical deep dive into rollout sequencing and readiness gates for production Kubernetes workloads.","tags":["kubernetes"],"article_markdown":"This article body is intentionally longer than one hundred characters so ParseBundle catches title-length violations before import jobs fail.","linkedin_post":"LinkedIn copy {{CANONICAL_URL}}","social_post":"Social copy {{CANONICAL_URL}}"}`
	if _, err := ParseBundle(raw); err != nil {
		t.Fatalf("ParseBundle() error = %v", err)
	}
}

func TestParseBundleRejectsUnicodeTitlePastCharacterLimit(t *testing.T) {
	title := strings.Repeat("界", 121)
	raw := `{"title":"` + title + `","slug":"unicode-title-over-limit","description":"A practical deep dive into rollout sequencing and readiness gates for production Kubernetes workloads.","tags":["kubernetes"],"article_markdown":"This article body is intentionally longer than one hundred characters so ParseBundle catches title-length violations before import jobs fail.","linkedin_post":"LinkedIn copy {{CANONICAL_URL}}","social_post":"Social copy {{CANONICAL_URL}}"}`
	if _, err := ParseBundle(raw); err == nil {
		t.Fatal("ParseBundle() expected an error for Unicode title longer than 120 characters")
	}
}