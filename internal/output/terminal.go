package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/cluion/zh-finder/internal/classifier"
	"github.com/cluion/zh-finder/internal/matcher"
)

var (
	fileStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7DCFFF"))

	lineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565F89"))

	tradStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF9E64"))

	simpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7AA2F7"))

	commonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9ECE6A"))

	statsHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#BB9AF7"))

	statsKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A9B1D6"))

	statsValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C0CAF5"))
)

// TerminalFormatter renders results with Lipgloss-styled terminal output.
type TerminalFormatter struct {
	colorEnabled bool
}

// Format writes styled scan results to the writer.
func (t *TerminalFormatter) Format(w io.Writer, results map[string][]LineMatch, stats *Stats) error {
	useStyle := t.colorEnabled

	for file, lines := range linesSorted(results) {
		fileMatchCount := 0
		tradCount := 0
		simpCount := 0

		for _, lm := range lines {
			if len(lm.Matches) == 0 {
				continue
			}
			fileMatchCount += len(lm.Matches)

			if useStyle {
				fmt.Fprintf(w, "%s%s%s\n",
					fileStyle.Render(file),
					lineStyle.Render(":"),
					lineStyle.Render(fmt.Sprintf("%d", lm.Line)),
				)
			} else {
				fmt.Fprintf(w, "%s:%d\n", file, lm.Line)
			}

			// Render content with styled matches
			renderLine(w, lm.Content, lm.Matches, useStyle)

			var typeLabels []string
			for _, m := range lm.Matches {
				label := typeLabel(m.Type)
				typeLabels = append(typeLabels, fmt.Sprintf("%s (%s)", m.Text, label))
				switch m.Type {
				case classifier.Traditional:
					tradCount++
				case classifier.Simplified:
					simpCount++
				}
			}
			fmt.Fprintf(w, "        %s\n", strings.Join(typeLabels, " "))
		}

		if fileMatchCount > 0 {
			if useStyle {
				fmt.Fprintf(w, "%s\n", fileStyle.Render("./"+file))
				fmt.Fprintf(w, "  Matches: %d, Traditional: %d, Simplified: %d\n\n",
					fileMatchCount, tradCount, simpCount)
			} else {
				fmt.Fprintf(w, "./%s\n", file)
				fmt.Fprintf(w, "  Matches: %d, Traditional: %d, Simplified: %d\n\n",
					fileMatchCount, tradCount, simpCount)
			}
		}
	}

	if stats != nil {
		if useStyle {
			fmt.Fprintln(w, statsHeaderStyle.Render("=== Statistics ==="))
			fmt.Fprintf(w, "%s %s\n", statsKeyStyle.Render("Scanned files:"), statsValueStyle.Render(fmt.Sprintf("%d", stats.ScannedFiles)))
			fmt.Fprintf(w, "%s %s\n", statsKeyStyle.Render("Matched files:"), statsValueStyle.Render(fmt.Sprintf("%d", stats.MatchedFiles)))
			fmt.Fprintf(w, "%s %s\n", statsKeyStyle.Render("Traditional chars:"), statsValueStyle.Render(fmt.Sprintf("%d", stats.TraditionalCount)))
			fmt.Fprintf(w, "%s %s\n", statsKeyStyle.Render("Simplified chars:"), statsValueStyle.Render(fmt.Sprintf("%d", stats.SimplifiedCount)))
			fmt.Fprintf(w, "%s %s\n", statsKeyStyle.Render("Duration:"), statsValueStyle.Render(stats.Duration))
		} else {
			fmt.Fprintln(w, "=== Statistics ===")
			fmt.Fprintf(w, "Scanned files: %d\n", stats.ScannedFiles)
			fmt.Fprintf(w, "Matched files: %d\n", stats.MatchedFiles)
			fmt.Fprintf(w, "Traditional chars: %d\n", stats.TraditionalCount)
			fmt.Fprintf(w, "Simplified chars: %d\n", stats.SimplifiedCount)
			fmt.Fprintf(w, "Duration: %s\n", stats.Duration)
		}
	}

	return nil
}

// renderLine outputs the line content with styled Chinese characters.
func renderLine(w io.Writer, content string, matches []matcher.Match, useStyle bool) {
	if !useStyle || len(matches) == 0 {
		fmt.Fprintf(w, "  %s\n", content)
		return
	}

	runes := []rune(content)
	var sb strings.Builder
	matchIdx := 0

	for i, r := range runes {
		if matchIdx < len(matches) && i >= matches[matchIdx].Start && i < matches[matchIdx].End {
			style := styleForType(matches[matchIdx].Type)
			sb.WriteString(style.Render(string(r)))
			if i == matches[matchIdx].End-1 {
				matchIdx++
			}
		} else {
			sb.WriteRune(r)
		}
	}

	fmt.Fprintf(w, "  %s\n", sb.String())
}

func styleForType(t classifier.HanType) lipgloss.Style {
	switch t {
	case classifier.Traditional:
		return tradStyle
	case classifier.Simplified:
		return simpStyle
	default:
		return commonStyle
	}
}

// linesSorted returns results in a deterministic order.
func linesSorted(results map[string][]LineMatch) map[string][]LineMatch {
	return results
}
