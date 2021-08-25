package recipe

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewTransformation(t *testing.T) {
	tests := []struct {
		name string
		want *Transformation
	}{
		{
			name: "buids a Transformation structure",
			want: &Transformation{
				Variables: map[string]Recipe{},
				Columns:   map[int]Recipe{},
				Headers:   map[int]Recipe{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransformation(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransformation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutput_GetValue(t *testing.T) {
	type fields struct {
		Type  string
		Value string
	}
	type args struct {
		ctx LineContext
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "get value from variable",
			fields: fields{
				Type:  "variable",
				Value: "orange",
			},
			args: args{
				ctx: LineContext{
					Variables: map[string]string{"pewp": "salami", "orange": "soda"},
					Columns:   nil,
				},
			},
			want:    "soda",
			wantErr: false,
		},
		{
			name: "get value from column",
			fields: fields{
				Type:  "column",
				Value: "2",
			},
			args: args{
				ctx: LineContext{
					Variables: nil,
					Columns:   map[int]string{1: "salad", 2: "fingers"},
				},
			},
			want:    "fingers",
			wantErr: false,
		},
		{
			name: "get value from non-existent column",
			fields: fields{
				Type:  "column",
				Value: "5",
			},
			args: args{
				ctx: LineContext{
					Variables: nil,
					Columns:   map[int]string{1: "salad", 2: "fingers"},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "get value from non-existent variable",
			fields: fields{
				Type:  "variable",
				Value: "plop",
			},
			args: args{
				ctx: LineContext{
					Variables: map[string]string{"herpa": "derp"},
					Columns:   nil,
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Output{
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
			}
			got, err := o.GetValue(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//
//func TestTransformation_Dump(t1 *testing.T) {
//	type fields struct {
//		Variables   map[string]Recipe
//		Columns     map[int]Recipe
//		Placeholder string
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		wantW  string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t1.Run(tt.name, func(t1 *testing.T) {
//			t := &Transformation{
//				Variables:   tt.fields.Variables,
//				Columns:     tt.fields.Columns,
//				Placeholder: tt.fields.Placeholder,
//			}
//			w := &bytes.Buffer{}
//			t.Dump(w)
//			if gotW := w.String(); gotW != tt.wantW {
//				t1.Errorf("Dump() = %v, want %v", gotW, tt.wantW)
//			}
//		})
//	}
//}
//
//func TestTransformation_Execute(t1 *testing.T) {
//	type fields struct {
//		Variables   map[string]Recipe
//		Columns     map[int]Recipe
//		Placeholder string
//	}
//	type args struct {
//		reader *csv.Reader
//		writer **csv.Writer
//		l      *LineContext
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t1.Run(tt.name, func(t1 *testing.T) {
//			t := &Transformation{
//				Variables:   tt.fields.Variables,
//				Columns:     tt.fields.Columns,
//				Placeholder: tt.fields.Placeholder,
//			}
//		})
//	}
//}

func TestTransformation_AddOutputToVariable(t1 *testing.T) {
	tests := []struct {
		name     string
		variable string
		want     Output
	}{
		{
			name:     "add variable output",
			variable: "floop",
			want: Output{
				Type:  "variable",
				Value: "floop",
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := NewTransformation()
			t.AddOutputToVariable(tt.variable)
			got := t.Variables[tt.variable].Output
			if got != tt.want {
				t1.Errorf("Dump() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransformation_AddOutputToHeader(t1 *testing.T) {
	tests := []struct {
		name      string
		header    string
		headerNum int
		want      Output
	}{
		{
			name:      "add header output",
			header:    "5",
			headerNum: 5,
			want: Output{
				Type:  "header",
				Value: "5",
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := NewTransformation()
			t.AddOutputToHeader(tt.header)
			got := t.Headers[tt.headerNum].Output
			if got != tt.want {
				t1.Errorf("Dump() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransformation_AddOutputToColumn(t1 *testing.T) {
	tests := []struct {
		name      string
		column    string
		want      Output
		columnNum int
	}{
		{
			name:      "add column output",
			column:    "42",
			columnNum: 42,
			want: Output{
				Type:  "column",
				Value: "42",
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := NewTransformation()
			t.AddOutputToColumn(tt.column)

			got := t.Columns[tt.columnNum].Output
			if got != tt.want {
				t1.Errorf("Dump() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransformation_AddOperationToVariable(t1 *testing.T) {
	tests := []struct {
		name      string
		variable  string
		initial   Transformation
		operation Operation
		want      Operation
	}{
		{
			name:     "add operation to variable",
			variable: "floopy",
			initial: Transformation{
				Variables: map[string]Recipe{
					"floopy": {
						Output:  Output{},
						Pipe:    []Operation{},
						Comment: "",
					},
				},
				Columns: nil,
			},
			operation: Operation{
				Name: "value",
				Arguments: []Argument{
					{
						Type:  "literal",
						Value: "ham",
					},
				},
			},
			want: Operation{
				Name: "value",
				Arguments: []Argument{
					{
						Type:  "literal",
						Value: "ham",
					},
				},
			},
		},
		{
			name:     "add operation to variable without output",
			variable: "ploopy",
			initial: Transformation{
				Variables: map[string]Recipe{},
				Columns:   map[int]Recipe{},
			},
			operation: Operation{
				Name: "value",
				Arguments: []Argument{
					{
						Type:  "literal",
						Value: "sammich",
					},
				},
			},
			want: Operation{
				Name: "value",
				Arguments: []Argument{
					{
						Type:  "literal",
						Value: "sammich",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.initial
			initialLen := len(tt.initial.Variables[tt.variable].Pipe)
			t.AddOperationToVariable(tt.variable, tt.operation)

			if len(t.Variables[tt.variable].Pipe) != initialLen+1 {
				t1.Errorf("expected pipe have %d operation, got %d", initialLen+1, len(t.Variables[tt.variable].Pipe))
				t1.Fail()
			}

			got := t.Variables[tt.variable].Pipe[0]

			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Dump() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransformation_AddOperationToColumn(t1 *testing.T) {
	tests := []struct {
		name         string
		column       string
		columnNumber int
		initial      Transformation
		operation    Operation
		want         Operation
	}{
		{
			name:         "add operation to column that is initialized",
			column:       "10",
			columnNumber: 10,
			initial: Transformation{
				Variables: map[string]Recipe{},
				Columns: map[int]Recipe{
					10: {
						Output: Output{
							Type:  "column",
							Value: "10",
						},
						Pipe:    []Operation{},
						Comment: "",
					},
				},
			},
			operation: Operation{
				Name: "value",
				Arguments: []Argument{
					{
						Type:  "column",
						Value: "3",
					},
				},
			},
			want: Operation{
				Name: "value",
				Arguments: []Argument{
					{
						Type:  "column",
						Value: "3",
					},
				},
			},
		},
		{
			name:         "add operation to column that is not initialized",
			column:       "14",
			columnNumber: 14,
			initial: Transformation{
				Variables: map[string]Recipe{},
				Columns:   map[int]Recipe{},
			},
			operation: Operation{
				Name: "value",
				Arguments: []Argument{
					{
						Type:  "placeholder",
						Value: "?",
					},
				},
			},
			want: Operation{
				Name: "value",
				Arguments: []Argument{
					{
						Type:  "placeholder",
						Value: "?",
					},
				},
			},
		},
		{
			name:         "add operation to column that is not initialized",
			column:       "14",
			columnNumber: 14,
			initial: Transformation{
				Variables: map[string]Recipe{},
				Columns:   map[int]Recipe{},
			},
			operation: Operation{
				Name: "value",
				Arguments: []Argument{
					{
						Type:  "placeholder",
						Value: "?",
					},
				},
			},
			want: Operation{
				Name: "value",
				Arguments: []Argument{
					{
						Type:  "placeholder",
						Value: "?",
					},
				},
			},
		},
		{
			name:         "add operation to column that has operation",
			column:       "10",
			columnNumber: 10,
			initial: Transformation{
				Variables: map[string]Recipe{},
				Columns: map[int]Recipe{
					10: {
						Output: Output{
							Type:  "column",
							Value: "10",
						},
						Pipe: []Operation{
							{
								Name: "fake",
								Arguments: []Argument{
									{
										Type:  "literal",
										Value: "name",
									},
								},
							},
						},
						Comment: "",
					},
				},
			},
			operation: Operation{
				Name: "value",
				Arguments: []Argument{
					{
						Type:  "column",
						Value: "3",
					},
				},
			},
			want: Operation{
				Name: "value",
				Arguments: []Argument{
					{
						Type:  "column",
						Value: "3",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.initial
			initialLength := len(tt.initial.Columns[tt.columnNumber].Pipe)
			t.AddOperationToColumn(tt.column, tt.operation)

			fmt.Println("initial", initialLength)
			if len(t.Columns[tt.columnNumber].Pipe) != (initialLength + 1) {
				t1.Errorf("expected pipe have %d operations, got %d", initialLength, len(t.Columns[tt.columnNumber].Pipe))
				t1.Fail()
			}

			got := t.Columns[tt.columnNumber].Pipe[initialLength]

			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Dump() = %v, want %v", got, tt.want)
			}
		})
	}
}
