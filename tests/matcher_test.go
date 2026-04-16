package tests

import (
	"testing"

	"github.com/cluion/zh-finder/internal/classifier"
	"github.com/cluion/zh-finder/internal/matcher"
)

func TestFindNoMatches(t *testing.T) {
	c := classifier.New("", "")
	m := matcher.New(c)

	results := m.Find("Hello World")
	if len(results) != 0 {
		t.Errorf("Expected 0 matches, got %d", len(results))
	}
}

func TestFindSingleMatch(t *testing.T) {
	c := classifier.New("", "")
	m := matcher.New(c)

	results := m.Find("Hello 世界")
	if len(results) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(results))
	}
	if results[0].Text != "世界" {
		t.Errorf("Expected '世界', got '%s'", results[0].Text)
	}
}

func TestFindMultipleMatches(t *testing.T) {
	c := classifier.New("", "")
	m := matcher.New(c)

	results := m.Find("這是繁體，这是简体")
	if len(results) < 2 {
		t.Errorf("Expected at least 2 matches, got %d", len(results))
	}
}

func TestFindCorrectPositions(t *testing.T) {
	c := classifier.New("", "")
	m := matcher.New(c)

	results := m.Find("abc測試def")
	if len(results) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(results))
	}
	if results[0].Start != 3 {
		t.Errorf("Expected start 3, got %d", results[0].Start)
	}
	if results[0].End != 5 {
		t.Errorf("Expected end 5, got %d", results[0].End)
	}
}
