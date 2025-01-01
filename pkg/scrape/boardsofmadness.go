package scrape

import (
	"log"

	"github.com/gocolly/colly/v2"
)

func ScrapeBoardsOfMadness() (map[string]any, []map[string]any, error) {
	store_id := int64(16)
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.AllowedDomains("boardsofmadness.com"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML("li.product", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".price > .amount > bdi")
		if raw_price == "" {
			raw_price = e.ChildText(".price > ins > .amount > bdi")
		}

		old_price := e.ChildText(".price > del > .amount > bdi")
		if old_price == "" {
			old_price = raw_price
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
			"name":           e.ChildText(".woocommerce-loop-product__title"),
			"store_id":       store_id,
			"store_thumb":    e.ChildAttr(".attachment-woocommerce_thumbnail", "src"),
			"stock":          stock,
			"price":          getPrice(raw_price),
			"original_price": getPrice(old_price),
			"url":            e.Request.AbsoluteURL(e.ChildAttr(".woocommerce-LoopProduct-link", "href")),
		}

		rs = append(rs, item)
	})

	collector.OnHTML("a.next", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))

		if Debug {
			log.Println("Visiting: " + link)
		}

		collector.Visit(link)
	})

	collector.Visit("https://boardsofmadness.com/product-category/epitrapezia-paixnidia/")
	collector.Wait()

	return map[string]interface{}{
		"name":    "Boards of Madness",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
