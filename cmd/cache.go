package cmd

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
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
	tracker := newProgressTracker("issues")
	issues, err := client.FetchAllIssues(repo.Owner, repo.Name, since, tracker.update)
	if err != nil {
		return fmt.Errorf("fetching issues: %w", err)
	}
	tracker.done()
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
	tracker = newProgressTracker("pull requests")
	prs, err := client.FetchAllPRs(repo.Owner, repo.Name, since, tracker.update)
	if err != nil {
		return fmt.Errorf("fetching pull requests: %w", err)
	}
	tracker.done()
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

// progressTracker adapts github.ProgressFunc callbacks onto a schollz
// progress bar. The bar is created lazily on the first update so it can adopt
// the server-reported total (determinate bar) or fall back to a spinner when
// the total is unknown (indeterminate PR delta scans).
type progressTracker struct {
	label  string
	writer io.Writer
	bar    *progressbar.ProgressBar
	maxSet bool
}

func newProgressTracker(label string) *progressTracker {
	return &progressTracker{label: label, writer: progressWriter()}
}

// update implements github.ProgressFunc.
func (t *progressTracker) update(current, total int) {
	if t.bar == nil {
		max := -1 // -1 selects spinner mode for an unknown length
		if total > 0 {
			max = total
		}
		t.bar = progressbar.NewOptions(max,
			progressbar.OptionSetDescription(t.label),
			progressbar.OptionSetWriter(t.writer),
			progressbar.OptionShowCount(),
			progressbar.OptionSetWidth(30),
			progressbar.OptionSetPredictTime(false),
			progressbar.OptionThrottle(50*time.Millisecond),
		)
	} else if total > 0 && !t.maxSet {
		t.bar.ChangeMax(total)
		t.maxSet = true
	}
	if total > 0 {
		t.maxSet = true
	}
	t.bar.Set(current)
}

// done finalises the bar: fills a determinate bar if it did not naturally
// complete, then emits a trailing newline.
func (t *progressTracker) done() {
	if t.bar == nil {
		return
	}
	if !t.bar.IsFinished() && t.bar.GetMax() > 0 {
		t.bar.Finish()
	}
	fmt.Fprintln(t.writer)
}

// progressWriter returns os.Stderr when it is a terminal device, otherwise
// io.Discard so progress bars never clutter captured or piped output.
func progressWriter() io.Writer {
	if fi, err := os.Stderr.Stat(); err == nil && fi.Mode()&os.ModeCharDevice != 0 {
		return os.Stderr
	}
	return io.Discard
}
