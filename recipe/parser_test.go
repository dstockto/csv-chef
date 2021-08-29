package recipe

import (
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
				Headers:   map[int]Recipe{},
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
				Headers:   map[int]Recipe{},
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
				Columns:       map[int]Recipe{},
				Headers:       map[int]Recipe{},
				VariableOrder: []string{"$foo"},
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
				Columns:       map[int]Recipe{},
				Headers:       map[int]Recipe{},
				VariableOrder: []string{"$lala"},
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
				Columns:       map[int]Recipe{},
				Headers:       map[int]Recipe{},
				VariableOrder: []string{"$lala"},
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
				Headers: map[int]Recipe{},
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
				Headers: map[int]Recipe{},
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
				Headers: map[int]Recipe{},
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
				Headers: map[int]Recipe{},
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
				Columns:       map[int]Recipe{},
				Headers:       map[int]Recipe{},
				VariableOrder: []string{"$foo"},
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
				Headers: map[int]Recipe{},
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
				Columns:       map[int]Recipe{},
				Headers:       map[int]Recipe{},
				VariableOrder: []string{"$foo"},
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
				Headers: map[int]Recipe{},
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
				Columns:       map[int]Recipe{},
				Headers:       map[int]Recipe{},
				VariableOrder: []string{"$big"},
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
								Name:      "today",
								Arguments: []Argument{},
							},
						},
						Comment: "store today's date in a variable",
					},
				},
				Columns:       map[int]Recipe{},
				Headers:       map[int]Recipe{},
				VariableOrder: []string{"$date"},
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
				Columns:       map[int]Recipe{},
				Headers:       map[int]Recipe{},
				VariableOrder: []string{"$name"},
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
				Headers: map[int]Recipe{},
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
				Columns:       map[int]Recipe{},
				Headers:       map[int]Recipe{},
				VariableOrder: []string{"$total"},
			},
			wantErr: false,
		},
		{
			name: "function with multiple args piped to another",
			args: args{source: strings.NewReader("$total <- add(2, $apples) -> normalize_date(\"Y-m-d\") #what even is this?")},
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
							{
								Name: "normalize_date",
								Arguments: []Argument{
									{
										Type:  "literal",
										Value: "Y-m-d",
									},
									{
										Type:  "placeholder",
										Value: "?",
									},
								},
							},
						},
						Comment: "what even is this?",
					},
				},
				Columns:       map[int]Recipe{},
				Headers:       map[int]Recipe{},
				VariableOrder: []string{"$total"},
			},
			wantErr: false,
		},
		{
			name: "literal with embedded quotes into variable",
			args: args{source: strings.NewReader("$foo <- \"this \\\" quote \"")},
			want: &Transformation{
				Variables: map[string]Recipe{
					"$foo": {
						Output: getOutputForVariable("$foo"),
						Pipe: []Operation{
							getLiteral("this \" quote "),
						},
					},
				},
				Columns:       map[int]Recipe{},
				Headers:       map[int]Recipe{},
				VariableOrder: []string{"$foo"},
			},
			wantErr: false,
		},
		{
			name: "literal with embedded quote and backslash into column",
			args: args{source: strings.NewReader("15 <- \"quote: \\\" bs: \\\\ !\"")},
			want: &Transformation{
				Variables: map[string]Recipe{},
				Columns: map[int]Recipe{
					15: {
						Output: getOutputForColumn("15"),
						Pipe: []Operation{
							getLiteral("quote: \" bs: \\ !"),
						},
					},
				},
				Headers: map[int]Recipe{},
			},
			wantErr: false,
		},
		{
			name: "column header recipe with literal",
			args: args{source: strings.NewReader("!2 <- \"col 2\"")},
			want: &Transformation{
				Variables: map[string]Recipe{},
				Columns:   map[int]Recipe{},
				Headers: map[int]Recipe{
					2: {
						Output: getOutputForHeader("2"),
						Pipe: []Operation{
							getLiteral("col 2"),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "column header recipe with column, joins, function, etc",
			args: args{source: strings.NewReader("!13 <- 12 + \" ham \" -> uppercase #silly example")},
			want: &Transformation{
				Variables: map[string]Recipe{},
				Columns:   map[int]Recipe{},
				Headers: map[int]Recipe{
					13: {
						Output: getOutputForHeader("13"),
						Pipe: []Operation{
							getColumn("12"),
							getJoinWithPlaceholder(),
							getLiteral(" ham "),
							getFunction("uppercase", []Argument{placeholderArg()}),
						},
						Comment: "silly example",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "parse header with joined with placeholder",
			args: args{source: strings.NewReader("!1 <- 1 + ?")},
			want: &Transformation{
				Variables: map[string]Recipe{},
				Columns:   map[int]Recipe{},
				Headers: map[int]Recipe{
					1: {
						Output: getOutputForHeader("1"),
						Pipe: []Operation{
							getColumn("1"),
							getJoinWithPlaceholder(),
							getPlaceholder(),
						},
						Comment: "",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
