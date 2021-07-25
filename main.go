package main

import (
	"compress/gzip"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set"
	wfa "github.com/stasundr/gomitohg/bridge"
	"github.com/stasundr/gomitohg/fasta"
	"github.com/urfave/cli/v2"
)

//go:embed data/RSRS.fa
var rsrs string

//go:embed data/phylotree17.json.gz
var phylotree string

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
			log.Println("No input fasta")
			return nil
		}
		_, err := os.Stat(c.String("input"))
		if os.IsNotExist(err) {
			log.Println("File does not exist")
			return nil
		}

		phylotreef := strings.NewReader(phylotree)
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

		reff := strings.NewReader(rsrs)
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

		mutations := mapset.NewSet()
		cigarLen := len(alignment.Opsc)
		refPos := 0
		seqPos := 0
		for i := 0; i < cigarLen; i++ {
			if alignment.Opsc[i] != 'D' {
				seqPos += alignment.Opsn[i]
			}
			if alignment.Opsc[i] != 'I' {
				refPos += alignment.Opsn[i]
			}

			// TODO: handle X D I
			switch alignment.Opsc[i] {
			case 'X':
				refL := string(ref[0].Sequence[refPos-1])
				seqL := string(seq[0].Sequence[seqPos-1])
				if refL != seqL && refL != "N" {
					// TODO: handle alphabet
					mutations.Add(fmt.Sprintf("%s%d%s", refL, refPos, seqL))
				}
			}
		}

		var resultHg string
		resultHgMutations := mapset.NewSet()
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
				resultHgMutations = hgMutations.Clone()
			}
		}

		re := regexp.MustCompile(`\d+`)
		var tail string
		ms := mutations.Difference(resultHgMutations).ToSlice()
		sort.SliceStable(ms, func(i, j int) bool {
			pi, _ := strconv.Atoi(re.FindString(ms[i].(string)))
			pj, _ := strconv.Atoi(re.FindString(ms[j].(string)))
			return pi < pj
		})
		for _, m := range ms {
			tail += " " + m.(string)
		}
		fmt.Println(resultHg + tail)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Println(err.Error())
	}
}
