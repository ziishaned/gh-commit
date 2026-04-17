package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/gh-commit/internal/types"

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

func NewClient() (*Client, error) {
	host, _ := auth.DefaultHost()
	token, _ := auth.TokenForHost(host)

	if token == "" {
		return nil, fmt.Errorf("no GitHub token found. Please run 'gh auth login' to authenticate")
	}

	return &Client{token: token}, nil
}

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
