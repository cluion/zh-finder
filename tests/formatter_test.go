package tests

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cluion/zh-finder/internal/classifier"
	"github.com/cluion/zh-finder/internal/matcher"
	"github.com/cluion/zh-finder/internal/output"
	"github.com/cluion/zh-finder/internal/scanner"
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

func TestFormatTerminalWithColor(t *testing.T) {
	f := output.New(true, "term")

	results := map[string][]output.LineMatch{
		"test.go": {
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
	if !strings.Contains(got, "test.go") {
		t.Error("Expected file name in colored output")
	}
	if !strings.Contains(got, "測試") {
		t.Error("Expected match text in colored output")
	}
}

func TestFormatTerminalSimplified(t *testing.T) {
	f := output.New(false, "term")

	results := map[string][]output.LineMatch{
		"test.go": {
			{
				Line:    1,
				Content: "// 简体",
				Matches: []matcher.Match{
					{Text: "简体", Type: classifier.Simplified, Start: 3, End: 5},
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
	if !strings.Contains(got, "简体 (simplified)") {
		t.Error("Expected simplified label in output")
	}
}

func TestFormatTerminalWithColorAndStats(t *testing.T) {
	f := output.New(true, "term")

	results := map[string][]output.LineMatch{
		"test.go": {
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
		ScannedFiles:     5,
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
		t.Error("Expected statistics in colored output")
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

func TestFormatJSONWithStats(t *testing.T) {
	f := output.New(false, "json")

	results := map[string][]output.LineMatch{
		"test.go": {
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
		ScannedFiles:     5,
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
	if !strings.Contains(got, `"summary"`) {
		t.Error("Expected summary field in JSON with stats")
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

func TestScannerSetExcludeDirs(t *testing.T) {
	s := scanner.New()
	s.SetExcludeDirs([]string{"custom"})
	if len(s.ExcludeDirs()) != 1 || s.ExcludeDirs()[0] != "custom" {
		t.Error("Expected custom exclude dir")
	}
}

func TestScannerAddExcludeDirs(t *testing.T) {
	s := scanner.New()
	original := len(s.ExcludeDirs())
	s.AddExcludeDirs([]string{"extra"})
	if len(s.ExcludeDirs()) != original+1 {
		t.Error("Expected additional exclude dir")
	}
}

func TestScannerSetMaxDepth(t *testing.T) {
	s := scanner.New()
	s.SetMaxDepth(5)
	if s.MaxDepth() != 5 {
		t.Error("Expected max depth 5")
	}
}

func TestScannerSetScanBinary(t *testing.T) {
	s := scanner.New()
	s.SetScanBinary(true)
	if !s.ScanBinary() {
		t.Error("Expected scan binary true")
	}
}

func TestMatcherMixedTraditionalSimplified(t *testing.T) {
	c := classifier.New("繁體", "简体")
	m := matcher.New(c)

	results := m.Find("繁體简体混合")
	if len(results) < 1 {
		t.Fatalf("Expected at least 1 match, got %d", len(results))
	}
	// Mixed traditional + simplified in one sequence should be classified as Common
	if results[0].Type != classifier.Common {
		t.Errorf("Expected Common for mixed, got %v", results[0].Type)
	}
}

func TestMatcherSimplifiedOnly(t *testing.T) {
	c := classifier.New("", "简体")
	m := matcher.New(c)

	results := m.Find("这是简体")
	if len(results) < 1 {
		t.Fatalf("Expected at least 1 match, got %d", len(results))
	}
	if results[0].Type != classifier.Simplified {
		t.Errorf("Expected Simplified, got %v", results[0].Type)
	}
}
