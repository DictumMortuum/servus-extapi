package scrape

import (
	"log"

	"github.com/gocolly/colly/v2"
)

func ScrapeGameTheory() (map[string]any, []map[string]any, error) {
	store_id := int64(40)
	rs := []map[string]any{}
	detected := 0

	collector := colly.NewCollector(
		colly.AllowedDomains("www.gametheory.gr"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML(".card-wrapper.product-card-wrapper", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".price")

		item := map[string]any{
			"name":           e.ChildText(".full-unstyled-link"),
			"store_id":       store_id,
			"store_thumb":    e.ChildAttr(".motion-reduce", "src"),
			"stock":          0,
			"price":          getPrice(raw_price),
			"original_price": getPrice(raw_price), // TODO
			"url":            e.Request.AbsoluteURL(e.ChildAttr(".full-unstyled-link", "href")),
		}

		rs = append(rs, item)
		detected++
	})

	collector.OnHTML("a.pagination__item--prev", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))

		if Debug {
			log.Println("Visiting: " + link)
		}

		collector.Visit(link)
	})

	collector.Visit("https://www.gametheory.gr/collections/%CE%B4%CE%B5%CF%82-%CF%84%CE%B1-%CF%8C%CE%BB%CE%B1?filter.v.availability=1&filter.v.price.gte=&filter.v.price.lte=&sort_by=best-selling")
	collector.Wait()

	return map[string]any{
		"name":    "Game Theory",
		"id":      store_id,
		"scraped": detected,
	}, rs, nil
}
