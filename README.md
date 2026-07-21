# project-parser

Cross-platform CLI that scans a codebase, detects stack / architecture / infrastructure via plugins, and emits structured reports.

Single static binary — no runtime dependencies beyond the scanned project itself.

## Requirements

- Go 1.24+

## Quick start

```bash
make build
make run-doctor
make run-report
make run-graph
```

Or without `make`:

```bash
go build -o bin/project-parser ./cmd/project-parser
./bin/project-parser report . --config configs/parser.yaml
./bin/project-parser graph . --config configs/parser.yaml
```

## Commands

| Command | Description |
|---------|-------------|
| `scan [path]` | Walk the tree and print file index summary |
| `report [path]` | Full analysis → JSON / Markdown / HTML |
| `graph [path]` | Architecture graph → SVG or JSON |
| `plugins` | List registered plugins |
| `version` | Print build version |
| `config init` | Create default `parser.yaml` |
| `doctor` | Environment and config diagnostics |

### Examples

```bash
# Scan current directory
project-parser scan . --config configs/parser.yaml

# Generate report (default: json + markdown)
project-parser report examples/sample-go --config configs/parser.yaml --format json,markdown,html

# Architecture graph as SVG
project-parser graph examples/sample-node --config configs/parser.yaml --format svg

# Graph as JSON
project-parser graph . --format json --out ./.project-parser
```

```bash
project-parser config init              # creates ./parser.yaml
project-parser config init --force      # overwrite existing file
project-parser config init --out ./cfg/parser.yaml
```

Output files are written to `.project-parser/` by default:

```text
.project-parser/
  report.json
  report.md
  report.html
  graph.svg
  graph.json
```

## Configuration

Copy [`configs/parser.yaml`](configs/parser.yaml) to your project root as `parser.yaml`, or pass `--config`.

Key options:

- `plugins.enabled` / `plugins.disabled` — control which analyzers run
- `scan.exclude` — glob patterns skipped during walk
- `scan.jobs` — worker pool size (`0` = NumCPU)
- `report.formats` — default output formats
- `graph.format` — default graph format (`svg`)

## Architecture

Clean Architecture layout:

```text
cmd/           → Cobra CLI (no business logic)
internal/app/  → use-cases (scan, report, graph)
internal/ports/→ interfaces (Plugin, Scanner, Renderer, …)
internal/domain/→ entities (Report, Finding, Graph)
internal/plugins/→ stack-specific analyzers
internal/scanner/→ filesystem walk + index
internal/analyzer/→ parallel plugin runner
internal/report/ → aggregation
internal/graph/  → graph builder
internal/output/ → JSON, Markdown, HTML, SVG renderers
```

Plugins are independent — the core never imports language-specific logic directly; plugins self-register via `internal/plugins/register.go`.

## Plugins (built-in)

`node`, `python`, `golang`, `rust`, `java`, `docker`, `git`, `githubactions`, `kubernetes`, `database`, `architecture`, `documentation`

## Sample projects

- [`examples/sample-go`](examples/sample-go) — Go module with Clean Architecture folders
- [`examples/sample-node`](examples/sample-node) — Node.js + TypeScript + React/Next

```bash
make demo-go
make demo-node
```

## Development

```bash
make test
make vet
make fmt
make tidy
```

## Roadmap

- PDF export (renderer stub exists)
- Tree-sitter / LSP integration
- GitHub / GitLab API
- Mermaid / PlantUML / Graphviz export
- Web UI / REST API

## License

See [LICENSE](LICENSE).
