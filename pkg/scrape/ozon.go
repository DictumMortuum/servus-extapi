package scrape

import (
	"log"

	"github.com/gocolly/colly/v2"
)

func ScrapeOzon() (map[string]any, []map[string]any, error) {
	store_id := int64(17)
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.AllowedDomains("www.ozon.gr"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML(".products-list div.col-xs-3", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".special-price")
		if raw_price == "" {
			raw_price = e.ChildText(".price")
		}

		old_price := e.ChildText(".old-price")
		if old_price == "" {
			old_price = raw_price
		}

		item := map[string]any{
			"name":           e.ChildText(".title"),
			"store_id":       store_id,
			"store_thumb":    e.ChildAttr(".image-wrapper img", "src"),
			"stock":          0,
			"price":          getPrice(raw_price),
			"original_price": getPrice(old_price), // TODO
			"url":            e.ChildAttr(".product-box", "href"),
		}

		rs = append(rs, item)
	})

	collector.OnHTML("a.next", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if link != "javascript:;" {
			if Debug {
				log.Println("Visiting: " + link)
			}

			collector.Visit(link)
		}
	})

	collector.Visit("https://www.ozon.gr/pazl-kai-paixnidia/epitrapezia-paixnidia")
	collector.Wait()

	return map[string]interface{}{
		"name":    "Ozon",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
