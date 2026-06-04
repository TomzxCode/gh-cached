package cmd

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/tomzxcode/gh-cached/internal/mockserver"
)

func TestBuildSimConfig_DefaultPreset(t *testing.T) {
	c, _, _ := makeMockServeCmd(t, []string{"--preset", "default"})
	cfg, err := buildSimConfig(c)
	if err != nil {
		t.Fatal(err)
	}
	defaultCfg := mockserver.DefaultConfig()
	if cfg.NumUsers != defaultCfg.NumUsers {
		t.Errorf("NumUsers: got %d, want %d", cfg.NumUsers, defaultCfg.NumUsers)
	}
	if cfg.Seed != defaultCfg.Seed {
		t.Errorf("Seed: got %d, want %d", cfg.Seed, defaultCfg.Seed)
	}
}

func TestBuildSimConfig_SmallPreset(t *testing.T) {
	c, _, _ := makeMockServeCmd(t, []string{"--preset", "small"})
	cfg, err := buildSimConfig(c)
	if err != nil {
		t.Fatal(err)
	}
	smallCfg := mockserver.SmallConfig()
	if cfg.NumUsers != smallCfg.NumUsers {
		t.Errorf("NumUsers: got %d, want %d", cfg.NumUsers, smallCfg.NumUsers)
	}
}

func TestBuildSimConfig_UnknownPreset(t *testing.T) {
	_, _, err := makeMockServeCmdWithErr(t, []string{"--preset", "bogus"})
	if err == nil {
		t.Fatal("expected error for unknown preset")
	}
}

func TestBuildSimConfig_FlagOverrides(t *testing.T) {
	c, _, _ := makeMockServeCmd(t, []string{
		"--preset", "default",
		"--users", "5",
		"--repos", "org/a,org/b",
		"--history", "48h",
		"--issues-per-repo", "10",
		"--prs-per-repo", "8",
		"--comments-per-issue", "3",
		"--comments-per-pr", "2",
		"--assignees-per-issue", "1",
		"--assignees-per-pr", "1",
		"--labels-per-item", "4",
		"--milestones-per-repo", "3",
		"--close-rate", "0.8",
		"--merge-rate", "0.9",
		"--draft-rate", "0.2",
		"--review-rate", "0.7",
		"--seed", "123",
		"--activity-bursts", "2",
	})
	cfg, err := buildSimConfig(c)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.NumUsers != 5 {
		t.Errorf("NumUsers: got %d, want 5", cfg.NumUsers)
	}
	if len(cfg.Repos) != 2 || cfg.Repos[0] != "org/a" || cfg.Repos[1] != "org/b" {
		t.Errorf("Repos: got %v", cfg.Repos)
	}
	if cfg.History != 48*time.Hour {
		t.Errorf("History: got %v, want 48h", cfg.History)
	}
	if cfg.IssuesPerRepo != 10 {
		t.Errorf("IssuesPerRepo: got %d, want 10", cfg.IssuesPerRepo)
	}
	if cfg.PRsPerRepo != 8 {
		t.Errorf("PRsPerRepo: got %d, want 8", cfg.PRsPerRepo)
	}
	if cfg.CommentsPerIssue != 3 {
		t.Errorf("CommentsPerIssue: got %d, want 3", cfg.CommentsPerIssue)
	}
	if cfg.CommentsPerPR != 2 {
		t.Errorf("CommentsPerPR: got %d, want 2", cfg.CommentsPerPR)
	}
	if cfg.AssigneesPerIssue != 1 {
		t.Errorf("AssigneesPerIssue: got %d, want 1", cfg.AssigneesPerIssue)
	}
	if cfg.AssigneesPerPR != 1 {
		t.Errorf("AssigneesPerPR: got %d, want 1", cfg.AssigneesPerPR)
	}
	if cfg.LabelsPerItem != 4 {
		t.Errorf("LabelsPerItem: got %d, want 4", cfg.LabelsPerItem)
	}
	if cfg.MilestonesPerRepo != 3 {
		t.Errorf("MilestonesPerRepo: got %d, want 3", cfg.MilestonesPerRepo)
	}
	if cfg.CloseRate != 0.8 {
		t.Errorf("CloseRate: got %f, want 0.8", cfg.CloseRate)
	}
	if cfg.MergeRate != 0.9 {
		t.Errorf("MergeRate: got %f, want 0.9", cfg.MergeRate)
	}
	if cfg.DraftRate != 0.2 {
		t.Errorf("DraftRate: got %f, want 0.2", cfg.DraftRate)
	}
	if cfg.ReviewRate != 0.7 {
		t.Errorf("ReviewRate: got %f, want 0.7", cfg.ReviewRate)
	}
	if cfg.Seed != 123 {
		t.Errorf("Seed: got %d, want 123", cfg.Seed)
	}
	if cfg.ActivityBursts != 2 {
		t.Errorf("ActivityBursts: got %d, want 2", cfg.ActivityBursts)
	}
}

