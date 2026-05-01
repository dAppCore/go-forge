package main

import (
	"context"
	"fmt"
	"os"

	forge "dappco.re/go/forge"
)

func main() {
	url := os.Getenv("FORGE_URL")
	token := os.Getenv("FORGE_TOKEN")
	if url == "" || token == "" {
		fmt.Println("skip: FORGE_URL and FORGE_TOKEN are required")
		return
	}

	f := forge.NewForge(url, token)
	repos, err := f.Repos.ListOrgRepos(context.Background(), "core")
	if err != nil {
		fmt.Fprintf(os.Stderr, "list core repos: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("core repos: %d\n", len(repos))
}
