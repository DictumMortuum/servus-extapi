package scrape

import (
	"log"

	"github.com/gocolly/colly/v2"
)

func ScrapeInnkeeper() (map[string]any, []map[string]any, error) {
	store_id := int64(30)
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.AllowedDomains("innkeeper.gr"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML(".product-grid-item", func(e *colly.HTMLElement) {
		var stock int

		if hasClass(e, "product_tag-pre-orders") {
			stock = 1
		} else if hasClass(e, "instock") {
			stock = 0
		} else {
			stock = 2
		}

		// log.Println(e.(".product-img-link img"))

		var url string
		if e.ChildAttr(".hover-img a", "href") != "" {
			url = e.ChildAttr(".hover-img a", "href")
		} else {
			url = e.ChildAttr("a.product-image-link", "href")
		}

		item := map[string]any{
			"name":        e.ChildText(".wd-entities-title"),
			"store_id":    store_id,
			"store_thumb": e.ChildAttr(".size-woocommerce_thumbnail", "data-lazy-src"),
			"stock":       stock,
			"price":       getPrice(e.ChildText(".woocommerce-Price-amount")),
			"url":         e.Request.AbsoluteURL(url),
		}

		rs = append(rs, item)
	})

	collector.OnHTML(".woocommerce-pagination a.next", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if Debug {
			log.Println("Visiting: " + link)
		}

		collector.Visit(link)
	})

	collector.Visit("https://innkeeper.gr/product-category/board-games/")
	collector.Wait()

	return map[string]interface{}{
		"name":    "Dragonphoenix",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
