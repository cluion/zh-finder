package output

import (
	"encoding/json"
	"io"
)

type jsonChar struct {
	Text  string `json:"text"`
	Type  string `json:"type"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

type jsonLineMatch struct {
	Line       int         `json:"line"`
	Content    string      `json:"content"`
	Characters []jsonChar  `json:"characters"`
}

type jsonResult struct {
	File    string          `json:"file"`
	Matches []jsonLineMatch `json:"matches"`
}

type jsonStats struct {
	ScannedFiles     int    `json:"scannedFiles"`
	MatchedFiles     int    `json:"matchedFiles"`
	TraditionalCount int    `json:"traditionalCount"`
	SimplifiedCount  int    `json:"simplifiedCount"`
	Duration         string `json:"duration"`
}

type jsonOutput struct {
	Results []jsonResult `json:"results"`
	Summary *jsonStats   `json:"summary,omitempty"`
}

// JSONFormatter renders results as pretty-printed JSON.
type JSONFormatter struct{}

// Format writes scan results as JSON to the writer.
func (j *JSONFormatter) Format(w io.Writer, results map[string][]LineMatch, stats *Stats) error {
	var output jsonOutput

	for file, lines := range results {
		var matches []jsonLineMatch
		for _, lm := range lines {
			var chars []jsonChar
			for _, m := range lm.Matches {
				chars = append(chars, jsonChar{
					Text:  m.Text,
					Type:  typeLabel(m.Type),
					Start: m.Start,
					End:   m.End,
				})
			}
			matches = append(matches, jsonLineMatch{
				Line:       lm.Line,
				Content:    lm.Content,
				Characters: chars,
			})
		}
		output.Results = append(output.Results, jsonResult{
			File:    file,
			Matches: matches,
		})
	}

	if stats != nil {
		output.Summary = &jsonStats{
			ScannedFiles:     stats.ScannedFiles,
			MatchedFiles:     stats.MatchedFiles,
			TraditionalCount: stats.TraditionalCount,
			SimplifiedCount:  stats.SimplifiedCount,
			Duration:         stats.Duration,
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}

type jsonCompactResult struct {
	File    string   `json:"file"`
	Line    int      `json:"line"`
	Matches []string `json:"matches"`
}

type jsonCompactOutput struct {
	Results []jsonCompactResult `json:"results"`
}

// JSONCompactFormatter renders results as compact JSON.
type JSONCompactFormatter struct{}

// Format writes scan results as compact JSON to the writer.
func (j *JSONCompactFormatter) Format(w io.Writer, results map[string][]LineMatch, _ *Stats) error {
	var output jsonCompactOutput

	for file, lines := range results {
		for _, lm := range lines {
			var matches []string
			for _, m := range lm.Matches {
				matches = append(matches, m.Text)
			}
			if len(matches) > 0 {
				output.Results = append(output.Results, jsonCompactResult{
					File:    file,
					Line:    lm.Line,
					Matches: matches,
				})
			}
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}
