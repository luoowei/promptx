package tui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	codeStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("0")).
			Foreground(lipgloss.Color("10")).
			Padding(0, 1)

	boldStyle = lipgloss.NewStyle().
			Bold(true)

	italicStyle = lipgloss.NewStyle().
			Italic(true)

	linkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).
			Underline(true)

	listStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("6"))

	codeBlockRegexp = regexp.MustCompile("```(?:\\w+)?\\s*\\n([\\s\\S]*?)```")
	inlineCodeRegexp = regexp.MustCompile("`([^`]+)`")
	boldRegexp      = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	italicRegexp    = regexp.MustCompile(`\*([^*]+)\*`)
	listItemRegexp  = regexp.MustCompile(`(?m)^[-*]\s+(.+)$`)
	numberedRegexp  = regexp.MustCompile(`(?m)^(\d+\.\s+.+)$`)
)

// renderMarkdown renders basic markdown formatting for the terminal
func renderMarkdown(text string, width int) string {
	if text == "" {
		return ""
	}

	// Extract and replace code blocks with placeholders
	codeBlocks := []string{}
	text = codeBlockRegexp.ReplaceAllStringFunc(text, func(match string) string {
		submatches := codeBlockRegexp.FindStringSubmatch(match)
		if len(submatches) > 1 {
			codeBlocks = append(codeBlocks, strings.TrimSpace(submatches[1]))
		}
		return fmt.Sprintf("{{CODEBLOCK_%d}}", len(codeBlocks)-1)
	})

	// Bold
	text = boldRegexp.ReplaceAllStringFunc(text, func(match string) string {
		submatches := boldRegexp.FindStringSubmatch(match)
		if len(submatches) > 1 {
			return boldStyle.Render(submatches[1])
		}
		return match
	})

	// Italic
	text = italicRegexp.ReplaceAllStringFunc(text, func(match string) string {
		submatches := italicRegexp.FindStringSubmatch(match)
		if len(submatches) > 1 {
			return italicStyle.Render(submatches[1])
		}
		return match
	})

	// Inline code
	text = inlineCodeRegexp.ReplaceAllStringFunc(text, func(match string) string {
		submatches := inlineCodeRegexp.FindStringSubmatch(match)
		if len(submatches) > 1 {
			return codeStyle.Render(submatches[1])
		}
		return match
	})

	// List items
	text = listItemRegexp.ReplaceAllString(text, listStyle.Render("  - ")+"$1")

	// Numbered items
	text = numberedRegexp.ReplaceAllString(text, listStyle.Render("$1"))

	// Restore code blocks
	for i, block := range codeBlocks {
		placeholder := fmt.Sprintf("{{CODEBLOCK_%d}}", i)
		rendered := codeStyle.Render(block)
		text = strings.Replace(text, placeholder, "\n"+rendered+"\n", 1)
	}

	// Wrap text to width
	return wordWrap(text, width)
}

// wordWrap wraps text to a maximum width
func wordWrap(text string, width int) string {
	if width <= 0 {
		return text
	}

	var result strings.Builder
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if len(line) == 0 {
			result.WriteString("\n")
			continue
		}

		if lipgloss.Width(line) <= width {
			result.WriteString(line)
			result.WriteString("\n")
			continue
		}

		words := strings.Fields(line)
		currentLine := ""

		for _, word := range words {
			testLine := currentLine
			if testLine != "" {
				testLine += " "
			}
			testLine += word

			if lipgloss.Width(testLine) <= width {
				currentLine = testLine
			} else {
				if currentLine != "" {
					result.WriteString(currentLine + "\n")
				}
				currentLine = word
			}
		}

		if currentLine != "" {
			result.WriteString(currentLine + "\n")
		}
	}

	return strings.TrimRight(result.String(), "\n")
}
