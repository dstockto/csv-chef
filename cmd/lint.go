/*
Copyright © 2021 David Stockton <dave@davidstockton.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/dstockto/csv-chef/recipe"
	"github.com/google/martian/log"
	"github.com/spf13/cobra"
)

var (
	lintRecipeFile string
	lintInputFile  string
)

// lintCmd represents the lint command
var lintCmd = &cobra.Command{
	Use:   "lint -r /path/to/recipe [-i /path/to/input.csv]",
	Short: "Validates a recipe without producing output",
	Long: `Lint parses and validates a recipe file without transforming any data. It
reports parse errors (such as unknown functions, unterminated literals or
incorrect argument counts) and recipe validation errors (such as missing
column definitions). If an input CSV is provided with -i, lint also checks
that the recipe does not reference an input column number greater than the
number of columns in the input file's header row.`,
	Run: runLint,
}

func runLint(cmd *cobra.Command, args []string) {
	if lintRecipeFile == "" {
		log.Errorf("Please specify a recipe file path with -r or --recipe")
		os.Exit(1)
	}

	recipeReader, err := os.Open(lintRecipeFile)
	if err != nil {
		log.Errorf("Unable to open recipe file: %v", err)
		os.Exit(2)
	}
	defer func() { _ = recipeReader.Close() }()

	transformer, err := recipe.Parse(recipeReader)
	if err != nil {
		log.Errorf("Error processing your recipe: %v", err)
		os.Exit(3)
	}

	if err := transformer.ValidateRecipe(); err != nil {
		log.Errorf("Recipe validation failed: %v", err)
		os.Exit(4)
	}

	if lintInputFile != "" {
		in, err := os.Open(lintInputFile)
		if err != nil {
			log.Errorf("Error opening input file: %v", err)
			os.Exit(5)
		}
		defer func() { _ = in.Close() }()

		header, err := csv.NewReader(in).Read()
		if err == io.EOF {
			log.Errorf("Input file %s is empty", lintInputFile)
			os.Exit(6)
		}
		if err != nil {
			log.Errorf("Error reading header row from input file: %v", err)
			os.Exit(6)
		}

		headerWidth := len(header)
		maxReferenced := transformer.MaxInputColumnReferenced()
		if maxReferenced > headerWidth {
			log.Errorf(
				"Recipe references input column %d, but input file %s only has %d column(s) in its header row",
				maxReferenced, lintInputFile, headerWidth,
			)
			os.Exit(7)
		}
	}

	fmt.Println("Recipe OK")
}

func init() {
	rootCmd.AddCommand(lintCmd)

	lintCmd.Flags().StringVarP(&lintRecipeFile, "recipe", "r", "", "-r /path/to/recipe.txt")
	lintCmd.Flags().StringVarP(&lintInputFile, "in", "i", "", "-i /path/to/input.csv")
	_ = lintCmd.MarkFlagRequired("recipe")
}
