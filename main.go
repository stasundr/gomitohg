package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/markbates/pkger"
	log "github.com/sirupsen/logrus"
	"github.com/stasundr/gomitohg/fasta"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "mitohg"
	app.Usage = "human mtDNA haplogroup classification tool"
	app.Action = func(c *cli.Context) error {
		muscle := getEnv("MUSCLE_BIN", "muscle")
		log.Info(muscle)

		rsrsf, err := pkger.Open("/data/RSRS.fa")
		if err != nil {
			return err
		}
		defer rsrsf.Close()

		f, err := fasta.Read(rsrsf)
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
