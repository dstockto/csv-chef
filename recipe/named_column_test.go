package recipe

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
)

func TestParse_NamedColumnReference(t *testing.T) {
	tr, err := Parse(strings.NewReader("1 <- {first_name}"))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	pipe := tr.Columns[1].Pipe
	if len(pipe) != 1 || len(pipe[0].Arguments) != 1 {
		t.Fatalf("unexpected pipe shape: %+v", pipe)
	}
	arg := pipe[0].Arguments[0]
	if arg.Type != NamedColumn {
		t.Errorf("arg type = %s, want NamedColumn", arg.Type)
	}
	if arg.Value != "first_name" {
		t.Errorf("arg value = %q, want %q", arg.Value, "first_name")
	}
}

func runBake(t *testing.T, recipeText, input string, processHeader bool) (string, error) {
	t.Helper()
	tr, err := Parse(strings.NewReader(recipeText))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	var buf bytes.Buffer
	_, err = tr.Execute(csv.NewReader(strings.NewReader(input)), csv.NewWriter(&buf), processHeader, -1, false)
	return buf.String(), err
}

func TestNamedColumn_ResolvesAgainstHeader(t *testing.T) {
	got, err := runBake(t, "1 <- {last}", "first,last\nAda,Lovelace\nGrace,Hopper\n", true)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	// Header row passes through (column 1 = "first"); data rows pull {last}.
	want := "first\nLovelace\nHopper\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNamedColumn_AsFunctionArgument(t *testing.T) {
	got, err := runBake(t, "1 <- uppercase({first})", "first,last\nada,lovelace\n", true)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	want := "first\nADA\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNamedColumn_UnknownNameErrors(t *testing.T) {
	_, err := runBake(t, "1 <- {nope}", "first,last\nAda,Lovelace\n", true)
	if err == nil {
		t.Fatalf("expected an error for an unknown header name, got none")
	}
	if !strings.Contains(err.Error(), "{nope}") {
		t.Errorf("error %q should mention the missing name", err.Error())
	}
}

func TestNamedColumn_RequiresHeaderProcessing(t *testing.T) {
	_, err := runBake(t, "1 <- {last}", "Ada,Lovelace\n", false)
	if err == nil {
		t.Fatalf("expected an error when header processing is disabled, got none")
	}
	if !strings.Contains(err.Error(), "header processing is disabled") {
		t.Errorf("error %q should explain that headers are required", err.Error())
	}
}

func TestNamedColumn_StripsBOMFromHeader(t *testing.T) {
	got, err := runBake(t, "1 <- {first}", "\ufefffirst,last\nAda,Lovelace\n", true)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	// The header row passes through untouched (BOM intact); {first} still
	// resolves because the lookup key is BOM-stripped.
	want := "\ufefffirst\nAda\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNamedColumn_TrimsWhitespaceOnBothSides(t *testing.T) {
	// The header cell has a trailing space and the reference has padding
	// inside the braces; both are trimmed before matching.
	got, err := runBake(t, "1 <- { last }", "first,last \nAda,Lovelace\n", true)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	want := "first\nLovelace\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNamedColumn_InVariableRecipe(t *testing.T) {
	got, err := runBake(t, "$v <- {first}\n1 <- $v", "first,last\nAda,Lovelace\n", true)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	want := "first\nAda\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNamedColumn_InHeaderRecipe(t *testing.T) {
	got, err := runBake(t, "1 <- 1\n!1 <- uppercase({first})", "first,last\nada,lovelace\n", true)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	want := "FIRST\nada\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNamedColumn_DuplicateHeaderLastWins(t *testing.T) {
	got, err := runBake(t, "1 <- {name}", "name,name\nleft,right\n", true)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	want := "name\nright\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestParse_UnterminatedNamedColumn(t *testing.T) {
	_, err := Parse(strings.NewReader("1 <- {oops"))
	if err == nil {
		t.Fatalf("expected a parse error for an unterminated named column, got none")
	}
	if !strings.Contains(err.Error(), "{oops") {
		t.Errorf("error %q should mention the unterminated reference", err.Error())
	}
}

func TestParse_UnterminatedNamedColumnAfterPipe(t *testing.T) {
	_, err := Parse(strings.NewReader("1 <- 2 -> {oops"))
	if err == nil {
		t.Fatalf("expected a parse error for an unterminated named column, got none")
	}
	if !strings.Contains(err.Error(), "{oops") {
		t.Errorf("error %q should mention the unterminated reference", err.Error())
	}
}

func TestParse_EmptyNamedColumn(t *testing.T) {
	_, err := Parse(strings.NewReader("1 <- {}"))
	if err == nil {
		t.Fatalf("expected a parse error for an empty named column, got none")
	}
	if !strings.Contains(err.Error(), "{}") {
		t.Errorf("error %q should mention the empty reference", err.Error())
	}
}

func TestParse_NamedColumnAsTargetRejected(t *testing.T) {
	_, err := Parse(strings.NewReader("{name} <- 1"))
	if err == nil {
		t.Fatalf("expected a parse error for a named column assignment target, got none")
	}
	if !strings.Contains(err.Error(), "expected column, header or variable") {
		t.Errorf("error %q should reject the named column target", err.Error())
	}
}
