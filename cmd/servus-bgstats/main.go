package main

import (
	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/urfave/cli/v2"

	"log"
	"os"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "map",
				Action: mapMissing,
			},
			{
				Name:   "integrity",
				Action: integrity,
			},
			{
				Name:   "load",
				Action: load,
			},
			{
				Name:   "cooperative",
				Action: cooperative,
			},
			{
				Name:   "score",
				Action: score,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
