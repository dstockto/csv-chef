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
	"github.com/dstockto/csv-chef/csv"
	"github.com/google/martian/log"
	"github.com/spf13/cobra"
	"math/rand"
	"syreclabs.com/go/faker"
	"time"
)

// writeCmd represents the write command
var writeCmd = &cobra.Command{
	Use:   "write <file> [-n=lines (default 100)",
	Short: "Writes a CSV",
	Long:  `Write a CSV file`,
	Run:   runWrite,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide file to write csv to")
		}
		return nil
	},
}

var lines int

func runWrite(cmd *cobra.Command, args []string) {
	output, closeFunc, err := csv.NewOutputSource(args[0])
	defer closeFunc()

	if err != nil {
		log.Errorf("%+v", err)
	}

	output.Write([]string{
		"voter_id",
		"first",
		"last",
		"address",
		"city",
		"state",
		"zipcode",
		"birthdate",
		"party",
		"sent",
		"email",
	})

	for i := 0; i < lines; i++ {
		var sent string
		if rand.Intn(100) < 10 {
			sent = faker.Date().Between(time.Now().AddDate(0, 0, -10), time.Now().AddDate(0, 0, 10)).Format("2006-01-02")
		}

		output.Write([]string{
			faker.Number().Between(100000, 99999999),
			faker.Name().FirstName(),
			faker.Name().LastName(),
			faker.Address().StreetAddress(),
			faker.Address().City(),
			faker.Address().State(),
			faker.Address().ZipCode(),
			faker.Date().Birthday(17, 99).Format("2006-01-02"),
			faker.RandomChoice([]string{
				"REP",
				"DEM",
				"",
				"IND",
				"GRN",
			}),
			sent,
			faker.Internet().Email(),
		})
	}
	output.Flush()
}

func init() {
	rootCmd.AddCommand(writeCmd)

	// Here you will define your flags and configuration settings.
	writeCmd.Flags().IntVarP(&lines, "lines", "n", 100, "Number of lines to write")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// writeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// writeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
