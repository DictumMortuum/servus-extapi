package main

import (
	"log"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/nas"
	"github.com/gocolly/colly/v2"
	"github.com/urfave/cli/v2"
)

func erotas(ctx *cli.Context) error {
	DB, err := db.DatabaseX()
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	episodes := scrapeErotasEpisodes("https://www.star.gr/tv/seires/erotas-me-diafora/tag_1224")

	if len(episodes) == 0 {
		return nil
	}

	for _, episode := range episodes {
		part := "https://www.star.gr" + episode

		_, err := nas.Insert(DB, "Erotas Me Diafora", part)
		if err != nil {
			return err
		}
	}

	return nil
}

func scrapeErotasEpisodes(url string) []string {
	urls := []string{}

	collector := colly.NewCollector()
	collector.OnHTML("#episodes_results a.swiper-slide__link", func(e *colly.HTMLElement) {
		urls = append(urls, e.Attr("href"))
	})

	collector.Visit(url)
	collector.Wait()

	return urls
}
