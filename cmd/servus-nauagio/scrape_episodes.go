package main

import (
	"github.com/gocolly/colly/v2"
)

func scrapeEpisodes(url string) []string {
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
