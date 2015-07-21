package parser

import "fmt"

// ParseError occurs when Parse can't build an AST of provided source.
type ParseError struct {
	line, position int
	message        string
}

// Error returns a text representation of ParseError
func (err ParseError) Error() string {
	return fmt.Sprintf("%s at line %d:%d",
		err.message, err.line, err.position)
}
