VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

.PHONY: build test lint clean install

build:
	go build -ldflags "$(LDFLAGS)" -o portpilot ./cmd/portpilot

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/portpilot

test:
	go test -v -race ./...

lint:
	golangci-lint run

clean:
	rm -f portpilot
