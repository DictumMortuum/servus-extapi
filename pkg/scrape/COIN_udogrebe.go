package scrape

import (
	"github.com/gocolly/colly/v2"
)

func ScrapeCOINUdo() (map[string]any, []map[string]any, error) {
	store_id := int64(45)
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.AllowedDomains("ugg2nd.de"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML(".productData", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".price")

		var stock int
		if childHasClass(e, ".stockFlag", "preorderStock") {
			stock = 1
		} else if childHasClass(e, ".stockFlag", "notOnStock") {
			stock = 2
		} else {
			stock = 0
		}

		item := map[string]any{
			"name":        e.ChildText(".title"),
			"store_id":    store_id,
			"store_thumb": e.ChildAttr(".pictureBox img", "src"),
			"stock":       stock,
			"price":       getPrice(raw_price),
			"url":         e.ChildAttr(".title", "href"),
			"tag":         "COIN",
		}

		rs = append(rs, item)
	})

	collector.Visit("https://ugg2nd.de/en/tag/COIN/")
	collector.Wait()

	return map[string]interface{}{
		"name":    "Udo Grebe",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
