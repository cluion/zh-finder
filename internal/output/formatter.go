package output

import (
	"io"

	"github.com/cluion/zh-finder/internal/classifier"
	"github.com/cluion/zh-finder/internal/matcher"
)

// LineMatch holds matches found on a single line.
type LineMatch struct {
	Line    int
	Content string
	Matches []matcher.Match
}

// FileMatch is an alias for a slice of LineMatch (per file).
type FileMatch = LineMatch

// Stats holds scan statistics.
type Stats struct {
	ScannedFiles     int
	MatchedFiles     int
	TraditionalCount int
	SimplifiedCount  int
	Duration         string
}

// Format renders scan results to the writer.
type Format interface {
	Format(w io.Writer, results map[string][]LineMatch, stats *Stats) error
}

// New returns the appropriate formatter based on format and color settings.
func New(colorEnabled bool, format string) Format {
	switch format {
	case "json", "json-verbose":
		return &JSONFormatter{}
	case "json-compact":
		return &JSONCompactFormatter{}
	default:
		return &TerminalFormatter{colorEnabled: colorEnabled}
	}
}

// typeLabel returns a human-readable label for a HanType.
func typeLabel(t classifier.HanType) string {
	switch t {
	case classifier.Traditional:
		return "traditional"
	case classifier.Simplified:
		return "simplified"
	default:
		return "common"
	}
}
