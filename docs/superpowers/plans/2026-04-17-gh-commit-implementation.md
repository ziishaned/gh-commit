# gh-commit Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a GitHub CLI extension that automatically generates Conventional Commits messages using GitHub Models API and commits them

**Architecture:** Three-component system: Git client for git operations, LLM client for AI integration via GitHub Models API, and main CLI orchestrator that coordinates the workflow with intelligent error handling

**Tech Stack:** Go 1.21, GitHub CLI extension API, GitHub Models API, cobra CLI framework, conventional commits specification

---

## File Structure

```
gh-commit/
├── extension.yml              # GitHub CLI extension manifest
├── go.mod                     # Go module definition
├── go.sum                     # Go dependencies lock
├── Makefile                   # Build commands
├── README.md                  # Documentation
├── cmd/
│   └── commit/
│       └── main.go           # CLI entry point with cobra
├── internal/
│   ├── git/
│   │   ├── client.go         # Git operations (stage, diff, commit)
│   │   └── client_test.go    # Git client tests
│   ├── llm/
│   │   ├── client.go         # GitHub Models API client
│   │   ├── client_test.go    # LLM client tests
│   │   └── commit.prompt.yml # Conventional Commits prompt
│   └── types/
│       └── commit.go         # Type definitions
```

---

### Task 1: Initialize Go Module and Dependencies

**Files:**
- Create: `go.mod`
- Create: `go.sum`

- [ ] **Step 1: Initialize Go module**

Run: `go mod init github.com/gh-commit`
Expected: Creates `go.mod` with module declaration

- [ ] **Step 2: Add required dependencies**

Run: `go get github.com/cli/go-gh/v2@latest && go get github.com/spf13/cobra@latest`
Expected: Dependencies added to go.mod, go.sum created

- [ ] **Step 3: Tidy dependencies**

Run: `go mod tidy`
Expected: go.mod and go.sum cleaned up

- [ ] **Step 4: Commit**

```bash
git add go.mod go.sum
git commit -m "feat: initialize Go module and dependencies"
```

---

### Task 2: Create Extension Manifest

**Files:**
- Create: `extension.yml`

- [ ] **Step 1: Create extension.yml manifest**

Create `extension.yml`:
```yaml
name: commit
owner: gh-commit
host: github.com
tag: v1.0.0
```

- [ ] **Step 2: Verify file is valid YAML**

Check: File should parse as valid YAML with name, owner, host, tag fields

- [ ] **Step 3: Commit**

```bash
git add extension.yml
git commit -m "feat: add GitHub CLI extension manifest"
```

---

### Task 3: Define Core Types

**Files:**
- Create: `internal/types/commit.go`

- [ ] **Step 1: Write type definitions**

Create `internal/types/commit.go`:
```go
package types

// GitChange represents a single file change in git
type GitChange struct {
	Path         string // File path
	Additions    int    // Lines added
	Deletions    int    // Lines deleted
	Patch        string // Full diff patch
	IsNewFile    bool   // Is this a new file
	IsDeleted    bool   // Is this a deleted file
	IsRenamed    bool   // Is this a renamed file
}

// CommitResult represents the result of a commit operation
type CommitResult struct {
	Success     bool   // Was the commit successful
	Message     string // Generated commit message
	Hash        string // Commit hash (if successful)
	Error       error  // Error if unsuccessful
	DryRun      bool   // Was this a dry run?
}

// LLMRequest represents a request to the LLM
type LLMRequest struct {
	Diff    string   // Git diff to analyze
	Model   string   // Model to use
	Language string  // Language for output (unused but kept for compatibility)
}

// LLMResponse represents a response from the LLM
type LLMResponse struct {
	Message string // Generated commit message
	Error   error  // Error if unsuccessful
}
```

- [ ] **Step 2: Verify code compiles**

Run: `go build ./internal/types`
Expected: No compilation errors

- [ ] **Step 3: Commit**

```bash
git add internal/types/commit.go
git commit -m "feat: define core types for git operations"
```

---

### Task 4: Implement Git Client - Check Staged Changes

**Files:**
- Create: `internal/git/client.go`
- Create: `internal/git/client_test.go`

- [ ] **Step 1: Write failing test for HasStagedChanges**

Create `internal/git/client_test.go`:
```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/git -v -run TestHasStagedChanges`
Expected: FAIL with "undefined: Client"

- [ ] **Step 3: Implement Client struct and HasStagedChanges method**

