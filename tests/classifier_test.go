package tests

import (
	"testing"

	"github.com/cluion/zh-finder/internal/classifier"
)

func TestClassifyTraditional(t *testing.T) {
	c := classifier.New("繁體測試", "")
	result := c.Classify('繁')
	if result != classifier.Traditional {
		t.Errorf("Expected Traditional, got %v", result)
	}
}

func TestClassifySimplified(t *testing.T) {
	c := classifier.New("", "简体测试")
	result := c.Classify('简')
	if result != classifier.Simplified {
		t.Errorf("Expected Simplified, got %v", result)
	}
}

func TestClassifyCommon(t *testing.T) {
	c := classifier.New("中", "中")
	result := c.Classify('中')
	if result != classifier.Common {
		t.Errorf("Expected Common, got %v", result)
	}
}

func TestClassifyNonHan(t *testing.T) {
	c := classifier.New("", "")
	result := c.Classify('a')
	if result != classifier.Common {
		t.Errorf("Expected Common for non-Han, got %v", result)
	}
}

func TestIsHan(t *testing.T) {
	tests := []struct {
		char     rune
		expected bool
	}{
		{'中', true},
		{'繁', true},
		{'简', true},
		{'a', false},
		{'1', false},
		{' ', false},
	}

	for _, tt := range tests {
		result := classifier.IsHan(tt.char)
		if result != tt.expected {
			t.Errorf("IsHan(%c) = %v, expected %v", tt.char, result, tt.expected)
		}
	}
}
