package scrape

import (
	"log"

	"github.com/gocolly/colly/v2"
)

func ScrapeKaissaChania() (map[string]any, []map[string]any, error) {
	store_id := int64(39)
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.AllowedDomains("kaissachania.gr"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML(".product-type-simple", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".price ins .amount")
		if raw_price == "" {
			raw_price = e.ChildText(".price .amount")
		}

		var stock int
		if e.ChildText(".ast-shop-product-out-of-stock") == "Εκτός αποθέματος" {
			stock = 2
		} else {
			stock = 0
		}

		item := map[string]any{
			"name":        e.ChildText(".woocommerce-loop-product__title"),
			"store_id":    store_id,
			"store_thumb": e.ChildAttr(".attachment-woocommerce_thumbnail", "src"),
			"stock":       stock,
			"price":       getPrice(raw_price),
			"url":         e.ChildAttr(".woocommerce-LoopProduct-link.woocommerce-loop-product__link", "href"),
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

	collector.Visit("https://kaissachania.gr/product-category/boardgames")
	collector.Wait()

	return map[string]interface{}{
		"name":    "Kaissa Chania",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
