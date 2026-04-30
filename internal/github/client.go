package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Client is a minimal GitHub GraphQL client.
type Client struct {
	token      string
	host       string
	httpClient *http.Client
}

type gqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type gqlError struct {
	Message string `json:"message"`
}

type gqlResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []gqlError      `json:"errors"`
}

// NewClient creates a client for the given GitHub host (e.g. "github.com").
func NewClient(host string) (*Client, error) {
	token := resolveToken(host)
	if token == "" {
		return nil, fmt.Errorf("no GitHub token found; set GH_TOKEN or GITHUB_TOKEN, or run `gh auth login`")
	}
	return &Client{
		token:      token,
		host:       host,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func resolveToken(host string) string {
	for _, env := range []string{"GH_TOKEN", "GITHUB_TOKEN"} {
		if t := os.Getenv(env); t != "" {
			return t
		}
	}
	// Try gh CLI as fallback.
	if out, err := exec.Command("gh", "auth", "token", "--hostname", host).Output(); err == nil {
		if t := strings.TrimSpace(string(out)); t != "" {
			return t
		}
	}
	if out, err := exec.Command("gh", "auth", "token").Output(); err == nil {
		return strings.TrimSpace(string(out))
	}
	return ""
}

func (c *Client) endpoint() string {
	if c.host == "github.com" {
		return "https://api.github.com/graphql"
	}
	return fmt.Sprintf("https://%s/api/graphql", c.host)
}

// Query executes a GraphQL query and unmarshals the data field into result.
func (c *Client) Query(query string, variables map[string]interface{}, result interface{}) error {
	body, err := json.Marshal(gqlRequest{Query: query, Variables: variables})
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", c.endpoint(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "gh-cached/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var gqlResp gqlResponse
	if err := json.Unmarshal(respBody, &gqlResp); err != nil {
		return fmt.Errorf("parsing GraphQL response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		msgs := make([]string, len(gqlResp.Errors))
		for i, e := range gqlResp.Errors {
			msgs[i] = e.Message
		}
		return fmt.Errorf("GraphQL errors: %s", strings.Join(msgs, "; "))
	}

	if result != nil {
		if err := json.Unmarshal(gqlResp.Data, result); err != nil {
			return fmt.Errorf("parsing response data: %w", err)
		}
	}
	return nil
}
