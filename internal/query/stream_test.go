package query

import (
	"bytes"
	"testing"
)

func TestStreamJSONArray(t *testing.T) {
	r := bytes.NewBufferString(`[{"a":1},{"a":2}]`)
	items, errs := StreamJSON(r)
	got := make([]interface{}, 0)
	for it := range items {
		got = append(got, it)
	}
	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}
}

func TestStreamJSONNDJSON(t *testing.T) {
	r := bytes.NewBufferString("{\"a\":1}\n{\"a\":2}\n")
	items, errs := StreamJSON(r)
	count := 0
	for range items {
		count++
	}
	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 items, got %d", count)
	}
}
