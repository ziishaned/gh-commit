package git

import (
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
	output, _ := cmd3.Output()
	untracked := len(strings.TrimSpace(string(output))) > 0

	return unstaged || staged || untracked
}
