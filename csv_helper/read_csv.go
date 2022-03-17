package csv_helper

import (
	"encoding/csv"
	"fmt"
	"os"
)

func ReadCsv(filename string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Println("err")
		return nil, err
	}
	var result [][]string

	for _, line := range lines[1:] {

		data := []string{line[0], line[1]}
		result = append(result, data)
	}

	return result, nil
}
