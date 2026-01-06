package query

import (
	"bufio"
	"encoding/json"
	"io"
)

// StreamBufferSize controls the buffer size of the items channel returned by StreamJSON.
// It can be adjusted by callers/tests for tuning.
var StreamBufferSize = 16

// StreamJSON decodes input as a stream of JSON values. It supports a top-level
// array (decodes each element) or multiple concatenated/NDJSON values.
func StreamJSON(r io.Reader) (<-chan interface{}, <-chan error) {
	items := make(chan interface{}, StreamBufferSize)
	errs := make(chan error, 1)
	go func() {
		defer close(items)
		defer close(errs)

		br := bufio.NewReader(r)
		// find first non-space byte
		var first byte
		for {
			b, err := br.ReadByte()
			if err != nil {
				if err == io.EOF {
					return
				}
				errs <- err
				return
			}
			if b == ' ' || b == '\n' || b == '\t' || b == '\r' {
				continue
			}
			first = b
			_ = br.UnreadByte()
			break
		}

		dec := json.NewDecoder(br)
		dec.UseNumber()

		if first == '[' {
			// consume '[' token
			if _, err := dec.Token(); err != nil {
				errs <- err
				return
			}
			for dec.More() {
				var v interface{}
				if err := dec.Decode(&v); err != nil {
					errs <- err
					return
				}
				items <- v
			}
			// consume closing ']' token
			if _, err := dec.Token(); err != nil {
				errs <- err
			}
			return
		}

		// otherwise decode successive values (NDJSON or concatenated JSON)
		for {
			var v interface{}
			if err := dec.Decode(&v); err != nil {
				if err == io.EOF {
					return
				}
				errs <- err
				return
			}
			items <- v
		}
	}()
	return items, errs
}
