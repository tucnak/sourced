package parser

import "fmt"

// Statement is a command with some arguments
type Statement struct {
	Command   string
	Arguments []Argument

	// Line in source
	line int
}

func (s Statement) String() (string, error) {
	self := s.Command

	if self == "with" {
		return generateForWith(&s)
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

func (s Statement) Undo() (string, error) {
	self := s.Command

	if self == "bind" {
		unbind, err := s.Arguments[0].String()
		if err != nil {
			return "", err
		}

		return "unbind " + unbind, nil
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

func generateForWith(st *Statement) (string, error) {
	if len(st.Arguments) != 2 {
		return "", &TypeError{
			Line:    st.line,
			Message: "Wrong amount of arguments for `with` (requires two)",
		}
	}

	if _, ok := st.Arguments[0].(*Sequence); ok {
		return "", &TypeError{
			Line:    st.line,
			Message: "First argument of `with` must be a keyboard key",
		}
	}

	key, err := st.Arguments[0].String()
	if err != nil {
		return "", err
	}

	seq, ok := st.Arguments[1].(*Sequence)
	if !ok {
		return "", &TypeError{
			Line:    st.line,
			Message: "Second argument of `with` must be a scope",
		}
	}

	enable, err := seq.String()
	if err != nil {
		return "", err
	}

	disable, err := seq.Undo()
	if err != nil {
		return "", err
	}

	alias_plus := fmt.Sprintf("alias +with_%s %s", key, enable)
	alias_minus := fmt.Sprintf("alias -with_%s %s", key, disable)
	bind := fmt.Sprintf("bind %s +with_%s", key, key)

	return fmt.Sprintf("%s\n%s\n%s", alias_plus, alias_minus, bind), nil
}
