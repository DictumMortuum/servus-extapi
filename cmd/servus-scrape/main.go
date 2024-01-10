package main

import (
	"log"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"

	rofi "github.com/DictumMortuum/gofi"
	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/scrape"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"
)

var wg sync.WaitGroup

func unique(col []map[string]any) []map[string]any {
	temp := map[string]map[string]any{}

	for _, item := range col {
		if val, ok := item["name"]; ok {
			if name, ok := val.(string); ok {
				name = strings.TrimSpace(name)

				if name == "" {
					continue
				}

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

func scrapeSingle(db *sqlx.DB, f func() (map[string]any, []map[string]any, error)) (int, error) {
	metadata, rs, err := f()
	if err != nil {
		return -1, err
	}

	temp := unique(rs)
	count := 0
	for _, item := range temp {
		id, err := Insert(db, item)
		if err != nil {
			return -1, err
		}

		if id != -1 {
			count++
		}
	}

	log.Println(metadata, len(temp), len(rs), count)

	return count, nil
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
						err := Stale(DB, scrape.IDs[val])
						if err != nil {
							return err
						}

						if f, ok := scrape.Scrapers[val].(func() (map[string]any, []map[string]any, error)); ok {
							_, err := scrapeSingle(DB, f)
							if err != nil {
								return err
							}
						}

						err = UpdateCounts(DB, scrape.IDs[val])
						if err != nil {
							return err
						}
					}

					return nil
				},
			},
			{
				Name: "scrapeall",
				Action: func(ctx *cli.Context) error {
					wg.Add(len(scrape.Scrapers))

					for key, scraper := range scrape.Scrapers {
						id := scrape.IDs[key]
						if f, ok := scraper.(func() (map[string]any, []map[string]any, error)); ok {
							go (func() {
								runtime.LockOSThread()

								err := Stale(DB, id)
								if err != nil {
									log.Println(err)
								}

								count, err := scrapeSingle(DB, f)
								if err != nil {
									log.Println(err)
								}

								if count > 0 {
									err = Cleanup(DB, id)
									if err != nil {
										log.Println(err)
									}
									log.Println("Cleaned up ", id, " with count ", count)
								}

								err = UpdateCounts(DB, id)
								if err != nil {
									log.Println(err)
								}

								defer wg.Done()
								runtime.UnlockOSThread()
							})()
						}
					}

					wg.Wait()

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
