package scrape

import (
	"log"

	"github.com/gocolly/colly/v2"
)

func ScrapeKaissaIoannina() (map[string]any, []map[string]any, error) {
	store_id := int64(38)
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.AllowedDomains("www.kaissa-ioannina.com"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML(".product-grid-item", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".price ins .amount")
		if raw_price == "" {
			raw_price = e.ChildText(".price .amount")
		}

		var stock int
		if e.ChildText(".out-of-stock") == "Sold out" {
			stock = 2
		} else {
			stock = 0
		}

		item := map[string]any{
			"name":           e.ChildText(".product-title"),
			"store_id":       store_id,
			"store_thumb":    e.ChildAttr(".attachment-woocommerce_thumbnail", "src"),
			"stock":          stock,
			"price":          getPrice(raw_price),
			"original_price": getPrice(raw_price), // TODO
			"url":            e.ChildAttr("a.product-image-link", "href"),
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

	collector.Visit("https://www.kaissa-ioannina.com/en/shop/category/epitrapezia-paixnidia-en/?per_page=36")
	collector.Wait()

	return map[string]interface{}{
		"name":    "Kaissa Ioannina",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
