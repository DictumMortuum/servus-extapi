package main

import (
	"log"
	"os"
	"sort"
	"strings"

	rofi "github.com/DictumMortuum/gofi"
	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/scrape"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"
)

func unique(col []map[string]any) []map[string]any {
	temp := map[string]map[string]any{}

	for _, item := range col {
		if val, ok := item["name"]; ok {
			if name, ok := val.(string); ok {
				name = strings.TrimSpace(name)
				temp[name] = item
			}
		}
	}

	rs := []map[string]any{}
	for _, val := range temp {
		rs = append(rs, val)
	}

	sort.Slice(rs, func(i int, j int) bool {
		return rs[i]["name"].(string) > rs[j]["name"].(string)
	})

	return rs
}

func scrapeSingle(db *sqlx.DB, f func() (map[string]any, []map[string]any, error)) error {
	metadata, rs, err := f()
	if err != nil {
		return err
	}

	temp := unique(rs)
	count := 0
	for _, item := range temp {
		id, err := Insert(db, item)
		if err != nil {
			return err
		}

		if id != -1 {
			count++
		}
	}

	log.Println(metadata, len(temp), len(rs), count)

	return nil
}

func listSingle(db *sqlx.DB, f func() (map[string]any, []map[string]any, error)) error {
	metadata, rs, err := f()
	if err != nil {
		return err
	}

	temp := unique(rs)
	for _, item := range temp {
		log.Println(item)
	}

	log.Println(metadata, len(temp), len(rs))

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
			{
				Name: "scrapeall",
				Action: func(ctx *cli.Context) error {
					for _, scraper := range scrape.Scrapers {
						if f, ok := scraper.(func() (map[string]any, []map[string]any, error)); ok {
							err := scrapeSingle(DB, f)
							if err != nil {
								return err
							}
						}
					}

					return nil
				},
			},
			{
				Name: "list",
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
							err := listSingle(DB, f)
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