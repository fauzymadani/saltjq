Saltjq — small jq-like MVP

A minimal prototype CLI for querying JSON with streaming-aware decoding and colored pretty printing.

Goals
- Fast, small, and safe for pipelines
- Familiar jq-like pipe expressions (small subset for MVP)
- Colored pretty-print and a simple table mode

Build

Run the build target in the repository root (requires Go toolchain):

```bash
make build  # builds ./bin/saltjq
```

or directly with go:

```bash
go build -o saltjq ./cmd/saltjq
```

Makefile targets

| Target  | Description |
|---------|-------------|
| build   | build the saltjq binary (into ./bin) |
| fmt     | run gofmt -w . |
| test    | run go test ./... |
| vet     | run go vet ./... |
| run     | build and run the binary (./bin/saltjq) |
| install | install binary to GOBIN |
| clean   | remove the built binary |

Flags

| Flag | Short | Description |
|------|-------|-------------|
| --expr | -e | Expression to run (subset supported, e.g. `.users[] | .name`) |
| --stream | -s | Stream top-level array elements (supports NDJSON or large arrays; values are decoded and evaluated one-by-one) |
| --raw | -r | Raw output for strings (no JSON quotes), useful for shell pipelines |
| --buffer-size |  | Buffer size for streaming items channel (default 16) |
| --table |  | Format arrays of objects as a table (note: `--table` expects a full array result; it is not automatically compatible with `-s` streaming mode unless you collect items) |
| --style |  | Choose output style: `clean`, `dev`, `viz` |
| --no-color |  | Disable color output (useful when piping to non-TTY) |

Examples

Pretty-print a field from the sample JSON (non-streaming):

```bash
./bin/saltjq -e '.users[] | .name' testdata/sample.json
```

Stream a large array or NDJSON and print a field as items arrive:

```bash
cat testdata/sample.json | ./bin/saltjq -s -e '.users[] | .name'
```

Control the stream buffer size (example: increase to 64):

```bash
cat testdata/sample.json | ./bin/saltjq -s --buffer-size 64 -e '.users[] | .name'
```

Print raw strings (no JSON quotes):

```bash
./bin/saltjq -e '.users[] | .name' -r testdata/sample.json
cat testdata/sample.json | ./bin/saltjq -s -e '.users[] | .name' -r
```

Notes
- This is an early MVP. The expression language is intentionally small. Improvements (streaming on large files without buffering, more builtins, better table formatting) are planned.
