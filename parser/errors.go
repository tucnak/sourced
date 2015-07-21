package parser

import "fmt"

type ParseError struct {
    line, position int
    message string
}

func (err ParseError) Error() string {
    return fmt.Sprintf("%s at line %d:%d",
        err.message, err.line, err.position)
}
