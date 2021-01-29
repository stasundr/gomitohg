package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/markbates/pkger"
	log "github.com/sirupsen/logrus"
	"github.com/stasundr/gomitohg/fasta"
	"github.com/urfave/cli"
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
	app.Action = func(c *cli.Context) error {
		rsrsf, err := pkger.Open("/data/RSRS.fa")
		if err != nil {
			return err
		}
		defer rsrsf.Close()

		r, err := fasta.Read(rsrsf)
		if err != nil {
			return err
		}

		sf, err := os.Open("/Users/me/dev/mtget/dryomov2015/KF874328.1_Chaedn23.fasta")
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

		log.Info(s[0].Name)

		sequence := C.CString(s[0].Sequence)
		defer C.free(unsafe.Pointer(sequence))

		fmt.Println(C.align(reference, sequence))

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
}
