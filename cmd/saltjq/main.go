package main

import (
	"flag"
	"io"
	"os"

	"myJsonParser/internal/errs"
	"myJsonParser/internal/printer"
	"myJsonParser/internal/query"
)

// helper to print a single value. If leadNL is true, write a leading newline before the value.
func emitValue(w io.Writer, v interface{}, style printer.Style, raw bool, leadNL bool) error {
	if leadNL {
		if _, err := io.WriteString(w, "\n"); err != nil {
			return err
		}
	}
	if raw {
		if s, ok := v.(string); ok {
			if _, err := io.WriteString(w, s); err != nil {
				return err
			}
			return nil
		}
	}
	return printer.PrintValue(w, v, style)
}

func main() {
	expr := flag.String("e", ".", "expression to run (simple subset: .field, .field1.field2, .[], pipes with |)")
	stream := flag.Bool("s", false, "stream top-level array elements")
	table := flag.Bool("table", false, "format arrays of objects as table")
	styleName := flag.String("style", "clean", "output style: clean|dev|viz")
	noColor := flag.Bool("no-color", false, "disable color output")
	raw := flag.Bool("r", false, "raw output for strings (no JSON quotes)")
	bufferSize := flag.Int("buffer-size", query.StreamBufferSize, "buffer size for streaming items channel")
	flag.Parse()

	// apply runtime buffer-size if provided
	if bufferSize != nil && *bufferSize > 0 {
		query.StreamBufferSize = *bufferSize
	}

	var r io.Reader
	if flag.NArg() > 0 {
		fpath := flag.Arg(0)
		f, err := os.Open(fpath)
		if err != nil {
			errs.Handle(errs.Wrap(err, 2, "failed to open file"))
		}
		defer func() {
			if err := f.Close(); err != nil {
				errs.Handle(errs.Wrap(err, 5, "failed to close file"))
			}
		}()
		r = f
	} else {
		r = os.Stdin
	}

	style := printer.GetStyle(*styleName)
	// if user explicitly asked to disable color, respect it
	if *noColor {
		style.NoColor = true
	} else {
		// auto-disable color when stdout is not a terminal
		if fi, err := os.Stdout.Stat(); err == nil {
			// if not a character device, assume non-tty (piped/redirected)
			if fi.Mode()&os.ModeCharDevice == 0 {
				style.NoColor = true
			}
		}
	}

	w := os.Stdout

	if *stream {
		items, errc := query.StreamJSON(r)
		firstOut := true
		for {
			select {
			case item, ok := <-items:
				if !ok {
					// items closed; check err channel
					if err := <-errc; err != nil {
						errs.Handle(errs.Wrap(err, 2, "failed to read json"))
					}
					return
				}
				results, err := query.Eval(item, *expr, true)
				if err != nil {
					errs.Handle(errs.Wrap(err, 3, "evaluation error"))
				}
				for i, v := range results {
					if err := emitValue(w, v, style, *raw, !firstOut || i > 0); err != nil {
						errs.Handle(errs.Wrap(err, 4, "write error"))
					}
					firstOut = false
				}
			case err := <-errc:
				if err != nil {
					errs.Handle(errs.Wrap(err, 2, "failed to read json"))
				}
				return
			}
		}
	}

	data, err := query.ReadAllJSON(r)
	if err != nil {
		errs.Handle(errs.Wrap(err, 2, "failed to read json"))
	}

	results, err := query.Eval(data, *expr, false)
	if err != nil {
		errs.Handle(errs.Wrap(err, 3, "evaluation error"))
	}

	if *table {
		if len(results) == 1 {
			if arr, ok := results[0].([]interface{}); ok {
				if err := printer.PrintTable(w, arr, style); err != nil {
					errs.Handle(errs.Wrap(err, 4, "write error"))
				}
				return
			}
		}
	}

	for i, v := range results {
		if err := emitValue(w, v, style, *raw, i > 0); err != nil {
			errs.Handle(errs.Wrap(err, 4, "write error"))
		}
		if _, err := io.WriteString(w, "\n"); err != nil {
			errs.Handle(errs.Wrap(err, 4, "write error"))
		}
	}
}
