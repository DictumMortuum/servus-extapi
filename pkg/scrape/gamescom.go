package scrape

import (
	"log"
	"net/http"

	"github.com/DictumMortuum/servus/pkg/w3m"
	"github.com/gocolly/colly/v2"
)

func ScrapeGamescom() (map[string]any, []map[string]any, error) {
	store_id := int64(18)
	rs := []map[string]any{}

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	collector := colly.NewCollector()
	collector.WithTransport(t)

	collector.OnHTML("div.col-tile", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".ty-price")

		var stock int

		if childHasClass(e, ".stock-block div", "block_avail_status_label") {
			stock = 0
		} else {
			stock = 2
		}

		item := map[string]any{
			"name":        e.ChildText(".product-title"),
			"store_id":    store_id,
			"store_thumb": e.ChildAttr(".cm-image", "src"),
			"stock":       stock,
			"price":       getPrice(raw_price),
			"url":         e.ChildAttr(".abt-single-image", "href"),
		}

		rs = append(rs, item)
	})

	collector.OnHTML("a.ty-pagination__next", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if Debug {
			log.Println("Visiting: " + link)
		}

		local_link, _ := w3m.BypassCloudflare(link)
		collector.Visit(local_link)
	})

	local, err := w3m.BypassCloudflare("https://www.gamescom.gr/epitrapezia-el")
	if err != nil {
		return nil, nil, err
	}

	collector.Visit(local)
	collector.Wait()

	return map[string]interface{}{
		"name":    "Gamescom",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
