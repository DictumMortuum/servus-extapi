package scrape

import (
	"log"

	"github.com/gocolly/colly/v2"
)

func ScrapeCOINMilan() (map[string]any, []map[string]any, error) {
	store_id := int64(47)
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.AllowedDomains("www.milan-spiele.de"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML(".master", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".priceRow span.price")

		var stock int
		if childHasClass(e, ".delivery", "yellow") {
			stock = 1
		} else if childHasClass(e, ".delivery", "green") {
			stock = 2
		} else {
			stock = 0
		}

		name := e.ChildText(".detail h1")

		item := map[string]any{
			"name":           name,
			"store_id":       store_id,
			"store_thumb":    "https://www.milan-spiele.de/images/gmt-games-logo.jpg",
			"stock":          stock,
			"price":          getPrice(raw_price),
			"original_price": getPrice(raw_price), // TODO
			"url":            e.Request.URL.String(),
			"tag":            "COIN",
		}

		log.Println(item)

		rs = append(rs, item)
	})

	URLs := []string{
		"https://www.milan-spiele.de/-p-28604.html",
		"https://www.milan-spiele.de/fire-lake-second-edition-engl-p-17569.html",
		"https://www.milan-spiele.de/-p-34822.html",
		"https://www.milan-spiele.de/-p-34842.html",
		"https://www.milan-spiele.de/-p-33628.html",
		"https://www.milan-spiele.de/-p-27345.html",
		"https://www.milan-spiele.de/-p-27751.html",
		"https://www.milan-spiele.de/-p-32602.html",
	}

	for _, item := range URLs {
		collector.Visit(item)
	}

	collector.Wait()

	return map[string]interface{}{
		"name":    "Milan",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
