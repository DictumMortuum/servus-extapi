package scrape

import (
	"log"

	"github.com/gocolly/colly/v2"
)

func ScrapeNoLabelX() (map[string]any, []map[string]any, error) {
	store_id := int64(32)
	rs := []map[string]any{}
	detected := 0

	collector := colly.NewCollector(
		colly.AllowedDomains("www.skroutz.gr"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML("li.cf.card.add-to-cart-cta", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".sku-link")

		var url string
		if e.ChildAttr("img[data-testid=sku-pic-img]", "src") != "" {
			url = "https" + e.ChildAttr("img[data-testid=sku-pic-img]", "src")
		} else {
			url = ""
		}

		item := map[string]any{
			"name":        e.ChildText(".card-content h2"),
			"store_id":    store_id,
			"store_thumb": url,
			"stock":       0,
			"price":       getPrice(raw_price),
			"url":         e.Request.AbsoluteURL(e.ChildAttr(".js-sku-link", "href")),
		}

		rs = append(rs, item)
		detected++
	})

	collector.OnHTML("ol li a", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))

		if Debug {
			log.Println("Visiting: " + link)
		}

		collector.Visit(link)
	})

	collector.Visit("https://www.skroutz.gr/shop/7101/No-Label-X/products.html")
	collector.Wait()

	return map[string]interface{}{
		"name":    "No Label X",
		"id":      store_id,
		"scraped": detected,
	}, rs, nil
}
