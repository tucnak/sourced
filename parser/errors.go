package parser

import "fmt"

// ParseError occurs when Parse can't build an AST of provided source.
type ParseError struct {
	Line, Position int
	Message        string
}

// Error returns a text representation of ParseError
func (err ParseError) Error() string {
	return fmt.Sprintf("%s at line %d:%d",
		err.Message, err.Line, err.Position)
}

// TypeError occurs when input source is not valid Source script.
type TypeError struct {
	Line    int
	Message string
}

// Error returns a text representation of TypeError
func (err TypeError) Error() string {
	return fmt.Sprintf("%s at line %d",
		err.Message, err.Line)
}
