package llm

import (
	"fmt"

	_ "embed"
	"gopkg.in/yaml.v3"
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

func loadPromptConfig() (*PromptConfig, error) {
	var config PromptConfig
	err := yaml.Unmarshal(commitPromptYAML, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse prompt configuration: %w", err)
	}
	return &config, nil
}
