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
				Name:   "sync",
				Action: sync,
			},
			{
				Name:   "cache",
				Action: cache,
			},
			{
				Name:   "guild",
				Action: guild,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
