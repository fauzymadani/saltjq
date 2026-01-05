Saltjq — small jq-like MVP

A minimal prototype CLI for querying JSON with streaming-aware decoding and colored pretty printing.

Goals
- Fast, small, and safe for pipelines
- Familiar jq-like pipe expressions (small subset for MVP)
- Colored pretty-print and a simple table mode

Build

Run the build target in the repository root (requires Go toolchain):

```bash
make build
```

or directly with go:

```bash
go build -o saltjq ./cmd/saltjq
```

Makefile targets

| Target  | Description |
|---------|-------------|
| build   | build the saltjq binary |
| fmt     | run gofmt -w . |
| test    | run go test ./... |
| vet     | run go vet ./... |
| run     | build and run the binary (./saltjq) |
| install | install binary to GOBIN |
| clean   | remove the built binary |

Flags

| Flag | Short | Description |
|------|-------|-------------|
| --expr | -e | Expression to run (subset supported, e.g. `.users[] | .name`) |
| --stream | -s | Stream top-level array elements (NDJSON or large arrays) |
| --table |  | Format arrays of objects as a table |
| --style |  | Choose output style: `clean`, `dev`, `viz` |
| --no-color |  | Disable color output (useful when piping to non-TTY) |

Examples

Pretty-print a field from the sample JSON:

```bash
cat testdata/sample.json | ./saltjq -e '.users[] | .name'
```

Print an array-of-objects as a table:

```bash
./saltjq -e '.users' --table testdata/sample.json
```

Build and run via Makefile:

```bash
make build
./saltjq -e '.users[] | .name' testdata/sample.json
```

Notes
- This is an early MVP. The expression language is intentionally small. Improvements (streaming on large files without buffering, more builtins, better table formatting) are planned.
