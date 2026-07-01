package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	cacheDuration int
	cacheForce    bool
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Fetch and cache issues and PRs (including comments)",
	Long: `Fetches all issues and pull requests (with comments) from GitHub and writes them
to ~/.cache/ghx/cache/<host>/<owner>/<repo>. Subsequent list/view commands will
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

	store := newStore()

	if !cacheForce {
		if fresh, _ := store.IsCacheFreshWithDuration(repo.Host, repo.Owner, repo.Name, cacheDuration); fresh {
			fmt.Printf("Cache is still fresh (within %d minutes). Use --force to refresh anyway.\n", cacheDuration)
			return nil
		}
	}

	// Determine the delta cutoff: the previous cache timestamp (if any).
	// --force bypasses the freshness check above but we still want to do a
	// full fetch, so only use the delta when not forcing.
	var since *time.Time
	if !cacheForce {
		if info, err := store.LoadCacheInfo(repo.Host, repo.Owner, repo.Name); err == nil {
			t := info.CachedAt
			since = &t
		}
	}

	client, err := newClient(repo.Host)
	if err != nil {
		return err
	}

	if since != nil {
		fmt.Printf("Fetching issues updated since %s for %s/%s...\n", since.Format("2006-01-02 15:04"), repo.Owner, repo.Name)
	} else {
		fmt.Printf("Caching issues for %s/%s...\n", repo.Owner, repo.Name)
	}
	issues, err := client.FetchAllIssues(repo.Owner, repo.Name, since)
	if err != nil {
		return fmt.Errorf("fetching issues: %w", err)
	}
	for _, issue := range issues {
		if err := store.SaveIssue(repo.Host, repo.Owner, repo.Name, issue); err != nil {
			return fmt.Errorf("saving issue #%d: %w", issue.Number, err)
		}
	}
	fmt.Printf("Cached %d issue(s).\n", len(issues))

	if since != nil {
		fmt.Printf("Fetching PRs updated since %s for %s/%s...\n", since.Format("2006-01-02 15:04"), repo.Owner, repo.Name)
	} else {
		fmt.Printf("Caching pull requests for %s/%s...\n", repo.Owner, repo.Name)
	}
	prs, err := client.FetchAllPRs(repo.Owner, repo.Name, since)
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
