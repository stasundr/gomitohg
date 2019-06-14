package main

import (
	"gomitohg/fasta"
	"os"

	"github.com/gobuffalo/packr"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "mitohg"
	app.Usage = "make an explosive entrance"
	app.Action = func(c *cli.Context) error {
		muscle := getEnv("MUSCLE_BIN", "muscle")
		box := packr.NewBox("./data")
		rsrsFasta, err := box.FindString("RSRS.fa")
		if err != nil {
			log.Error(err)

			return err
		}

		log.Info(rsrsFasta)
		log.Info(muscle)

		f, err := fasta.Read("./data/RSRS.fa")
		if err != nil {
			log.Error(err)

			return err
		}

		if len(f) > 0 {
			log.Info(len(f[0].Sequence))
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
