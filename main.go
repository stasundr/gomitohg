package main

import (
	"os"

	"github.com/gobuffalo/packr"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "mitohg"
	app.Usage = "make an explosive entrance"
	app.Action = func(c *cli.Context) error {
		box := packr.NewBox("./data")
		rsrsFasta, err := box.FindString("RSRS.fa")
		if err != nil {
			log.Error(err)

			return err
		}

		log.Info(rsrsFasta)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
}
