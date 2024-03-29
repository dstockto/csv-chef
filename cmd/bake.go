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
	"github.com/dstockto/csv-chef/recipe"
	"github.com/google/martian/log"
	"os"

	"github.com/spf13/cobra"
)

var (
	transformLines int
	disableHeader  bool
	forceOverwrite bool
	inputFile      string
	outputFile     string
	recipeFile     string
)

// bakeCmd represents the bake command
var bakeCmd = &cobra.Command{
	Use:   "bake -i /path/to/input.csv -o /path/to/output.csv -r /path/to/recipe",
	Short: "Bake uses a recipe to transform a CSV file",
	Long: `Using a recipe file, bake allows you to transform an input file to another
output file where each output line can be manipulated according to the rules you've 
created in the recipe file. Please see the README for how to make recipes. The -f flag can be used to
overwrite the output file if it exists. The -d flag will disable processing of headers with header rules 
for the first line of the file. The -n flag can tag a number representing the maximum number of lines
to process from the input file. This can be helpful if you are testing a recipe and the input file is large.'`,
	Run: runBake,
}

func runBake(cmd *cobra.Command, args []string) {
	if inputFile == "" {
		log.Errorf("Please specify an input file path with -i or --in")
		os.Exit(1)
	}
	if outputFile == "" {
		log.Errorf("Please specify an output file path with -o or --out")
		os.Exit(1)
	}
	if recipeFile == "" {
		log.Errorf("Please specify a recipe file path with -r -or --recipe")
		os.Exit(1)
	}
	parseErrIsError, err := cmd.Flags().GetBool("parseErrorIsError")
	if err != nil {
		log.Errorf("Error reading parseErrIsError flag: %s\n", err)
		os.Exit(1)
	}

	in, err := os.Open(inputFile)
	if err != nil {
		log.Errorf("Error opening input file: %v", err)
		os.Exit(1)
	}
	defer in.Close()

	// ensure output doesn't exist, or force is specified
	if _, err := os.Stat(outputFile); err == nil && !forceOverwrite {
		log.Errorf("Output file already exists: %s", output)
		os.Exit(5)
	}

	out, err := os.Create(outputFile)
	if err != nil {
		log.Errorf("Error creating output file: %v", err)
		os.Exit(6)
	}
	defer out.Close()

	recipeFile, err := os.Open(recipeFile)
	if err != nil {
		log.Errorf("Unable to open recipe file: %v", err)
		os.Exit(6)
	}
	defer recipeFile.Close()

	transformer, err := recipe.Parse(recipeFile)
	if err != nil {
		log.Errorf("Error processing your recipe: %v", err)
		os.Exit(7)
	}

	// Don't count the header
	if transformLines > 0 && !disableHeader {
		transformLines++
	}

	result, err := transformer.Execute(csv.NewReader(in), csv.NewWriter(out), !disableHeader, transformLines, parseErrIsError)
	if err != nil {
		log.Errorf("Error during baking: %v", err)
		os.Exit(8)
	}

	fmt.Printf("Baking complete. Your output is here: %s\n\n", outputFile)
	fmt.Printf("Processed %d header lines and %d input lines\n", result.HeaderLines, result.Lines)
}

func init() {
	rootCmd.AddCommand(bakeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bakeCmd.PersistentFlags().String("foo", "", "A help for foo")
	bakeCmd.Flags().IntVarP(&transformLines, "lines", "n", -1, "-n 100")
	bakeCmd.Flags().BoolVarP(&disableHeader, "no-header", "d", false, "--no-header")
	bakeCmd.Flags().BoolVarP(&forceOverwrite, "force", "f", false, "--force (force output)")
	bakeCmd.Flags().StringVarP(&inputFile, "in", "i", "", "-i /path/to/input.csv")
	bakeCmd.Flags().StringVarP(&outputFile, "out", "o", "", "-o /path/to/output.csv")
	bakeCmd.Flags().StringVarP(&recipeFile, "recipe", "r", "", "-r /path/to/recipe.txt")
	bakeCmd.Flags().BoolP("parseErrorIsError", "p", false, "-p")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bakeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
