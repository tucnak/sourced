package parser

import (
	"strconv"
	"unicode"
)

// Build tries to generate a legit script source
// for Source engine.
func Build(input []byte) ([]byte, error) {
	ast, err := Parse(input)
	if err != nil {
		return []byte{}, err
	}

	output, err := ast.String()
	if err != nil {
		return []byte{}, err
	}

	return []byte(output), nil
}

//go:generate stringer -type=scope
type scope int

const (
	OUTSCOPED scope = iota
	IN_COMMENT
	READING_COMMAND
	READING_ARGUMENTS
	READING_NUMBER
	READING_WORD
	READING_STRING
)

type parsingCtx struct {
	line, pos int
	input     *[]byte
	index     int
	symbol    rune

	scope scope

	head, root *Sequence

	command, word, argstring, number string
}

func Parse(input []byte) (*Sequence, error) {
	ctx := &parsingCtx{}
	ctx.input = &input
	ctx.root = &Sequence{}
	ctx.head = ctx.root

	ctx.line = 1
	ctx.pos = 1
	ctx.index = -1

	for _, char := range input {
		ctx.index++

		symbol := rune(char)
		ctx.symbol = symbol

		if ctx.scope == IN_COMMENT && symbol != '\n' {
			ctx.pos++
			continue
		}

		switch symbol {
		case '\n':
			if ctx.scope == READING_COMMAND {
				statement := ctx.head.Add(ctx.command)
				statement.line = ctx.line

				ctx.scope = OUTSCOPED
				ctx.line++
				ctx.pos = 0

				continue
			}

			ctx.line++
			ctx.pos = 0

			err := finishReadingArgument(ctx)
			if err != nil {
				return ctx.root, err
			}

			if ctx.scope == IN_COMMENT {
				ctx.scope = OUTSCOPED
				continue
			}

			if ctx.scope != READING_STRING {
				ctx.scope = OUTSCOPED
				continue
			}

			continue

		case ' ':
			if ctx.scope == OUTSCOPED {
				break
			}

			if ctx.scope == READING_COMMAND {
				ctx.scope = READING_ARGUMENTS
				ctx.head.Add(ctx.command)
				break
			}

			if ctx.scope == READING_ARGUMENTS {
				break
			}

			if ctx.scope == READING_STRING {
				err := parseWhileReadingString(ctx)
				if err != nil {
					return ctx.root, err
				}

				break
			}

			err := finishReadingArgument(ctx)
			if err != nil {
				return ctx.root, err
			}

		case '\t':
			continue

		case '{':
			if ctx.scope != READING_ARGUMENTS {
				return ctx.root, &ParseError{
					Line:     ctx.line,
					Position: ctx.pos,
					Message:  "Opening a scope in inappropriate place",
				}
			}

			last := ctx.head.Last()
			scope := &Sequence{Parent: ctx.head}
			last.AddArgument(scope)
			ctx.head = scope

		case '}':
			if ctx.head.Parent == nil {
				return ctx.root, &ParseError{
					Line:     ctx.line,
					Position: ctx.pos,
					Message:  "Trying to close a non-existing scope",
				}
			}

			ctx.head = ctx.head.Parent

		default:
			if symbol == '/' && relativeTo(ctx, 1) == '/' {
				ctx.scope = IN_COMMENT
				break
			}

			if ctx.scope == OUTSCOPED {
				if !isLegitCharForName(symbol) {
					return ctx.root, &ParseError{
						Line:     ctx.line,
						Position: ctx.pos,
						Message:  "Command starts with a non-letter character",
					}
				}

				ctx.scope = READING_COMMAND
				ctx.command = ""
				ctx.command += string(symbol)
				break
			}

			if ctx.scope == READING_COMMAND {
				ctx.command += string(symbol)

				break
			}

			if ctx.scope == READING_ARGUMENTS {
				if symbol == '"' {
					ctx.scope = READING_STRING
					ctx.argstring = ""
					break
				}

				if unicode.IsDigit(symbol) {
					ctx.scope = READING_NUMBER
					ctx.number = string(symbol)
					break
				}

				if isLegitCharForName(symbol) {
					ctx.scope = READING_WORD
					ctx.word = string(symbol)
					break
				}

				break
			}

			if ctx.scope == READING_STRING {
				err := parseWhileReadingString(ctx)
				if err != nil {
					return ctx.root, err
				}

				break
			}

			if ctx.scope == READING_NUMBER {
				if unicode.IsDigit(symbol) {
					ctx.number += string(symbol)
					break
				}

				return ctx.root, &ParseError{
					Line:     ctx.line,
					Position: ctx.pos,
					Message:  "Unexpected non-digit character in the number parameter",
				}
			}

			if ctx.scope == READING_WORD {
				if isLegitCharForName(symbol) {
					ctx.word += string(symbol)
					break
				}

				return ctx.root, &ParseError{
					Line:     ctx.line,
					Position: ctx.pos,
					Message:  "Unexpected non-letter character in the word parameter",
				}
			}
		}

		ctx.pos++
	}

	return ctx.root, nil
}

// isLegitCharForName checks whether symbol is a legit "word" character.
func isLegitCharForName(symbol rune) bool {
	isLetter := unicode.IsLetter(symbol)
	isUnderscore := (symbol == '_')
	isSwitch := ((symbol == '+') || (symbol == '-'))
	isNumber := unicode.IsNumber(symbol)

	if isLetter || isUnderscore || isSwitch || isNumber {
		return true
	}

	return false
}

// relativeTo returns a symbol with relative offset to current,
// or zero-rune if it doesn't exist.
func relativeTo(ctx *parsingCtx, offset int) rune {
	i := ctx.index
	length := len(*ctx.input)

	if i+offset < 0 {
		return rune(0)
	}

	if i+offset >= length {
		return rune(0)
	}

	return rune((*ctx.input)[i+offset])
}

// parseWhileReadingString parses a character in string scope.
func parseWhileReadingString(ctx *parsingCtx) *ParseError {
	symbol := ctx.symbol

	if symbol == '\\' {
		if relativeTo(ctx, -1) != '\\' {
			return nil
		}

		ctx.argstring += "\\"
		return nil
	}

	if ctx.index > 0 && relativeTo(ctx, -1) == '\\' && relativeTo(ctx, -2) != '\\' {
		if symbol == '"' {
			ctx.argstring += "\""
		} else if symbol == 'n' {
			ctx.argstring += "\n"
		} else if symbol == '/' {
			ctx.argstring += "/"
		}

		return nil
	}

	if symbol == '"' {
		ctx.scope = READING_ARGUMENTS
		if statement := ctx.head.Last(); statement != nil {
			statement.AddArgument(&StringArgument{ctx.argstring})
			return nil
		}
	}

	ctx.argstring += string(symbol)

	return nil
}

// finishReadingArgument finishes a word/number argument
// with READING_ARGUMENTS scope
func finishReadingArgument(ctx *parsingCtx) *ParseError {
	if ctx.scope == READING_WORD {
		ctx.scope = READING_ARGUMENTS
		if statement := ctx.head.Last(); statement != nil {
			statement.AddArgument(&WordArgument{ctx.word})

			return nil
		}
	}

	if ctx.scope == READING_NUMBER {
		ctx.scope = READING_ARGUMENTS
		if statement := ctx.head.Last(); statement != nil {
			number, _ := strconv.Atoi(ctx.number)
			statement.AddArgument(&NumberArgument{number})

			return nil
		}
	}

	return nil
}
