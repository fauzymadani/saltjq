package errs

import (
	"errors"
	"fmt"
	"os"
)

// CLIError represents an error with an associated exit code and optional message.
type CLIError struct {
	Err     error
	Code    int
	Message string
}

func (e *CLIError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err == nil {
		return e.Message
	}
	if e.Message == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *CLIError) Unwrap() error { return e.Err }

// New creates a CLIError with no wrapped error.
func New(code int, msg string) error {
	return &CLIError{Err: nil, Code: code, Message: msg}
}

// Wrap wraps an existing error with a CLI exit code and message. If err is nil returns nil.
func Wrap(err error, code int, msg string) error {
	if err == nil {
		return nil
	}
	// If err already has a code, preserve the inner code unless overridden by non-zero code
	var inner *CLIError
	if errors.As(err, &inner) {
		// If caller provided code 0, keep inner code
		if code == 0 {
			code = inner.Code
		}
	}
	return &CLIError{Err: err, Code: code, Message: msg}
}

// Wrapf wraps err with formatted message.
func Wrapf(err error, code int, format string, a ...interface{}) error {
	return Wrap(err, code, fmt.Sprintf(format, a...))
}

// ExitCode extracts the exit code from err (defaults to 1).
func ExitCode(err error) int {
	var e *CLIError
	if errors.As(err, &e) {
		if e.Code != 0 {
			return e.Code
		}
	}
	return 1
}

// Handle prints the error to stderr and exits with the proper exit code. If err is nil, does nothing.
func Handle(err error) {
	if err == nil {
		return
	}
	if _, writeErr := fmt.Fprintln(os.Stderr, err.Error()); writeErr != nil {
		// If we can't write the error message, exit with a fallback code
		os.Exit(1)
	}
	os.Exit(ExitCode(err))
}
