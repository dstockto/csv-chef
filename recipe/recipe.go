package recipe

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type Output struct {
	Type  string
	Value string
}

func (o *Output) GetValue(ctx LineContext) (string, error) {
	if o.Type == "variable" {
		value, ok := ctx.Variables[o.Value]
		if !ok {
			return "", errors.New("unrecognized variable")
		}
		return value, nil
	}
	if o.Type == "column" {
		column, _ := strconv.Atoi(o.Value)
		value, ok := ctx.Columns[column]
		if !ok {
			return "", errors.New("unrecognized/unfilled column number")
		}
		return value, nil
	}

	return "", errors.New("unknown column type")
}

type Argument struct {
	Type  string
	Value string
}

type Operation struct {
	Name      string
	Arguments []Argument
}

type Recipe struct {
	Output  Output
	Pipe    []Operation
	Comment string
}

type Transformation struct {
	Variables map[string]Recipe
	Columns   map[int]Recipe
	Headers   map[int]Recipe
}

func (t *Transformation) Dump(w io.Writer) {
	_, _ = fmt.Fprintln(w, "Headers: \n=====")
	for _, h := range t.Headers {
		_, _ = fmt.Fprintf(w, "Header: %s\n", h.Output.Value)
		_, _ = fmt.Fprintf(w, "pipe: ")
		for _, p := range h.Pipe {
			_, _ = fmt.Fprintf(w, p.Name+"(")
			for _, a := range p.Arguments {
				_, _ = fmt.Fprintf(w, "%s: %s, ", a.Type, a.Value)
			}
			_, _ = fmt.Fprintf(w, ") -> ")
		}
		_, _ = fmt.Fprintln(w)
		_, _ = fmt.Fprintf(w, "Comment: # %s\n---\n", h.Comment)
	}

	_, _ = fmt.Fprintln(w, "Variables: \n======")
	for _, v := range t.Variables {
		_, _ = fmt.Fprintf(w, "Var: %s\n", v.Output.Value)
		_, _ = fmt.Fprint(w, "pipe: ")
		for _, p := range v.Pipe {
			_, _ = fmt.Fprint(w, p.Name+"(")
			for _, a := range p.Arguments {
				_, _ = fmt.Fprintf(w, "%s: %s, ", a.Type, a.Value)
			}
			_, _ = fmt.Fprintf(w, ") -> ")
		}
		_, _ = fmt.Fprintln(w)
		_, _ = fmt.Fprintf(w, "Comment: %s\n---\n", v.Comment)
	}

	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(w, "Columns: \n======")
	for _, c := range t.Columns {
		_, _ = fmt.Fprintf(w, "Column: %s\n", c.Output.Value)
		_, _ = fmt.Fprint(w, "pipe: ")
		for _, p := range c.Pipe {
			_, _ = fmt.Fprintf(w, p.Name+"(")
			for _, a := range p.Arguments {
				_, _ = fmt.Fprintf(w, "%s: %s, ", a.Type, a.Value)
			}
			_, _ = fmt.Fprint(w, ") -> ")
		}
		_, _ = fmt.Fprintln(w)
		_, _ = fmt.Fprintf(w, "Comment: %s\n---\n", c.Comment)
	}
}

func (t *Transformation) AddOutputToVariable(variable string) {
	t.Variables[variable] = Recipe{Output: getOutputForVariable(variable)}
}

func (t *Transformation) AddOutputToColumn(column string) {
	output := getOutputForColumn(column)
	columnNum, _ := strconv.Atoi(column)
	t.Columns[columnNum] = Recipe{Output: output}
}

func (t *Transformation) AddOutputToHeader(header string) {
	output := getOutputForHeader(header)
	headerNum, _ := strconv.Atoi(header)
	t.Headers[headerNum] = Recipe{Output: output}
}

func (t *Transformation) Execute(reader *csv.Reader, writer *csv.Writer, processHeader bool) error {
	defer writer.Flush()

	numColumns := len(t.Columns)

	if err := t.ValidateRecipe(); err != nil {
		return err
	}
	var linesRead int

	// TODO process header

	// TODO process columns
	for {
		var context = LineContext{
			Variables: make(map[string]string),
			Columns:   make(map[int]string),
		}
		var output = make(map[int]string)
		// load line into context
		row, err := reader.Read()
		if err != nil {
			return fmt.Errorf("processing line %d error: %v", linesRead, err)
		}
		linesRead++
		for col, cell := range row {
			context.Columns[col+1] = cell
		}

		// junk
		output[1] = "column 1"
		output[2] = "col 2"

		// TODO process variables
		// TODO process columns
		// TODO write output
		var outputStrings []string
		for c := 1; c <= numColumns; c++ {
			value, ok := output[c]
			if !ok {
				return fmt.Errorf("internal logic error: expected to find output for column %d, but nothing was found", c)
			}
			outputStrings = append(outputStrings, value)
		}

		err = writer.Write(outputStrings)
		if err != nil {
			return fmt.Errorf("output error: %v", err)
		}

		//row, err := reader.Read()
		break
	}

	writer.Write([]string{"fruit", "veggie", "mineral"})
	writer.Write([]string{"apple", "carrot", "rock"})
	return nil
}

func (t *Transformation) AddOperationToVariable(variable string, operation Operation) {
	recipe, ok := t.Variables[variable]
	if !ok {
		t.AddOutputToVariable(variable)
		recipe = t.Variables[variable]
	}
	pipe := recipe.Pipe
	if pipe == nil {
		pipe = []Operation{}
	}
	pipe = append(pipe, operation)
	recipe.Pipe = pipe
	t.Variables[variable] = recipe
}

func (t *Transformation) AddOperationToColumn(column string, operation Operation) {
	columnNumber, _ := strconv.Atoi(column)
	recipe, ok := t.Columns[columnNumber]
	if !ok {
		t.AddOutputToColumn(column)
		recipe = t.Columns[columnNumber]
	}
	pipe := recipe.Pipe
	if pipe == nil {
		pipe = []Operation{}
	}
	pipe = append(pipe, operation)
	recipe.Pipe = pipe
	t.Columns[columnNumber] = recipe
}

func (t *Transformation) AddOperationToHeader(header string, operation Operation) {
	headerNumber, _ := strconv.Atoi(header)
	recipe, ok := t.Headers[headerNumber]
	if !ok {
		t.AddOutputToHeader(header)
		recipe = t.Headers[headerNumber]
	}
	pipe := recipe.Pipe
	if pipe == nil {
		pipe = []Operation{}
	}
	pipe = append(pipe, operation)
	recipe.Pipe = pipe
	t.Headers[headerNumber] = recipe
}

func (t *Transformation) AddOperationByType(targetType string, target string, operation Operation) {
	switch targetType {
	case "variable":
		t.AddOperationToVariable(target, operation)
	case "column":
		t.AddOperationToColumn(target, operation)
	case "header":
		t.AddOperationToHeader(target, operation)
	}
}

func (t *Transformation) ValidateRecipe() error {
	numColumns := len(t.Columns)

	// recipe with no columns is pointless/invalid
	if numColumns == 0 {
		return errors.New("no column recipes provided")
	}

	// validate all columns are specified
	for c := 1; c <= numColumns; c++ {
		if _, ok := t.Columns[c]; !ok {
			return fmt.Errorf("missing column definition for column #%d", c)
		}
	}

	return nil
}

type LineContext struct {
	Variables map[string]string
	Columns   map[int]string
}

func NewTransformation() *Transformation {
	return &Transformation{
		Variables: make(map[string]Recipe),
		Columns:   make(map[int]Recipe),
		Headers:   make(map[int]Recipe),
	}
}
