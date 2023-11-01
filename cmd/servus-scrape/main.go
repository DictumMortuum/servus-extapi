package main

import (
	"log"
	"os"

	rofi "github.com/DictumMortuum/gofi"
	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/scrape"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"
)

func scrapeSingle(db *sqlx.DB, f func() (map[string]any, []map[string]any, error)) error {
	metadata, rs, err := f()
	if err != nil {
		return err
	}

	for _, item := range rs {
		log.Println(item)

		_, err := Insert(db, item)
		if err != nil {
			return err
		}
	}

	log.Println(metadata)

	return nil
}

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
				Name: "scrape",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "store",
						Value: "",
					},
				},
				Action: func(ctx *cli.Context) error {
					var scrapers []string
					scraper := ctx.String("store")
					if scraper != "" {
						scrapers = []string{scraper}
					} else {
						opts := rofi.GofiOptions{
							Description: "scraper",
						}

						scrapers, err = rofi.FromInterface(&opts, scrape.Scrapers)
						if err != nil {
							return err
						}
					}

					for _, val := range scrapers {
						if f, ok := scrape.Scrapers[val].(func() (map[string]any, []map[string]any, error)); ok {
							err := scrapeSingle(DB, f)
							if err != nil {
								return err
							}
						}
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
