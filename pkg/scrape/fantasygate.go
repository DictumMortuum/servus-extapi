package scrape

import (
	"log"

	"github.com/gocolly/colly/v2"
)

func ScrapeFantasyGate() (map[string]any, []map[string]any, error) {
	store_id := int64(2)
	rs := []map[string]any{}
	detected := 0

	collector := colly.NewCollector(
		colly.AllowedDomains("fantasygate.gr"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML(".sblock4", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".jshop_price")

		var stock int
		if e.ChildText(".not_available") != "" {
			stock = 2
		} else {
			stock = 0
		}

		item := map[string]any{
			"name":        e.ChildText(".name"),
			"store_id":    store_id,
			"store_thumb": e.ChildAttr(".jshop_img", "src"),
			"stock":       stock,
			"price":       getPrice(raw_price),
			"url":         e.Request.AbsoluteURL(e.ChildAttr(".image_block a", "href")),
		}

		rs = append(rs, item)
		detected++
	})

	collector.OnHTML("a.pagenav", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))

		if Debug {
			log.Println("Visiting: " + link)
		}

		collector.Visit(link)
	})

	collector.Visit("https://fantasygate.gr/shop/ellinika-epitrapezia")
	collector.Visit("https://fantasygate.gr/strategygames")
	collector.Visit("https://fantasygate.gr/fantasygames")
	collector.Visit("https://fantasygate.gr/family-games")
	collector.Visit("https://fantasygate.gr/2-paiktes")
	collector.Visit("https://fantasygate.gr/party-games")
	collector.Visit("https://fantasygate.gr/cardgames")
	collector.Visit("https://fantasygate.gr/miniature-games")
	collector.Wait()

	return map[string]any{
		"name":    "Fantasy Gate",
		"id":      store_id,
		"scraped": detected,
	}, rs, nil
}
