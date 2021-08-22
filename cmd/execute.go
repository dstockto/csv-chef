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
	"github.com/dstockto/csv-transform/recipe"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// ploopCmd represents the ploop command
var ploopCmd = &cobra.Command{
	Use:   "execute",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		vars := make(map[string]recipe.Recipe)
		vars["foobar"] = recipe.Recipe{
			Output: recipe.Output{
				Type:  "variable",
				Value: "foobar",
			},
			Pipe: []recipe.Operation{
				{
					Name: "value",
					Arguments: []recipe.Argument{
						{Type: "literal", Value: "Zamboni"},
					},
				},
				{
					Name: "join",
				},
				{
					Name: "value",
					Arguments: []recipe.Argument{
						{Type: "column", Value: "1"},
					},
				},
				{
					Name: "lowercase",
					Arguments: []recipe.Argument{
						{Type: "placeholder", Value: "?"},
					},
				},
			},
			Comment: "What the poop!?",
		}

		transform := recipe.Transformation{
			Variables: vars,
			Columns:   nil,
		}

		vars["pewp"] = recipe.Recipe{
			Output: recipe.Output{
				Type:  "variable",
				Value: "pewp",
			},
			Pipe: []recipe.Operation{
				{
					Name: "value",
					Arguments: []recipe.Argument{
						{
							Type:  "column",
							Value: "2",
						},
					},
				},
			},
			Comment: "And another one",
		}

		transform.Dump(os.Stdout)
		var buffer []byte

		reader := csv.NewReader(strings.NewReader("alpha,banana,carrot,delta,3billion\ndomain,elephant,fart,salami,donkey\n"))
		//writer := csv.NewWriter(bytes.NewBuffer(buffer))
		writer := csv.NewWriter(os.Stdout)
		context := recipe.LineContext{
			Variables: make(map[string]string),
			Columns:   make(map[int]string),
		}

		transform.Execute(reader, &writer, &context)
		writer.Flush()
		fmt.Println(buffer)
		fmt.Printf("%+v\n", context)
	},
}

func init() {
	rootCmd.AddCommand(ploopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ploopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ploopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
