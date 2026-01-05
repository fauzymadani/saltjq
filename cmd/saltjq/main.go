package main

import (
	"flag"
	"io"
	"os"

	"myJsonParser/internal/errs"
	"myJsonParser/internal/printer"
	"myJsonParser/internal/query"
)

func main() {
	expr := flag.String("e", ".", "expression to run (simple subset: .field, .field1.field2, .[], pipes with |)")
	stream := flag.Bool("s", false, "stream top-level array elements")
	table := flag.Bool("table", false, "format arrays of objects as table")
	styleName := flag.String("style", "clean", "output style: clean|dev|viz")
	noColor := flag.Bool("no-color", false, "disable color output")
	flag.Parse()

	// Input source: optional file argument or stdin
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

	data, err := query.ReadAllJSON(r)
	if err != nil {
		errs.Handle(errs.Wrap(err, 2, "failed to read json"))
	}

	style := printer.GetStyle(*styleName)
	if *noColor {
		style.NoColor = true
	}

	// Evaluate expression
	results, err := query.Eval(data, *expr, *stream)
	if err != nil {
		errs.Handle(errs.Wrap(err, 3, "evaluation error"))
	}

	w := os.Stdout
	if *table {
		// If a single array-of-objects returned, print as table
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
		// separate multiple results with newline
		if i > 0 {
			if _, err := io.WriteString(w, "\n"); err != nil {
				errs.Handle(errs.Wrap(err, 4, "write error"))
			}
		}
		if err := printer.PrintValue(w, v, style); err != nil {
			errs.Handle(errs.Wrap(err, 4, "write error"))
		}
		if _, err := io.WriteString(w, "\n"); err != nil {
			errs.Handle(errs.Wrap(err, 4, "write error"))
		}
	}
}
