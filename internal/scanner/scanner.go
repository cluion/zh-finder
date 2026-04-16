package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

var DefaultExcludes = []string{
	".git", ".svn", ".hg",
	"node_modules", "vendor",
	"dist", "build", "out", "target",
	".cache", ".npm", ".yarn",
	"__pycache__", ".pytest_cache",
}

type FileInfo struct {
	Path         string
	RelativePath string
	Extension    string
	IsBinary     bool
}

type Scanner struct {
	extensions  []string
	excludeDirs []string
	maxDepth    int
	scanBinary  bool
}

func New() *Scanner {
	return &Scanner{
		excludeDirs: DefaultExcludes,
		maxDepth:    -1,
	}
}

func (s *Scanner) Scan(root string) <-chan FileInfo {
	out := make(chan FileInfo)

	go func() {
		defer close(out)

		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.Mode().IsRegular() {
				return nil
			}

			relPath, _ := filepath.Rel(root, path)
			if s.isExcluded(relPath) {
				return nil
			}

			ext := strings.ToLower(filepath.Ext(path))
			if len(ext) > 0 {
				ext = ext[1:] // remove dot
			}

			if len(s.extensions) > 0 && !s.contains(s.extensions, ext) {
				return nil
			}

			isBinary := s.IsBinary(path)
			if isBinary && !s.scanBinary {
				return nil
			}

			out <- FileInfo{
				Path:         path,
				RelativePath: relPath,
				Extension:    ext,
				IsBinary:     isBinary,
			}
			return nil
		})
	}()

	return out
}

func (s *Scanner) isExcluded(relPath string) bool {
	parts := strings.Split(filepath.Dir(relPath), string(os.PathSeparator))
	for _, part := range parts {
		if s.contains(s.excludeDirs, part) {
			return true
		}
	}
	return false
}

func (s *Scanner) IsBinary(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return true
	}
	defer f.Close()

	buf := make([]byte, 8000)
	n, _ := f.Read(buf)
	for i := 0; i < n; i++ {
		if buf[i] == 0 {
			return true
		}
	}
	return false
}

func (s *Scanner) contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func (s *Scanner) SetExtensions(exts []string) {
	s.extensions = exts
}

func (s *Scanner) SetExcludeDirs(dirs []string) {
	s.excludeDirs = dirs
}

func (s *Scanner) AddExcludeDirs(dirs []string) {
	s.excludeDirs = append(s.excludeDirs, dirs...)
}

func (s *Scanner) SetMaxDepth(depth int) {
	s.maxDepth = depth
}

func (s *Scanner) SetScanBinary(scan bool) {
	s.scanBinary = scan
}
