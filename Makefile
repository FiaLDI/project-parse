APP := project-parser
MODULE := github.com/FiaLDI/project-parse
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo 0.1.0-dev)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X $(MODULE)/internal/version.Version=$(VERSION) \
	-X $(MODULE)/internal/version.Commit=$(COMMIT) \
	-X $(MODULE)/internal/version.BuildDate=$(BUILD_DATE)

.PHONY: build test vet fmt tidy run-version run-doctor

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(APP) ./cmd/project-parser

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

run-version: build
	./bin/$(APP) version

run-doctor: build
	./bin/$(APP) doctor --config configs/parser.yaml
