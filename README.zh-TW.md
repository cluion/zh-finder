# zh-finder

> 掃描檔案中中文字（繁體/簡體）的 CLI 工具

[English](README.md)

## 特色

- 掃描檔案中的中文字（繁體/簡體）
- Lipgloss 樣式化的彩色輸出
- 多種輸出格式（終端機、JSON、JSON-Compact）
- 彈性的排除/包含規則
- 統計資訊顯示

## 安裝

### Homebrew

```bash
brew install cluion/tap/zh-finder
```

### Go Install

```bash
go install github.com/cluion/zh-finder/cmd/zh-finder@latest
```

### 從原始碼建置

```bash
git clone https://github.com/cluion/zh-finder.git
cd zh-finder
make install
```

## 使用方式

```bash
# 基本掃描
zh-finder scan ./src

# 顯示統計資訊
zh-finder scan ./src --stats

# 篩選副檔名
zh-finder scan ./src --ext=go,js,ts

# JSON 輸出
zh-finder scan ./src --format=json

# 篩選繁體/簡體
zh-finder scan ./src --type=traditional

# 排除目錄
zh-finder scan ./src --exclude=node_modules,dist

# 顯示版本
zh-finder version
```

## 參數

| 參數 | 預設值 | 說明 |
|------|--------|------|
| `--ext` | 全部 | 只掃描指定副檔名（逗號分隔） |
| `--exclude` | 內建 | 排除目錄（逗號分隔） |
| `--exclude-add` | | 額外排除的目錄 |
| `--no-exclude` | false | 停用預設排除規則 |
| `--format` | term | 輸出格式：term、json、json-compact |
| `--stats` | false | 顯示統計資訊 |
| `--max-depth` | 無限制 | 最大遞迴深度 |
| `--binary` | false | 掃描二進位檔案 |
| `--no-color` | false | 停用彩色輸出 |
| `--type` | all | 篩選：all、traditional、simplified |

## 開發

```bash
make build          # 建置執行檔
make test           # 執行測試（含覆蓋率）
make lint           # 執行 go vet
make clean          # 清理建置產物
```

## 授權

[MIT License](LICENSE)

## 致謝

- 繁簡字資料來自 [OpenCC](https://github.com/BYVoid/OpenCC)
