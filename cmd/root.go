package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tomzxcode/gh-cached/internal/cache"
	"github.com/tomzxcode/gh-cached/internal/gitremote"
	"github.com/tomzxcode/gh-cached/internal/github"
	"github.com/tomzxcode/gh-cached/internal/version"
)

var (
	repoFlag   string
	apiURLFlag string
	cacheDir   string
)

var rootCmd = &cobra.Command{
	Use:           "gh-cached",
	Short:         "GitHub CLI with local caching",
	Version:       version.Get(),
	SilenceUsage:  true,
	SilenceErrors: true,
	Long: `gh-cached is a GitHub CLI that caches issues, pull requests, and their comments
locally to minimise API calls. The cache lives at ~/.cache/gh-cached/<host>/<owner>/<repo>.`,
}

// Execute is the entry point called from main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&repoFlag, "repo", "", "Repository in [HOST/]OWNER/REPO format")
	rootCmd.PersistentFlags().StringVar(&apiURLFlag, "api-url", "", "Override the GitHub GraphQL API endpoint URL (for testing)")
	rootCmd.PersistentFlags().StringVar(&cacheDir, "cache-dir", "", "Override the cache directory path")

	rootCmd.AddCommand(issueCmd)
	rootCmd.AddCommand(prCmd)
	rootCmd.AddCommand(cacheCmd)
	rootCmd.AddCommand(repoCmd)
}

// getRepo resolves the target repository from the --repo flag or the current
// directory's git remote.
func getRepo() (*gitremote.Repo, error) {
	if repoFlag != "" {
		return gitremote.ParseRepo(repoFlag)
	}
	return gitremote.DetectRepo()
}

// newClient creates a GitHub client, using --api-url if provided.
func newClient(host string) (*github.Client, error) {
	if apiURLFlag != "" {
		return github.NewClientWithURL(apiURLFlag, "test-token", host)
	}
	return github.NewClient(host)
}

// newStore creates a cache store, using --cache-dir if provided.
func newStore() *cache.Store {
	if cacheDir != "" {
		return cache.NewStoreWithPath(cacheDir)
	}
	return cache.NewStore()
}
