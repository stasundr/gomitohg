package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/markbates/pkger"
	log "github.com/sirupsen/logrus"
	"github.com/stasundr/gomitohg/fasta"
	"github.com/urfave/cli/v2"
)

// #cgo CFLAGS: -Iwfa_bridge -I../WFA/gap_affine
// #cgo LDFLAGS: -Lwfa_bridge -lwfabridge -L../WFA/build -lwfa
// #include <stdlib.h>
// #include <wfa_bridge/wfa_bridge.h>
import "C"

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

		rsrsf, err := pkger.Open("/data/RSRS.fa")
		if err != nil {
			return err
		}
		defer rsrsf.Close()

		r, err := fasta.Read(rsrsf)
		if err != nil {
			return err
		}

		sf, err := os.Open(c.String("input"))
		if err != nil {
			return err
		}
		defer sf.Close()

		s, err := fasta.Read(sf)
		if err != nil {
			return err
		}

		reference := C.CString(r[0].Sequence)
		defer C.free(unsafe.Pointer(reference))

		sequence := C.CString(s[0].Sequence)
		defer C.free(unsafe.Pointer(sequence))

		fmt.Println(C.GoString(C.align(reference, sequence)))

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
}
