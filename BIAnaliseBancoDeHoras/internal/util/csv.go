package util

import (
	"encoding/csv"
	"io"
	"os"
)

// ReadCSVAll lê um CSV genérico e devolve headers + linhas.
func ReadCSVAll(path string, sep rune) ([]string, [][]string, error) {
	f, err := os.Open(path)
	if err != nil { return nil, nil, err }
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = sep
	r.LazyQuotes = true

	head, err := r.Read()
	if err != nil { return nil, nil, err }

	var rows [][]string
	for {
		rec, err := r.Read()
		if err == io.EOF { break }
		if err != nil { return nil, nil, err }
		rows = append(rows, rec)
	}
	return head, rows, nil
}