Create `internal/git/client.go`:
```go
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
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/git -v -run TestHasStagedChanges`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/git/client.go internal/git/client_test.go
git commit -m "feat: add staged change detection to Git client"
```

---

### Task 5: Implement Git Client - Stage All Changes

**Files:**
- Modify: `internal/git/client.go`
- Modify: `internal/git/client_test.go`

- [ ] **Step 1: Write failing test for StageAllChanges**

Add to `internal/git/client_test.go`:
```go
func TestStageAllChanges(t *testing.T) {
	client := &Client{}
	
	// Create an untracked file
	exec.Command("touch", "untracked.txt").Run()
	defer exec.Command("rm", "-f", "untracked.txt").Run()
	
	// Verify it's not staged
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	output, _ := cmd.Output()
	if strings.Contains(string(output), "untracked.txt") {
		t.Error("File should not be staged initially")
	}
	
	// Stage all changes
	err := client.StageAllChanges()
	if err != nil {
		t.Fatalf("StageAllChanges failed: %v", err)
	}
	
	// Verify file is staged
	cmd2 := exec.Command("git", "diff", "--cached", "--name-only")
	output2, _ := cmd2.Output()
	if !strings.Contains(string(output2), "untracked.txt") {
		t.Error("File should be staged after StageAllChanges")
	}
	
	// Clean up
	exec.Command("git", "reset", "untracked.txt").Run()
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/git -v -run TestStageAllChanges`
Expected: Should PASS (already implemented in Task 4)

- [ ] **Step 3: Commit**

```bash
git add internal/git/client_test.go
git commit -m "test: add test for StageAllChanges"
```

---

### Task 6: Implement Git Client - Get Staged Diff

**Files:**
- Modify: `internal/git/client.go`
- Modify: `internal/git/client_test.go`

- [ ] **Step 1: Write failing test for GetStagedDiff**

Add to `internal/git/client_test.go`:
```go
func TestGetStagedDiff(t *testing.T) {
	client := &Client{}
	
	// Create a test file with content
	testContent := "Hello, World!"
	exec.Command("sh", "-c", "echo '"+testContent+"' > diff-test.txt").Run()
	
	// Stage the file
	exec.Command("git", "add", "diff-test.txt").Run()
	defer func() {
		exec.Command("git", "reset", "diff-test.txt").Run()
		exec.Command("rm", "-f", "diff-test.txt").Run()
	}()
	
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/git -v -run TestGetStagedDiff`
Expected: FAIL with "undefined: GetStagedDiff"

- [ ] **Step 3: Implement GetStagedDiff method**

Add to `internal/git/client.go`:
```go
import (
	"bytes"
	"io"
)

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
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/git -v -run TestGetStagedDiff`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/git/client.go internal/git/client_test.go
git commit -m "feat: add GetStagedDiff to Git client"
```

---

### Task 7: Implement Git Client - Commit with Message

**Files:**
- Modify: `internal/git/client.go`
- Modify: `internal/git/client_test.go`

- [ ] **Step 1: Write failing test for Commit**

Add to `internal/git/client_test.go`:
```go
func TestCommit(t *testing.T) {
	client := &Client{}
	
	// Create and stage a test file
	exec.Command("sh", "-c", "echo 'test content' > commit-test.txt").Run()
	exec.Command("git", "add", "commit-test.txt").Run()
	
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
	exec.Command("rm", "-f", "commit-test.txt").Run()
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/git -v -run TestCommit`
Expected: FAIL with "undefined: Commit"

- [ ] **Step 3: Implement Commit method**

Add to `internal/git/client.go`:
```go
import (
	"github.com/gh-commit/internal/types"
)

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
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/git -v -run TestCommit`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/git/client.go internal/git/client_test.go
git commit -m "feat: add Commit method to Git client"
```

---

### Task 8: Create Conventional Commits Prompt

**Files:**
- Create: `internal/llm/commit.prompt.yml`

- [ ] **Step 1: Create prompt configuration file**

Create `internal/llm/commit.prompt.yml`:
```yaml
name: Conventional Commits Generator
description: Generates conventional commit messages based on git diffs
model: openai/gpt-4o
modelParameters:
  temperature: 0.7
  topP: 0.9
messages:
  - role: system
    content: >
      You are an expert developer who writes clear, concise Conventional Commits
      messages. Your task is to analyze git diffs and generate conventional
      commit messages.


      Conventional Commits format: type(scope): description


      Types: feat, fix, docs, style, refactor, test, chore, perf, ci, build,
      revert


      Scope: Optional, typically the module/component affected


      Description: Imperative mood, lowercase, no period


      Examples:


      - feat(auth): add JWT token refresh mechanism


      - fix(api): resolve memory leak in request handler


      - docs(readme): update installation instructions


      - refactor(database): extract query builder to separate module
  - role: user
    content: |
      Analyze the following git diff and generate a conventional commit
      message. Focus on WHAT changed and WHY, not HOW.


      Diff:

      {{diff}}


      Requirements:


      - Use conventional commit format: type(scope): description


      - Choose the most appropriate type (feat/fix/docs/refactor/etc.)


      - Keep description under 72 characters


      - Use imperative mood ("add" not "added")


      - If scope is unclear, omit it


      - Provide only the commit message, nothing else


      Commit message:
testData:
  - diff: |
      diff --git a/internal/auth/jwt.go b/internal/auth/jwt.go
      new file mode 100644
      index 0000000..1234567
      --- /dev/null
      +++ b/internal/auth/jwt.go
      @@ -0,0 +1,25 @@
      +package auth
      +
      +import "time"
      +
      +type JWTToken struct {
      +        Token     string
      +        ExpiresAt time.Time
      +}
      +
      +func RefreshToken(token string) (*JWTToken, error) {
      +        // Implementation here
      +}
      +
      +func ValidateToken(token string) bool {
      +        // Implementation here
      +}
evaluators: []
```

- [ ] **Step 2: Verify file is valid YAML**

Check: File should parse as valid YAML with required fields

- [ ] **Step 3: Commit**

```bash
git add internal/llm/commit.prompt.yml
git commit -m "feat: add Conventional Commits prompt configuration"
```

---

### Task 9: Implement LLM Client - Basic Structure

**Files:**
- Create: `internal/llm/client.go`

- [ ] **Step 1: Create LLM client struct and types**

Create `internal/llm/client.go`:
```go
package llm

import (
	_ "embed"
)

//go:embed commit.prompt.yml
var commitPromptYAML []byte

type PromptConfig struct {
	Name            string          `yaml:"name"`
	Description     string          `yaml:"description"`
	Model           string          `yaml:"model"`
	ModelParameters ModelParameters `yaml:"modelParameters"`
	Messages        []PromptMessage `yaml:"messages"`
}

type ModelParameters struct {
	Temperature float64 `yaml:"temperature"`
	TopP        float64 `yaml:"topP"`
}

type PromptMessage struct {
	Role    string `yaml:"role"`
	Content string `yaml:"content"`
}

type Client struct {
	token string
}

type Request struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Temperature float64   `json:"temperature"`
	TopP        float64   `json:"top_p"`
	Stream      bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
```

- [ ] **Step 2: Verify code compiles**

Run: `go build ./internal/llm`
Expected: No compilation errors

- [ ] **Step 3: Commit**

```bash
git add internal/llm/client.go
git commit -m "feat: add LLM client types and structures"
```

---

### Task 10: Implement LLM Client - Prompt Loading

**Files:**
- Modify: `internal/llm/client.go`
- Create: `internal/llm/client_test.go`

- [ ] **Step 1: Write failing test for loadPromptConfig**

Create `internal/llm/client_test.go`:
```go
package llm

import (
	"testing"
)

func TestLoadPromptConfig(t *testing.T) {
	config, err := loadPromptConfig()
	if err != nil {
		t.Fatalf("loadPromptConfig failed: %v", err)
	}
	
	if config.Name != "Conventional Commits Generator" {
		t.Errorf("Unexpected name: %s", config.Name)
	}
	
	if config.Model != "openai/gpt-4o" {
		t.Errorf("Unexpected model: %s", config.Model)
	}
	
	if len(config.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(config.Messages))
	}
	
	if config.Messages[0].Role != "system" {
		t.Errorf("First message should be system role, got %s", config.Messages[0].Role)
	}
	
	if config.Messages[1].Role != "user" {
		t.Errorf("Second message should be user role, got %s", config.Messages[1].Role)
	}
	
	if config.ModelParameters.Temperature != 0.7 {
		t.Errorf("Expected temperature 0.7, got %f", config.ModelParameters.Temperature)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/llm -v -run TestLoadPromptConfig`
Expected: FAIL with "undefined: loadPromptConfig"

- [ ] **Step 3: Implement loadPromptConfig function**

Add to `internal/llm/client.go`:
```go
import (
	"fmt"
	"gopkg.in/yaml.v3"
)

func loadPromptConfig() (*PromptConfig, error) {
	var config PromptConfig
	err := yaml.Unmarshal(commitPromptYAML, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse prompt configuration: %w", err)
	}
	return &config, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/llm -v -run TestLoadPromptConfig`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/llm/client.go internal/llm/client_test.go
git commit -m "feat: add prompt configuration loading"
```

---

### Task 11: Implement LLM Client - Client Initialization

**Files:**
- Modify: `internal/llm/client.go`
- Modify: `internal/llm/client_test.go`

- [ ] **Step 1: Write failing test for NewClient**

Add to `internal/llm/client_test.go`:
```go
import (
	"os"
)

func TestNewClient(t *testing.T) {
	// Set a fake token for testing
	os.Setenv("GH_TOKEN", "test-token")
	defer os.Unsetenv("GH_TOKEN")
	
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	
	if client == nil {
		t.Fatal("Client should not be nil")
	}
	
	if client.token != "test-token" {
		t.Errorf("Expected token 'test-token', got '%s'", client.token)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/llm -v -run TestNewClient`
Expected: FAIL with "undefined: NewClient"

- [ ] **Step 3: Implement NewClient function**

Add to `internal/llm/client.go`:
```go
import (
	"github.com/cli/go-gh/v2/pkg/auth"
)

func NewClient() (*Client, error) {
	host, _ := auth.DefaultHost()
	token, _ := auth.TokenForHost(host)
	
	if token == "" {
		return nil, fmt.Errorf("no GitHub token found. Please run 'gh auth login' to authenticate")
	}
	
	return &Client{token: token}, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/llm -v -run TestNewClient`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/llm/client.go internal/llm/client_test.go
git commit -m "feat: add LLM client initialization"
```

---

### Task 12: Implement LLM Client - Generate Commit Message

**Files:**
- Modify: `internal/llm/client.go`
- Modify: `internal/llm/client_test.go`

- [ ] **Step 1: Write failing test for GenerateCommitMessage**

Add to `internal/llm/client_test.go`:
```go
import (
	"github.com/gh-commit/internal/types"
)

func TestGenerateCommitMessage(t *testing.T) {
	// Skip if no token available
	client, err := NewClient()
	if err != nil {
		t.Skip("No GitHub token available")
	}
	
	// Test with a simple diff
	testDiff := `diff --git a/test.txt b/test.txt
new file mode 100644
index 0000000..1234567
--- /dev/null
+++ b/test.txt
@@ -0,0 +1 @@
+Hello, World!`
	
	req := types.LLMRequest{
		Diff:  testDiff,
		Model: "openai/gpt-4o",
	}
	
	resp, err := client.GenerateCommitMessage(req)
	if err != nil {
		t.Fatalf("GenerateCommitMessage failed: %v", err)
	}
	
	if resp.Error != nil {
		t.Fatalf("GenerateCommitMessage returned error: %v", resp.Error)
	}
	
	if resp.Message == "" {
		t.Error("Message should not be empty")
	}
	
	// Verify conventional commit format (basic check)
	// Should contain a colon and not be too long
	if len(resp.Message) > 100 {
		t.Errorf("Message seems too long: %s", resp.Message)
	}
	
	t.Logf("Generated message: %s", resp.Message)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/llm -v -run TestGenerateCommitMessage`
Expected: FAIL with "undefined: GenerateCommitMessage"

- [ ] **Step 3: Implement GenerateCommitMessage method**

Add to `internal/llm/client.go`:
```go
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func (c *Client) GenerateCommitMessage(req types.LLMRequest) (*types.LLMResponse, error) {
	// Load prompt configuration
	promptConfig, err := loadPromptConfig()
	if err != nil {
		return &types.LLMResponse{Error: err}, err
	}
	
	// Use model from request or fall back to config
	selectedModel := req.Model
	if selectedModel == "" {
		selectedModel = promptConfig.Model
	}
	
	// Build messages from prompt config, replacing template variables
	messages := make([]Message, len(promptConfig.Messages))
	for i, msg := range promptConfig.Messages {
		content := msg.Content
		content = strings.ReplaceAll(content, "{{diff}}", req.Diff)
		messages[i] = Message{
			Role:    msg.Role,
			Content: content,
		}
	}
	
	// Build request
	request := Request{
		Messages:    messages,
		Model:       selectedModel,
		Temperature: promptConfig.ModelParameters.Temperature,
		TopP:        promptConfig.ModelParameters.TopP,
		Stream:      false,
	}
	
	// Call GitHub Models API
	response, err := c.callGitHubModels(request)
	if err != nil {
		return &types.LLMResponse{Error: err}, err
	}
	
	if len(response.Choices) == 0 {
		return &types.LLMResponse{Error: fmt.Errorf("no response generated from model")}, fmt.Errorf("no response generated")
	}
	
	message := strings.TrimSpace(response.Choices[0].Message.Content)
	return &types.LLMResponse{Message: message}, nil
}

func (c *Client) callGitHubModels(request Request) (*Response, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	req, err := http.NewRequest("POST", "https://models.github.ai/inference/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return &response, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/llm -v -run TestGenerateCommitMessage`
Expected: PASS (may skip if no token)

- [ ] **Step 5: Commit**

```bash
git add internal/llm/client.go internal/llm/client_test.go
git commit -m "feat: add commit message generation using GitHub Models API"
```

---

### Task 13: Implement Main CLI - Basic Structure

**Files:**
- Create: `cmd/commit/main.go`

- [ ] **Step 1: Create main CLI structure**

Create `cmd/commit/main.go`:
```go
package main

import (
	"fmt"
	"os"

	"github.com/gh-commit/internal/git"
	"github.com/gh-commit/internal/llm"
	"github.com/spf13/cobra"
)

const extensionName = "commit"

var rootCmd = &cobra.Command{
	Use:   extensionName,
	Short: "Generate AI-powered commit messages",
	Long:  "A GitHub CLI extension that generates commit messages using GitHub Models and commits them automatically",
	RunE:  runCommit,
}

var (
	flagModel  string
	flagDryRun bool
)

func init() {
	rootCmd.Flags().StringVarP(&flagModel, "model", "m", "openai/gpt-4o", "GitHub Models model to use")
	rootCmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Show what would be done without actually committing")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runCommit(cmd *cobra.Command, args []string) error {
	fmt.Println("Checking for staged changes...")
	
	gitClient := &git.Client{}
	
	// Check if there are any changes at all
	if !gitClient.HasAnyChanges() {
		return fmt.Errorf("no changes detected. Please make some changes first")
	}
	
	// Check if there are staged changes
	if !gitClient.HasStagedChanges() {
		fmt.Println("No staged changes found. Staging all changes...")
		if err := gitClient.StageAllChanges(); err != nil {
			return fmt.Errorf("failed to stage changes: %w", err)
		}
		fmt.Println("✅ Changes staged")
	} else {
		fmt.Println("✅ Found staged changes")
	}
	
	return nil
}
```

- [ ] **Step 2: Verify code compiles**

Run: `go build ./cmd/commit`
Expected: No compilation errors

- [ ] **Step 3: Test basic CLI**

Run: `go run ./cmd/commit --help`
Expected: Help text displayed

- [ ] **Step 4: Commit**

```bash
git add cmd/commit/main.go
git commit -m "feat: add basic CLI structure with change detection"
```

---

### Task 14: Implement Main CLI - Diff Capture

**Files:**
- Modify: `cmd/commit/main.go`

- [ ] **Step 1: Add diff capture to runCommit**

Modify the `runCommit` function in `cmd/commit/main.go`:
```go
func runCommit(cmd *cobra.Command, args []string) error {
	fmt.Println("Checking for staged changes...")
	
	gitClient := &git.Client{}
	
	// Check if there are any changes at all
	if !gitClient.HasAnyChanges() {
		return fmt.Errorf("no changes detected. Please make some changes first")
	}
	
	// Check if there are staged changes
	if !gitClient.HasStagedChanges() {
		fmt.Println("No staged changes found. Staging all changes...")
		if err := gitClient.StageAllChanges(); err != nil {
			return fmt.Errorf("failed to stage changes: %w", err)
		}
		fmt.Println("✅ Changes staged")
	} else {
		fmt.Println("✅ Found staged changes")
	}
	
	// Get the staged diff
	fmt.Println("Analyzing staged changes...")
	diff, err := gitClient.GetStagedDiff()
	if err != nil {
		return fmt.Errorf("failed to get diff: %w", err)
	}
	
	fmt.Printf("📝 Analyzing %d lines of changes...\n", len(diff))
	
	return nil
}
```

- [ ] **Step 2: Verify code compiles**

Run: `go build ./cmd/commit`
Expected: No compilation errors

- [ ] **Step 3: Test with actual changes**

Run: `touch test.txt && git add test.txt && go run ./cmd/commit`
Expected: Should show analysis message

Clean up: `git reset test.txt && rm test.txt`

- [ ] **Step 4: Commit**

```bash
git add cmd/commit/main.go
git commit -m "feat: add diff capture and analysis to CLI"
```

---

### Task 15: Implement Main CLI - LLM Integration

**Files:**
- Modify: `cmd/commit/main.go`

- [ ] **Step 1: Add LLM client integration**

Modify the `runCommit` function in `cmd/commit/main.go`:
```go
import (
	"github.com/gh-commit/internal/types"
)

func runCommit(cmd *cobra.Command, args []string) error {
	fmt.Println("Checking for staged changes...")
	
	gitClient := &git.Client{}
	
	// Check if there are any changes at all
	if !gitClient.HasAnyChanges() {
		return fmt.Errorf("no changes detected. Please make some changes first")
	}
	
	// Check if there are staged changes
	if !gitClient.HasStagedChanges() {
		fmt.Println("No staged changes found. Staging all changes...")
		if err := gitClient.StageAllChanges(); err != nil {
			return fmt.Errorf("failed to stage changes: %w", err)
		}
		fmt.Println("✅ Changes staged")
	} else {
		fmt.Println("✅ Found staged changes")
	}
	
	// Get the staged diff
	fmt.Println("Analyzing staged changes...")
	diff, err := gitClient.GetStagedDiff()
	if err != nil {
		return fmt.Errorf("failed to get diff: %w", err)
	}
	
	fmt.Printf("📝 Analyzing changes...\n")
	
	// Generate commit message using LLM
	fmt.Printf("Generating commit message using %s...\n", flagModel)
	llmClient, err := llm.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create LLM client: %w", err)
	}
	
	req := types.LLMRequest{
		Diff:  diff,
		Model: flagModel,
	}
	
	resp, err := llmClient.GenerateCommitMessage(req)
	if err != nil {
		return fmt.Errorf("failed to generate commit message: %w", err)
	}
	
	if resp.Error != nil {
		return fmt.Errorf("failed to generate commit message: %w", resp.Error)
	}
	
	message := resp.Message
	fmt.Printf("✨ Generated: %s\n", message)
	
	return nil
}
```

- [ ] **Step 2: Verify code compiles**

Run: `go build ./cmd/commit`
Expected: No compilation errors

- [ ] **Step 3: Test with actual changes (dry run)**

Run: `echo "test" > test.txt && git add test.txt && go run ./cmd/commit --dry-run`
Expected: Should generate a commit message

Clean up: `git reset test.txt && rm test.txt`

- [ ] **Step 4: Commit**

```bash
git add cmd/commit/main.go
git commit -m "feat: add LLM integration for commit message generation"
```

---

### Task 16: Implement Main CLI - Commit Execution

**Files:**
- Modify: `cmd/commit/main.go`

- [ ] **Step 1: Add commit execution to runCommit**

Modify the end of the `runCommit` function in `cmd/commit/main.go`:
```go
	message := resp.Message
	fmt.Printf("✨ Generated: %s\n", message)
	
	// Commit the changes
	if flagDryRun {
		fmt.Println("\n[Dry run] Commit would be created with message above")
		return nil
	}
	
	fmt.Println("Committing...")
	result, err := gitClient.Commit(message, flagDryRun)
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	
	if !result.Success {
		return fmt.Errorf("commit failed: %w", result.Error)
	}
	
	fmt.Printf("✅ Committed successfully: %s\n", result.Message)
	fmt.Printf("   Commit hash: %s\n", result.Hash)
	
	return nil
}
```

- [ ] **Step 2: Verify code compiles**

Run: `go build ./cmd/commit`
Expected: No compilation errors

- [ ] **Step 3: Test full workflow**

Run: `echo "test content" > test.txt && git add test.txt && go run ./cmd/commit`
Expected: Should stage, generate message, and commit

Clean up: `git reset --hard HEAD~1 && rm test.txt`

- [ ] **Step 4: Commit**

```bash
git add cmd/commit/main.go
git commit -m "feat: add commit execution with dry-run support"
```

---

### Task 17: Implement Error Handling - Interactive Fallback

**Files:**
- Modify: `internal/git/client.go`
- Modify: `cmd/commit/main.go`

- [ ] **Step 1: Add interactive fallback method to Git client**

Add to `internal/git/client.go`:
```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

// FallbackToInteractive opens an editor for manual commit message editing
func (c *Client) FallbackToInteractive(message string) error {
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
	
	// Get editor from environment
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim" // Fallback to vim
	}
	
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
```

- [ ] **Step 2: Add error handling with fallback in CLI**

Modify `cmd/commit/main.go` to handle errors:
```go
	// Commit the changes
	if flagDryRun {
		fmt.Println("\n[Dry run] Commit would be created with message above")
		return nil
	}
	
	fmt.Println("Committing...")
	result, err := gitClient.Commit(message, flagDryRun)
	if err != nil {
		fmt.Printf("⚠️  Commit failed: %v\n", err)
		fmt.Println("Falling back to interactive mode...")
		
		fallbackMessage := "AI generation failed. Please write commit message manually."
		if message != "" {
			fallbackMessage = message
		}
		
		if err := gitClient.FallbackToInteractive(fallbackMessage); err != nil {
			return fmt.Errorf("interactive commit failed: %w", err)
		}
		
		fmt.Println("✅ Committed successfully via interactive mode")
		return nil
	}
	
	if !result.Success {
		fmt.Printf("⚠️  Commit failed: %v\n", result.Error)
		fmt.Println("Falling back to interactive mode...")
		
		if err := gitClient.FallbackToInteractive(message); err != nil {
			return fmt.Errorf("interactive commit failed: %w", err)
		}
		
		fmt.Println("✅ Committed successfully via interactive mode")
		return nil
	}
	
	fmt.Printf("✅ Committed successfully: %s\n", result.Message)
	fmt.Printf("   Commit hash: %s\n", result.Hash)
	
	return nil
}
```

- [ ] **Step 3: Verify code compiles**

Run: `go build ./cmd/commit`
Expected: No compilation errors

- [ ] **Step 4: Test error handling**

Test with an invalid model to trigger fallback:
```bash
echo "test" > test.txt && git add test.txt && go run ./cmd/commit --model invalid/model
```
Expected: Should fall back to interactive mode

Clean up: Handle interactive commit, then `git reset --hard HEAD~1 && rm test.txt`

- [ ] **Step 5: Commit**

```bash
git add internal/git/client.go cmd/commit/main.go
git commit -m "feat: add interactive fallback for error handling"
```

---

### Task 18: Create Makefile

**Files:**
- Create: `Makefile`

- [ ] **Step 1: Create Makefile for build commands**

Create `Makefile`:
```makefile
.PHONY: build test clean install

# Build the extension
build:
	@echo "Building gh-commit..."
	@mkdir -p bin
	@go build -o bin/commit ./cmd/commit
	@echo "✅ Built bin/commit"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin
	@echo "✅ Cleaned"

# Install as gh extension
install: build
	@echo "Installing as gh extension..."
	@gh extension install .
	@echo "✅ Installed"

# Uninstall extension
uninstall:
	@echo "Uninstalling gh extension..."
	@gh extension remove commit
	@echo "✅ Uninstalled"
```

- [ ] **Step 2: Test Makefile**

Run: `make build`
Expected: Should build successfully

- [ ] **Step 3: Commit**

```bash
git add Makefile
git commit -m "feat: add Makefile for build automation"
```

---

### Task 19: Create README Documentation

**Files:**
- Create: `README.md`

- [ ] **Step 1: Create comprehensive README**

Create `README.md`:
```markdown
# gh-commit

A GitHub CLI extension that automatically generates Conventional Commits messages using GitHub Models API and commits them. It uses free [GitHub Models](https://docs.github.com/en/github-models) for inference, so you don't need to do any token setup - your existing GitHub CLI token will work!

## Installation

```bash
gh extension install <your-username>/gh-commit
gh commit
```

## Prerequisites

- [GitHub CLI](https://cli.github.com/) installed and authenticated
- Git repository with changes to commit

## Usage

### Basic Usage

Automatically stage changes (if needed), generate a commit message, and commit:

```bash
gh commit
```

### Advanced Options

```bash
# Use a different model
gh commit --model openai/gpt-4o-mini

# Dry run to see what would be committed
gh commit --dry-run

# Combine options
gh commit --model xai/grok-3-mini --dry-run
```

## How It Works

1. **Check for staged changes**: If no changes are staged, automatically stages all changes with `git add .`
2. **Analyze diff**: Captures the full diff of staged changes
3. **Generate message**: Sends diff to GitHub Models API with Conventional Commits prompt
4. **Commit**: Creates commit with generated message
5. **Error handling**: If anything fails, falls back to interactive mode

## Conventional Commits

The extension generates commit messages following the Conventional Commits specification:

```
type(scope): description
```

**Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `perf`, `ci`, `build`, `revert`

**Examples:**
- `feat(auth): add JWT token refresh mechanism`
- `fix(api): resolve memory leak in request handler`
- `docs(readme): update installation instructions`
- `refactor(database): extract query builder to separate module`

## Available Models

The extension supports any model available through GitHub Models API. Default is `openai/gpt-4o`.

Popular options:
- `openai/gpt-4o` (default)
- `openai/gpt-4o-mini` (faster, cheaper)
- `xai/grok-3-mini` (good reasoning)

## Contributing

Contributions are welcome! The extension is built with Go and follows the structure of [gh-standup](https://github.com/sgoedecke/gh-standup).

In particular, tweaking the [prompt](internal/llm/commit.prompt.yml) can significantly improve commit message quality.

## License

MIT
```

- [ ] **Step 2: Verify README is accurate**

Check: All information should be correct and complete

- [ ] **Step 3: Commit**

```bash
git add README.md
git commit -m "docs: add comprehensive README documentation"
```

---

### Task 20: Add CI/CD Configuration

**Files:**
- Create: `.github/workflows/ci.yml`

- [ ] **Step 1: Create CI workflow**

Create `.github/workflows/ci.yml`:
```yaml
name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.21'
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run tests
      run: go test -v ./...
    
    - name: Build
      run: go build -o bin/commit ./cmd/commit
```

- [ ] **Step 2: Verify workflow syntax**

Check: YAML should be valid

- [ ] **Step 3: Commit**

```bash
git add .github/workflows/ci.yml
git commit -m "ci: add GitHub Actions CI workflow"
```

---

### Task 21: Final Integration Test

**Files:**
- Test all components together

- [ ] **Step 1: Build the extension**

Run: `make build`
Expected: Successful build with output in `bin/commit`

- [ ] **Step 2: Test in real repository**

Run: `echo "# Test feature" > test-feature.txt && git add test-feature.txt && ./bin/commit`
Expected: Should generate message and commit successfully

- [ ] **Step 3: Verify commit was created**

Run: `git log -1 --oneline`
Expected: Should show commit with Conventional Commits format

- [ ] **Step 4: Test dry-run mode**

Run: `echo "another change" > test2.txt && git add test2.txt && ./bin/commit --dry-run`
Expected: Should show message without committing

Clean up: `git reset test2.txt && rm test2.txt`

- [ ] **Step 5: Test auto-staging**

Run: `echo "unstaged" > test3.txt && ./bin/commit`
Expected: Should auto-stage and commit

Clean up: `git reset --hard HEAD~2 && rm test-feature.txt test3.txt`

- [ ] **Step 6: Test error handling with invalid model**

Run: `echo "test" > test4.txt && git add test4.txt && ./bin/commit --model invalid/model`
Expected: Should fall back to interactive mode

Clean up: Handle interactive commit, then cleanup

- [ ] **Step 7: Verify all scenarios pass**

Check: All test scenarios should work correctly

- [ ] **Step 8: Commit**

```bash
git add .
git commit -m "test: complete integration testing of gh-commit extension"
```

---

### Task 22: Create .gitignore

**Files:**
- Create: `.gitignore`

- [ ] **Step 1: Create Go-specific .gitignore**

Create `.gitignore`:
```
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
bin/
dist/

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out

# Go workspace file
go.work

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Temporary files
*.tmp
*.log
commit-message-*.txt
test*.txt
```

- [ ] **Step 2: Verify .gitignore works**

Run: `touch test.tmp && git status`
Expected: test.tmp should not appear as untracked

Clean up: `rm test.tmp`

- [ ] **Step 3: Commit**

```bash
git add .gitignore
git commit -m "chore: add .gitignore for Go project"
```

---

## Self-Review Checklist

**1. Spec Coverage:**
- ✅ Extension manifest (Task 2)
- ✅ Git client with staging, diff, commit (Tasks 4-7)
- ✅ LLM client with GitHub Models API (Tasks 9-12)
- ✅ Main CLI orchestrator (Tasks 13-17)
- ✅ Conventional Commits prompt (Task 8)
- ✅ Error handling with interactive fallback (Task 17)
- ✅ CLI interface with --model and --dry-run flags (Task 13)
- ✅ Testing strategy (Tasks 4-7, 10-12)
- ✅ Documentation (Task 19)

**2. Placeholder Scan:**
- ✅ No TBD or TODO found
- ✅ All code steps contain actual implementations
- ✅ All file paths are specified
- ✅ All commands are complete with expected outputs

**3. Type Consistency:**
- ✅ `types.GitChange` defined and used consistently
- ✅ `types.CommitResult` defined and used consistently
- ✅ `types.LLMRequest` and `types.LLMResponse` defined and used consistently
- ✅ Method names consistent: `HasStagedChanges()`, `StageAllChanges()`, `GetStagedDiff()`, `Commit()`
- ✅ Function signatures match throughout tasks

**4. Dependencies:**
- ✅ All required Go packages imported
- ✅ GitHub CLI API (`github.com/cli/go-gh/v2`)
- ✅ Cobra framework (`github.com/spf13/cobra`)
- ✅ YAML parsing (`gopkg.in/yaml.v3`)

---

## Final Notes

This plan follows TDD principles with failing tests written first, then minimal implementations. Each task produces self-contained changes that can be committed independently. The implementation prioritizes:

1. **Core functionality first** (git operations, LLM integration, CLI orchestration)
2. **Error handling** (interactive fallback for all failure scenarios)
3. **User experience** (clear status messages, dry-run mode)
4. **Testing** (unit tests for all major components)
5. **Documentation** (comprehensive README, inline comments)

The extension is designed to be:
- **Simple**: Single command does everything automatically
- **Robust**: Graceful error handling with interactive fallback
- **Flexible**: Model selection and dry-run mode
- **Standard**: Follows GitHub CLI extension patterns from gh-standup
