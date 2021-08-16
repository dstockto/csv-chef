package csv

import (
	"encoding/csv"
	"os"
)

func NewCsvSource(filename string) (*csv.Reader, func() error, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	reader := csv.NewReader(file)

	return reader, file.Close, nil
}

func NewOutputSource(filename string) (*csv.Writer, func() error, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, nil, err
	}
	writer := csv.NewWriter(file)
	return writer, file.Close, nil
}
