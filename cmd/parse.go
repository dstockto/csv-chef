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
	"errors"
	"github.com/dstockto/csv-chef/recipe"
	"github.com/google/martian/log"
	"github.com/spf13/cobra"
	"os"
)

// parseCmd represents the parse command
var parseFakeCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parses a given recipe file",
	Long:  `Tests the parser by ensuring that it is reading the instructions as expected`,
	Run:   runParse,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("please provide recipe file")
		}
		return nil
	},
}

func runParse(cmd *cobra.Command, args []string) {
	recipeFile, err := os.Open(args[0])
	if err != nil {
		log.Errorf("%+v\n", err)
		os.Exit(2)
	}
	transformation, err := recipe.Parse(recipeFile)
	if err != nil {
		log.Errorf("%+v\n", err)
		os.Exit(10)
	}

	transformation.Dump(os.Stdout)
}

func init() {
	rootCmd.AddCommand(parseFakeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// parseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// parseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
