package parser

import "strconv"

// Argument represents a single command argument.
type Argument interface {
	// Actual representation, compatible with Source engine.
	String() (string, error)
}

// NumberArgument represents an integer argument, like 42.
type NumberArgument struct {
	value int
}

// String returns a number which represents a corresponding argument.
func (arg NumberArgument) String() (string, error) {
	return strconv.Itoa(arg.value), nil
}

// WordArgument represents a single word argument (like "one_two")
type WordArgument struct {
	value string
}

// String returns a single word argument.
func (arg WordArgument) String() (string, error) {
	return arg.value, nil
}

// StringArgument represents a string argument, like "hey, im johny"
type StringArgument struct {
	value string
}

// String returns a string of corresponding argument
func (arg StringArgument) String() (string, error) {
	return "\"" + arg.value + "\"", nil
}
