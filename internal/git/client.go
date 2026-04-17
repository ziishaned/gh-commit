package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Client struct{}

// HasStagedChanges checks if there are any staged changes
func (c *Client) HasStagedChanges() bool {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	err := cmd.Run()
	return err != nil
}

// StageAllChanges stages all unstaged changes
func (c *Client) StageAllChanges() error {
	cmd := exec.Command("git", "add", ".")
	return cmd.Run()
}

// HasAnyChanges checks if there are any changes (staged or unstaged)
func (c *Client) HasAnyChanges() bool {
	// Check for unstaged changes
	cmd1 := exec.Command("git", "diff", "--quiet")
	unstaged := cmd1.Run() != nil

	// Check for staged changes
	cmd2 := exec.Command("git", "diff", "--cached", "--quiet")
	staged := cmd2.Run() != nil

	// Check for untracked files
	cmd3 := exec.Command("git", "ls-files", "--others", "--exclude-standard")
	output, err := cmd3.Output()
	if err != nil {
		// If git command fails, log the error but don't fail the entire check
		// We'll conservatively assume there might be changes
		return true
	}
	untracked := len(strings.TrimSpace(string(output))) > 0

	return unstaged || staged || untracked
}

// GetStagedDiff returns the full diff of staged changes
func (c *Client) GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git diff failed: %w: %s", err, stderr.String())
	}

	diff := stdout.String()
	if strings.TrimSpace(diff) == "" {
		return "", fmt.Errorf("no staged changes found")
	}

	return diff, nil
}
