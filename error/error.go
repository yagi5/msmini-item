package error

import (
	"bytes"
	"fmt"
)

const (
	// EINTERNAL is internal error
	EINTERNAL code = "internal"

	// EINVALID is invalid request error
	EINVALID code = "invalid"

	// ENOTFOUND is not found error
	ENOTFOUND code = "not_found"
)

// Code is machine-readable error code
type code string

// Error represents application error
type Error struct {
	Code    code
	Message string
	Op      string
	Err     error
}

// Error returns the string representation of the error message.
func (e *Error) Error() string {
	var buf bytes.Buffer

	if e.Op != "" {
		fmt.Fprintf(&buf, "%s: ", e.Op)
	}

	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		if e.Code != "" {
			fmt.Fprintf(&buf, "<%s> ", e.Code)
		}
		buf.WriteString(e.Message)
	}
	return buf.String()
}

// Code returns the code of the root error, if available. Otherwise returns EINTERNAL.
func Code(err error) string {
	if err == nil {
		return ""
	} else if e, ok := err.(*Error); ok && e.Code != "" {
		return string(e.Code)
	} else if ok && e.Err != nil {
		return Code(e.Err)
	}
	return string(EINTERNAL)
}

// Message returns the human-readable message of the error, if available.
func Message(err error) string {
	if err == nil {
		return ""
	} else if e, ok := err.(*Error); ok && e.Message != "" {
		return e.Message
	} else if ok && e.Err != nil {
		return Message(e.Err)
	}
	return "An internal error has occurred"
}
