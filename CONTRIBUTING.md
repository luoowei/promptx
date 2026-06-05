# Contributing to PromptX

Thanks for your interest in contributing! Here's how to get started.

## Getting Started

```bash
git clone https://github.com/luoowei/promptx.git
cd promptx
go mod download
go build ./cmd/px
./px ask "hello world"
```

## Development

### Prerequisites
- Go 1.22+
- An API key for testing (OpenAI or Anthropic)

### Project Structure
```
promptx/
├── cmd/px/main.go          # Entry point
├── internal/
�?  ├── ai/                 # AI provider clients
�?  ├── cli/                # Cobra CLI commands
�?  ├── config/             # Configuration management
�?  ├── shell/              # Shell command generation
�?  └── tui/                # Bubble Tea UI
└── website/                # Landing page
```

### Adding a New Provider

1. Create `internal/ai/yourprovider.go`
2. Implement the `Client` interface
3. Add to the factory in `internal/ai/client.go`
4. Add config defaults in `internal/config/config.go`

### Running Tests
```bash
go test ./...
```

### Code Style
- Follow standard Go conventions (`go fmt`, `go vet`)
- Use descriptive variable names
- Add comments for exported functions
- Keep functions small and focused

## Pull Request Process

1. Fork the repo
2. Create a feature branch (`git checkout -b feature/amazing`)
3. Make your changes
4. Verify builds: `go build ./...`
5. Commit with a clear message
6. Push and open a PR

## Issue Guidelines

When reporting bugs:
- Include your OS and shell
- Include exact error messages
- Include steps to reproduce
- Share what you expected to happen

## Code of Conduct

Be kind, be respectful, be constructive.

## Questions?

Open a [Discussion](https://github.com/luoowei/promptx/discussions) or an [Issue](https://github.com/luoowei/promptx/issues).
