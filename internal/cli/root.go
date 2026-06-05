package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/luoowei/promptx/internal/ai"
	"github.com/luoowei/promptx/internal/config"
	"github.com/luoowei/promptx/internal/shell"
	"github.com/luoowei/promptx/internal/tui"
	"github.com/spf13/cobra"
)

var (
	cfgFile   string
	provider  string
	model     string
	cfg       *config.Config
)

// rootCmd is the base command
var rootCmd = &cobra.Command{
	Use:   "px",
	Short: "PromptX - AI-powered terminal assistant",
	Long: `PromptX is an AI-powered terminal assistant that helps you:
  - Get shell command suggestions
  - Debug errors and explain code
  - Automate development workflows
  - Chat with AI models directly from your terminal

Run without arguments to start interactive chat mode.
Use 'px ask "your question"' for one-shot queries.
Use 'px sh "describe what you want"' for shell command generation.`,
	Run: func(cmd *cobra.Command, args []string) {
		runInteractive()
	},
}

// askCmd handles one-shot questions
var askCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask a one-shot question",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		question := strings.Join(args, " ")
		runAsk(question)
	},
}

// shCmd generates shell commands
var shCmd = &cobra.Command{
	Use:   "sh [description]",
	Short: "Generate shell commands from natural language",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		description := strings.Join(args, " ")
		runShell(description)
	},
}

// configCmd manages configuration
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage PromptX configuration",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			return
		}
		fmt.Printf("Default Provider: %s\n", c.DefaultProvider)
		fmt.Printf("Theme: %s\n", c.Theme)
		fmt.Printf("Max Tokens: %d\n", c.MaxTokens)
		fmt.Println("\nProviders:")
		for name, p := range c.Providers {
			keyStatus := "not set"
			if p.APIKey != "" {
				keyStatus = "********"
			}
			envKey := config.GetAPIKey(name)
			if envKey != "" {
				keyStatus = "(from env) ********"
			}
			fmt.Printf("  %s: model=%s, api_key=%s, base_url=%s\n",
				name, p.Model, keyStatus, p.BaseURL)
		}
	},
}

var configSetProviderCmd = &cobra.Command{
	Use:   "set-provider [name]",
	Short: "Set the default AI provider",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		c.DefaultProvider = args[0]
		if err := config.Save(c); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		fmt.Printf("Default provider set to: %s\n", args[0])
	},
}

var configSetKeyCmd = &cobra.Command{
	Use:   "set-key [provider] [api-key]",
	Short: "Set API key for a provider",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		p, ok := c.Providers[args[0]]
		if !ok {
			p = config.Provider{Name: args[0]}
		}
		p.APIKey = args[1]
		c.Providers[args[0]] = p
		if err := config.Save(c); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		fmt.Printf("API key set for: %s\n", args[0])
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&provider, "provider", "", "AI provider to use (openai, anthropic, ollama)")
	rootCmd.PersistentFlags().StringVar(&model, "model", "", "Model to use (overrides default)")

	rootCmd.AddCommand(askCmd)
	rootCmd.AddCommand(shCmd)
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetProviderCmd)
	configCmd.AddCommand(configSetKeyCmd)
}

// Execute runs the CLI
func Execute() error {
	return rootCmd.Execute()
}

// loadClient creates an AI client from configuration
func loadClient() (ai.Client, error) {
	var err error
	cfg, err = config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	selectedProvider := cfg.DefaultProvider
	if provider != "" {
		selectedProvider = provider
	}

	p, ok := cfg.Providers[selectedProvider]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s (available: %v)",
			selectedProvider, providerNames(cfg))
	}

	apiKey := config.GetAPIKey(selectedProvider)
	if apiKey == "" && selectedProvider != "ollama" {
		return nil, fmt.Errorf("no API key found for %s. Set %s_API_KEY env var or run 'px config set-key %s <key>'",
			selectedProvider, strings.ToUpper(selectedProvider), selectedProvider)
	}

	selectedModel := p.Model
	if model != "" {
		selectedModel = model
	}

	return ai.NewClient(ai.Config{
		Provider: selectedProvider,
		APIKey:   apiKey,
		BaseURL:  p.BaseURL,
		Model:    selectedModel,
	})
}

func providerNames(cfg *config.Config) []string {
	var names []string
	for k := range cfg.Providers {
		names = append(names, k)
	}
	return names
}

func runInteractive() {
	client, err := loadClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	t := tui.New(client)
	if _, err := t.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runAsk(question string) {
	client, err := loadClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	handle, _ := shell.NewCommandHandler(client)
	response, err := handle.Ask(question)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(response)
}

func runShell(description string) {
	client, err := loadClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	handle, _ := shell.NewCommandHandler(client)
	cmds, explanation, err := handle.GenerateCommands(description)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(explanation)
	fmt.Println()
	fmt.Println("Suggested commands:")
	for i, cmd := range cmds {
		fmt.Printf("  %d. %s\n", i+1, cmd)
	}
}
