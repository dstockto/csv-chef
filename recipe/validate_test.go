package recipe

import "testing"

func TestMaxInputColumnReferenced(t *testing.T) {
	tests := []struct {
		name           string
		transformation *Transformation
		want           int
	}{
		{
			name:           "no column references returns 0",
			transformation: NewTransformation(),
			want:           0,
		},
		{
			name: "literal only references no columns",
			transformation: func() *Transformation {
				tr := NewTransformation()
				tr.AddOperationToColumn("1", Operation{
					Name:      "value",
					Arguments: []Argument{{Type: Literal, Value: "hello"}},
				})
				return tr
			}(),
			want: 0,
		},
		{
			name: "tracks max column across recipes",
			transformation: func() *Transformation {
				tr := NewTransformation()
				tr.AddOperationToColumn("1", Operation{
					Name:      "value",
					Arguments: []Argument{{Type: Column, Value: "3"}},
				})
				tr.AddOperationToColumn("2", Operation{
					Name: "join",
					Arguments: []Argument{
						{Type: Column, Value: "7"},
						{Type: Literal, Value: "x"},
					},
				})
				tr.AddOperationToVariable("foo", Operation{
					Name:      "value",
					Arguments: []Argument{{Type: Column, Value: "5"}},
				})
				tr.AddOperationToHeader("1", Operation{
					Name:      "value",
					Arguments: []Argument{{Type: Column, Value: "2"}},
				})
				return tr
			}(),
			want: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.transformation.MaxInputColumnReferenced(); got != tt.want {
				t.Errorf("MaxInputColumnReferenced() = %d, want %d", got, tt.want)
			}
		})
	}
}
