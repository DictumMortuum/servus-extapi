package scrape

import (
	"log"

	"github.com/gocolly/colly/v2"
)

func ScrapeEfantasyCrete() (map[string]any, []map[string]any, error) {
	store_id := int64(33)
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.AllowedDomains("efantasy-crete.gr"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML("div.product-layout", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".price-new")
		if raw_price == "" {
			raw_price = e.ChildText(".price")
		}

		item := map[string]any{
			"name":        e.ChildText(".caption h4"),
			"store_id":    store_id,
			"store_thumb": e.ChildAttr(".img-primary", "src"),
			"stock":       0,
			"price":       getPrice(raw_price),
			"url":         e.ChildAttr(".image a", "href"),
		}

		rs = append(rs, item)
	})

	collector.OnHTML(".pagination li a", func(e *colly.HTMLElement) {
		if e.Text != ">" {
			return
		}

		link := e.Attr("href")

		if Debug {
			log.Println("Visiting: " + link)
		}

		collector.Visit(link)
	})

	collector.Visit("https://efantasy-crete.gr/index.php?route=product/category&path=1")
	collector.Wait()

	return map[string]interface{}{
		"name":    "eFantasy Crete",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
