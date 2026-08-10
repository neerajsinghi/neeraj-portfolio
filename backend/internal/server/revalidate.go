package server

import (
	"context"
	"net/http"
	"os"
	"time"
)

var revalidateHTTPClient = &http.Client{Timeout: 5 * time.Second}

// notifyFrontendRevalidate asks the frontend to bust its blog fetch cache
// immediately after a publish-affecting change, instead of relying solely on
// the time-based revalidate window. Best effort: never fails the caller, and
// uses its own budget independent of the caller's (possibly near-expired)
// request context.
func notifyFrontendRevalidate() {
	url := os.Getenv("FRONTEND_REVALIDATE_URL")
	token := os.Getenv("FRONTEND_REVALIDATE_TOKEN")
	if url == "" || token == "" {
		return
	}
	postRevalidate(revalidateHTTPClient, url, token)
}

func postRevalidate(client *http.Client, url, token string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return
	}
	request.Header.Set("Authorization", "Bearer "+token)
	response, err := client.Do(request)
	if err != nil {
		return
	}
	_ = response.Body.Close()
}
