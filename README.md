# project-parser

Cross-platform CLI that scans a codebase, detects stack / architecture / infrastructure via plugins, and emits structured reports.

> Status: **Stage 0–1** — bootstrap, domain model, ports, Cobra CLI, config, registry. Scanner and plugins are not wired yet.

## Requirements

- Go 1.24+

## Build

```bash
make build
./bin/project-parser version
./bin/project-parser doctor --config configs/parser.yaml
```

## Commands

```text
project-parser scan .
project-parser report .
project-parser graph .
project-parser plugins
project-parser version
project-parser doctor
```

## Layout

See architecture discussion: Clean Architecture with `domain`, `ports`, `app`, plugins, and output renderers. Business logic stays out of `cmd/`.

## Config

Default example: [`configs/parser.yaml`](configs/parser.yaml)
