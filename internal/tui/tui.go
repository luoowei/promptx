package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luoowei/promptx/internal/ai"
)

type streamTickMsg struct{}
type streamCompleteMsg struct{ err error }
type errMsg struct{ err error }

type Model struct {
	viewport  viewport.Model
	textarea  textarea.Model
	messages  []ai.Message
	client    ai.Client
	streaming bool
	response  strings.Builder
	err       error
	width     int
	height    int
	ready     bool
}

var (
	senderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true)
	userStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	infoStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	borderStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("5")).Padding(0, 1)
	promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
)

func New(client ai.Client) *Model {
	ta := textarea.New()
	ta.Placeholder = "Ask anything... (Ctrl+D to quit, Enter to send)"
	ta.Prompt = "> "
	ta.ShowLineNumbers = false
	ta.SetHeight(3)
	ta.SetWidth(80)
	ta.Focus()
	ta.CharLimit = 4000

	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().Padding(0, 1)

	m := &Model{
		viewport: vp,
		textarea: ta,
		client:   client,
		messages: []ai.Message{},
	}

	m.messages = append(m.messages, ai.Message{
		Role: "assistant",
		Content: fmt.Sprintf(
			"Welcome to **PromptX**!\n\nProvider: %s | Model: %s\n\nType your question and press Enter.\n\n/clear - Clear  |  /model - Info  |  /help - Help",
			client.Name(), client.Model(),
		),
	})

	return m
}

func (m *Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.ready {
		if wm, ok := msg.(tea.WindowSizeMsg); ok {
			m.width = wm.Width
			m.height = wm.Height
			m.viewport = viewport.New(wm.Width-4, wm.Height-8)
			m.viewport.YPosition = 0
			m.textarea.SetWidth(wm.Width - 6)
			m.ready = true
			m.updateViewport()
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 8
		m.textarea.SetWidth(msg.Width - 6)

	case streamTickMsg:
		if m.streaming {
			if len(m.messages) > 0 {
				m.messages[len(m.messages)-1].Content = m.response.String()
			}
			m.updateViewport()
			return m, tickCmd()
		}

	case streamCompleteMsg:
		m.streaming = false
		if msg.err != nil {
			m.err = msg.err
			if len(m.messages) > 0 {
				m.messages[len(m.messages)-1].Content = fmt.Sprintf("Error: %v", msg.err)
			}
		} else {
			if len(m.messages) > 0 {
				m.messages[len(m.messages)-1].Content = m.response.String()
			}
		}
		m.updateViewport()

	case errMsg:
		m.err = msg.err
		m.streaming = false

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyCtrlD:
			return m, tea.Quit

		case tea.KeyEnter:
			if m.streaming {
				return m, nil
			}
			input := strings.TrimSpace(m.textarea.Value())
			if input == "" {
				return m, nil
			}
			if strings.HasPrefix(input, "/") {
				m.handleCommand(input)
				m.textarea.Reset()
				m.updateViewport()
				return m, nil
			}
			m.messages = append(m.messages, ai.Message{Role: "user", Content: input})
			m.streaming = true
			m.response.Reset()
			m.messages = append(m.messages, ai.Message{Role: "assistant", Content: "Thinking..."})
			m.textarea.Reset()
			m.updateViewport()
			return m, tea.Batch(m.streamRequest(input), tickCmd())
		}
	}

	var tiCmd tea.Cmd
	m.textarea, tiCmd = m.textarea.Update(msg)
	var vpCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m *Model) View() string {
	if !m.ready {
		return "  Initializing..."
	}
	var msgView strings.Builder
	for _, msg := range m.messages {
		switch msg.Role {
		case "user":
			msgView.WriteString(userStyle.Render("You") + "\n")
			msgView.WriteString(wordWrap(msg.Content, m.width-8) + "\n\n")
		case "assistant":
			msgView.WriteString(senderStyle.Render("PromptX") + "\n")
			msgView.WriteString(renderMarkdown(msg.Content, m.width-8) + "\n\n")
		}
	}
	m.viewport.SetContent(msgView.String())
	m.viewport.GotoBottom()

	status := fmt.Sprintf(" %s | %s ", m.client.Name(), m.client.Model())
	if m.streaming {
		status += "| Thinking..."
	}
	if m.err != nil {
		status += fmt.Sprintf("| Error: %v", m.err)
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		borderStyle.Render(m.viewport.View()),
		infoStyle.Render(status),
		promptStyle.Render(m.textarea.View()),
	)
}

func (m *Model) handleCommand(input string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}
	switch parts[0] {
	case "/clear":
		m.messages = []ai.Message{}
		m.messages = append(m.messages, ai.Message{Role: "assistant", Content: "Conversation cleared."})
	case "/model":
		m.messages = append(m.messages, ai.Message{Role: "assistant", Content: fmt.Sprintf("Provider: %s | Model: %s", m.client.Name(), m.client.Model())})
	case "/help":
		m.messages = append(m.messages, ai.Message{Role: "assistant", Content: "/clear - Clear history\n/model - Show model info\n/help - This help\nCtrl+D - Quit"})
	}
}

func (m *Model) streamRequest(input string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
		defer cancel()
		messages := make([]ai.Message, len(m.messages)-1)
		copy(messages, m.messages[:len(m.messages)-1])
		response, err := m.client.Chat(ctx, messages)
		if err != nil {
			return streamCompleteMsg{err: err}
		}
		m.response.Reset()
		m.response.WriteString(response)
		return streamCompleteMsg{}
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return streamTickMsg{}
	})
}

func (m *Model) updateViewport() {
	var msgView strings.Builder
	for _, msg := range m.messages {
		switch msg.Role {
		case "user":
			msgView.WriteString(userStyle.Render("You") + "\n")
			msgView.WriteString(wordWrap(msg.Content, m.width-8) + "\n\n")
		case "assistant":
			msgView.WriteString(senderStyle.Render("PromptX") + "\n")
			msgView.WriteString(renderMarkdown(msg.Content, m.width-8) + "\n\n")
		}
	}
	m.viewport.SetContent(msgView.String())
	m.viewport.GotoBottom()
}

func (m *Model) Run() (*tea.Program, error) {
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return nil, fmt.Errorf("run TUI: %w", err)
	}
	return p, nil
}
