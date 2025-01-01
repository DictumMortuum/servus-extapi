package scrape

import (
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
)

func ScrapeAvalon() (map[string]any, []map[string]any, error) {
	store_id := int64(25)
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.AllowedDomains("avalongames.gr"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML(".product-layout", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".price-normal")
		old_price := e.ChildText(".price-old")

		if raw_price == "" {
			raw_price = e.ChildText(".price-new")
		}

		if old_price == "" {
			old_price = raw_price
		}

		var stock int

		if !hasClass(e, ".out-of-stock") {
			stock = 0
		} else {
			stock = 2
		}

		item := map[string]any{
			"name":           e.ChildText(".name"),
			"store_id":       store_id,
			"store_thumb":    e.ChildAttr(".product-img div img", "src"),
			"stock":          stock,
			"price":          getPrice(raw_price),
			"original_price": getPrice(old_price),
			"url":            e.Request.AbsoluteURL(e.ChildAttr(".name a", "href")),
		}

		rs = append(rs, item)
	})

	collector.OnHTML(".pagination-results", func(e *colly.HTMLElement) {
		pageCount := getPages(e.Text)
		for i := 2; i <= pageCount; i++ {
			link := fmt.Sprintf("https://avalongames.gr/index.php?route=product/category&path=59&limit=100&page=%d", i)

			if Debug {
				log.Println("Visiting: ", link)
			}

			collector.Visit(link)
		}
	})

	collector.Visit("https://avalongames.gr/index.php?route=product/category&path=59&limit=100")
	collector.Wait()

	return map[string]interface{}{
		"name":    "Avalon",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
