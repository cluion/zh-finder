package main

import (
	"embed"

	"github.com/cluion/zh-finder/internal/cli"
)

//go:embed data/traditional.txt data/simplified.txt
var dataFS embed.FS

func main() {
	cli.SetDataFS(dataFS)
	cli.Execute()
}
