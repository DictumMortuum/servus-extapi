package scrape

import (
	"log"

	"github.com/gocolly/colly/v2"
)

func ScrapeMyFriendsGames() (map[string]any, []map[string]any, error) {
	store_id := int64(46)
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.AllowedDomains("myfriendsgames.com"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML("li.ast-grid-common-col", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".price")

		var stock int
		if e.ChildText(".ast-shop-product-out-of-stock") == "Out of stock" {
			stock = 2
		} else if e.ChildText(".berocket_better_labels_position") == "PreOrder" {
			stock = 1
		} else {
			stock = 0
		}

		item := map[string]any{
			"name":        e.ChildText(".woocommerce-loop-product__title"),
			"store_id":    store_id,
			"store_thumb": e.ChildAttr(".woocommerce-loop-product__link img", "src"),
			"stock":       stock,
			"price":       getPrice(raw_price),
			"url":         e.ChildAttr(".ast-loop-product__link", "href"),
		}

		rs = append(rs, item)
	})

	collector.OnHTML("a.next", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))

		if Debug {
			log.Println("Visiting: " + link)
		}

		collector.Visit(link)
	})

	collector.Visit("https://myfriendsgames.com/product-category/%ce%b5%cf%80%ce%b9%cf%84%cf%81%ce%b1%cf%80%ce%b5%ce%b6%ce%b9%ce%b1-%cf%80%ce%b1%ce%b9%cf%87%ce%bd%ce%b9%ce%b4%ce%b9%ce%b1/")
	collector.Wait()

	return map[string]interface{}{
		"name":    "MyFriendsGames",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
