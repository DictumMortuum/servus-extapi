package scrape

import (
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

func ScrapeGamesUniverse() (map[string]any, []map[string]any, error) {
	store_id := int64(20)
	rs := []map[string]any{}
	detected := 0

	collector := colly.NewCollector(
		colly.AllowedDomains("gamesuniverse.gr"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML("article.product-miniature", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".product-price")

		old_price := e.ChildText(".regular-price")
		if old_price == "" {
			old_price = raw_price
		}

		url := e.ChildAttr(".product-thumbnail", "href")
		if strings.Contains(url, "paidika") || strings.Contains(url, "ekpaideftika") || strings.Contains(url, "trapoules") {
			return
		}

		item := map[string]any{
			"name":           e.ChildText(".product-title"),
			"store_id":       store_id,
			"store_thumb":    e.ChildAttr(".thumbnail img", "data-src"),
			"stock":          0,
			"price":          getPrice(raw_price),
			"original_price": getPrice(old_price),
			"url":            url,
		}

		rs = append(rs, item)
		detected++
	})

	collector.OnHTML("a.next", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if Debug {
			log.Println("Visiting: " + link)
		}

		collector.Visit(link)
	})

	collector.Visit("https://gamesuniverse.gr/el/10-epitrapezia")
	collector.Wait()

	return map[string]any{
		"name":    "Games Universe",
		"id":      store_id,
		"scraped": detected,
	}, rs, nil
}
