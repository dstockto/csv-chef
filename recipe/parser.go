package recipe

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var allFuncs = map[string][]int{
	"uppercase":    {1},
	"lowercase":    {1},
	"join":         {1},
	"add":          {2},
	"subtract":     {2},
	"multiply":     {2},
	"divide":       {2},
	"change":       {3},
	"changei":      {3},
	"ifempty":      {3},
	"isempty":      {3}, // alias for ifempty
	"numberformat": {2},
	"lineno":       {0},
	"removedigits": {1},
	"onlydigits":   {1},
	"mod":          {2},
	"trim":         {1},
	"firstchars":   {1},
	"lastchars":    {1},
	"repeat":       {2},
	"replace":      {3},
	"today":        {0},
	"now":          {0},
	"formatdate":   {2},
	"readdate":     {2},
	"smartdate":    {1},
	"ispast":       {3},
	"isfuture":     {3},

	"normalize_date": {1, 1},
	"fake":           {1},
}

func Parse(source io.Reader) (*Transformation, error) {
	transformation := NewTransformation()

	// split by newlines
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(source)
	s := buf.String()
	lines := strings.Split(s, "\n")

	for lineNo, l := range lines {
		p := NewParser(strings.NewReader(l))

		// Full Line Comment
		tok, lit := p.scanIgnoreWhitespace()
		if tok == COMMENT {
			p.scanComment()
			continue
		}
		if tok == EOF {
			break
		}

		if tok != COLUMN_ID && tok != VARIABLE && tok != HEADER {
			return transformation, fmt.Errorf("expected column, header or variable on line %d, but found %s", lineNo, lit)
		}

		// Found column or variable to assign result to
		target := lit
		var targetType string
		if tok == COLUMN_ID {
			transformation.AddOutputToColumn(lit)
			targetType = "column"
		} else if tok == VARIABLE {
			transformation.AddOutputToVariable(lit)
			transformation.VariableOrder = append(transformation.VariableOrder, lit)
			targetType = "variable"
		} else if tok == HEADER {
			transformation.AddOutputToHeader(lit)
			targetType = "header"
		}

		// After column or variable, we need the assignment <- operator
		if err := consumeAssignment(p); err != nil {
			return nil, err
		}

		// grab first pipe piece - literal, column, variable, function, function w/ args
		tok, lit = p.scanIgnoreWhitespace()
		switch tok {
		case COLUMN_ID:
			transformation.AddOperationByType(targetType, target, getColumn(lit))
		case LITERAL:
			transformation.AddOperationByType(targetType, target, getLiteral(lit))
		case VARIABLE:
			transformation.AddOperationByType(targetType, target, getVariable(lit))
		case FUNCTION:
			function := lit
			operation, err := consumeFunctionArgs(p, function)
			if err != nil {
				return nil, err
			}
			transformation.AddOperationByType(targetType, target, operation)
		default:
			return nil, fmt.Errorf("unexpected token [%d] %s\n", tok, lit)
		}

	LOOPSCAN:
		for {
			tok, lit := p.scanIgnoreWhitespace()
			switch tok {
			case EOF:
				break LOOPSCAN
			case PIPE:
				break
			case PLUS:
				transformation.AddOperationByType(targetType, target, getJoinWithPlaceholder())
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
				if targetType == "header" {
					headerNum, _ := strconv.Atoi(target)
					recipe := transformation.Headers[headerNum]
					recipe.Comment = lit
					transformation.Headers[headerNum] = recipe
				}
				break LOOPSCAN
			default:
				break
			}

			// After connection scan stuff we can do (column, variable, literal, function)
			// Comments or EOL are no bueno here like 1 <- 2 + # comment <- what??
			tok, lit = p.scanIgnoreWhitespace()
			switch tok {
			case COLUMN_ID:
				transformation.AddOperationByType(targetType, target, getColumn(lit))
			case VARIABLE:
				transformation.AddOperationByType(targetType, target, getVariable(lit))
			case LITERAL:
				transformation.AddOperationByType(targetType, target, getLiteral(lit))
			case FUNCTION:
				function := lit
				operation, err := consumeFunctionArgs(p, function)
				if err != nil {
					return nil, err
				}
				transformation.AddOperationByType(targetType, target, operation)
			case PLACEHOLDER:
				transformation.AddOperationByType(targetType, target, getPlaceholder())
			default:
				return nil, fmt.Errorf("unexpected token [%d]-'%s' in parse loop", tok, lit)
			}
		}
	}

	return transformation, nil
}

func getLiteral(lit string) Operation {
	return Operation{
		Name: "value",
		Arguments: []Argument{
			literalArg(lit),
		},
	}
}

func getColumn(lit string) Operation {
	op := Operation{
		Name: "value",
		Arguments: []Argument{
			columnArg(lit),
		},
	}
	return op
}

