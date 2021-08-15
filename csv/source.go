package csv

import (
	"encoding/csv"
	"os"
)

func NewCsvSource(filename string) (*csv.Reader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(file)

	return reader, nil
}
