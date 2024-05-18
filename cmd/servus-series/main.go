package main

import (
	"log"
	"os"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/urfave/cli/v2"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "nauagio",
				Action: nauagio,
			},
			{
				Name:   "magissa",
				Action: magissa,
			},
			{
				Name:   "ioulia",
				Action: ioulia,
			},
			{
				Name:   "erotas",
				Action: erotas,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
