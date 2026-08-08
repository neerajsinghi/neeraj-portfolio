package content

import "testing"

func TestParseBundleAcceptsFencedJSON(t *testing.T) {
	raw := "```json\n{\"title\":\"Reliable Go Services\",\"slug\":\"reliable-go-services\",\"description\":\"A practical guide to reliability patterns for production Go services and the tradeoffs behind them.\",\"tags\":[\"go\",\"backend\"],\"article_markdown\":\"# Reliable Go Services\\n\\nContent\",\"linkedin_post\":\"A practical reliability guide {{CANONICAL_URL}}\",\"social_post\":\"Reliable Go services: {{CANONICAL_URL}}\"}\n```"

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