package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/tomzxcode/gh-cached/internal/cache"
	"github.com/tomzxcode/gh-cached/internal/github"
)

// ---------------------------------------------------------------------------
// Command tree
// ---------------------------------------------------------------------------

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Work with GitHub pull requests",
}

// pr list flags
var (
	prListApp      string
	prListAssignee string
	prListAuthor   string
	prListBase     string
	prListDraft    bool
	prListHead     string
	prListLabels   []string
	prListLimit    int
	prListSearch   string
	prListState    string
)

// pr view flags
var prViewComments bool

var prListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pull requests",
	RunE:  runPRList,
}

var prViewCmd = &cobra.Command{
	Use:   "view <number>",
	Short: "View a pull request",
	Args:  cobra.ExactArgs(1),
	RunE:  runPRView,
}

func init() {
	prCmd.AddCommand(prListCmd)
	prCmd.AddCommand(prViewCmd)

	prListCmd.Flags().StringVar(&prListApp, "app", "", "Filter by GitHub App author")
	prListCmd.Flags().StringVarP(&prListAssignee, "assignee", "a", "", "Filter by assignee")
	prListCmd.Flags().StringVarP(&prListAuthor, "author", "A", "", "Filter by author")
	prListCmd.Flags().StringVarP(&prListBase, "base", "B", "", "Filter by base branch")
	prListCmd.Flags().BoolVarP(&prListDraft, "draft", "d", false, "Filter by draft state")
	prListCmd.Flags().StringVarP(&prListHead, "head", "H", "", "Filter by head branch")
	prListCmd.Flags().StringSliceVarP(&prListLabels, "label", "l", nil, "Filter by label")
	prListCmd.Flags().IntVarP(&prListLimit, "limit", "L", 30, "Maximum number of items to fetch")
	prListCmd.Flags().StringVarP(&prListSearch, "search", "S", "", "Search pull requests with query")
	prListCmd.Flags().StringVarP(&prListState, "state", "s", "open", "Filter by state: {open|closed|merged|all}")

	prViewCmd.Flags().BoolVarP(&prViewComments, "comments", "c", false, "View pull request comments")
}

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------

func runPRList(cmd *cobra.Command, args []string) error {
	repo, err := getRepo()
	if err != nil {
		return err
	}

	store := cache.NewStore()

	// Try serving from cache when it is fresh.
	if fresh, _ := store.IsCacheFresh(repo.Host, repo.Owner, repo.Name); fresh {
		if prs, err := store.LoadAllPRs(repo.Host, repo.Owner, repo.Name); err == nil {
			filtered := filterPRs(prs, prListState, prListAssignee, prListAuthor,
				prListLabels, prListBase, prListHead, prListApp, prListSearch, prListDraft)
			sort.Slice(filtered, func(i, j int) bool {
				return filtered[i].UpdatedAt.After(filtered[j].UpdatedAt)
			})
			if prListLimit > 0 && len(filtered) > prListLimit {
				filtered = filtered[:prListLimit]
			}
			printPRList(filtered)
			return nil
		}
	}

	// Fall back to the GitHub API.
	client, err := github.NewClient(repo.Host)
	if err != nil {
		return err
	}

	opts := github.PRListOptions{
		Limit:    prListLimit,
		State:    prListState,
		Assignee: prListAssignee,
		Author:   prListAuthor,
		Labels:   prListLabels,
		Base:     prListBase,
		Head:     prListHead,
		Draft:    prListDraft,
		Search:   prListSearch,
		App:      prListApp,
	}

	prs, err := client.ListPRs(repo.Owner, repo.Name, opts)
	if err != nil {
		return err
	}

	printPRList(prs)
	return nil
}

func runPRView(cmd *cobra.Command, args []string) error {
	number, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid pull request number: %s", args[0])
	}

	repo, err := getRepo()
	if err != nil {
		return err
	}

	store := cache.NewStore()
	if pr, mtime, err := store.LoadPR(repo.Host, repo.Owner, repo.Name, number); err == nil {
		if time.Since(mtime) < 60*time.Minute {
			printPRView(pr, prViewComments)
			return nil
		}
	}

	client, err := github.NewClient(repo.Host)
	if err != nil {
		return err
	}

	pr, err := client.GetPR(repo.Owner, repo.Name, number)
	if err != nil {
		return err
	}

	// Persist to cache for future use.
	_ = store.SavePR(repo.Host, repo.Owner, repo.Name, pr)

	printPRView(pr, prViewComments)
	return nil
}

