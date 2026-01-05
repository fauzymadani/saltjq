package printer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
)

// PrintValue prints a value with pretty JSON and color according to style.
// Returns an error if writing to the provided writer fails.
func PrintValue(w io.Writer, v interface{}, style Style) error {
	var buf strings.Builder

	// For primitives, build colored value string
	switch vv := v.(type) {
	case nil:
		buf.WriteString("null")
	case bool:
		if style.NoColor {
			buf.WriteString(fmt.Sprint(vv))
		} else {
			buf.WriteString(fmt.Sprintf("%s%v%s", style.BoolColor, vv, style.Reset))
		}
	case string:
		if style.NoColor {
			buf.WriteString(fmt.Sprintf("%q", vv))
		} else {
			buf.WriteString(fmt.Sprintf("%s%q%s", style.StringColor, vv, style.Reset))
		}
	case json.Number:
		if style.NoColor {
			buf.WriteString(vv.String())
		} else {
			buf.WriteString(fmt.Sprintf("%s%s%s", style.NumberColor, vv.String(), style.Reset))
		}
	case float64, int, int64:
		if style.NoColor {
			buf.WriteString(fmt.Sprint(vv))
		} else {
			buf.WriteString(fmt.Sprintf("%s%v%s", style.NumberColor, vv, style.Reset))
		}
	default:
		// For composite types, marshal with indentation then colorize keys/strings/numbers
		js, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			// fallback to default fmt
			buf.WriteString(fmt.Sprint(v))
		} else {
			out := string(js)
			if style.NoColor {
				buf.WriteString(out)
			} else {
				// Colorize keys: "key":
				keyRe := regexp.MustCompile(`"([^"\\]+)":`)
				out = keyRe.ReplaceAllString(out, style.KeyColor+`"$1"`+style.Reset+":")

				// Colorize string values (naive): : "value"
				strRe := regexp.MustCompile(`: "([^"\\]*)"`)
				out = strRe.ReplaceAllString(out, ": "+style.StringColor+`"$1"`+style.Reset)

				// Colorize numbers (naive)
				numRe := regexp.MustCompile(`: ([-]?[0-9]+\.?[0-9]*)`)
				out = numRe.ReplaceAllString(out, ": "+style.NumberColor+`$1`+style.Reset)

				buf.WriteString(out)
			}
		}
	}

	// write once to provided writer
	if _, err := io.WriteString(w, buf.String()); err != nil {
		return err
	}
	return nil
}

// PrintTable prints an array of objects as a simple ASCII table.
// Returns an error if writing to the provided writer fails.
func PrintTable(w io.Writer, arr []interface{}, style Style) error {
	var buf strings.Builder
	if len(arr) == 0 {
		buf.WriteString("(empty)\n")
		_, err := io.WriteString(w, buf.String())
		return err
	}
	// collect keys
	keySet := map[string]struct{}{}
	rows := make([]map[string]string, 0, len(arr))
	for _, it := range arr {
		m, ok := it.(map[string]interface{})
		row := map[string]string{}
		if ok {
			for k, v := range m {
				keySet[k] = struct{}{}
				// marshal value to compact string
				b, err := json.Marshal(v)
				if err != nil {
					row[k] = fmt.Sprintf("<marshal error: %v>", err)
				} else {
					row[k] = string(b)
				}
			}
		}
		rows = append(rows, row)
	}
	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// compute widths
	widths := make([]int, len(keys))
	for i, k := range keys {
		widths[i] = len(k)
	}
	for _, r := range rows {
		for i, k := range keys {
			if len(r[k]) > widths[i] {
				widths[i] = len(r[k])
			}
		}
	}

	// print header
	var header bytes.Buffer
	header.WriteString("+")
	for i := range keys {
		header.WriteString(strings.Repeat("-", widths[i]+2))
		header.WriteString("+")
	}
	header.WriteString("\n")
	buf.WriteString(header.String())
	// header row
	buf.WriteString("|")
	for _, k := range keys {
		if style.NoColor {
			buf.WriteString(fmt.Sprintf(" %s ", k))
		} else {
			buf.WriteString(fmt.Sprintf(" %s%v%s ", style.KeyColor, k, style.Reset))
		}
		buf.WriteString("|")
	}
	buf.WriteString("\n")
	buf.WriteString(header.String())

	// rows
	for _, r := range rows {
		buf.WriteString("|")
		for i, k := range keys {
			cell := r[k]
			buf.WriteString(fmt.Sprintf(" %s ", cell))
			// pad
			pad := widths[i] - len(cell)
			if pad > 0 {
				buf.WriteString(strings.Repeat(" ", pad))
			}
			buf.WriteString("|")
		}
		buf.WriteString("\n")
	}
	buf.WriteString(header.String())

	_, err := io.WriteString(w, buf.String())
	return err
}
