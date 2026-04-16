package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cluion/zh-finder/internal/scanner"
)

func TestScanReturnsFiles(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "zh-finder-test")
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "test.txt")
	os.WriteFile(testFile, []byte("test content"), 0644)

	s := scanner.New()
	var files []scanner.FileInfo
	for f := range s.Scan(tempDir) {
		files = append(files, f)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}
}

func TestScanWithExtensionFilter(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "zh-finder-test")
	defer os.RemoveAll(tempDir)

	os.WriteFile(filepath.Join(tempDir, "test.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tempDir, "test.go"), []byte("package main"), 0644)

	s := scanner.New()
	s.SetExtensions([]string{"go"})

	var files []scanner.FileInfo
	for f := range s.Scan(tempDir) {
		files = append(files, f)
	}

	for _, f := range files {
		if f.Extension != "go" {
			t.Errorf("Expected go extension, got %s", f.Extension)
		}
	}
}

func TestIsBinaryDetectsBinary(t *testing.T) {
	tempFile, _ := os.CreateTemp("", "test")
	defer os.Remove(tempFile.Name())

	tempFile.Write([]byte("text\x00binary"))
	tempFile.Close()

	s := scanner.New()
	if !s.IsBinary(tempFile.Name()) {
		t.Error("Expected binary detection")
	}
}

func TestIsBinaryDetectsText(t *testing.T) {
	tempFile, _ := os.CreateTemp("", "test")
	defer os.Remove(tempFile.Name())

	tempFile.Write([]byte("plain text"))
	tempFile.Close()

	s := scanner.New()
	if s.IsBinary(tempFile.Name()) {
		t.Error("Expected text detection")
	}
}
