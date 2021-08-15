package main

import (
	"fmt"
	"github.com/dstockto/csv-transform/recipe"
)



func main() {
	fmt.Println("CSV Transform")

	recipe := recipe.Recipe{
		Output: recipe.Output{
			Type:  "column",
			Value: "1",
		},
		Pipe: []recipe.Operation{
			{Name: "lower", Arguments: []recipe.Argument{
				{
					Type:  "column",
					Value: "1",
				},
			}},
		},
	}

	csvInput, err := NewCsvInput("foo.csv")
	if err != nil {

	}

	csvOutput, err := NewCsvOutput("output.csv")

	Transform(recipe, csvInput, csvOutput)
}
