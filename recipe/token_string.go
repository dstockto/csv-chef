// Code generated by "stringer -type=Token"; DO NOT EDIT.

package recipe

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ILLEGAL-0]
	_ = x[EOF-1]
	_ = x[WS-2]
	_ = x[NEWLINE-3]
	_ = x[COLUMN_ID-4]
	_ = x[ASSIGNMENT-5]
	_ = x[PIPE-6]
	_ = x[COMMENT-7]
	_ = x[PLACEHOLDER-8]
	_ = x[PLUS-9]
	_ = x[LITERAL-10]
	_ = x[VARIABLE-11]
	_ = x[FUNCTION-12]
	_ = x[OPEN_PAREN-13]
	_ = x[CLOSE_PAREN-14]
	_ = x[COMMA-15]
	_ = x[HEADER-16]
}

const _Token_name = "ILLEGALEOFWSNEWLINECOLUMN_IDASSIGNMENTPIPECOMMENTPLACEHOLDERPLUSLITERALVARIABLEFUNCTIONOPEN_PARENCLOSE_PARENCOMMAHEADER"

var _Token_index = [...]uint8{0, 7, 10, 12, 19, 28, 38, 42, 49, 60, 64, 71, 79, 87, 97, 108, 113, 119}

func (i Token) String() string {
	if i < 0 || i >= Token(len(_Token_index)-1) {
		return "Token(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Token_name[_Token_index[i]:_Token_index[i+1]]
}
