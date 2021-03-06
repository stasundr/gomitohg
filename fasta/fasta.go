package fasta

import (
	"io"
	"io/ioutil"
	"strings"

	mapset "github.com/deckarep/golang-set"
)

type Record struct {
	Name     string
	Sequence string
}

type Fasta []Record

func Read(r io.Reader) (Fasta, error) {
	allowed := mapset.NewSet("A", "a", "C", "c", "G", "g", "T", "t", "U", "u", "R", "r", "Y", "y", "K", "k", "M", "m", "S", "s", "W", "w", "B", "b", "D", "d", "H", "h", "V", "v", "N", "n", "-", ">", "\n", " ")
	fasta := Fasta{}

	file, err := ioutil.ReadAll(r)
	if err != nil {
		return fasta, err
	}

	data := strings.Split(string(file), ">")
	for _, rawEntry := range data[1:] {
		entry := strings.Split(rawEntry, "\n")
		var sequence string
		s := strings.Join(entry[1:], "")
		l := len(s)
		for i := 0; i < l; i++ {
			c := string(s[i])
			if allowed.Contains(c) {
				sequence += c
			}
		}

		record := Record{
			Name:     entry[0],
			Sequence: sequence,
		}
		fasta = append(fasta, record)
	}

	return fasta, nil
}
