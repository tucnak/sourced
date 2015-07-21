package parser

// Statement is a command with some arguments
type Statement struct {
	Command   string
	Arguments []Argument
}

func (s Statement) Source() string {
	self := s.Command

	for _, argument := range s.Arguments {
		self += " " + argument.Source()
	}

	return self
}

func (s *Statement) AddArgument(arg Argument) *Argument {
	s.Arguments = append(s.Arguments, arg)

	return &s.Arguments[len(s.Arguments)-1]
}
