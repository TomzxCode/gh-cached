package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tomzxcode/gh-cached/internal/cache"
	"github.com/tomzxcode/gh-cached/internal/github"
)

var (
	cacheDuration int
	cacheForce    bool
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Fetch and cache issues and PRs (including comments)",
	Long: `Fetches all issues and pull requests (with comments) from GitHub and writes them
to ~/.cache/gh-cached/<host>/<owner>/<repo>. Subsequent list/view commands will
serve results from this cache until it expires.`,
	RunE: runCache,
}

func init() {
	cacheCmd.Flags().IntVar(&cacheDuration, "cache-duration", 60, "Cache duration in minutes")
	cacheCmd.Flags().BoolVar(&cacheForce, "force", false, "Re-fetch even if the cache is still fresh")
}

func runCache(cmd *cobra.Command, args []string) error {
	repo, err := getRepo()
	if err != nil {
		return err
	}

	store := cache.NewStore()

	if !cacheForce {
		if fresh, _ := store.IsCacheFreshWithDuration(repo.Host, repo.Owner, repo.Name, cacheDuration); fresh {
			fmt.Printf("Cache is still fresh (within %d minutes). Use --force to refresh anyway.\n", cacheDuration)
			return nil
		}
	}

	client, err := github.NewClient(repo.Host)
	if err != nil {
		return err
	}

	fmt.Printf("Caching issues for %s/%s...\n", repo.Owner, repo.Name)
	issues, err := client.FetchAllIssues(repo.Owner, repo.Name)
	if err != nil {
		return fmt.Errorf("fetching issues: %w", err)
	}
	for _, issue := range issues {
		if err := store.SaveIssue(repo.Host, repo.Owner, repo.Name, issue); err != nil {
			return fmt.Errorf("saving issue #%d: %w", issue.Number, err)
		}
	}
	fmt.Printf("Cached %d issue(s).\n", len(issues))

	fmt.Printf("Caching pull requests for %s/%s...\n", repo.Owner, repo.Name)
	prs, err := client.FetchAllPRs(repo.Owner, repo.Name)
	if err != nil {
		return fmt.Errorf("fetching pull requests: %w", err)
	}
	for _, pr := range prs {
		if err := store.SavePR(repo.Host, repo.Owner, repo.Name, pr); err != nil {
			return fmt.Errorf("saving PR #%d: %w", pr.Number, err)
		}
	}
	fmt.Printf("Cached %d pull request(s).\n", len(prs))

	if err := store.SaveCacheInfo(repo.Host, repo.Owner, repo.Name, cacheDuration); err != nil {
		return fmt.Errorf("saving cache info: %w", err)
	}

	fmt.Printf("Cache updated. Valid for %d minute(s).\n", cacheDuration)
	return nil
}
