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
	"errors"
	"fmt"
	"github.com/google/martian/log"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var (
	withHeaders bool
	output      string
)

// identityCmd represents the identity command
var identityCmd = &cobra.Command{
	Use:   "identity",
	Short: "Creates a recipe from the given input file",
	Long: `The identity command creates a recipe that will read in and write out a file
unchanged from a given input file. This can then be used to build the recipe you want without
needing to specify all the columns (and headers, optionally) through typing. The intent is to
save you some time. To save it to a file, redirect the output to a file or provide the -o flag.`,
	Run: runIdentity,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("please provide input file")
		}
		return nil
	},
}

func runIdentity(cmd *cobra.Command, args []string) {
	// try to read file
	in, err := os.Open(args[0])
	if err != nil {
		log.Errorf("Unable to read input file: %v", err)
		os.Exit(3)
	}

	csvReader := csv.NewReader(in)
	row, err := csvReader.Read()
	if err == io.EOF {
		log.Errorf("Input CSV was empty")
		os.Exit(2)
	}
	if err != nil {
		log.Errorf("Error reading a line from CSV: %v", err)
		os.Exit(4)
	}

	var w io.Writer

	if output != "" {
		// check for existence
		if _, err := os.Stat(output); err == nil {
			log.Errorf("Output file already exists: %s", output)
			os.Exit(5)
		}

		f, err := os.Create(output)
		if err != nil {
			log.Errorf("Unable to open output file: %v", err)
			os.Exit(1)
		}
		defer f.Close()
		w = f
	} else {
		w = os.Stdout
	}

	for zeroIndex, column := range row {
		num := zeroIndex + 1
		if withHeaders {
			_, err = fmt.Fprintf(w, "!%d <- %d # %s header\n", num, num, column)
			_, err = fmt.Fprintf(w, "%d <- %d # %s\n", num, num, column)
			if err != nil {
				log.Errorf("Error writing header recipe: %v", err)
				os.Exit(10)
			}
		} else {
			_, err = fmt.Fprintf(w, "%d <- %d\n", num, num)
			if err != nil {
				log.Errorf("Error writing recipe line: %v", err)
				os.Exit(11)
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(identityCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// identityCmd.PersistentFlags().String("foo", "", "A help for foo")
	identityCmd.Flags().BoolVarP(&withHeaders, "with-headers", "w", false, "--with-headers")
	identityCmd.Flags().StringVarP(&output, "output", "o", "", "-o /path/to/output.csv")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// identityCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