func getPlaceholder() Operation {
	return Operation{
		Name: "value",
		Arguments: []Argument{
			placeholderArg(),
		},
	}
}

func getFunction(name string, args []Argument) Operation {
	return Operation{
		Name:      name,
		Arguments: args,
	}
}

func getVariable(lit string) Operation {
	return Operation{
		Name: "value",
		Arguments: []Argument{
			variableArg(lit),
		},
	}
}

func getJoinWithPlaceholder() Operation {
	return Operation{
		Name: "join",
		Arguments: []Argument{
			placeholderArg(),
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

func getOutputForHeader(h string) Output {
	return Output{
		Type:  "header",
		Value: h,
	}
}

func consumeAssignment(p *Parser) error {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != ASSIGNMENT {
		return fmt.Errorf("expected assignment ( <- ) but found [%s] (%d) instead.\n", lit, tok)
	}
	return nil
}

func consumeFunctionArgs(p *Parser, name string) (Operation, error) {
	// check if the function even exists
	funcArgs, ok := allFuncs[strings.ToLower(name)]
	if !ok {
		return Operation{}, fmt.Errorf("unrecognized function %s", name)
	}
	var totalArgs int

	for _, count := range funcArgs {
		totalArgs += count
	}

	// look for paren
	tok, _ := p.scan()

	var operation Operation

	operation.Name = name
	// if not a paren, then we have the default placeholder arg
	if tok != OPEN_PAREN {
		// Whatever it was, we can let it get parsed elsewhere
		p.s.unread()
		operation.Arguments = []Argument{}
		for i := 0; i < totalArgs; i++ {
			operation.Arguments = append(operation.Arguments, placeholderArg())
		}

		return operation, nil
	}

	var gotPlaceholder bool // track if the placeholder was explicitly provided or not
	var args []Argument
ARGLOOP:
	for {
		tok, lit := p.scanIgnoreWhitespace()
		switch tok {
		case EOF:
			return operation, fmt.Errorf("expected function args for %s. found EOF", name)
		case LITERAL:
			args = append(args, literalArg(lit))
		case PLACEHOLDER:
			gotPlaceholder = true
			args = append(args, placeholderArg())
		case COLUMN_ID:
			args = append(args, columnArg(lit))
		case VARIABLE:
			args = append(args, variableArg(lit))
		case COMMA:
			break
		case CLOSE_PAREN:
			break ARGLOOP
		default:
			return operation, fmt.Errorf("expected function args, got [%d] - %s", tok, lit)
		}
	}

	if !gotPlaceholder || len(args) == 0 {
		args = append(args, placeholderArg())
	}
	//var args []Argument
	// must now get args until we get a close paren
	operation.Arguments = args

	return operation, nil
}

func variableArg(lit string) Argument {
	return Argument{
		Type:  "variable",
		Value: lit,
	}
}

func columnArg(lit string) Argument {
	return Argument{
		Type:  "column",
		Value: lit,
	}
}

func literalArg(lit string) Argument {
	return Argument{
		Type:  "literal",
		Value: lit,
	}
}

func placeholderArg() Argument {
	return Argument{
		Type:  "placeholder",
		Value: "?",
	}
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
	VARIABLE          //11 - starts w/ $
	FUNCTION          //12 - letters
	OPEN_PAREN        //13 - (
	CLOSE_PAREN       //14 - )
	COMMA             //15 - ,
	HEADER            //16 - !<digits>
	//ARGUMENT			// unknown if needed

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
		return s.scanFunction()
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
	} else if ch == '!' {
		_, lit := s.scanColumn()
		return HEADER, lit
	} else if ch == '#' {
		return s.scanComment()
	}

	// Otherwise read the individual character.
	switch ch {
	case '\n':
		return NEWLINE, ""
	case eof:
		return EOF, "EOF"
	case '?':
		return PLACEHOLDER, string(ch)
	case '+':
		return PLUS, string(ch)
	case '(':
		return OPEN_PAREN, string(ch)
	case ')':
		return CLOSE_PAREN, string(ch)
	case ',':
		return COMMA, string(ch)
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

func (s *Scanner) scanFunction() (Token, string) {
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

	// Otherwise return as a regular identifier.
	return FUNCTION, buf.String()
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

	s.read() // eat the leading quote

	for {
		ch := s.read()
		if ch == '\\' {
			_, _ = buf.WriteRune(s.read())
		} else if ch != '"' {
			_, _ = buf.WriteRune(ch)
		} else {
			break
		}
	}

	return LITERAL, buf.String()
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
		} else if isLetter(ch) {
			_, _ = buf.WriteRune(ch)
		} else {
			s.unread()
			break
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
