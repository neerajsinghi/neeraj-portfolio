package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

// Repo is a public GitHub repository.
type Repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	Language    string `json:"language"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	Fork        bool   `json:"fork"`
	Updated     string `json:"updated_at"`
}

var (
	repoMu     sync.Mutex
	repoCache  []Repo
	repoCached time.Time
)

// User returns the configured GitHub username.
func User() string {
	if u := os.Getenv("GITHUB_USER"); u != "" {
		return u
	}
	return "neerajsinghi"
}

// FetchRepos returns public, non-fork repos sorted by stars, cached for 10 minutes.
func FetchRepos() ([]Repo, error) {
	repoMu.Lock()
	defer repoMu.Unlock()
	if repoCache != nil && time.Since(repoCached) < 10*time.Minute {
		return repoCache, nil
	}

	url := "https://api.github.com/users/" + User() + "/repos?per_page=100&sort=updated"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "neeraj-portfolio")
	if tok := os.Getenv("GITHUB_TOKEN"); tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}

	resp, err := (&http.Client{Timeout: 8 * time.Second}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github %d: %s", resp.StatusCode, string(body))
	}

	var all []Repo
	if err := json.NewDecoder(resp.Body).Decode(&all); err != nil {
		return nil, err
	}
	out := all[:0]
	for _, r := range all {
		if !r.Fork {
			out = append(out, r)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Stars > out[j].Stars })
	repoCache, repoCached = out, time.Now()
	return out, nil
}
