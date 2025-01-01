// https://kaissapagrati.gr/product-category/boardgames/

package scrape

import (
	"log"

	"github.com/gocolly/colly/v2"
)

func ScrapeMythicVault() (map[string]any, []map[string]any, error) {
	store_id := int64(37)
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.AllowedDomains("mythicvault.com"),
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
	})

	collector.OnHTML("a.next", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if Debug {
			log.Println("Visiting: " + link)
		}

		collector.Visit(link)
	})

	collector.Visit("https://mythicvault.com/product-category/board-games/?per_page=36")
	collector.Wait()

	return map[string]interface{}{
		"name":    "Mythic Vault",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
