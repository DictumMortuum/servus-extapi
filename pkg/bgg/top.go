package bgg

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/DictumMortuum/servus-extapi/pkg/util"
	"github.com/gocolly/colly/v2"
)

func GetAllBoardgames() ([]map[string]any, error) {
	tmp := map[string]map[string]any{}
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
			tmp[bgg_id[2]] = item
		}
	})

	for i := range 1711 {
		log.Println(i)
		collector.Visit(fmt.Sprintf("https://boardgamegeek.com/browse/boardgame/page/%d", i))
	}

	collector.Wait()

	for _, val := range tmp {
		rs = append(rs, val)
	}

	return rs, nil
}

func GetAllBoardgameIds() ([]string, error) {
	rs := []string{}

	collector := colly.NewCollector(
	// colly.CacheDir("/tmp/bla"),
	)

	collector.OnHTML("tr", func(e *colly.HTMLElement) {
		url := e.ChildAttr("td a.primary", "href")
		bgg_id := strings.Split(url, "/")
		log.Println(bgg_id)

		if len(bgg_id) > 3 {
			rs = append(rs, bgg_id[2])
		}
	})

	collector.OnHTML("a", func(e *colly.HTMLElement) {
		log.Println(e.Text)
	})

	for i := range 1711 {
		log.Println(i)
		time.Sleep(2000 * time.Millisecond)
		collector.Visit(fmt.Sprintf("https://boardgamegeek.com/browse/boardgame/page/%d", i))
	}

	collector.Wait()

	return util.Unique(rs), nil
}
