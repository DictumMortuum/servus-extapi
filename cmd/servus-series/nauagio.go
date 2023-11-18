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
		// rs, _ := nas.Exists(DB, episode)
		// if rs == nil {
		_, err := nas.Insert(DB, "nauagio", episode)
		if err != nil {
			return err
		}

		// payload := map[string]any{
		// 	"url":   episode,
		// 	"path":  "/volume1/plex/greek series/Nauagio/",
		// 	"owner": "dimitris@dictummortuum.com",
		// 	"group": "dimitris@dictummortuum.com",
		// }

		// err = nas.YoutubeDL(payload)
		// if err != nil {
		// 	return err
		// }

		// log.Println(payload)
		// }
	}

	return nil
}

func scrapeNauagioEpisodes(url string) []string {
	urls := []string{}

	collector := colly.NewCollector(
		colly.CacheDir("/tmp"),
	)

	collector.OnHTML("#ShowEpisodes a.prel.relative-post.blocked", func(e *colly.HTMLElement) {
		urls = append(urls, e.Attr("href"))
	})

	collector.Visit(url)
	collector.Wait()

	return urls
}
