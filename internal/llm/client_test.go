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
