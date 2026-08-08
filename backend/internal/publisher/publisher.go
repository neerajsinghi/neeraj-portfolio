package publisher

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"neeraj-portfolio/backend/internal/content"
)

type Result struct {
	Platform string `json:"platform"`
	ID       string `json:"id,omitempty"`
	URL      string `json:"url,omitempty"`
	Draft    bool   `json:"draft"`
}

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{httpClient: httpClient}
}

func (c *Client) PublishDEV(ctx context.Context, bundle content.Bundle, apiKey string, publish bool) (Result, error) {
	if apiKey == "" {
		return Result{}, errors.New("DEVTO_API_KEY is required")
	}
	article := map[string]any{
		"title": bundle.Title, "body_markdown": bundle.Article, "published": publish,
		"description": bundle.Description, "tags": strings.Join(bundle.Tags, ","),
	}
	if bundle.CanonicalURL != "" {
		article["canonical_url"] = bundle.CanonicalURL
	}
	payload := map[string]any{"article": article}
	request, err := newJSONRequest(ctx, http.MethodPost, "https://dev.to/api/articles", payload)
	if err != nil {
		return Result{}, err
	}
	request.Header.Set("api-key", apiKey)
	request.Header.Set("accept", "application/vnd.forem.api-v1+json")

	var response struct {
		ID  int    `json:"id"`
		URL string `json:"url"`
	}
	if err := c.do(request, &response); err != nil {
		return Result{}, fmt.Errorf("publish DEV article: %w", err)
	}
	return Result{Platform: "devto", ID: fmt.Sprint(response.ID), URL: response.URL, Draft: !publish}, nil
}

func (c *Client) PublishLinkedIn(ctx context.Context, bundle content.Bundle, token, authorURN, version string) (Result, error) {
	if token == "" || authorURN == "" || version == "" {
		return Result{}, errors.New("LINKEDIN_ACCESS_TOKEN, LINKEDIN_AUTHOR_URN, and LINKEDIN_VERSION are required")
	}
	if strings.Contains(bundle.LinkedIn, "{{CANONICAL_URL}}") {
		return Result{}, errors.New("linkedin post still contains {{CANONICAL_URL}}; publish the article first or set canonical_url")
	}
	payload := map[string]any{
		"author": authorURN, "commentary": bundle.LinkedIn, "visibility": "PUBLIC",
		"distribution": map[string]any{"feedDistribution": "MAIN_FEED", "targetEntities": []any{}, "thirdPartyDistributionChannels": []any{}},
		"lifecycleState": "PUBLISHED", "isReshareDisabledByAuthor": false,
	}
	request, err := newJSONRequest(ctx, http.MethodPost, "https://api.linkedin.com/rest/posts", payload)
	if err != nil {
		return Result{}, err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Linkedin-Version", version)
	request.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return Result{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		raw, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		return Result{}, fmt.Errorf("publish LinkedIn post: HTTP %d: %s", response.StatusCode, strings.TrimSpace(string(raw)))
	}
	id := response.Header.Get("x-restli-id")
	return Result{Platform: "linkedin", ID: id, URL: "https://www.linkedin.com/feed/update/" + id}, nil
}

func newJSONRequest(ctx context.Context, method, url string, payload any) (*http.Request, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	return request, nil
}

func (c *Client) do(request *http.Request, output any) error {
	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		raw, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		return fmt.Errorf("HTTP %d: %s", response.StatusCode, strings.TrimSpace(string(raw)))
	}
	return json.NewDecoder(response.Body).Decode(output)
}