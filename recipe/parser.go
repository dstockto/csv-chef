package recipe

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func Parse(source io.Reader) (*Transformation, error) {
	transformation := NewTransformation()

	// split by newlines
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(source)
	s := buf.String()
	lines := strings.Split(s, "\n")

	for lineNo, l := range lines {
		fmt.Printf("Parsing line #%d\n", lineNo)
		p := NewParser(strings.NewReader(l))

		// Full Line Comment
		tok, lit := p.scanIgnoreWhitespace()
		if tok == COMMENT {
			comment := p.scanComment()
			fmt.Printf("Ignore comment: # %s\n", comment)
			continue
		}

		if tok != COLUMN_ID && tok != VARIABLE {
			return transformation, fmt.Errorf("expected column or variable on line %d, but found %s", lineNo, lit)
		}

		// Found column or variable to assign result to
		target := lit
		var targetType string
		if tok == COLUMN_ID {
			transformation.AddOutputToColumn(lit)
			targetType = "column"
		} else if tok == VARIABLE {
			transformation.AddOutputToVariable(lit)
			targetType = "variable"
		}

		// After column or variable, we need the assignment <- operator
		if err := consumeAssignment(p); err != nil {
			return nil, err
		}

		// grab first pipe piece - literal, column, variable, function, function w/ args
		tok, lit = p.scanIgnoreWhitespace()
		var operation Operation
		switch tok {
		case COLUMN_ID:
			operation = getColumn(lit)
		case LITERAL:
			operation = getLiteral(lit)
		case VARIABLE:
			operation = getVariable(lit)
		default:
			return nil, fmt.Errorf("unexpected token [%d] %s\n", tok, lit)
		}

		// we have the operation, add it
		switch targetType {
		case "variable":
			transformation.AddOperationToVariable(target, operation)
		case "column":
			transformation.AddOperationToColumn(target, operation)
		}

		// now we're in a loop because we've got at least one thing
		// sanity kill
		killCount := 0
		for {
			var operation Operation
			killCount++
			if killCount >= 10 {
				break
			}
			tok, lit := p.scanIgnoreWhitespace()
			fmt.Printf("%d %s\n", tok, lit)
			switch tok {
			case EOF:
				break
			case LITERAL:
				operation = Operation{
					Name: "value",
					Arguments: []Argument{
						{Type: "literal", Value: lit},
					},
				}
			case COMMENT:
				if targetType == "variable" {
					recipe := transformation.Variables[target]
					recipe.Comment = lit
					transformation.Variables[target] = recipe
				}
				if targetType == "column" {
					columnNum, _ := strconv.Atoi(target)
					recipe := transformation.Columns[columnNum]
					recipe.Comment = lit
					transformation.Columns[columnNum] = recipe
				}
				break
			default:
				break
			}
			if operation.Name != "" {
				if targetType == "variable" {
					transformation.AddOperationToVariable(target, operation)
				} else if targetType == "column" {
					transformation.AddOperationToColumn(target, operation)
				}
			}
		}
	}

	return transformation, nil
}

func getLiteral(lit string) Operation {
	return Operation{
		Name: "value",
		Arguments: []Argument{
			{
				Type:  "literal",
				Value: lit,
			},
		},
	}
}

func getColumn(lit string) Operation {
	op := Operation{
		Name: "value",
		Arguments: []Argument{
			{
				Type:  "column",
				Value: lit,
			},
		},
	}
	return op
}

func getVariable(lit string) Operation {
	return Operation{
		Name: "value",
		Arguments: []Argument{
			{
				Type:  "variable",
				Value: lit,
			},
		},
	}
}

func getJoinWithPlaceholder() Operation {
	return Operation{
		Name: "join",
		Arguments: []Argument{
			{
				Type:  "placeholder",
				Value: "?",
			},
		},
	}
}

func getOutputForColumn(col string) Output {
	return Output{
		Type:  "column",
		Value: col,
	}
}

func getOutputForVariable(v string) Output {
	return Output{
		Type:  "variable",
		Value: v,
	}
}

