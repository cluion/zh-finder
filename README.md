# zh-finder

> A CLI tool to find Chinese characters (Traditional/Simplified) in files

[中文](README.zh-TW.md)

## Features

- Scan files for Chinese characters (Traditional/Simplified)
- Color-coded output with Lipgloss styling
- Multiple output formats (Terminal, JSON, JSON-Compact)
- Flexible exclude/include rules
- Statistics display

## Install

### Homebrew

```bash
brew install cluion/tap/zh-finder
```

### Go Install

```bash
go install github.com/cluion/zh-finder/cmd/zh-finder@latest
```

### Build from Source

```bash
git clone https://github.com/cluion/zh-finder.git
cd zh-finder
make install
```

## Usage

```bash
# Basic scan
zh-finder scan ./src

# With statistics
zh-finder scan ./src --stats

# Filter by file extensions
zh-finder scan ./src --ext=go,js,ts

# JSON output
zh-finder scan ./src --format=json

# Filter by type
zh-finder scan ./src --type=traditional

# Exclude directories
zh-finder scan ./src --exclude=node_modules,dist

# Show version
zh-finder version
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--ext` | all | Only scan specific extensions (comma-separated) |
| `--exclude` | built-in | Exclude directories (comma-separated) |
| `--exclude-add` | | Additional directories to exclude |
| `--no-exclude` | false | Disable default excludes |
| `--format` | term | Output format: term, json, json-compact |
| `--stats` | false | Show statistics |
| `--max-depth` | unlimited | Maximum recursion depth |
| `--binary` | false | Scan binary files |
| `--no-color` | false | Disable color output |
| `--type` | all | Filter: all, traditional, simplified |

## Development

```bash
make build          # Build binary
make test           # Run tests with coverage
make lint           # Run go vet
make clean          # Clean build artifacts
```

## License

[MIT License](LICENSE)

## Acknowledgments

- Traditional/Simplified character data from [OpenCC](https://github.com/BYVoid/OpenCC)
