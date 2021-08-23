package recipe

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		source io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *Transformation
		wantErr bool
	}{
		{
			name: "Full line comment",
			args: args{
				source: strings.NewReader("# full line comment"),
			},
			want: &Transformation{
				Variables: map[string]Recipe{},
				Columns:   map[int]Recipe{},
			},
			wantErr: false,
		},
		{
			name: "Full line comment no space",
			args: args{
				source: strings.NewReader("#full line comment"),
			},
			want: &Transformation{
				Variables: map[string]Recipe{},
				Columns:   map[int]Recipe{},
			},
			wantErr: false,
		},
		{
			name: "Variable loaded from column",
			args: args{
				source: strings.NewReader("$foo <- 103"),
			},
			want: &Transformation{
				Variables: map[string]Recipe{
					"$foo": {
						Output: getOutputForVariable("$foo"),
						Pipe: []Operation{
							getColumn("103"),
						},
						Comment: "",
					},
				},
				Columns: map[int]Recipe{},
			},
			wantErr: false,
		},
		{
			name: "Variable loaded from column with comment",
			args: args{
				source: strings.NewReader("$lala <- 101 # put column 101 into $lala"),
			},
			want: &Transformation{
				Variables: map[string]Recipe{
					"$lala": {
						Output: getOutputForVariable("$lala"),
						Pipe: []Operation{
							getColumn("101"),
						},
						Comment: "put column 101 into $lala",
					},
				},
				Columns: map[int]Recipe{},
			},
			wantErr: false,
		},
		{
			name: "Variable loaded from column with comment trims spaces from ends",
			args: args{
				source: strings.NewReader("$lala <- 101 #   put column 101 into $lala  "),
			},
			want: &Transformation{
				Variables: map[string]Recipe{
					"$lala": {
						Output: getOutputForVariable("$lala"),
						Pipe: []Operation{
							getColumn("101"),
						},
						Comment: "put column 101 into $lala",
					},
				},
				Columns: map[int]Recipe{},
			},
			wantErr: false,
		},
		{
			name: "column loaded from column",
			args: args{
				source: strings.NewReader("1 <- 2"),
			},
			want: &Transformation{
				Variables: map[string]Recipe{},
				Columns: map[int]Recipe{
					1: {
						Output: getOutputForColumn("1"),
						Pipe: []Operation{
							getColumn("2"),
						},
						Comment: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "column loaded from column with comment",
			args: args{
				source: strings.NewReader("1 <- 2 # move col 1 to 2"),
			},
			want: &Transformation{
				Variables: map[string]Recipe{},
				Columns: map[int]Recipe{
					1: {
						Output: getOutputForColumn("1"),
						Pipe: []Operation{
							getColumn("2"),
						},
						Comment: "move col 1 to 2",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "column from literal with comment",
			args: args{
				source: strings.NewReader("3 <- \"banana\" # col 3 is always banana"),
			},
			want: &Transformation{
				Variables: map[string]Recipe{},
				Columns: map[int]Recipe{
					3: {
						Output: getOutputForColumn("3"),
						Pipe: []Operation{
							getLiteral("banana"),
						},
						Comment: "col 3 is always banana",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "assign variable to column with comment",
			args: args{
				source: strings.NewReader("13 <- $salad # put whatever $salad has into column 13"),
			},
			want: &Transformation{
				Variables: map[string]Recipe{},
				Columns: map[int]Recipe{
					13: {
						Output: getOutputForColumn("13"),
						Pipe: []Operation{
							getVariable("$salad"),
						},
						Comment: "put whatever $salad has into column 13",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "assign variable to variable",
			args: args{
				source: strings.NewReader("$foo <- $bar #assign $bar to $foo"),
			},
			want: &Transformation{
				Variables: map[string]Recipe{
					"$foo": {
						Output: getOutputForVariable("$foo"),
						Pipe: []Operation{
							getVariable("$bar"),
						},
						Comment: "assign $bar to $foo",
					},
				},
				Columns: map[int]Recipe{},
			},
			wantErr: false,
		},
		{
			name: "join 2 columns with comment",
			args: args{
				source: strings.NewReader("1 <- 2 + 3 #concat columns 2 and 3 into 1"),
			},
			want: &Transformation{
				Variables: map[string]Recipe{},
				Columns: map[int]Recipe{
					1: {
						Output: getOutputForColumn("1"),
						Pipe: []Operation{
							getColumn("2"),
							getJoinWithPlaceholder(),
							getColumn("3"),
						},
						Comment: "concat columns 2 and 3 into 1",
					},
				},
				Placeholder: "",
			},
			wantErr: false,
		},
		{
			name: "join two literals",
			args: args{source: strings.NewReader("$foo <- \"foo\" + \"bar\"")},
			want: &Transformation{
				Variables: map[string]Recipe{
					"$foo": {
						Output: getOutputForVariable("$foo"),
						Pipe: []Operation{
							getLiteral("foo"),
							getJoinWithPlaceholder(),
							getLiteral("bar"),
						},
					},
				},
				Columns: map[int]Recipe{},
			},
			wantErr: false,
		},
		{
			name: "prepend column with variable",
			args: args{source: strings.NewReader("12 <- 12 + $foo # tack $foo on the back of column 12")},
			want: &Transformation{
				Variables: map[string]Recipe{},
				Columns: map[int]Recipe{
					12: {
						Output: getOutputForColumn("12"),
						Pipe: []Operation{
							getColumn("12"),
							getJoinWithPlaceholder(),
							getVariable("$foo"),
						},
						Comment: "tack $foo on the back of column 12",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "Just a column is an error",
			args:    args{source: strings.NewReader("5")},
			wantErr: true,
		},
		{
			name:    "Just a variable is an error",
			args:    args{source: strings.NewReader("$foo")},
			wantErr: true,
		},
		{
			name:    "column must be followed by assign or error",
			args:    args{source: strings.NewReader("4 = 3")},
			wantErr: true,
		},
		{
			name:    "variable must be followed by assign or error",
			args:    args{source: strings.NewReader("$foo = $bar")},
			wantErr: true,
		},
		{
			name: "uppercase a column into a variable",
			args: args{source: strings.NewReader("$big <- 6 -> uppercase")},
			want: &Transformation{
				Variables: map[string]Recipe{
					"$big": {
						Output: getOutputForVariable("$big"),
						Pipe: []Operation{
							getColumn("6"),
							getFunction("uppercase", []Argument{
								{
									Type:  "placeholder",
									Value: "?",
								},
							}),
						},
					},
				},
				Columns: map[int]Recipe{},
			},
			wantErr: false,
		},
		{
			name:    "error when parsing function but no closing paren",
			args:    args{source: strings.NewReader("4 <- error(")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error when parsing function with comment in parens",
			args:    args{source: strings.NewReader("5 <- nope(#this does not work")},
			want:    nil,
			wantErr: true,
		},
		{
			name: "function with no args to column with comment",
			args: args{source: strings.NewReader("$date <- today #store today's date in a variable")},
			want: &Transformation{
				Variables: map[string]Recipe{
					"$date": {
						Output: getOutputForVariable("$date"),
						Pipe: []Operation{
							{
								Name: "today",
								Arguments: []Argument{
									{
										Type:  "placeholder",
										Value: "?",
									},
								},
							},
						},
						Comment: "store today's date in a variable",
					},
				},
				Columns: map[int]Recipe{},
			},
			wantErr: false,
		},
		{
			name: "function with args feeds to variable",
			args: args{source: strings.NewReader("$name <- fake(\"name\") # random name goes in")},
			want: &Transformation{
				Variables: map[string]Recipe{
					"$name": {
						Output: getOutputForVariable("$name"),
						Pipe: []Operation{
							getFunction("fake", []Argument{
								{
									Type:  "literal",
									Value: "name",
								},
								{
									Type:  "placeholder",
									Value: "?",
								},
							}),
						},
						Comment: "random name goes in",
					},
				},
				Columns: map[int]Recipe{},
			},
			wantErr: false,
		},
		{
			name: "function with explicit placeholder",
			args: args{source: strings.NewReader("13 <- fake(?)")},
			want: &Transformation{
				Variables: map[string]Recipe{},
				Columns: map[int]Recipe{
					13: {
						Output: getOutputForColumn("13"),
						Pipe: []Operation{
							{
								Name: "fake",
								Arguments: []Argument{
									{
										Type:  "placeholder",
										Value: "?",
									},
								},
							},
						},
						Comment: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "function with multiple args",
			args: args{source: strings.NewReader("$total <- add(2, $apples)")},
			want: &Transformation{
				Variables: map[string]Recipe{
					"$total": {
						Output: getOutputForVariable("$total"),
						Pipe: []Operation{
							{
								Name: "add",
								Arguments: []Argument{
									{
										Type:  "column",
										Value: "2",
									},
									{
										Type:  "variable",
										Value: "$apples",
									},
									{
										Type:  "placeholder",
										Value: "?",
									},
								},
							},
						},
						Comment: "",
					},
				},
				Columns:     map[int]Recipe{},
				Placeholder: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.name)
			got, err := Parse(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
