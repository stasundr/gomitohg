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
		log.WithFields(log.Fields{
			"animal": "walrus",
		}).Info("A walrus appears")

		// set up a new box by giving it a (relative) path to a folder on disk:
		box := packr.NewBox("./templates")

		// Get the string representation of a file, or an error if it doesn't exist:
		html, err := box.FindString("index.html")
		if err != nil {
			log.WithFields(log.Fields{
				"animal": "walrus",
			}).Info("A walrus appears")

			return err
		}

		log.WithFields(log.Fields{
			"animal": "walrus",
		}).Info(html)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.WithFields(log.Fields{
			"animal": "walrus",
		}).Info("A walrus appears")
	}
}
