package parser

// Sequence represents a sequence of commands
type Sequence struct {
	Content []Statement
	Parent  *Sequence // nil for root sequence
}

// Source returns an actual Source-compatible script source.
func (s Sequence) Source() string {
	var content string

	if s.Parent != nil {
		// For inner scopes

		length := len(s.Content)
		if length > 0 {
			content += s.Content[0].Source()
		}

		for i := 1; i < length; i++ {
			content += ";" + s.Content[i].Source()
		}

		content = "\"" + content + "\""
	} else {
		// For outer scopes

		for _, statement := range s.Content {
			content += statement.Source() + "\n"
		}
	}

	return content
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
