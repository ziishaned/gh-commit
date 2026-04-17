package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestRepo(t *testing.T) string {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Configure git user
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to configure git email: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to configure git name: %v", err)
	}

	// Create initial commit
	testFile := filepath.Join(tempDir, "initial.txt")
	if err := os.WriteFile(testFile, []byte("initial content"), 0644); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create initial file: %v", err)
	}

	cmd = exec.Command("git", "add", "initial.txt")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to add initial file: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	return tempDir
}

func TestHasStagedChanges(t *testing.T) {
	tempDir := setupTestRepo(t)
	defer os.RemoveAll(tempDir)

	// Change to test directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	client := &Client{}

	// Test with no staged changes
	hasStaged := client.HasStagedChanges()
	if hasStaged {
		t.Error("Expected no staged changes, got true")
	}

	// Create a test file and stage it
	testFile := filepath.Join(tempDir, "test-file.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd := exec.Command("git", "add", "test-file.txt")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage test file: %v", err)
	}

	// Test with staged changes
	client2 := &Client{}
	hasStaged2 := client2.HasStagedChanges()
	if !hasStaged2 {
		t.Error("Expected staged changes, got false")
	}
}

func TestStageAllChanges(t *testing.T) {
	tempDir := setupTestRepo(t)
	defer os.RemoveAll(tempDir)

	// Change to test directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	client := &Client{}

	// Create untracked file
	testFile := filepath.Join(tempDir, "untracked.txt")
	if err := os.WriteFile(testFile, []byte("untracked content"), 0644); err != nil {
		t.Fatalf("Failed to create untracked file: %v", err)
	}

	// Create unstaged changes
	initialFile := filepath.Join(tempDir, "initial.txt")
	if err := os.WriteFile(initialFile, []byte("modified content"), 0644); err != nil {
		t.Fatalf("Failed to modify initial file: %v", err)
	}

	// Verify changes are not staged
	if client.HasStagedChanges() {
		t.Error("Expected no staged changes before StageAllChanges")
	}

	// Stage all changes
	if err := client.StageAllChanges(); err != nil {
		t.Fatalf("StageAllChanges failed: %v", err)
	}

	// Verify changes are now staged
	if !client.HasStagedChanges() {
		t.Error("Expected staged changes after StageAllChanges")
	}
}

func TestHasAnyChanges(t *testing.T) {
	tempDir := setupTestRepo(t)
	defer os.RemoveAll(tempDir)

	// Change to test directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	client := &Client{}

	// Test with no changes
	if client.HasAnyChanges() {
		t.Error("Expected no changes, got true")
	}

	// Test with unstaged changes
	initialFile := filepath.Join(tempDir, "initial.txt")
	if err := os.WriteFile(initialFile, []byte("modified content"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	if !client.HasAnyChanges() {
		t.Error("Expected changes with unstaged modifications, got false")
	}

	// Stage the changes
	cmd := exec.Command("git", "add", "initial.txt")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage changes: %v", err)
	}

	if !client.HasAnyChanges() {
		t.Error("Expected changes with staged modifications, got false")
	}

	// Commit staged changes
	cmd = exec.Command("git", "commit", "-m", "Test commit")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	if client.HasAnyChanges() {
		t.Error("Expected no changes after commit, got true")
	}

	// Test with untracked files
	testFile := filepath.Join(tempDir, "untracked.txt")
	if err := os.WriteFile(testFile, []byte("untracked content"), 0644); err != nil {
		t.Fatalf("Failed to create untracked file: %v", err)
	}

	if !client.HasAnyChanges() {
		t.Error("Expected changes with untracked files, got false")
	}

	// Clean up untracked file for next test
	if err := os.Remove(testFile); err != nil {
		t.Fatalf("Failed to remove untracked file: %v", err)
	}

	// Test edge case: file in .gitignore
	gitignoreFile := filepath.Join(tempDir, ".gitignore")
	if err := os.WriteFile(gitignoreFile, []byte("*.log\n"), 0644); err != nil {
		t.Fatalf("Failed to create .gitignore: %v", err)
	}

	logFile := filepath.Join(tempDir, "test.log")
	if err := os.WriteFile(logFile, []byte("log content"), 0644); err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}

	// Untracked files matching .gitignore should not count as changes
	// Only the .gitignore file itself should count
	if !client.HasAnyChanges() {
		t.Error("Expected changes with untracked .gitignore file, got false")
	}
}

func TestGetStagedDiff(t *testing.T) {
	tempDir := setupTestRepo(t)
	defer os.RemoveAll(tempDir)

	// Change to test directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	client := &Client{}

	// Create a test file with content
	testContent := "Hello, World!"
	testFile := filepath.Join(tempDir, "diff-test.txt")
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Stage the file
	cmd := exec.Command("git", "add", "diff-test.txt")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage test file: %v", err)
	}

	// Get diff
	diff, err := client.GetStagedDiff()
	if err != nil {
		t.Fatalf("GetStagedDiff failed: %v", err)
	}

	// Verify diff contains our content
	if !strings.Contains(diff, "Hello, World!") {
		t.Errorf("Diff should contain test content, got: %s", diff)
	}

	// Verify diff is not empty
	if strings.TrimSpace(diff) == "" {
		t.Error("Diff should not be empty when there are staged changes")
	}
}

