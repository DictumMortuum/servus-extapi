package scrape

import (
	"log"

	"github.com/gocolly/colly/v2"
)

func ScrapePlayce() (map[string]any, []map[string]any, error) {
	store_id := int64(35)
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.AllowedDomains("shop.playce.gr"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML("div.product-grid-item", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".woocommerce-Price-amount")

		var stock int
		if hasClass(e, "outofstock") {
			stock = 2
		} else {
			stock = 0
		}

		item := map[string]any{
			"name":           e.ChildText(".wd-entities-title"),
			"store_id":       store_id,
			"store_thumb":    e.ChildAttr(".product-image-link img", "data-src"),
			"stock":          stock,
			"price":          getPrice(raw_price),
			"original_price": getPrice(raw_price), // TODO
			"url":            e.ChildAttr(".product-image-link", "href"),
		}

		rs = append(rs, item)
	})

	collector.OnHTML(".next", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if Debug {
			log.Println("Visiting: " + link)
		}

		collector.Visit(link)
	})

	collector.Visit("https://shop.playce.gr/shop/page/1/?per_page=24")
	collector.Wait()

	return map[string]interface{}{
		"name":    "Playce",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
