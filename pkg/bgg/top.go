package bgg

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

func GetAllBoardgames() ([]map[string]any, error) {
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.CacheDir("/tmp/bla"),
	)

	collector.OnHTML("tr", func(e *colly.HTMLElement) {
		name := e.ChildText("td a.primary")
		rank := e.ChildText("td.collection_rank")
		url := e.ChildAttr("td a.primary", "href")
		bgg_id := strings.Split(url, "/")

		if name != "" {
			item := map[string]any{
				"name": name,
				"rank": rank,
				"url":  url,
				"id":   bgg_id[2],
			}
			rs = append(rs, item)
		}
	})

	collector.Visit("https://boardgamegeek.com/browse/boardgame")
	for i := range 1657 {
		log.Println(i)
		collector.Visit(fmt.Sprintf("https://boardgamegeek.com/browse/boardgame/%d", i))
	}

	collector.Wait()
	return rs, nil
}
