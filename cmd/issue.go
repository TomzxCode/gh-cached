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

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Work with GitHub issues",
}

// issue list flags
var (
	issueListApp       string
	issueListAssignee  string
	issueListAuthor    string
	issueListLabels    []string
	issueListLimit     int
	issueListMention   string
	issueListMilestone string
	issueListSearch    string
	issueListState     string
)

// issue view flags
var issueViewComments bool

var issueListCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues",
	RunE:  runIssueList,
}

var issueViewCmd = &cobra.Command{
	Use:   "view <number>",
	Short: "View an issue",
	Args:  cobra.ExactArgs(1),
	RunE:  runIssueView,
}

func init() {
	issueCmd.AddCommand(issueListCmd)
	issueCmd.AddCommand(issueViewCmd)

	issueListCmd.Flags().StringVar(&issueListApp, "app", "", "Filter by GitHub App author")
	issueListCmd.Flags().StringVarP(&issueListAssignee, "assignee", "a", "", "Filter by assignee")
	issueListCmd.Flags().StringVarP(&issueListAuthor, "author", "A", "", "Filter by author")
	issueListCmd.Flags().StringSliceVarP(&issueListLabels, "label", "l", nil, "Filter by label")
	issueListCmd.Flags().IntVarP(&issueListLimit, "limit", "L", 30, "Maximum number of issues to fetch")
	issueListCmd.Flags().StringVar(&issueListMention, "mention", "", "Filter by mention")
	issueListCmd.Flags().StringVarP(&issueListMilestone, "milestone", "m", "", "Filter by milestone number or title")
	issueListCmd.Flags().StringVarP(&issueListSearch, "search", "S", "", "Search issues with query")
	issueListCmd.Flags().StringVarP(&issueListState, "state", "s", "open", "Filter by state: {open|closed|all}")

	issueViewCmd.Flags().BoolVarP(&issueViewComments, "comments", "c", false, "View issue comments")
}

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------

func runIssueList(cmd *cobra.Command, args []string) error {
	repo, err := getRepo()
	if err != nil {
		return err
	}

	store := cache.NewStore()

	// Try serving from cache when it is fresh.
	if fresh, _ := store.IsCacheFresh(repo.Host, repo.Owner, repo.Name); fresh {
		if issues, err := store.LoadAllIssues(repo.Host, repo.Owner, repo.Name); err == nil {
			filtered := filterIssues(issues, issueListState, issueListAssignee, issueListAuthor,
				issueListLabels, issueListMilestone, issueListMention, issueListApp, issueListSearch)
			sort.Slice(filtered, func(i, j int) bool {
				return filtered[i].UpdatedAt.After(filtered[j].UpdatedAt)
			})
			if issueListLimit > 0 && len(filtered) > issueListLimit {
				filtered = filtered[:issueListLimit]
			}
			printIssueList(filtered)
			return nil
		}
	}

	// Fall back to the GitHub API.
	client, err := github.NewClient(repo.Host)
	if err != nil {
		return err
	}

	opts := github.IssueListOptions{
		Limit:     issueListLimit,
		State:     issueListState,
		Assignee:  issueListAssignee,
		Author:    issueListAuthor,
		Labels:    issueListLabels,
		Milestone: issueListMilestone,
		Mention:   issueListMention,
		Search:    issueListSearch,
		App:       issueListApp,
	}

	issues, err := client.ListIssues(repo.Owner, repo.Name, opts)
	if err != nil {
		return err
	}

	printIssueList(issues)
	return nil
}

func runIssueView(cmd *cobra.Command, args []string) error {
	number, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid issue number: %s", args[0])
	}

	repo, err := getRepo()
	if err != nil {
		return err
	}

	store := cache.NewStore()
	if issue, mtime, err := store.LoadIssue(repo.Host, repo.Owner, repo.Name, number); err == nil {
		if time.Since(mtime) < 60*time.Minute {
			printIssueView(issue, issueViewComments)
			return nil
		}
	}

	client, err := github.NewClient(repo.Host)
	if err != nil {
		return err
	}

	issue, err := client.GetIssue(repo.Owner, repo.Name, number)
	if err != nil {
		return err
	}

	// Persist to cache for future use.
	_ = store.SaveIssue(repo.Host, repo.Owner, repo.Name, issue)

	printIssueView(issue, issueViewComments)
	return nil
}