func consumeAssignment(p *Parser) error {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != ASSIGNMENT {
		return fmt.Errorf("expected assignment ( <- ) but found [%s] (%d) instead.\n", lit, tok)
	}
	return nil
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

type Token int

const (
	ILLEGAL     Token = iota
	EOF               //1 - end of file
	WS                //2 - space, tab, newline
	NEWLINE           //3 - \n (probably not needed)
	COLUMN_ID         //4 - digits
	ASSIGNMENT        //5 - <-
	PIPE              //6 - ->
	COMMENT           //7 - # ...
	PLACEHOLDER       //8 - ?
	PLUS              //9 - +
	LITERAL           //10 - "quoted"
	IDENT             //11 - all letters
	VARIABLE          //12 - starts w/ $
	COMMA             //13 - ,
	//ARGUMENT			// unknown if needed
	//FUNCTION          // need to figure out - IDENT may be this

	//	column_id <- [0-9]+ | p + column_id
	//column_assign <- "<-"
	//pipe <- "->"
	//comment <- "#"
	//placeholder <- "?"
	//variable <- [a-zA-z_][a-zA-Z_0-9]*
	//expression <- column_id | variable | processed_column | function
	//identifier <- placeholder | variable | column_id
	//argument_list <- identifier | identifier + "," + argument_list
	//function <- function_identifier + "(" + argument_list + ")"
)

func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

var eof = rune(0)

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

type Parser struct {
	s   *Scanner
	buf struct {
		tok Token
		lit string
		n   int
	}
}

// read reads the next rune from the buffered reader
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// Scan returns the next token and literal value
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune
	ch := s.read()

	// If we see whitespace then we consume all contiguous whitespace
	// If we see a letter then consume as an ident or reserved word.
	if isWhiteSpace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	} else if ch == '<' {
		s.unread()
		return s.scanAssignment()
	} else if ch == '-' {
		s.unread()
		return s.scanPipe()
	} else if isDigit(ch) {
		s.unread()
		return s.scanColumn()
	} else if ch == '$' {
		s.unread()
		return s.scanVariable()
	} else if ch == '"' {
		s.unread()
		return s.scanLiteral()
	} else if ch == '#' {
		return s.scanComment()
	}

	// Otherwise read the individual character.
	switch ch {
	case '\n':
		return NEWLINE, ""
	case eof:
		return EOF, ""
	case '?':
		return PLACEHOLDER, string(ch)
	case '+':
		return PLUS, string(ch)
	}

	return ILLEGAL, string(ch)
}

func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhiteSpace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

func (s *Scanner) scanIdent() (Token, string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// If the string matches a keyword then return that keyword.
	switch strings.ToUpper(buf.String()) {
	case "->":
		return PIPE, buf.String()
	case "<-":
		return ASSIGNMENT, buf.String()
	}

	// Otherwise return as a regular identifier.
	return IDENT, buf.String()
}

func (s *Scanner) scanAssignment() (Token, string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	ch := s.read()
	if ch != '-' {
		return ILLEGAL, string(ch)
	}

	return ASSIGNMENT, "<-"
}

func (s *Scanner) scanPipe() (Token, string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	ch := s.read()
	if ch != '>' {
		return ILLEGAL, string(ch)
	}

	return PIPE, "->"
}

func (s *Scanner) scanColumn() (Token, string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		ch := s.read()
		if isDigit(ch) {
			_, _ = buf.WriteRune(ch)
		} else {
			s.unread()
			break
		}
	}

	return COLUMN_ID, buf.String()
}

func (s *Scanner) scanLiteral() (Token, string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		ch := s.read()
		if ch != '"' {
			_, _ = buf.WriteRune(ch)
		} else {
			_, _ = buf.WriteRune(ch)
			break
		}
	}

	return LITERAL, strings.Trim(buf.String(), "\"")
}

func (s *Scanner) scanVariable() (Token, string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		ch := s.read()
		if isWhiteSpace(ch) || ch == eof {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return VARIABLE, buf.String()
}

func (s *Scanner) scanComment() (Token, string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		ch := s.read()
		if ch == '\n' || ch == eof {
			break
		}
		buf.WriteRune(ch)
	}

	return COMMENT, strings.TrimSpace(buf.String())
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS || tok == NEWLINE {
		tok, lit = p.scan()
	}

	return
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan it later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

func (p *Parser) scanComment() string {
	var tok Token
	var lit string
	var comment strings.Builder

	for {
		tok, lit = p.scan()
		if tok == EOF || tok == NEWLINE {
			break
		}
		comment.WriteString(lit)
	}

	return strings.TrimSpace(comment.String())
}
