package parser

// Sequence represents a sequence of commands
type Sequence struct {
	Content []Statement
	Parent  *Sequence // nil for root sequence
}

// String returns an actual Source-compatible script source.
func (s Sequence) String() (string, error) {
	var content string

	if s.Parent != nil {
		// For inner scopes

		length := len(s.Content)
		if length > 0 {
			repr, err := s.Content[0].String()
			if err != nil {
				return "", err
			}

			content += repr
		}

		for i := 1; i < length; i++ {
			repr, err := s.Content[i].String()
			if err != nil {
				return "", err
			}
			content += ";" + repr
		}

		content = "\"" + content + "\""
	} else {
		// For outer scopes

		for _, statement := range s.Content {
			repr, err := statement.String()
			if err != nil {
				return "", err
			}

			content += repr + "\n"
		}
	}

	return content, nil
}

func (s Sequence) Undo() (string, error) {
	var content string

	length := len(s.Content)
	if length > 0 {
		repr, err := s.Content[0].Undo()
		if err != nil {
			return "", err
		}

		content += repr
	}

	for i := 1; i < length; i++ {
		repr, err := s.Content[i].Undo()
		if err != nil {
			return "", err
		}
		content += ";" + repr
	}

	content = "\"" + content + "\""

	return content, nil
}

// Last returns the last statement in the chain.
func (s *Sequence) Last() *Statement {
	if len(s.Content) == 0 {
		return nil
	}

	return &s.Content[len(s.Content)-1]
}

// Add pushes a new command (without arguments) back and returns its statement.
func (s *Sequence) Add(command string) *Statement {
	s.Content = append(s.Content, Statement{Command: command})

	return &s.Content[len(s.Content)-1]
}
