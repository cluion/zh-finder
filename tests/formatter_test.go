package tests

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cluion/zh-finder/internal/classifier"
	"github.com/cluion/zh-finder/internal/matcher"
	"github.com/cluion/zh-finder/internal/output"
)

func TestFormatTerminalOutput(t *testing.T) {
	f := output.New(false, "term")

	results := map[string][]output.LineMatch{
		"test.php": {
			{
				Line:    1,
				Content: "// 測試",
				Matches: []matcher.Match{
					{Text: "測試", Type: classifier.Traditional, Start: 3, End: 5},
				},
			},
		},
	}

	var buf bytes.Buffer
	err := f.Format(&buf, results, nil)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "test.php:1") {
		t.Error("Expected file:line in output")
	}
	if !strings.Contains(got, "測試") {
		t.Error("Expected match text in output")
	}
}

func TestFormatJSONOutput(t *testing.T) {
	f := output.New(false, "json")

	results := map[string][]output.LineMatch{
		"test.php": {
			{
				Line:    1,
				Content: "// 測試",
				Matches: []matcher.Match{
					{Text: "測試", Type: classifier.Traditional, Start: 3, End: 5},
				},
			},
		},
	}

	var buf bytes.Buffer
	err := f.Format(&buf, results, nil)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `"file": "test.php"`) {
		t.Error("Expected file field in JSON")
	}
	if !strings.Contains(got, `"results"`) {
		t.Error("Expected results field in JSON")
	}
}

func TestFormatJSONCompactOutput(t *testing.T) {
	f := output.New(false, "json-compact")

	results := map[string][]output.LineMatch{
		"test.php": {
			{
				Line:    1,
				Content: "// 測試",
				Matches: []matcher.Match{
					{Text: "測試", Type: classifier.Traditional, Start: 3, End: 5},
				},
			},
		},
	}

	var buf bytes.Buffer
	err := f.Format(&buf, results, nil)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `"matches"`) {
		t.Error("Expected matches field in compact JSON")
	}
}

func TestFormatWithStats(t *testing.T) {
	f := output.New(false, "term")

	results := map[string][]output.LineMatch{
		"test.php": {
			{
				Line:    1,
				Content: "// 測試",
				Matches: []matcher.Match{
					{Text: "測試", Type: classifier.Traditional, Start: 3, End: 5},
				},
			},
		},
	}

	stats := &output.Stats{
		MatchedFiles:     1,
		TraditionalCount: 1,
		SimplifiedCount:  0,
		Duration:         "0.01s",
	}

	var buf bytes.Buffer
	err := f.Format(&buf, results, stats)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "Statistics") {
		t.Error("Expected statistics in output")
	}
}
