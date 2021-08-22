package recipe

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
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
	Variables   map[string]Recipe
	Columns     map[int]Recipe
	Placeholder string
}

func (t *Transformation) Dump(w io.Writer) {
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
	output := Output{
		Type:  "variable",
		Value: variable,
	}
	t.Variables[variable] = Recipe{Output: output}
}

func (t *Transformation) AddOutputToColumn(column string) {
	output := Output{
		Type:  "column",
		Value: column,
	}
	columnNum, _ := strconv.Atoi(column)
	t.Columns[columnNum] = Recipe{Output: output}
}

func (t *Transformation) Execute(reader *csv.Reader, writer **csv.Writer, l *LineContext) {
	for {
		record, err := reader.Read()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalln(err)
		}
		l.Variables = make(map[string]string)
		l.Columns = make(map[int]string)
		for i, r := range record {
			l.Columns[i+1] = r
		}

		for v := range t.Variables {
			name := t.Variables[v].Output.Value
			pipe := t.Variables[v].Pipe
			var placeholder string
			mode := "replace"

			for _, op := range pipe {
				opName := op.Name
				if opName == "value" {
					arg := op.Arguments[0]
					if arg.Type == "literal" {
						if mode == "replace" {
							placeholder = arg.Value
						} else {
							placeholder += arg.Value
							mode = "replace"
						}
					} else if arg.Type == "column" {
						column, _ := strconv.Atoi(arg.Value)
						value := l.Columns[column]
						//fmt.Println("Got col value", value)
						if mode == "replace" {
							placeholder = value
						} else {
							placeholder += value
							mode = "replace"
						}
					} else {
						fmt.Println("implement", arg.Type, "for value")
					}
				}
				if opName == "join" {
					mode = "join"
				}
				if opName == "lowercase" {
					placeholder = strings.ToLower(placeholder)
				}
				fmt.Println(opName, placeholder)
			}
			l.Variables[name] = placeholder
			(*writer).Write([]string{"Variable:", name, "=", placeholder})
			fmt.Println("Got name: ", name, placeholder)
		}
	}
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

type LineContext struct {
	Variables map[string]string
	Columns   map[int]string
}

func NewTransformation() *Transformation {
	return &Transformation{
		Variables: make(map[string]Recipe),
		Columns:   make(map[int]Recipe),
	}
}
