package scrape

import (
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

func ScrapeCrystalLotus2() (map[string]any, []map[string]any, error) {
	store_id := int64(24)
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.AllowedDomains("crystallotus.eu"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML(".grid__item", func(e *colly.HTMLElement) {
		link := e.ChildAttr(".motion-reduce", "src")
		if strings.HasPrefix(link, "//") {
			link = "https:" + link
		}

		raw_price := e.ChildText(".price__sale")
		item := map[string]any{
			"name":        e.ChildText(".card-information__text"),
			"store_id":    store_id,
			"store_thumb": link,
			"stock":       0,
			"price":       getPrice(raw_price),
			"url":         "https://crystallotus.eu" + e.ChildAttr("a.card-information__text", "href"),
		}

		rs = append(rs, item)
	})

	collector.OnHTML(".pagination__list li:last-child a", func(e *colly.HTMLElement) {
		link := "https://crystallotus.eu" + e.Attr("href")

		if Debug {
			log.Println("Visiting: " + link)
		}

		collector.Visit(link)
	})

	collector.Visit("https://crystallotus.eu/collections/tabletop-games/")
	collector.Wait()

	return map[string]interface{}{
		"name":    "Crystal Lotus",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
