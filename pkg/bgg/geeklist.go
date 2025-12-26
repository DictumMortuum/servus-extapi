package bgg

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

func GetGeeklist(id int64) ([]map[string]any, error) {
	rs := []map[string]any{}
	url := fmt.Sprintf("https://boardgamegeek.com/geeklist/%d", id)

	collector := colly.NewCollector(
		colly.CacheDir("/tmp/bla"),
	)
	collector.OnHTML("gg-geeklist-item-ui", func(e *colly.HTMLElement) {
		name := e.ChildText("h2")
		url := e.ChildAttr("h2 a", "href")
		bgg_id := strings.Split(url, "/")

		log.Println(name)

		if name != "" {
			item := map[string]any{
				"name": name,
				"url":  url,
				"id":   bgg_id[2],
			}
			rs = append(rs, item)
		}
	})

	collector.Visit(url)
	collector.Wait()
	return rs, nil
}