func TestBuildSimConfig_PartialOverrides(t *testing.T) {
	c, _, _ := makeMockServeCmd(t, []string{
		"--preset", "default",
		"--users", "12",
		"--seed", "99",
	})
	cfg, err := buildSimConfig(c)
	if err != nil {
		t.Fatal(err)
	}
	defaultCfg := mockserver.DefaultConfig()

	if cfg.NumUsers != 12 {
		t.Errorf("NumUsers: got %d, want 12", cfg.NumUsers)
	}
	if cfg.Seed != 99 {
		t.Errorf("Seed: got %d, want 99", cfg.Seed)
	}
	if cfg.IssuesPerRepo != defaultCfg.IssuesPerRepo {
		t.Errorf("IssuesPerRepo should remain default: got %d, want %d", cfg.IssuesPerRepo, defaultCfg.IssuesPerRepo)
	}
}

func TestBuildSimConfig_NonePresetRequiresFlags(t *testing.T) {
	c, _, _ := makeMockServeCmd(t, []string{
		"--preset", "none",
		"--users", "3",
		"--repos", "org/test",
	})
	cfg, err := buildSimConfig(c)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.NumUsers != 3 {
		t.Errorf("NumUsers: got %d, want 3", cfg.NumUsers)
	}
}

func TestBuildSimConfig_NonePresetMissingRepos(t *testing.T) {
	_, _, err := makeMockServeCmdWithErr(t, []string{"--preset", "none", "--users", "3"})
	if err == nil {
		t.Fatal("expected error for missing repos")
	}
}

func TestBuildSimConfig_NonePresetMissingUsers(t *testing.T) {
	_, _, err := makeMockServeCmdWithErr(t, []string{"--preset", "none", "--repos", "org/test"})
	if err == nil {
		t.Fatal("expected error for missing users")
	}
}

func TestBuildSimConfig_StatsFlag(t *testing.T) {
	c, _, _ := makeMockServeCmd(t, []string{"--preset", "small", "--stats"})
	stats, _ := c.Flags().GetBool("stats")
	if !stats {
		t.Error("stats flag should be true")
	}
}

func TestBuildSimConfig_GeneratesValidScenario(t *testing.T) {
	c, _, _ := makeMockServeCmd(t, []string{
		"--preset", "small",
		"--seed", "42",
	})
	cfg, err := buildSimConfig(c)
	if err != nil {
		t.Fatal(err)
	}

	scenario := mockserver.Generate(cfg)
	if len(scenario.Issues) == 0 {
		t.Error("expected at least one issue")
	}
	if len(scenario.PRs) == 0 {
		t.Error("expected at least one PR")
	}

	summary := scenario.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func makeMockServeCmd(t *testing.T, args []string) (*cobra.Command, []string, error) {
	t.Helper()
	c := &cobra.Command{Use: "serve"}
	addMockServeFlags(c)
	c.SetArgs(args)
	err := c.ParseFlags(args)
	return c, args, err
}

func makeMockServeCmdWithErr(t *testing.T, args []string) (*cobra.Command, []string, error) {
	t.Helper()
	c := &cobra.Command{Use: "serve"}
	addMockServeFlags(c)
	c.SetArgs(args)
	err := c.ParseFlags(args)
	if err != nil {
		return c, args, err
	}
	_, err2 := buildSimConfig(c)
	return c, args, err2
}

func addMockServeFlags(c *cobra.Command) {
	c.Flags().IntVar(&mockServeUsers, "users", 0, "")
	c.Flags().StringSliceVar(&mockServeRepos, "repos", nil, "")
	c.Flags().DurationVar(&mockServeHistory, "history", 0, "")
	c.Flags().IntVar(&mockServeIssuesPerRepo, "issues-per-repo", 0, "")
	c.Flags().IntVar(&mockServePRsPerRepo, "prs-per-repo", 0, "")
	c.Flags().IntVar(&mockServeCommentsIssue, "comments-per-issue", 0, "")
	c.Flags().IntVar(&mockServeCommentsPR, "comments-per-pr", 0, "")
	c.Flags().IntVar(&mockServeAssigneesIssue, "assignees-per-issue", 0, "")
	c.Flags().IntVar(&mockServeAssigneesPR, "assignees-per-pr", 0, "")
	c.Flags().IntVar(&mockServeLabels, "labels-per-item", 0, "")
	c.Flags().IntVar(&mockServeMilestones, "milestones-per-repo", 0, "")
	c.Flags().Float64Var(&mockServeCloseRate, "close-rate", -1, "")
	c.Flags().Float64Var(&mockServeMergeRate, "merge-rate", -1, "")
	c.Flags().Float64Var(&mockServeDraftRate, "draft-rate", -1, "")
	c.Flags().Float64Var(&mockServeReviewRate, "review-rate", -1, "")
	c.Flags().Int64Var(&mockServeSeed, "seed", 0, "")
	c.Flags().IntVar(&mockServeBursts, "activity-bursts", -1, "")
	c.Flags().StringVar(&mockServePreset, "preset", "default", "")
	c.Flags().BoolVar(&mockServeStats, "stats", false, "")
}
