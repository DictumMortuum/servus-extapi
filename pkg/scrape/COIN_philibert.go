package scrape

import (
	"github.com/gocolly/colly/v2"
)

func ScrapeCOINPhilibertnet() (map[string]any, []map[string]any, error) {
	store_id := int64(41)
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.AllowedDomains("www.philibertnet.com"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML(".ajax_block_product", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".price")

		var stock int
		if e.ChildAttr(".ajax_add_to_cart_button", "disabled") == "disabled" {
			stock = 2
		} else {
			stock = 0
		}

		item := map[string]any{
			"name":        e.ChildText(".s_title_block"),
			"store_id":    store_id,
			"store_thumb": e.ChildAttr(".product_img_link img", "src"),
			"stock":       stock,
			"price":       getPrice(raw_price),
			"url":         e.ChildAttr(".s_title_block a", "href"),
			"tag":         "COIN",
		}

		rs = append(rs, item)
	})

	collector.Visit("https://www.philibertnet.com/en/12004-coin")
	collector.Wait()

	return map[string]interface{}{
		"name":    "Philibert",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
