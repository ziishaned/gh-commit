package git

import (
	"os"
	"os/exec"
	"path/filepath"
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
