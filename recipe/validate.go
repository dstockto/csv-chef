package recipe

import (
	"sort"
	"strconv"
)

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

// NamedColumnsReferenced scans every recipe (Variables, Columns and Headers)
// and returns the unique header names referenced by NamedColumn arguments
// like {first_name}, sorted alphabetically. It returns an empty slice when
// no named references are used.
func (t *Transformation) NamedColumnsReferenced() []string {
	seen := map[string]bool{}

	collect := func(r Recipe) {
		for _, op := range r.Pipe {
			for _, arg := range op.Arguments {
				if arg.Type == NamedColumn {
					seen[arg.Value] = true
				}
			}
		}
	}

	for _, r := range t.Variables {
		collect(r)
	}
	for _, r := range t.Columns {
		collect(r)
	}
	for _, r := range t.Headers {
		collect(r)
	}

	names := make([]string, 0, len(seen))
	for n := range seen {
		names = append(names, n)
	}
	sort.Strings(names)

	return names
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
