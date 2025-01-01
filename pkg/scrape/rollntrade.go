package scrape

import (
	"log"

	"github.com/gocolly/colly/v2"
)

func ScrapeRollntrade() (map[string]any, []map[string]any, error) {
	store_id := int64(36)
	rs := []map[string]any{}
	detected := 0

	collector := colly.NewCollector(
		colly.AllowedDomains("rollntrade.com"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML(".product-grid-item", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".price > .amount > bdi")
		if raw_price == "" {
			raw_price = e.ChildText(".price > ins > .amount > bdi")
		}

		var stock int

		if hasClass(e, "instock") {
			stock = 0
		} else if hasClass(e, "onbackorder") {
			stock = 1
		} else if hasClass(e, "outofstock") {
			stock = 2
		}

		item := map[string]any{
			"name":           e.ChildText(".wd-entities-title a"),
			"store_id":       store_id,
			"store_thumb":    e.ChildAttr(".attachment-woocommerce_thumbnail", "src"),
			"stock":          stock,
			"price":          getPrice(raw_price),
			"original_price": getPrice(raw_price), // TODO
			"url":            e.ChildAttr(".wd-entities-title a", "href"),
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

	collector.Visit("https://rollntrade.com/product-category/board-games/")
	collector.Wait()

	return map[string]interface{}{
		"name":    "Roll 'n' Trade",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
