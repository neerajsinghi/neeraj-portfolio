package publisher

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"neeraj-portfolio/backend/internal/content"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (function roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return function(request)
}

func TestPublishDEVDefaultsToDraftPayload(t *testing.T) {
	client := NewClient(&http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		body, _ := io.ReadAll(request.Body)
		if !strings.Contains(string(body), `"published":false`) {
			t.Fatalf("request body did not create a draft: %s", body)
		}
		return &http.Response{StatusCode: http.StatusCreated, Body: io.NopCloser(strings.NewReader(`{"id":42,"url":"https://dev.to/example"}`)), Header: make(http.Header)}, nil
	})})

	result, err := client.PublishDEV(context.Background(), content.Bundle{Title: "Title", Article: "Body", Tags: []string{"go"}}, "key", false)
	if err != nil {
		t.Fatalf("PublishDEV() error = %v", err)
	}
	if !result.Draft {
		t.Fatal("PublishDEV() expected a draft result")
	}
}