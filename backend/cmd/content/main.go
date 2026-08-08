package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"neeraj-portfolio/backend/internal/content"
	llmanthropic "neeraj-portfolio/backend/internal/llm/anthropic"
	"neeraj-portfolio/backend/internal/publisher"
)

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "content:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New("usage: content <generate|publish> [options]")
	}
	switch args[0] {
	case "generate":
		return generate(ctx, args[1:])
	case "publish":
		return publish(ctx, args[1:])
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func generate(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("generate", flag.ContinueOnError)
	topic := flags.String("topic", "", "article topic")
	audience := flags.String("audience", "", "target audience")
	canonicalBase := flags.String("canonical-base", "", "canonical site base URL when the article already exists there")
	output := flags.String("output", "content/drafts/draft.json", "output JSON path")
	if err := flags.Parse(args); err != nil {
		return err
	}

	bundle, err := content.NewGenerator(llmanthropic.New()).Generate(ctx, *topic, *audience, *canonicalBase)
	if err != nil {
		return err
	}
	return writeJSON(*output, bundle)
}

func publish(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("publish", flag.ContinueOnError)
	file := flags.String("file", "", "approved bundle JSON path")
	platforms := flags.String("platforms", "devto", "comma-separated: devto,linkedin")
	approve := flags.Bool("approve", false, "confirm public publishing; without this flag DEV receives a draft and LinkedIn is skipped")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *file == "" {
		return errors.New("--file is required")
	}

	var bundle content.Bundle
	raw, err := os.ReadFile(*file)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, &bundle); err != nil {
		return err
	}

	client := publisher.NewClient(&http.Client{Timeout: 30 * time.Second})
	results := make([]publisher.Result, 0)
	for _, platform := range strings.Split(*platforms, ",") {
		switch strings.TrimSpace(platform) {
		case "devto":
			result, err := client.PublishDEV(ctx, bundle, os.Getenv("DEVTO_API_KEY"), *approve)
			if err != nil {
				return err
			}
			results = append(results, result)
			if bundle.CanonicalURL == "" && result.URL != "" {
				bundle.CanonicalURL = result.URL
				bundle.LinkedIn = strings.ReplaceAll(bundle.LinkedIn, "{{CANONICAL_URL}}", result.URL)
				bundle.Social = strings.ReplaceAll(bundle.Social, "{{CANONICAL_URL}}", result.URL)
			}
		case "linkedin":
			if !*approve {
				fmt.Fprintln(os.Stderr, "content: skipping LinkedIn because --approve was not provided")
				continue
			}
			result, err := client.PublishLinkedIn(ctx, bundle, os.Getenv("LINKEDIN_ACCESS_TOKEN"), os.Getenv("LINKEDIN_AUTHOR_URN"), os.Getenv("LINKEDIN_VERSION"))
			if err != nil {
				return err
			}
			results = append(results, result)
		case "":
		default:
			return fmt.Errorf("unsupported platform %q", platform)
		}
	}
	return json.NewEncoder(os.Stdout).Encode(results)
}

func writeJSON(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, append(raw, '\n'), 0o644); err != nil {
		return err
	}
	fmt.Println(path)
	return nil
}