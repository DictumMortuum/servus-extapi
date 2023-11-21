package main

import (
	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/nas"
	"github.com/gocolly/colly/v2"
	"github.com/urfave/cli/v2"
)

func nauagio(ctx *cli.Context) error {
	episodes := scrapeNauagioEpisodes("https://www.megatv.com/ekpompes/1110551/to-nayagio/")

	if len(episodes) == 0 {
		return nil
	}

	DB, err := db.DatabaseX()
	if err != nil {
		return err
	}
	defer DB.Close()

	for _, episode := range episodes {
		_, err := nas.Insert(DB, "Nauagio", episode)
		if err != nil {
			return err
		}
	}

	return nil
}

func scrapeNauagioEpisodes(url string) []string {
	urls := []string{}

	collector := colly.NewCollector()
	collector.OnHTML("#ShowEpisodes a.prel.relative-post.blocked", func(e *colly.HTMLElement) {
		urls = append(urls, e.Attr("href"))
	})

	collector.Visit(url)
	collector.Wait()

	return urls
}
