package printer

// Style holds ANSI color codes for token types. Minimal for MVP.
type Style struct {
	KeyColor    string
	StringColor string
	NumberColor string
	BoolColor   string
	PunctColor  string
	Reset       string
	NoColor     bool
}

// GetStyle returns a style by name. Default "clean".
func GetStyle(name string) Style {
	switch name {
	case "dev":
		return Style{KeyColor: "\x1b[94m", StringColor: "\x1b[92m", NumberColor: "\x1b[93m", BoolColor: "\x1b[91m", PunctColor: "\x1b[90m", Reset: "\x1b[0m"}
	case "viz":
		return Style{KeyColor: "\x1b[36m", StringColor: "\x1b[36m", NumberColor: "\x1b[33m", BoolColor: "\x1b[35m", PunctColor: "\x1b[90m", Reset: "\x1b[0m"}
	default:
		return Style{KeyColor: "\x1b[36m", StringColor: "\x1b[32m", NumberColor: "\x1b[33m", BoolColor: "\x1b[35m", PunctColor: "\x1b[90m", Reset: "\x1b[0m"}
	}
}
