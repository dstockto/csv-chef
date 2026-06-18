package recipe

import "strconv"

// MaxInputColumnReferenced scans every recipe (Variables, Columns and
// Headers), every Operation in each recipe's Pipe and every Argument. For
// arguments whose Type is Column, it parses the Value as an integer and
// tracks the highest column number referenced. It returns 0 if no input
// columns are referenced.
func (t *Transformation) MaxInputColumnReferenced() int {
	max := 0

	check := func(recipes map[int]Recipe) {
		for _, r := range recipes {
			scanRecipe(r, &max)
		}
	}

	for _, r := range t.Variables {
		scanRecipe(r, &max)
	}
	check(t.Columns)
	check(t.Headers)

	return max
}

func scanRecipe(r Recipe, max *int) {
	for _, op := range r.Pipe {
		for _, arg := range op.Arguments {
			if arg.Type != Column {
				continue
			}
			col, err := strconv.Atoi(arg.Value)
			if err != nil {
				continue
			}
			if col > *max {
				*max = col
			}
		}
	}
}