// ---------------------------------------------------------------------------
// Filtering
// ---------------------------------------------------------------------------

func filterPRs(prs []*github.PullRequest, state, assignee, author string,
	labels []string, base, head, app, search string, draft bool) []*github.PullRequest {

	var result []*github.PullRequest
	for _, pr := range prs {
		if state != "all" && state != "" {
			if !strings.EqualFold(pr.State, state) {
				continue
			}
		}
		if assignee != "" {
			found := false
			for _, a := range pr.Assignees {
				if strings.EqualFold(a.Login, assignee) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if author != "" && !strings.EqualFold(pr.Author.Login, author) {
			continue
		}
		if len(labels) > 0 {
			if !hasAllLabelsPR(pr.Labels, labels) {
				continue
			}
		}
		if base != "" && !strings.EqualFold(pr.BaseRefName, base) {
			continue
		}
		if head != "" && !strings.EqualFold(pr.HeadRefName, head) {
			continue
		}
		if draft && !pr.IsDraft {
			continue
		}
		// app cannot be verified from cached data; skip.
		if search != "" {
			q := strings.ToLower(search)
			if !strings.Contains(strings.ToLower(pr.Title), q) &&
				!strings.Contains(strings.ToLower(pr.Body), q) {
				continue
			}
		}
		result = append(result, pr)
	}
	return result
}

func hasAllLabelsPR(prLabels []github.Label, wantLabels []string) bool {
	for _, want := range wantLabels {
		found := false
		for _, l := range prLabels {
			if strings.EqualFold(l.Name, want) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// Display
// ---------------------------------------------------------------------------

func printPRList(prs []*github.PullRequest) {
	if len(prs) == 0 {
		fmt.Println("No pull requests found.")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	defer w.Flush()
	for _, pr := range prs {
		labels := make([]string, len(pr.Labels))
		for i, l := range pr.Labels {
			labels[i] = l.Name
		}
		draft := ""
		if pr.IsDraft {
			draft = " [draft]"
		}
		branch := fmt.Sprintf("%s → %s", pr.HeadRefName, pr.BaseRefName)
		fmt.Fprintf(w, "#%d\t%s%s\t%s\t%s\n",
			pr.Number,
			truncate(pr.Title, 55),
			draft,
			branch,
			pr.UpdatedAt.Format("2006-01-02"),
		)
	}
}

func printPRView(pr *github.PullRequest, showComments bool) {
	fmt.Printf("#%d %s\n", pr.Number, pr.Title)

	status := strings.ToUpper(pr.State)
	draftTag := ""
	if pr.IsDraft {
		draftTag = " • draft"
	}
	fmt.Printf("%s%s • opened by %s • %d comment(s)\n",
		status, draftTag, pr.Author.Login, len(pr.Comments))
	fmt.Println()

	fmt.Printf("Branch:    %s → %s\n", pr.HeadRefName, pr.BaseRefName)
	if len(pr.Labels) > 0 {
		names := make([]string, len(pr.Labels))
		for i, l := range pr.Labels {
			names[i] = l.Name
		}
		fmt.Printf("Labels:    %s\n", strings.Join(names, ", "))
	}
	if len(pr.Assignees) > 0 {
		logins := make([]string, len(pr.Assignees))
		for i, a := range pr.Assignees {
			logins[i] = a.Login
		}
		fmt.Printf("Assignees: %s\n", strings.Join(logins, ", "))
	}
	if pr.Milestone != nil {
		fmt.Printf("Milestone: %s\n", pr.Milestone.Title)
	}
	fmt.Printf("Created:   %s\n", pr.CreatedAt.Format("2006-01-02 15:04"))
	if pr.MergedAt != nil {
		fmt.Printf("Merged:    %s\n", pr.MergedAt.Format("2006-01-02 15:04"))
	} else if pr.ClosedAt != nil {
		fmt.Printf("Closed:    %s\n", pr.ClosedAt.Format("2006-01-02 15:04"))
	}
	fmt.Printf("URL:       %s\n", pr.URL)
	fmt.Println()
	fmt.Println(pr.Body)

	if showComments && len(pr.Comments) > 0 {
		for i, c := range pr.Comments {
			fmt.Printf("\n── Comment %d by %s (%s) ──\n\n",
				i+1, c.Author.Login, c.CreatedAt.Format("2006-01-02 15:04"))
			fmt.Println(c.Body)
		}
	}
}

