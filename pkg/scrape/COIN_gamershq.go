package scrape

import (
	"strings"

	"github.com/gocolly/colly/v2"
)

func ScrapeCOINGamersHQ() (map[string]any, []map[string]any, error) {
	store_id := int64(43)
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.AllowedDomains("gamers-hq.de"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML("div.product--box", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".product--price")

		var stock int
		if e.ChildText(".buy-btn--cart-text") == "cart" {
			stock = 0
		} else {
			stock = 2
		}

		var img string
		tmp := strings.Split(e.ChildAttr(".image--element img", "data-srcset"), ",")
		img = tmp[0]

		item := map[string]any{
			"name":        e.ChildAttr(".product--info a", "title"),
			"store_id":    store_id,
			"store_thumb": img,
			"stock":       stock,
			"price":       getPrice(raw_price),
			"url":         e.ChildAttr(".product--info a", "href"),
			"tag":         "COIN",
		}

		rs = append(rs, item)
	})

	collector.Visit("https://gamers-hq.de/en/cosim-wargames/gmt-cosim/gmt-coin-series/")
	collector.Wait()

	return map[string]interface{}{
		"name":    "Gamers HQ",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
