package ai

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`    // "system", "user", "assistant"
	Content string `json:"content"` // message content
}

// StreamChunk represents a chunk of streaming response
type StreamChunk struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
	Error   error  `json:"error,omitempty"`
}

// Client defines the AI provider interface
type Client interface {
	// Chat sends messages and returns the response
	Chat(ctx context.Context, messages []Message) (string, error)

	// ChatStream sends messages and streams the response
	ChatStream(ctx context.Context, messages []Message) (<-chan StreamChunk, error)

	// Name returns the provider name
	Name() string

	// Model returns the current model name
	Model() string
}

// Config holds client configuration
type Config struct {
	Provider string
	APIKey   string
	BaseURL  string
	Model    string
	Timeout  time.Duration
}

// NewClient creates an AI client for the given configuration
func NewClient(cfg Config) (Client, error) {
	switch strings.ToLower(cfg.Provider) {
	case "openai":
		return NewOpenAIClient(cfg)
	case "anthropic":
		return NewAnthropicClient(cfg)
	case "ollama":
		return NewOllamaClient(cfg)
	default:
		return nil, fmt.Errorf("unsupported provider: %s (supported: openai, anthropic, ollama)", cfg.Provider)
	}
}

// SystemPrompt is the default system prompt for PromptX
const SystemPrompt = `You are PromptX, an AI assistant that lives in the terminal.
You help developers with:
- Shell commands and scripting
- Code explanations and debugging
- System administration
- Development workflows

Rules:
1. Be concise - terminal users value brevity
2. Use markdown for code blocks with language hints
3. For shell commands, provide the exact command with brief explanation
4. When fixing errors, explain the root cause first, then the fix
5. When generating commands, ensure they are safe (no destructive operations without warning)
6. For multi-step solutions, number the steps clearly`
