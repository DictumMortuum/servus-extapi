package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/urfave/cli/v2"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	DB, err := db.DatabaseX()
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "list",
				Action: func(ctx *cli.Context) error {
					rs, err := NetworkPlayersInYear(DB)
					if err != nil {
						return err
					}

					for _, item := range rs {
						fmt.Println(item.Name, item.Email, item.Count)
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
