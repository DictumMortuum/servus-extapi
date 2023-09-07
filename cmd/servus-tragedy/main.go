package main

import (
	"github.com/urfave/cli/v2"

	"log"
	"os"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "scrape",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "author_id",
						Value: 150,
					},
				},
				Action: func(c *cli.Context) error {
					parts, err := scrapeParts(c)
					if err != nil {
						return err
					}

					for _, part := range parts {
						log.Println(part, part.Browse(), part.GetText())
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
