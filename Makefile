VERSION ?= dev
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo "unknown")
BINARY  := bin/zh-finder
LDFLAGS := -X github.com/cluion/zh-finder/internal/cli.version=$(VERSION) \
           -X github.com/cluion/zh-finder/internal/cli.commit=$(COMMIT) \
           -X github.com/cluion/zh-finder/internal/cli.date=$(DATE)

.PHONY: build test lint clean install uninstall

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/zh-finder/

test:
	go test -race -cover ./...

lint:
	go vet ./...

clean:
	rm -rf bin/ dist/

install: build
	cp $(BINARY) /usr/local/bin/zh-finder

uninstall:
	@rm -f /usr/local/bin/zh-finder
	@echo "zh-finder has been uninstalled from /usr/local/bin/zh-finder"
