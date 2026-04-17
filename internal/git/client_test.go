package git

import (
	"os/exec"
	"testing"
)

func TestHasStagedChanges(t *testing.T) {
	// This test requires a git repository
	// We'll use git commands to set up test state

	client := &Client{}

	// Test with no staged changes
	hasStaged := client.HasStagedChanges()
	if hasStaged {
		t.Error("Expected no staged changes, got true")
	}

	// Create a test file and stage it
	exec.Command("touch", "test-file.txt").Run()
	exec.Command("git", "add", "test-file.txt").Run()
	defer exec.Command("git", "rm", "-f", "test-file.txt").Run()

	// Test with staged changes
	client2 := &Client{}
	hasStaged2 := client2.HasStagedChanges()
	if !hasStaged2 {
		t.Error("Expected staged changes, got false")
	}
}
