package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// ParseError is returned for parsing reader errors.
// The first line is 1.
type ParseError struct {
	Err error // The actual error
}

func (e ParseError) Error() string {
	return fmt.Sprintf("%T %s", e.Err, e.Err)
}

// ErrorList represents an array of errors which is also an error itself.
type ErrorList []error

// Add appends err onto the ErrorList. Errors are kept in order.
func (e *ErrorList) Add(err error) {
	*e = append(*e, err)
}

// Err returns the first error (or nil).
func (e ErrorList) Err() error {
	if e == nil || len(e) == 0 {
		return nil
	}
	return e[0]
}

// Error implements the error interface
func (e ErrorList) Error() string {
	if len(e) == 0 {
		return "<nil>"
	}
	var buf bytes.Buffer
	e.Print(&buf)
	return buf.String()
}

// Print formats the ErrorList into a string written to w.
// If ErrorList contains multiple errors those after the first
// are indented.
func (e ErrorList) Print(w io.Writer) {
	if w == nil || len(e) == 0 {
		fmt.Fprintf(w, "<nil>")
		return
	}

	fmt.Fprintf(w, "%s", e[0])
	if len(e) > 1 {
		fmt.Fprintf(w, "\n")
	}

	for i := 1; i < len(e); i++ {
		fmt.Fprintf(w, "  %s", e[i])
		if i < len(e)-1 { // don't add \n to last error
			fmt.Fprintf(w, "\n")
		}
	}
}

// Empty no errors to return
func (e ErrorList) Empty() bool {
	return e == nil || len(e) == 0
}

// MarshalJSON marshals error list
func (e ErrorList) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Error())
}
