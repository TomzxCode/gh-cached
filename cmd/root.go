package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tomzxcode/gh-cached/internal/gitremote"
)

var repoFlag string

var rootCmd = &cobra.Command{
	Use:   "gh-cached",
	Short: "GitHub CLI with local caching",
	Long: `gh-cached is a GitHub CLI that caches issues, pull requests, and their comments
locally to minimise API calls. The cache lives at ~/.cache/gh-cached/<host>/<owner>/<repo>.`,
}

// Execute is the entry point called from main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&repoFlag, "repo", "", "Repository in [HOST/]OWNER/REPO format")

	rootCmd.AddCommand(issueCmd)
	rootCmd.AddCommand(prCmd)
	rootCmd.AddCommand(cacheCmd)
}

// getRepo resolves the target repository from the --repo flag or the current
// directory's git remote.
func getRepo() (*gitremote.Repo, error) {
	if repoFlag != "" {
		return gitremote.ParseRepo(repoFlag)
	}
	return gitremote.DetectRepo()
}
