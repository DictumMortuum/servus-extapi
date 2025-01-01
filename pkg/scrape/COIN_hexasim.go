package scrape

import (
	"github.com/gocolly/colly/v2"
)

func ScrapeCOINHexasim() (map[string]any, []map[string]any, error) {
	store_id := int64(44)
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.AllowedDomains("www.hexasim.com"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML("li", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".prix")

		var stock int
		if e.ChildText(".stock") == "Rupture de stock temporaire" {
			stock = 2
		} else if e.ChildText(".stock") == "Epuisé" {
			stock = 2
		} else {
			stock = 0
		}

		name := e.ChildText(".titre_jeu")

		item := map[string]any{
			"name":           name,
			"store_id":       store_id,
			"store_thumb":    "https://www.hexasim.com/2-16-COIN/" + e.ChildAttr(".lien_img img", "src"),
			"stock":          stock,
			"price":          getPrice(raw_price),
			"original_price": getPrice(raw_price), // TODO
			"url":            "https://www.hexasim.com/" + e.ChildAttr(".lien_img", "href"),
			"tag":            "COIN",
		}

		if name != "" {
			rs = append(rs, item)
		}
	})

	collector.Visit("https://www.hexasim.com/2-16-COIN")
	collector.Wait()

	return map[string]interface{}{
		"name":    "Hexasim",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
