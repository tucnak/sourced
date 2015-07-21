package parser

import "strconv"

// Argument represents a single command argument.
type Argument interface {
    // Actual representation, compatible with Source engine.
    Source() string
}

// NumberArgument represents an integer argument, like 42.
type NumberArgument struct {
    value int
}

// Source returns a number which represents a corresponding argument.
func (arg NumberArgument) Source() string {
    return strconv.Itoa(arg.value)
}

// WordArgument represents a single word argument (like "one_two")
type WordArgument struct {
    value string
}

// Source returns a single word argument.
func (arg WordArgument) Source() string {
    return arg.value
}

// StringArgument represents a string argument, like "hey, im johny"
type StringArgument struct {
    value string
}

// Source returns a string of corresponding argument
func (arg StringArgument) Source() string {
    return "\""+arg.value+"\""
}
