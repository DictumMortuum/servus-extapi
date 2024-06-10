package scrape

import (
	"github.com/gocolly/colly/v2"
)

func ScrapeCOINFanen() (map[string]any, []map[string]any, error) {
	store_id := int64(42)
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.AllowedDomains("www.fanen.com"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML(".articles tr", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".cx[align=right] p")

		var stock int
		if childHasClass(e, ".cn p", "vergriffen") {
			stock = 2
		} else {
			stock = 0
		}

		name := e.ChildText(".cx b a")
		item := map[string]any{
			"name":        name,
			"store_id":    store_id,
			"store_thumb": "https://www.fanen.com" + e.ChildAttr(".cx a img", "src"),
			"stock":       stock,
			"price":       getPrice(raw_price),
			"url":         "https://www.fanen.com" + e.ChildAttr(".cx a", "href"),
			"tag":         "COIN",
		}

		if name != "" {
			rs = append(rs, item)
		}
	})

	collector.Visit("https://www.fanen.com/allgemeines/s--1_1500_-1_-1/katalog/2702540553/coin-series.html")
	collector.Wait()

	return map[string]interface{}{
		"name":    "Fanen",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
