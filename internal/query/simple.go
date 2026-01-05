package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ReadAllJSON reads all JSON from r. It decodes a single JSON value.
func ReadAllJSON(r io.Reader) (interface{}, error) {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	var v interface{}
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}

// Eval evaluates a very small expression language over the JSON value.
// Supported forms (MVP):
//   - "." root
//   - ".field" nested via dots
//   - ".field[]" iterate over array in field
//   - ".[]" iterate over root array
//   - pipes: expr | expr
func Eval(root interface{}, expr string, stream bool) ([]interface{}, error) {
	stages := parsePipe(expr)
	if len(stages) == 0 {
		return nil, errors.New("empty expression")
	}

	current := []interface{}{root}
	for _, st := range stages {
		var next []interface{}
		for _, item := range current {
			res, err := applyStage(item, st)
			if err != nil {
				return nil, err
			}
			// flatten
			for _, x := range res {
				next = append(next, x)
			}
		}
		current = next
	}
	return current, nil
}

func parsePipe(expr string) []string {
	parts := strings.Split(expr, "|")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// applyStage applies a single stage like ".a.b" or ".a[]" or ".[]"
func applyStage(v interface{}, stage string) ([]interface{}, error) {
	if stage == "." {
		return []interface{}{v}, nil
	}
	if !strings.HasPrefix(stage, ".") {
		return nil, fmt.Errorf("unsupported stage: %s", stage)
	}
	// Trim leading dot
	body := stage[1:]
	// handle [] suffix
	if body == "[]" {
		// iterate over v if array
		if arr, ok := v.([]interface{}); ok {
			out := make([]interface{}, 0, len(arr))
			for _, e := range arr {
				out = append(out, e)
			}
			return out, nil
		}
		return nil, nil // no results
	}

	// field access possibly with [] suffix
	if strings.HasSuffix(body, "[]") {
		field := body[:len(body)-2]
		val, err := getByPath(v, field)
		if err != nil || val == nil {
			return nil, nil
		}
		if arr, ok := val.([]interface{}); ok {
			out := make([]interface{}, 0, len(arr))
			for _, e := range arr {
				out = append(out, e)
			}
			return out, nil
		}
		return nil, nil
	}

	// simple field or dotted path
	val, err := getByPath(v, body)
	if err != nil {
		return nil, err
	}
	return []interface{}{val}, nil
}

// getByPath supports dotted keys: a.b.c
func getByPath(v interface{}, path string) (interface{}, error) {
	if path == "" {
		return v, nil
	}
	parts := strings.Split(path, ".")
	cur := v
	for _, p := range parts {
		switch m := cur.(type) {
		case map[string]interface{}:
			val, ok := m[p]
			if !ok {
				return nil, nil
			}
			cur = val
		default:
			return nil, nil
		}
	}
	return cur, nil
}
