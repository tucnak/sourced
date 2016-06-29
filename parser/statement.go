package parser

import (
	"fmt"
	"strings"
)

// Statement is a command with some arguments
type Statement struct {
	Command   string
	Arguments []Argument

	// Line in source
	line int
}

func (s *Statement) String() (string, error) {
	self := s.Command

	if self == "bind" {
		return generateForBind(s)
	}

	for _, argument := range s.Arguments {
		repr, err := argument.String()
		if err != nil {
			return "", err
		}

		self += " " + repr
	}

	return self, nil
}

func (s *Statement) Undo() (string, error) {
	self := s.Command

	if self == "bind" {
		unbind, err := s.Arguments[0].String()
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("unbind %s;bind %s bind_%s",
			unbind, unbind, keyRepresentation(unbind)), nil
	}

	if self == "alias" {
		unalias, err := s.Arguments[0].String()
		if err != nil {
			return "", err
		}

		return "unalias " + unalias, nil
	}

	return "", nil
}

func (s *Statement) AddArgument(arg Argument) *Argument {
	s.Arguments = append(s.Arguments, arg)

	return &s.Arguments[len(s.Arguments)-1]
}

func keyRepresentation(key string) string {
	key = strings.Trim(key, "\"'")

	if len(key) > 1 {
		return key
	}

	if repr, ok := Keys[rune(key[0])]; ok {
		return repr
	} else {
		return key
	}
}

func generateForBind(st *Statement) (string, error) {
	if len(st.Arguments) != 2 {
		return "", &TypeError{
			Line:    st.line,
			Message: "Wrong amount of arguments for `bind` (requires two)",
		}
	}

	if _, ok := st.Arguments[0].(*Sequence); ok {
		return "", &TypeError{
			Line:    st.line,
			Message: "First argument of `bind` must be a keyboard key",
		}
	}

	key, err := st.Arguments[0].String()
	if err != nil {
		return "", err
	}

	key_repr := keyRepresentation(key)

	clause, err := st.Arguments[1].String()
	if err != nil {
		return "", err
	}

	key_alias := fmt.Sprintf("alias bind_%s %s", key_repr, clause)
	bind := fmt.Sprintf("bind %s bind_%s", key, key_repr)

	return fmt.Sprintf("%s;%s", key_alias, bind), nil
}
