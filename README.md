Saltjq - small jq-like MVP

This is a minimal prototype CLI for JSON querying with streaming and colored pretty printing.

Build:

```bash
go build ./...
```

Example:

```bash
cat testdata/sample.json | ./saltjq -e '.users[] | .name'
```

Flags:
- -e expr : expression (simple subset)
- -s      : streaming (handles top-level arrays)
- --table : format array-of-objects as table
- --style : style name (clean|dev|viz)
- --no-color : disable color

