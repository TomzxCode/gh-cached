package gitremote

import (
	"fmt"
	"os/exec"
	"strings"
)

// Repo holds the parsed components of a GitHub repository reference.
type Repo struct {
	Host  string
	Owner string
	Name  string
}

func (r *Repo) String() string {
	return fmt.Sprintf("%s/%s/%s", r.Host, r.Owner, r.Name)
}

// ParseRepo parses a [HOST/]OWNER/REPO string.
func ParseRepo(s string) (*Repo, error) {
	parts := strings.Split(s, "/")
	switch len(parts) {
	case 2:
		return &Repo{Host: "github.com", Owner: parts[0], Name: parts[1]}, nil
	case 3:
		return &Repo{Host: parts[0], Owner: parts[1], Name: parts[2]}, nil
	default:
		return nil, fmt.Errorf("invalid repo format %q: expected [HOST/]OWNER/REPO", s)
	}
}

// DetectRepo reads the origin remote URL from the current directory's git repo.
func DetectRepo() (*Repo, error) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return nil, fmt.Errorf("could not get git remote origin: %w", err)
	}
	return parseRemoteURL(strings.TrimSpace(string(out)))
}

func parseRemoteURL(rawURL string) (*Repo, error) {
	rawURL = strings.TrimSuffix(rawURL, ".git")

	// SSH: git@github.com:owner/repo
	if strings.HasPrefix(rawURL, "git@") {
		rawURL = strings.TrimPrefix(rawURL, "git@")
		parts := strings.SplitN(rawURL, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid SSH remote URL: %s", rawURL)
		}
		ownerRepo := strings.SplitN(parts[1], "/", 2)
		if len(ownerRepo) != 2 {
			return nil, fmt.Errorf("invalid SSH remote URL: %s", rawURL)
		}
		return &Repo{Host: parts[0], Owner: ownerRepo[0], Name: ownerRepo[1]}, nil
	}

	// HTTPS: https://github.com/owner/repo
	if strings.HasPrefix(rawURL, "https://") || strings.HasPrefix(rawURL, "http://") {
		rawURL = strings.TrimPrefix(rawURL, "https://")
		rawURL = strings.TrimPrefix(rawURL, "http://")
		parts := strings.SplitN(rawURL, "/", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid HTTPS remote URL: %s", rawURL)
		}
		return &Repo{Host: parts[0], Owner: parts[1], Name: parts[2]}, nil
	}

	return nil, fmt.Errorf("unrecognized remote URL format: %s", rawURL)
}
