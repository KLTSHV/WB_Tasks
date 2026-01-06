package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"mirror/internal/crawl"
	"mirror/internal/fetch"
	"mirror/internal/urlmap"
)

func main() {
	var (
		rawURL   = flag.String("url", "", "Start URL (required)")
		outDir   = flag.String("out", "./out", "Output directory")
		depth    = flag.Int("depth", 1, "Recursion depth for page links a[href]")
		parallel = flag.Int("parallel", max(4, runtime.NumCPU()), "Max parallel downloads")
		timeout  = flag.Duration("timeout", 20*time.Second, "HTTP timeout per request")
		maxBytes = flag.Int64("max-bytes", 20<<20, "Max bytes per response (0=unlimited)")
	)
	flag.Parse()

	if *rawURL == "" {
		log.Fatal("missing -url")
	}

	mapper, err := urlmap.New(*rawURL, *outDir)
	if err != nil {
		log.Fatalf("urlmap init: %v", err)
	}

	httpClient := &http.Client{
		Timeout:   *timeout,
		Transport: fetch.NewTransport(*parallel),
	}

	cfg := crawl.Settings{
		MaxDepth:   *depth,
		Parallel:   *parallel,
		MaxBytes:   *maxBytes,
		UserAgent:  "mini-mirror/0.1",
		HTTPClient: httpClient,
		URLMap:     mapper,
	}

	ctx := context.Background()
	if err := crawl.Run(ctx, cfg, *rawURL); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done:", *outDir)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
