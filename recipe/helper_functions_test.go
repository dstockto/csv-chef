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
