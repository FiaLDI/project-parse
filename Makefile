APP := project-parser
MODULE := github.com/FiaLDI/project-parse
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo 0.1.0-dev)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X $(MODULE)/internal/version.Version=$(VERSION) \
	-X $(MODULE)/internal/version.Commit=$(COMMIT) \
	-X $(MODULE)/internal/version.BuildDate=$(BUILD_DATE)
CONFIG ?= configs/parser.yaml
GO ?= go

.PHONY: build test vet fmt tidy clean install run-version run-doctor run-scan run-report run-graph demo-go demo-node

build:
	$(GO) build -ldflags "$(LDFLAGS)" -o bin/$(APP) ./cmd/project-parser

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

fmt:
	$(GO) fmt ./...

tidy:
	$(GO) mod tidy

clean:
	rm -rf bin/ .project-parser/

install: build
	install -d $(DESTDIR)/usr/local/bin
	install -m 755 bin/$(APP) $(DESTDIR)/usr/local/bin/$(APP)

run-version: build
	./bin/$(APP) version

run-doctor: build
	./bin/$(APP) doctor --config $(CONFIG)

run-scan: build
	./bin/$(APP) scan . --config $(CONFIG)

run-report: build
	./bin/$(APP) report . --config $(CONFIG)

run-graph: build
	./bin/$(APP) graph . --config $(CONFIG)

config-init: build
	./bin/$(APP) config init --config $(CONFIG)

demo-go: build
	./bin/$(APP) report examples/sample-go --config $(CONFIG)
	./bin/$(APP) graph examples/sample-go --config $(CONFIG)

demo-node: build
	./bin/$(APP) report examples/sample-node --config $(CONFIG)
	./bin/$(APP) graph examples/sample-node --config $(CONFIG)