func TestCommit(t *testing.T) {
	tempDir := setupTestRepo(t)
	defer os.RemoveAll(tempDir)

	// Change to test directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	client := &Client{}

	// Create and stage a test file
	testFile := filepath.Join(tempDir, "commit-test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd := exec.Command("git", "add", "commit-test.txt")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage test file: %v", err)
	}

	// Commit with a test message
	testMessage := "test: automated commit message"
	result, err := client.Commit(testMessage, false)
	if err != nil {
		t.Fatalf("Commit failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected success, got failure: %v", result.Error)
	}

	if result.Message != testMessage {
		t.Errorf("Message mismatch: expected %s, got %s", testMessage, result.Message)
	}

	if result.Hash == "" {
		t.Error("Commit hash should not be empty on success")
	}

	// Clean up
	exec.Command("git", "reset", "--hard", "HEAD~1").Run()
}

func TestFallbackToInteractiveDryRun(t *testing.T) {
	tempDir := setupTestRepo(t)
	defer os.RemoveAll(tempDir)

	// Change to test directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	client := &Client{}

	// Create and stage a test file
	testFile := filepath.Join(tempDir, "interactive-test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd := exec.Command("git", "add", "interactive-test.txt")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage test file: %v", err)
	}

	// Create a simple script that acts as a non-interactive editor
	scriptPath := filepath.Join(tempDir, "fake-editor.sh")
	scriptContent := "#!/bin/sh\nexit 0\n"
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("Failed to create fake editor script: %v", err)
	}

	// Set EDITOR to our fake editor
	oldEditor := os.Getenv("EDITOR")
	os.Setenv("EDITOR", scriptPath)
	defer os.Setenv("EDITOR", oldEditor)

	// Test dry run mode - should not create actual commit
	testMessage := "test: dry run message"
	err = client.FallbackToInteractive(testMessage, true)
	if err != nil {
		t.Fatalf("FallbackToInteractive dry run failed: %v", err)
	}

	// Verify no commit was created
	cmd2 := exec.Command("git", "log", "--oneline")
	output, err := cmd2.Output()
	if err != nil {
		t.Fatalf("Failed to get git log: %v", err)
	}

	logLines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(logLines) != 1 {
		t.Errorf("Expected 1 commit (initial only), got %d: %s", len(logLines), string(output))
	}
}

func TestValidateEditor(t *testing.T) {
	client := &Client{}

	// Test with EDITOR set to a common editor
	commonEditors := []string{"vi", "vim", "nano"}
	foundEditor := false
	for _, editor := range commonEditors {
		if _, err := exec.LookPath(editor); err == nil {
			foundEditor = true
			// Set EDITOR environment variable
			oldEditor := os.Getenv("EDITOR")
			os.Setenv("EDITOR", editor)
			defer os.Setenv("EDITOR", oldEditor)

			validatedEditor, err := client.validateEditor()
			if err != nil {
				t.Errorf("validateEditor failed for %s: %v", editor, err)
			}
			if validatedEditor != editor {
				t.Errorf("Expected editor %s, got %s", editor, validatedEditor)
			}
			break
		}
	}

	if !foundEditor {
		t.Skip("No common editor found for testing")
	}

	// Test with EDITOR set to non-existent editor
	os.Setenv("EDITOR", "nonexistent-editor-xyz")
	_, err := client.validateEditor()
	if err == nil {
		t.Error("Expected error for non-existent editor, got nil")
	}

	// Test with EDITOR unset and no common editors available
	// This is hard to test reliably as most systems have at least one editor
	// So we'll just verify the function doesn't crash
	os.Unsetenv("EDITOR")
	editor, err := client.validateEditor()
	// On most systems, this will find a common editor and succeed
	// On systems without any common editor, it should return an error
	if err != nil && editor == "" {
		// This is acceptable - no editor found
		t.Logf("No editor found (acceptable on minimal systems): %v", err)
	} else if err == nil && editor != "" {
		// This is also acceptable - found a common editor
		t.Logf("Found common editor: %s", editor)
	}
}

func TestCommitDryRun(t *testing.T) {
	tempDir := setupTestRepo(t)
	defer os.RemoveAll(tempDir)

	// Change to test directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}

	client := &Client{}

	// Create and stage a test file
	testFile := filepath.Join(tempDir, "dryrun-test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd := exec.Command("git", "add", "dryrun-test.txt")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage test file: %v", err)
	}

	// Test dry run
	testMessage := "test: dry run message"
	result, err := client.Commit(testMessage, true)
	if err != nil {
		t.Fatalf("Commit dry run failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected success in dry run mode, got failure: %v", result.Error)
	}

	if result.Message != testMessage {
		t.Errorf("Message mismatch: expected %s, got %s", testMessage, result.Message)
	}

	if result.Hash != "" {
		t.Error("Commit hash should be empty in dry run mode")
	}

	// Verify no commit was actually created
	cmd2 := exec.Command("git", "log", "--oneline")
	output, err := cmd2.Output()
	if err != nil {
		t.Fatalf("Failed to get git log: %v", err)
	}

	logLines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(logLines) != 1 {
		t.Errorf("Expected 1 commit (initial only) after dry run, got %d", len(logLines))
	}
}
