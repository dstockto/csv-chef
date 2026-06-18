package recipe

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
)

// TestStrictMode_ReadDate verifies that readDate passes unparseable input
// through by default but errors when the transformation is in strict mode.
func TestStrictMode_ReadDate(t *testing.T) {
	const recipeText = "1 <- readDate(\"2006-01-02\", 1)"
	const input = "not-a-date\n"

	// Default (lenient): the unparseable value is passed through unchanged.
	lenient, err := Parse(strings.NewReader(recipeText))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	var buf bytes.Buffer
	if _, err := lenient.Execute(csv.NewReader(strings.NewReader(input)), csv.NewWriter(&buf), false, -1, false); err != nil {
		t.Fatalf("lenient execute: unexpected error: %v", err)
	}
	if got := strings.TrimSpace(buf.String()); got != "not-a-date" {
		t.Errorf("lenient: got %q, want %q", got, "not-a-date")
	}

	// Strict: the unparseable value stops processing with an error.
	strict, err := Parse(strings.NewReader(recipeText))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	strict.Strict = true
	var strictBuf bytes.Buffer
	if _, err := strict.Execute(csv.NewReader(strings.NewReader(input)), csv.NewWriter(&strictBuf), false, -1, false); err == nil {
		t.Errorf("strict: expected an error for an unparseable date, got none")
	}
}
