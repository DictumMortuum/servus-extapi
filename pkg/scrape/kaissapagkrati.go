// https://kaissapagrati.gr/product-category/boardgames/

package scrape

import (
	"log"

	"github.com/gocolly/colly/v2"
)

func ScrapeKaissaPagkrati() (map[string]any, []map[string]any, error) {
	store_id := int64(31)
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.AllowedDomains("kaissapagrati.gr"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML(".product-grid-container", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".woocommerce-Price-amount")

		// var stock int

		// if e.ChildText(".release-date") != "" {
		// 	stock = 1
		// } else {
		// 	if !childHasClass(e, "div.stock", "unavailable") {
		// 		stock = 0
		// 	} else {
		// 		stock = 2
		// 	}
		// }

		item := map[string]any{
			"name":           e.ChildText(".product-grid-title"),
			"store_id":       store_id,
			"store_thumb":    e.ChildAttr(".product-grid-image img", "data-src"),
			"stock":          -1,
			"price":          getPrice(raw_price),
			"original_price": getPrice(raw_price), // TODO
			"url":            e.Attr("href"),
		}

		rs = append(rs, item)
	})

	collector.OnHTML("a.next", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if Debug {
			log.Println("Visiting: " + link)
		}

		collector.Visit(link)
	})

	collector.Visit("https://kaissapagrati.gr/product-category/boardgames/")
	collector.Wait()

	return map[string]interface{}{
		"name":    "Kaissa Pagkrati",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