// ---------------------------------------------------------------------------
// Filtering
// ---------------------------------------------------------------------------

func filterIssues(issues []*github.Issue, state, assignee, author string,
	labels []string, milestone, mention, app, search string) []*github.Issue {

	var result []*github.Issue
	for _, issue := range issues {
		if state != "all" && state != "" {
			if !strings.EqualFold(issue.State, state) {
				continue
			}
		}
		if assignee != "" {
			found := false
			for _, a := range issue.Assignees {
				if strings.EqualFold(a.Login, assignee) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if author != "" && !strings.EqualFold(issue.Author.Login, author) {
			continue
		}
		if len(labels) > 0 {
			if !hasAllLabels(issue.Labels, labels) {
				continue
			}
		}
		if milestone != "" {
			if issue.Milestone == nil {
				continue
			}
			if !strings.EqualFold(issue.Milestone.Title, milestone) &&
				strconv.Itoa(issue.Milestone.Number) != milestone {
				continue
			}
		}
		// mention and app cannot be verified from cached data; skip filtering on them.
		if search != "" {
			q := strings.ToLower(search)
			if !strings.Contains(strings.ToLower(issue.Title), q) &&
				!strings.Contains(strings.ToLower(issue.Body), q) {
				continue
			}
		}
		result = append(result, issue)
	}
	return result
}

func hasAllLabels(issueLabels []github.Label, wantLabels []string) bool {
	for _, want := range wantLabels {
		found := false
		for _, l := range issueLabels {
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

func printIssueList(issues []*github.Issue) {
	if len(issues) == 0 {
		fmt.Println("No issues found.")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	defer w.Flush()
	for _, issue := range issues {
		labels := labelNames(issue.Labels)
		count := issue.CommentCount
		if count == 0 {
			count = len(issue.Comments) // fallback for cache files written before CommentCount was added
		}
		fmt.Fprintf(w, "#%d\t%s\t%s\t%d\t%s\n",
			issue.Number,
			truncate(issue.Title, 60),
			strings.Join(labels, ", "),
			count,
			issue.UpdatedAt.Format("2006-01-02"),
		)
	}
}

func printIssueView(issue *github.Issue, showComments bool) {
	fmt.Printf("#%d %s\n", issue.Number, issue.Title)

	status := strings.ToLower(issue.State)
	fmt.Printf("%s • opened by %s • %d comment(s)\n",
		status, issue.Author.Login, len(issue.Comments))
	fmt.Println()

	if len(issue.Labels) > 0 {
		fmt.Printf("Labels:    %s\n", strings.Join(labelNames(issue.Labels), ", "))
	}
	if len(issue.Assignees) > 0 {
		logins := make([]string, len(issue.Assignees))
		for i, a := range issue.Assignees {
			logins[i] = a.Login
		}
		fmt.Printf("Assignees: %s\n", strings.Join(logins, ", "))
	}
	if issue.Milestone != nil {
		fmt.Printf("Milestone: %s\n", issue.Milestone.Title)
	}
	fmt.Printf("Created:   %s\n", issue.CreatedAt.Format("2006-01-02 15:04"))
	if issue.ClosedAt != nil {
		fmt.Printf("Closed:    %s\n", issue.ClosedAt.Format("2006-01-02 15:04"))
	}
	fmt.Printf("URL:       %s\n", issue.URL)
	fmt.Println()
	fmt.Println(issue.Body)

	if showComments && len(issue.Comments) > 0 {
		for i, c := range issue.Comments {
			fmt.Printf("\n── Comment %d by %s (%s) ──\n\n",
				i+1, c.Author.Login, c.CreatedAt.Format("2006-01-02 15:04"))
			fmt.Println(c.Body)
		}
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func labelNames(labels []github.Label) []string {
	names := make([]string, len(labels))
	for i, l := range labels {
		names[i] = l.Name
	}
	return names
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
