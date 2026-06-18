package recipe

import (
	"reflect"
	"testing"
)

func TestLowercase(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "all caps", args: args{"BOZO"}, want: "bozo"},
		{name: "numbers", args: args{"1234"}, want: "1234"},
		{name: "mixed", args: args{"Banana"}, want: "banana"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Lowercase(tt.args.s); got != tt.want {
				t.Errorf("Lowercase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMassProcess(t *testing.T) {
	type args struct {
		incoming  []string
		processor Processor
	}
	tests := []struct {
		name    string
		args    args
		wantOut []string
	}{
		{name: "single", args: args{
			incoming:  []string{"BLOB"},
			processor: Lowercase,
		}, wantOut: []string{"blob"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := MassProcess(tt.args.incoming, tt.args.processor); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("MassProcess() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestUppercase(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "all lower", args: args{"bozo"}, want: "BOZO"},
		{name: "numbers", args: args{"1234"}, want: "1234"},
		{name: "mixed", args: args{"Banana"}, want: "BANANA"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Uppercase(tt.args.s); got != tt.want {
				t.Errorf("Lowercase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoinFunc(t *testing.T) {
	type args struct {
		p string // placeholder
		s string // string to join
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "hob",
			args: args{
				p: "hob",
				s: "goblin",
			},
			want: "hobgoblin",
		},
		{
			name: "gob",
			args: args{
				p: "gob",
				s: "hoblin",
			},
			want: "gobhoblin",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinFunc(tt.args.p)(tt.args.s); got != tt.want {
				t.Errorf("Lowercase() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestNoDigits(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "no digits", args: args{"bozo"}, want: "bozo"},
		{name: "numbers", args: args{"1234"}, want: ""},
		{name: "mixed", args: args{"a1b2c3"}, want: "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NoDigits(tt.args.s); got != tt.want {
				t.Errorf("Lowercase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSanitizeField(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"leading equals", "=1+1", "'=1+1"},
		{"leading plus", "+1", "'+1"},
		{"leading minus", "-1", "'-1"},
		{"leading at", "@SUM(A1)", "'@SUM(A1)"},
		{"leading tab", "\tvalue", "'\tvalue"},
		{"leading carriage return", "\rvalue", "'\rvalue"},
		{"normal text", "hello", "hello"},
		{"empty string", "", ""},
		{"equals not at start", "a=b", "a=b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeField(tt.input); got != tt.want {
				t.Errorf("SanitizeField(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestCoalesce(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want string
	}{
		{"first non-empty", "apple", "banana", "apple"},
		{"first empty", "", "banana", "banana"},
		{"both empty", "", "", ""},
		{"first whitespace kept", " ", "banana", " "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Coalesce(tt.a, tt.b)
			if err != nil {
				t.Fatalf("Coalesce() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Coalesce(%q, %q) = %q, want %q", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestNth(t *testing.T) {
	tests := []struct {
		name      string
		delimiter string
		index     string
		input     string
		want      string
		wantErr   bool
	}{
		{name: "first field", delimiter: ",", index: "1", input: "a,b,c", want: "a"},
		{name: "middle field", delimiter: ",", index: "2", input: "a,b,c", want: "b"},
		{name: "last field", delimiter: ",", index: "3", input: "a,b,c", want: "c"},
		{name: "beyond fields", delimiter: ",", index: "5", input: "a,b,c", want: ""},
		{name: "multichar delimiter", delimiter: "::", index: "2", input: "a::b::c", want: "b"},
		{name: "empty delimiter", delimiter: "", index: "1", input: "abc", wantErr: true},
		{name: "non-int index", delimiter: ",", index: "x", input: "a,b", wantErr: true},
		{name: "zero index", delimiter: ",", index: "0", input: "a,b", wantErr: true},
		{name: "negative index", delimiter: ",", index: "-1", input: "a,b", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Nth(tt.delimiter, tt.index, tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Nth() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Nth(%q, %q, %q) = %q, want %q", tt.delimiter, tt.index, tt.input, got, tt.want)
			}
		})
	}
}

func TestPadLeft(t *testing.T) {
	tests := []struct {
		name    string
		width   string
		pad     string
		input   string
		want    string
		wantErr bool
	}{
		{name: "pad with zeros", width: "5", pad: "0", input: "42", want: "00042"},
		{name: "already long enough", width: "2", pad: "0", input: "42", want: "42"},
		{name: "longer than width", width: "1", pad: "0", input: "42", want: "42"},
		{name: "multichar pad trimmed left", width: "4", pad: "ab", input: "x", want: "babx"},
		{name: "runes", width: "4", pad: "é", input: "ab", want: "ééab"},
		{name: "zero width", width: "0", pad: "0", input: "ab", want: "ab"},
		{name: "non-int width", width: "x", pad: "0", input: "ab", wantErr: true},
		{name: "negative width", width: "-1", pad: "0", input: "ab", wantErr: true},
		{name: "empty pad", width: "5", pad: "", input: "ab", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PadLeft(tt.width, tt.pad, tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("PadLeft() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("PadLeft(%q, %q, %q) = %q, want %q", tt.width, tt.pad, tt.input, got, tt.want)
			}
		})
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		name    string
		width   string
		pad     string
		input   string
		want    string
		wantErr bool
	}{
		{name: "pad with spaces", width: "5", pad: " ", input: "ab", want: "ab   "},
		{name: "already long enough", width: "2", pad: "0", input: "ab", want: "ab"},
		{name: "longer than width", width: "1", pad: "0", input: "ab", want: "ab"},
		{name: "multichar pad trimmed right", width: "5", pad: "ab", input: "x", want: "xabab"},
		{name: "runes", width: "4", pad: "é", input: "ab", want: "abéé"},
		{name: "zero width", width: "0", pad: "0", input: "ab", want: "ab"},
		{name: "non-int width", width: "x", pad: "0", input: "ab", wantErr: true},
		{name: "negative width", width: "-1", pad: "0", input: "ab", wantErr: true},
		{name: "empty pad", width: "5", pad: "", input: "ab", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PadRight(tt.width, tt.pad, tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("PadRight() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("PadRight(%q, %q, %q) = %q, want %q", tt.width, tt.pad, tt.input, got, tt.want)
			}
		})
	}
}

func TestTitleCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "simple", input: "hello world", want: "Hello World"},
		{name: "mixed case", input: "hELLO wORLD", want: "Hello World"},
		{name: "all caps", input: "HELLO WORLD", want: "Hello World"},
		{name: "collapses whitespace", input: "  hello   world  ", want: "Hello World"},
		{name: "empty", input: "", want: ""},
		{name: "single word", input: "banana", want: "Banana"},
		{name: "unicode", input: "éric", want: "Éric"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TitleCase(tt.input)
			if err != nil {
				t.Fatalf("TitleCase() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("TitleCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRegexReplace(t *testing.T) {
	tests := []struct {
		name        string
		pattern     string
		replacement string
		input       string
		want        string
		wantErr     bool
	}{
		{name: "digits to hash", pattern: "[0-9]+", replacement: "#", input: "a1b22c", want: "a#b#c"},
		{name: "capture group", pattern: "(\\w+)@(\\w+)", replacement: "$2.$1", input: "user@host", want: "host.user"},
		{name: "no match", pattern: "z+", replacement: "x", input: "abc", want: "abc"},
		{name: "invalid pattern", pattern: "(", replacement: "x", input: "abc", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RegexReplace(tt.pattern, tt.replacement, tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RegexReplace() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("RegexReplace(%q, %q, %q) = %q, want %q", tt.pattern, tt.replacement, tt.input, got, tt.want)
			}
		})
	}
}

func TestSubstring(t *testing.T) {
	tests := []struct {
		name    string
		start   string
		length  string
		input   string
		want    string
		wantErr bool
	}{
		{name: "middle", start: "2", length: "3", input: "abcdef", want: "bcd"},
		{name: "from start", start: "1", length: "2", input: "abcdef", want: "ab"},
		{name: "length beyond end clamps", start: "4", length: "10", input: "abcdef", want: "def"},
		{name: "start beyond end", start: "10", length: "3", input: "abc", want: ""},
		{name: "zero length", start: "1", length: "0", input: "abc", want: ""},
		{name: "runes", start: "1", length: "2", input: "éîça", want: "éî"},
		{name: "non-int start", start: "x", length: "2", input: "abc", wantErr: true},
		{name: "start zero", start: "0", length: "2", input: "abc", wantErr: true},
		{name: "negative start", start: "-1", length: "2", input: "abc", wantErr: true},
		{name: "non-int length", start: "1", length: "x", input: "abc", wantErr: true},
		{name: "negative length", start: "1", length: "-2", input: "abc", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Substring(tt.start, tt.length, tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Substring() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Substring(%q, %q, %q) = %q, want %q", tt.start, tt.length, tt.input, got, tt.want)
			}
		})
	}
}
