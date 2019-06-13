package fasta

import (
	"io/ioutil"
	"strings"
)

type Record struct {
	Name     string
	Sequence string
}

type Fasta []Record

func Read(filename string) (Fasta, error) {
	fasta := Fasta{}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return fasta, err
	}

	// todo: trim string
	data := strings.Split(string(file), ">")

	for _, rawEntry := range data[1:] {
		entry := strings.Split(rawEntry, "\n")
		record := Record{
			Name: entry[0],
			// todo: check sequence for junk
			Sequence: strings.Join(entry[1:], ""),
		}

		fasta = append(fasta, record)
	}

	return fasta, nil
}
