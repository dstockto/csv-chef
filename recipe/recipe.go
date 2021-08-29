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

type JoinMode int

const (
	Replace JoinMode = iota
	Join
)

func (j *JoinMode) String() string {
	switch *j {
	case Replace:
		return "replace"
	case Join:
		return "join"
	default:
		return "unknown join type"
	}
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

func (a *Argument) GetValue(context LineContext, placeholder string) (string, error) {
	var value string
	switch a.Type {
	case "column":
		colNum, _ := strconv.Atoi(a.Value)
		colValue, ok := context.Columns[colNum]
		if !ok {
			return "", fmt.Errorf("column %d referenced but it does not exist in input file", colNum)
		}
		value = colValue
	case "variable":
		varValue, ok := context.Variables[a.Value]
		if !ok {
			return "", fmt.Errorf("variable '%s' referenced, but it is not defined", a.Value)
		}
		value = varValue
	case "literal":
		return a.Value, nil
	case "placeholder":
		return placeholder, nil
	default:
		return "", fmt.Errorf("argument GetValue not implemented for type %s", a.Type)
	}

	return value, nil
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
	Variables     map[string]Recipe
	Columns       map[int]Recipe
	Headers       map[int]Recipe
	VariableOrder []string
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

	for {

		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		linesRead++

		var context = LineContext{
			Variables: map[string]string{},
			Columns:   map[int]string{},
		}
		// Load context with all the columns
		for i, v := range row {
			context.Columns[i+1] = v
		}

		// process variables
		for _, v := range t.VariableOrder {
			variableName := t.Variables[v].Output.Value
			variableRecipe := t.Variables[v]
			placeholder, err := t.processRecipe("variable", variableRecipe, context)
			if err != nil {
				return err
			}
			context.Variables[variableName] = placeholder
		}

		if processHeader && linesRead == 1 {
			// Load existing headers up to size of output
			var output = make(map[int]string)
			for i := 1; i <= numColumns; i++ {
				output[i] = row[i-1]
			}

			for h := range t.Headers {
				headerRecipe := t.Headers[h]
				placeholder, err := t.processRecipe("header", headerRecipe, context)
				if err != nil {
					return err
				}
				output[h] = placeholder
			}

			err := t.outputCsvRow(numColumns, output, writer)
			if err != nil {
				return err
			}
		}

		if !processHeader || linesRead > 1 {
			var output = make(map[int]string)

			for c := range t.Columns {
				columnRecipe := t.Columns[c]
				placeholder, err := t.processRecipe("column", columnRecipe, context)
				if err != nil {
					return err
				}
				output[c] = placeholder
			}

			err := t.outputCsvRow(numColumns, output, writer)
			if err != nil {
				return err
			}
		}

		if linesRead%100 == 0 {
			writer.Flush()
		}
	}

	return nil
}

func (t *Transformation) outputCsvRow(numColumns int, output map[int]string, writer *csv.Writer) error {
	var outputRow []string
	for i := 1; i <= numColumns; i++ {
		outputRow = append(outputRow, output[i])
	}
	err := writer.Write(outputRow)
	if err != nil {
		return err
	}
	return nil
}

func (t *Transformation) processRecipe(recipeType string, variable Recipe, context LineContext) (string, error) {
	var placeholder string
	var value string
	mode := Replace
	for _, o := range variable.Pipe {
		switch o.Name {
		case "value":
			firstArg := o.Arguments[0]
			argValue, err := firstArg.GetValue(context, placeholder)
			if err != nil {
				return "", err
			}
			value = argValue
		case "join":
			firstArg := o.Arguments[0]
			mode = Join
			argValue, err := firstArg.GetValue(context, placeholder)
			if err != nil {
				return "", err
			}
			value = argValue
			continue
		case "uppercase":
			firstArg, err := o.Arguments[0].GetValue(context, placeholder)
			if err != nil {
				return "", fmt.Errorf("uppercase: error evaluating arg - %v", err)
			}
			value = Uppercase(firstArg)
		case "lowercase":
			firstArg, err := o.Arguments[0].GetValue(context, placeholder)
			if err != nil {
				return "", fmt.Errorf("lowercase: error evaluating arg - %v", err)
			}
			value = Lowercase(firstArg)
			// TODO make function calling more smart, using the allFuncs thing
		default:
			return "", fmt.Errorf("error: processing variable, unimplemented operation %s", o.Name)
		}

		switch mode {
		case Replace:
			placeholder = value
		case Join:
			placeholder += value
			mode = Replace
		default:
			return "", fmt.Errorf("invalid join mode %d", mode)
		}
	}
	return placeholder, nil
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

	// ensure there are not header recipes for a column we don't have
	for h := range t.Headers {
		if _, ok := t.Columns[h]; !ok {
			return fmt.Errorf("found header for column %d, but no recipe for column %d", h, h)
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
