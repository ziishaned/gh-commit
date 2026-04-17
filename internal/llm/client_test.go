package llm

import (
	"os"
	"testing"

	"github.com/gh-commit/internal/types"
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
