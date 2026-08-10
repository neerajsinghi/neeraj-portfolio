package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostRevalidateSendsBearerToken(t *testing.T) {
	var gotAuth string
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	postRevalidate(testServer.Client(), testServer.URL, "revalidate-secret")

	if gotAuth != "Bearer revalidate-secret" {
		t.Fatalf("Authorization header = %q, want Bearer revalidate-secret", gotAuth)
	}
}

func TestNotifyFrontendRevalidateSkipsWhenUnconfigured(t *testing.T) {
	t.Setenv("FRONTEND_REVALIDATE_URL", "")
	t.Setenv("FRONTEND_REVALIDATE_TOKEN", "")
	notifyFrontendRevalidate() // must not panic or block when unconfigured
}
