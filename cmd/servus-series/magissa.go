package main

import (
	"log"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/nas"
	"github.com/gocolly/colly/v2"
	"github.com/urfave/cli/v2"
)

func magissa(ctx *cli.Context) error {
	DB, err := db.DatabaseX()
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	episodes := scrapeMagissaEpisodes("https://www.antenna.gr/magissa")

	if len(episodes) == 0 {
		return nil
	}

	for _, episode := range episodes {
		part := "https://www.antenna.gr" + episode

		_, err := nas.Insert(DB, "Magissa", part)
		if err != nil {
			return err
		}
	}

	return nil
}

func scrapeMagissaEpisodes(url string) []string {
	urls := []string{}

	collector := colly.NewCollector()
	collector.OnHTML("#showepisodes a", func(e *colly.HTMLElement) {
		urls = append(urls, e.Attr("href"))
	})

	collector.Visit(url)
	collector.Wait()

	return urls
}
