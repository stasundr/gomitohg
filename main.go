package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	mapset "github.com/deckarep/golang-set"
	"github.com/markbates/pkger"
	log "github.com/sirupsen/logrus"
	wfa "github.com/stasundr/gomitohg/bridge"
	"github.com/stasundr/gomitohg/fasta"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "mitohg"
	app.Usage = "human mtDNA haplogroup classification tool"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "input",
			Aliases: []string{"i"},
			Value:   "",
			Usage:   "Input fasta `FILE`",
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.String("input") == "" {
			log.Warn("No input fasta")
			return nil
		}
		_, err := os.Stat(c.String("input"))
		if os.IsNotExist(err) {
			log.Warn("File does not exist")
			return nil
		}

		phylotreef, err := pkger.Open("/data/phylotree17.json.gz")
		if err != nil {
			return err
		}
		defer phylotreef.Close()
		phylotreer, err := gzip.NewReader(phylotreef)
		if err != nil {
			return err
		}
		phylotreeb, err := ioutil.ReadAll(phylotreer)
		if err != nil {
			return err
		}

		var phylotree []struct {
			Haplogroup string   `json:"haplogroup"`
			Haplotype  []string `json:"haplotype"`
		}
		err = json.Unmarshal(phylotreeb, &phylotree)
		if err != nil {
			return err
		}

		reff, err := pkger.Open("/data/RSRS.fa")
		if err != nil {
			return err
		}
		defer reff.Close()

		ref, err := fasta.Read(reff)
		if err != nil {
			return err
		}

		seqf, err := os.Open(c.String("input"))
		if err != nil {
			return err
		}
		defer seqf.Close()

		seq, err := fasta.Read(seqf)
		if err != nil {
			return err
		}

		alignment, err := wfa.AffineWaveformAlign(ref[0].Sequence, seq[0].Sequence)
		if err != nil {
			return err
		}

		var position, insertionPosition int
		var currentInsertion string
		mutations := mapset.NewSet()
		for i := 0; i < len(alignment.Reference); i++ {
			r := alignment.Reference[i]
			s := alignment.Sequence[i]

			if s == '-' {
				position++
				if r != 'N' {
					mutations.Add(`${r}${position}D`)
				}
			} else if r == '-' {
				if currentInsertion == "" {
					insertionPosition = i
				}
				currentInsertion += string(s)
			} else {
				if currentInsertion != "" {
					relativeIndex := 1
					if currentInsertion[0] != alignment.Reference[insertionPosition+1] {
						relativeIndex = 2
					}
					ins := fmt.Sprintf("%d.%d%s", insertionPosition, relativeIndex, currentInsertion)
					mutations.Add(ins)
					currentInsertion = ""
				}
				position++

				if r != s {
					switch s {
					case 'M':
						mutations.Add(fmt.Sprintf("%s%dA", string(r), position))
						mutations.Add(fmt.Sprintf("%s%dC", string(r), position))
					case 'R':
						mutations.Add(fmt.Sprintf("%s%dA", string(r), position))
						mutations.Add(fmt.Sprintf("%s%dG", string(r), position))
					case 'W':
						mutations.Add(fmt.Sprintf("%s%dA", string(r), position))
						mutations.Add(fmt.Sprintf("%s%dT", string(r), position))
					case 'S':
						mutations.Add(fmt.Sprintf("%s%dC", string(r), position))
						mutations.Add(fmt.Sprintf("%s%dG", string(r), position))
					case 'Y':
						mutations.Add(fmt.Sprintf("%s%dC", string(r), position))
						mutations.Add(fmt.Sprintf("%s%dT", string(r), position))
					case 'K':
						mutations.Add(fmt.Sprintf("%s%dG", string(r), position))
						mutations.Add(fmt.Sprintf("%s%dT", string(r), position))
					case 'V':
						mutations.Add(fmt.Sprintf("%s%dA", string(r), position))
						mutations.Add(fmt.Sprintf("%s%dC", string(r), position))
						mutations.Add(fmt.Sprintf("%s%dG", string(r), position))
					case 'H':
						mutations.Add(fmt.Sprintf("%s%dA", string(r), position))
						mutations.Add(fmt.Sprintf("%s%dC", string(r), position))
						mutations.Add(fmt.Sprintf("%s%dT", string(r), position))
					case 'D':
						mutations.Add(fmt.Sprintf("%s%dA", string(r), position))
						mutations.Add(fmt.Sprintf("%s%dG", string(r), position))
						mutations.Add(fmt.Sprintf("%s%dT", string(r), position))
					case 'B':
						mutations.Add(fmt.Sprintf("%s%dC", string(r), position))
						mutations.Add(fmt.Sprintf("%s%dG", string(r), position))
						mutations.Add(fmt.Sprintf("%s%dT", string(r), position))
					case 'X':
						mutations.Add(fmt.Sprintf("%s%dA", string(r), position))
						mutations.Add(fmt.Sprintf("%s%dC", string(r), position))
						mutations.Add(fmt.Sprintf("%s%dG", string(r), position))
						mutations.Add(fmt.Sprintf("%s%dT", string(r), position))
					default:
						mutations.Add(fmt.Sprintf("%s%d%s", string(r), position, string(s)))
					}
				}
			}
		}

		var resultHg string
		resultScore := 100000
		for _, hg := range phylotree {
			hgMutations := mapset.NewSet()
			for _, m := range hg.Haplotype {
				hgMutations.Add(m)
			}
			intersection := mutations.Intersect(hgMutations)
			union := mutations.Union(hgMutations)
			score := len(union.ToSlice()) - len(intersection.ToSlice())
			if resultScore > score {
				resultScore = score
				resultHg = hg.Haplogroup
			}
		}

		log.Info(resultHg)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
}
