package shell

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/luoowei/promptx/internal/ai"
)

// Handler generates and executes shell commands
type Handler struct {
	client  ai.Client
	history []string
}

// NewCommandHandler creates a new shell command handler
func NewCommandHandler(client ai.Client) (*Handler, error) {
	return &Handler{
		client:  client,
		history: make([]string, 0),
	}, nil
}

// Ask sends a general question to the AI
func (h *Handler) Ask(question string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	messages := []ai.Message{
		{Role: "system", Content: ai.SystemPrompt},
		{Role: "user", Content: question},
	}

	return h.client.Chat(ctx, messages)
}

// GenerateCommands generates shell commands from natural language
func (h *Handler) GenerateCommands(description string) ([]string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	systemPrompt := fmt.Sprintf(`You are a shell command generator. The user is running on %s with shell: %s.
Generate the exact shell command(s) to accomplish the user's task.
Output format:
---EXPLANATION---
Brief explanation of what the commands do
---COMMANDS---
command1
command2

Rules:
- Output ONLY in the format above
- Each command should be on its own line
- Commands should be safe (no rm -rf without confirmation, no destructive operations)
- Use cross-platform commands when possible
- For Windows PowerShell, use appropriate cmdlets
- Do NOT use && to chain unless necessary`, runtime.GOOS, detectShell())

	messages := []ai.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: description},
	}

	response, err := h.client.Chat(ctx, messages)
	if err != nil {
		return nil, "", fmt.Errorf("generate commands: %w", err)
	}

	return parseCommandResponse(response)
}

// ExplainError gets an explanation for an error message
func (h *Handler) ExplainError(errorMsg string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	messages := []ai.Message{
		{Role: "system", Content: "You are an expert at debugging. Explain the error in simple terms, then provide the fix."},
		{Role: "user", Content: fmt.Sprintf("Explain this error and how to fix it:\n%s", errorMsg)},
	}

	return h.client.Chat(ctx, messages)
}

// DetectShell returns the current shell name
func detectShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		if runtime.GOOS == "windows" {
			shell = os.Getenv("COMSPEC")
			if shell == "" {
				return "PowerShell"
			}
			if strings.Contains(strings.ToLower(shell), "powershell") {
				return "PowerShell"
			}
			return "cmd.exe"
		}
		return "bash"
	}
	return shell
}

// IsSafeCommand checks if a command is potentially dangerous
func IsSafeCommand(cmd string) bool {
	dangerous := []string{
		"rm -rf /",
		"dd if=",
		"mkfs.",
		"> /dev/sda",
		"format c:",
		"del /f /s",
		"DROP TABLE",
		"DROP DATABASE",
	}

	cmdLower := strings.ToLower(cmd)
	for _, d := range dangerous {
		if strings.Contains(cmdLower, strings.ToLower(d)) {
			return false
		}
	}
	return true
}

// ExecuteCommand runs a shell command and returns output
func ExecuteCommand(command string) (string, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-Command", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func parseCommandResponse(response string) ([]string, string, error) {
	var explanation string
	var commands []string

	lines := strings.Split(response, "\n")
	inCommands := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(trimmed, "---COMMANDS---"):
			inCommands = true
			continue
		case strings.HasPrefix(trimmed, "---EXPLANATION---"):
			continue
		}

		if inCommands {
			if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
				commands = append(commands, trimmed)
			}
		} else {
			explanation += line + "\n"
		}
	}

	explanation = strings.TrimSpace(explanation)
	return commands, explanation, nil
}
