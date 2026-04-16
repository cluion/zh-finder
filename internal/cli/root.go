package cli

import (
	"embed"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cluion/zh-finder/internal/classifier"
	"github.com/cluion/zh-finder/internal/matcher"
	"github.com/cluion/zh-finder/internal/output"
	"github.com/cluion/zh-finder/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	extsString    string
	excludeString string
	excludeAdd    string
	noExclude     bool
	format        string
	showStats     bool
	maxDepth      int
	scanBinary    bool
	noColor       bool
	typeFilter    string
	dataFS        embed.FS
)

var rootCmd = &cobra.Command{
	Use:   "zh-finder",
	Short: "Find Chinese characters (Traditional/Simplified) in files",
	Long:  "A CLI tool to scan files for Chinese characters, distinguishing between Traditional and Simplified Chinese.",
	SilenceUsage:  true,
	SilenceErrors: true,
}

var scanCmd = &cobra.Command{
	Use:   "scan <path>",
	Short: "Scan files for Chinese characters",
	Args:  cobra.ExactArgs(1),
	RunE:  runScan,
}

func init() {
	scanCmd.Flags().StringVar(&extsString, "ext", "", "Only scan specific extensions (comma-separated)")
	scanCmd.Flags().StringVar(&excludeString, "exclude", "", "Exclude directories (comma-separated)")
	scanCmd.Flags().StringVar(&excludeAdd, "exclude-add", "", "Additional directories to exclude")
	scanCmd.Flags().BoolVar(&noExclude, "no-exclude", false, "Disable default excludes")
	scanCmd.Flags().StringVar(&format, "format", "term", "Output format: term, json, json-compact")
	scanCmd.Flags().BoolVar(&showStats, "stats", false, "Show statistics")
	scanCmd.Flags().IntVar(&maxDepth, "max-depth", -1, "Maximum recursion depth")
	scanCmd.Flags().BoolVar(&scanBinary, "binary", false, "Scan binary files")
	scanCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable color output")
	scanCmd.Flags().StringVar(&typeFilter, "type", "all", "Filter: all, traditional, simplified")

	rootCmd.AddCommand(scanCmd)
}

// SetDataFS injects the embedded data filesystem from main.
func SetDataFS(fs embed.FS) {
	dataFS = fs
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runScan(cmd *cobra.Command, args []string) error {
	path := args[0]
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}

	// Load classifier from embedded data
	tradData, err := dataFS.ReadFile("data/traditional.txt")
	if err != nil {
		return fmt.Errorf("loading traditional data: %w", err)
	}
	simpData, err := dataFS.ReadFile("data/simplified.txt")
	if err != nil {
		return fmt.Errorf("loading simplified data: %w", err)
	}

	c := classifier.New(string(tradData), string(simpData))
	m := matcher.New(c)
	s := scanner.New()

	// Configure scanner
	if extsString != "" {
		s.SetExtensions(strings.Split(extsString, ","))
	}
	if noExclude {
		s.SetExcludeDirs([]string{})
	} else if excludeString != "" {
		s.SetExcludeDirs(strings.Split(excludeString, ","))
	}
	if excludeAdd != "" {
		s.AddExcludeDirs(strings.Split(excludeAdd, ","))
	}
	if maxDepth >= 0 {
		s.SetMaxDepth(maxDepth)
	}
	s.SetScanBinary(scanBinary)

	// Filter type
	var filterType classifier.HanType
	switch typeFilter {
	case "traditional":
		filterType = classifier.Traditional
	case "simplified":
		filterType = classifier.Simplified
	default:
		filterType = classifier.AllHanType
	}

	// Scan
	start := time.Now()
	results := make(map[string][]output.LineMatch)
	scannedCount := 0

	for file := range s.Scan(path) {
		scannedCount++
		content, err := os.ReadFile(file.Path)
		if err != nil {
			continue
		}

		lines := strings.Split(string(content), "\n")
		var fileResults []output.LineMatch

		for lineNum, line := range lines {
			matches := m.Find(line)
			if filterType != classifier.AllHanType {
				var filtered []matcher.Match
				for _, match := range matches {
					if match.Type == filterType {
						filtered = append(filtered, match)
					}
				}
				matches = filtered
			}
			if len(matches) > 0 {
				fileResults = append(fileResults, output.LineMatch{
					Line:    lineNum + 1,
					Content: line,
					Matches: matches,
				})
			}
		}

		if len(fileResults) > 0 {
			results[file.RelativePath] = fileResults
		}
	}

	duration := time.Since(start).Seconds()

	// Build stats
	var statsData *output.Stats
	if showStats || format == "json-verbose" {
		tradCount, simpCount := 0, 0
		for _, lines := range results {
			for _, lm := range lines {
				for _, m := range lm.Matches {
					if m.Type == classifier.Traditional {
						tradCount++
					} else if m.Type == classifier.Simplified {
						simpCount++
					}
				}
			}
		}
		statsData = &output.Stats{
			ScannedFiles:     scannedCount,
			MatchedFiles:     len(results),
			TraditionalCount: tradCount,
			SimplifiedCount:  simpCount,
			Duration:         fmt.Sprintf("%.2fs", duration),
		}
	}

	// Format output
	fmt := output.New(!noColor, format)
	return fmt.Format(cmd.OutOrStdout(), results, statsData)
}
