<p align="center">
  <a href="README.md"><img src="https://img.shields.io/badge/English-EN-blue" alt="English"></a>
  <a href="README_CN.md"><img src="https://img.shields.io/badge/简体中文-CN-red" alt="简体中文"></a>
</p>

<h1 align="center">PromptX</h1>
<p align="center"><strong>AI-Powered Terminal Assistant — Ask, Generate, Debug.</strong></p>
<p align="center">
  <a href="README.md">English</a> | <a href="README_CN.md">简体中文</a>
</p>
<br>
<p align="center">
  <a href="https://github.com/luoowei/promptx/stargazers"><img src="https://img.shields.io/github/stars/luoowei/promptx?style=social" alt="Stars"></a>
  <a href="https://github.com/luoowei/promptx/blob/main/LICENSE"><img src="https://img.shields.io/github/license/luoowei/promptx" alt="License"></a>
  <a href="https://pkg.go.dev/github.com/luoowei/promptx"><img src="https://pkg.go.dev/badge/github.com/luoowei/promptx.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/luoowei/promptx"><img src="https://goreportcard.com/badge/github.com/luoowei/promptx" alt="Go Report Card"></a>
</p>

---

**PromptX** brings AI superpowers directly into your terminal. Ask questions, generate shell commands, debug errors, and chat with AI models — without leaving the command line.

- **Single binary** — no Python, no Node.js, no Docker. Just download and run.
- **Multi-model** — OpenAI, Anthropic, Ollama, OpenRouter, Groq, and any OpenAI-compatible API.
- **Beautiful TUI** — interactive chat mode with markdown rendering and syntax highlighting.
- **Privacy-first** — API keys never leave your machine. Use local models for 100% offline use.
- **Shell integration** — pipe errors directly: `command 2>&1 | px ask "fix this"`

## Quick Start

```bash
# macOS / Linux
curl -fsSL https://raw.githubusercontent.com/luoowei/promptx/main/install.sh | bash

# Homebrew
brew install luoowei/tap/promptx

# Go
go install github.com/luoowei/promptx/cmd/px@latest

# Windows
iwr https://raw.githubusercontent.com/luoowei/promptx/main/install.ps1 | iex
```

Set your API key and start:

```bash
export OPENAI_API_KEY="sk-..."     # or ANTHROPIC_API_KEY
px ask "how do I find the largest files in a directory?"
```

## Usage

### Interactive Chat Mode
```bash
px
# Opens beautiful TUI — chat with AI in your terminal
```

### One-Shot Questions
```bash
px ask "explain this regex: ^(?=.*[A-Z])(?=.*[a-z])(?=.*\d).{8,}$"
px ask "what's the git command to undo the last commit but keep changes?"
```

### Shell Command Generation
```bash
px sh "find all .log files modified in the last 24 hours and compress them"
px sh "create a new user with sudo privileges on ubuntu"
```

### Debug Errors
```bash
# Pipe errors directly
make build 2>&1 | px ask "what's wrong and how do I fix it?"

# Or paste the error
px ask "what causes 'bind: address already in use' and how to fix it?"
```

### Configuration
```bash
px config show                          # View current config
px config set-provider anthropic        # Switch to Claude
px config set-key openai sk-...         # Store API key
px --provider ollama --model llama3.2   # Use local Ollama model
```

## Supported Providers

| Provider | Model Examples | Setup |
|----------|---------------|-------|
| **OpenAI** | gpt-4o, gpt-4.1 | `export OPENAI_API_KEY=sk-...` |
| **Anthropic** | claude-opus-4-8, claude-sonnet-4-6 | `export ANTHROPIC_API_KEY=sk-ant-...` |
| **Ollama** | llama3.2, mistral, qwen | `ollama pull llama3.2` (no key needed) |
| **OpenRouter** | 250+ models | `export OPENROUTER_API_KEY=...` |
| **Groq** | llama-3.3-70b | `export GROQ_API_KEY=...` |
| **Any OpenAI-compatible** | Custom endpoint | `px config set-key custom <key>` |

## How It Works

```
  Terminal     PromptX CLI      AI API
  (You)   -->  (Go Binary)  --> (OpenAI/
                |               Anthropic/
                v               Ollama)
              Response <--  Backend
```

PromptX is a statically-compiled Go binary that:
1. Reads your input from CLI args, pipes, or interactive TUI
2. Sends the request to your configured AI provider
3. Streams the response back with markdown rendering

**Zero dependencies at runtime.** No Python venv. No npm install. No Docker pull.

## Why PromptX?

| | PromptX | LLM CLI wrappers | Web Chat UIs |
|---|:---:|:---:|:---:|
| **Single binary** | Yes | No (Python/Node) | No (Browser only) |
| **Pipe support** | Yes | Partial | No |
| **Shell integration** | Yes | Partial | No |
| **Offline (local models)** | Yes | Yes | No |
| **Multi-provider** | Yes | Partial | Yes |
| **TUI mode** | Yes | Yes | No |
| **Fast startup** | <50ms | 200-2000ms | 2-5s |
| **Binary size** | ~15MB | 50-500MB+ | N/A |

## Build from Source

```bash
# Requires Go 1.22+
git clone https://github.com/luoowei/promptx.git
cd promptx
go build -o px ./cmd/px
./px ask "hello world"
```

## Project Structure

```
promptx/
├── cmd/px/main.go          # Entry point
├── internal/
│   ├── ai/                 # AI provider clients
│   │   ├── client.go       # Client interface + factory
│   │   ├── openai.go       # OpenAI-compatible client
│   │   └── anthropic.go    # Anthropic client
│   ├── cli/root.go         # Cobra CLI commands
│   ├── config/config.go    # Configuration management
│   ├── shell/handler.go    # Shell command generation
│   └── tui/
│       ├── tui.go          # Bubble Tea UI model
│       └── renderer.go     # Markdown renderer
├── docs/                   # Landing page (GitHub Pages)
└── install.sh              # One-line installer
```

## Roadmap

- [ ] Streaming responses in TUI mode
- [ ] Shell auto-completion scripts (bash, zsh, fish, powershell)
- [ ] Session history with search
- [ ] Custom system prompts
- [ ] MCP server integration
- [ ] Plugin system for additional providers
- [ ] Windows Terminal theming

## Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT © PromptX Contributors — see [LICENSE](LICENSE)

---

<p align="center">
  <b>Star this repo</b> if you find it useful!<br>
  <sub>Built with Go, Bubble Tea, and love for the terminal</sub>
</p>
