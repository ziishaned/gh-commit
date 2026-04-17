package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gh-commit/internal/types"
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

// Commit creates a commit with the given message
func (c *Client) Commit(message string, dryRun bool) (*types.CommitResult, error) {
	result := &types.CommitResult{
		Message: message,
		DryRun:  dryRun,
	}

	if dryRun {
		result.Success = true
		return result, nil
	}

	// Create commit with message
	cmd := exec.Command("git", "commit", "-m", message)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		result.Success = false
		result.Error = fmt.Errorf("git commit failed: %w: %s", err, stderr.String())
		return result, result.Error
	}

	// Get commit hash
	cmd2 := exec.Command("git", "rev-parse", "HEAD")
	hashOutput, err := cmd2.Output()
	if err != nil {
		result.Success = false
		result.Error = fmt.Errorf("failed to get commit hash: %w", err)
		return result, result.Error
	}

	result.Success = true
	result.Hash = strings.TrimSpace(string(hashOutput))
	return result, nil
}

// FallbackToInteractive opens an editor for manual commit message editing
func (c *Client) FallbackToInteractive(message string, dryRun bool) error {
	// Validate and check editor availability
	editor, err := c.validateEditor()
	if err != nil {
		return fmt.Errorf("editor validation failed: %w", err)
	}

	// Create temp file with message
	tmpFile, err := os.CreateTemp("", "commit-message-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write pre-filled message
	if _, err := tmpFile.WriteString(message); err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}
	tmpFile.Close()

	// Open editor
	editCmd := exec.Command(editor, tmpFile.Name())
	editCmd.Stdin = os.Stdin
	editCmd.Stdout = os.Stdout
	editCmd.Stderr = os.Stderr

	if err := editCmd.Run(); err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}

	// Read edited message
	editedContent, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("failed to read edited message: %w", err)
	}

	editedMessage := string(editedContent)
	if strings.TrimSpace(editedMessage) == "" {
		return fmt.Errorf("empty commit message, aborting")
	}

	// If dry run, just show what would be committed
	if dryRun {
		fmt.Printf("\n[Dry run] Would commit with message:\n%s\n", editedMessage)
		return nil
	}

	// Commit with edited message
	cmd := exec.Command("git", "commit", "-m", editedMessage)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}

	return nil
}

// validateEditor checks if EDITOR is set and the editor exists
func (c *Client) validateEditor() (string, error) {
	editor := os.Getenv("EDITOR")

	// If EDITOR is not set, try common editors
	if editor == "" {
		commonEditors := []string{"vim", "nano", "vi"}
		for _, e := range commonEditors {
			if _, err := exec.LookPath(e); err == nil {
				editor = e
				break
			}
		}

		// If no common editor found, return error
		if editor == "" {
			return "", fmt.Errorf("no editor found. Please set EDITOR environment variable or ensure vim/nano/vi is available")
		}

		fmt.Printf("EDITOR not set, using '%s'\n", editor)
	}

	// Check if the editor exists
	if _, err := exec.LookPath(editor); err != nil {
		return "", fmt.Errorf("editor '%s' not found: %w. Please set EDITOR to a valid editor", editor, err)
	}

	return editor, nil
}
