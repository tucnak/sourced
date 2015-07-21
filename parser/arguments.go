package parser

import "strconv"

// Argument represents a single command argument
type Argument interface {
    String() string
}

type NumberArgument struct {
    value int
}

func (arg NumberArgument) String() string {
    return strconv.Itoa(arg.value)
}

type WordArgument struct {
    value string
}

func (arg WordArgument) String() string {
    return arg.value
}

type StringArgument struct {
    value string
}

func (arg StringArgument) String() string {
    return "\""+arg.value+"\""
}
