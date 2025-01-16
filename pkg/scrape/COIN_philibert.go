package scrape

import (
	"log"

	"github.com/gocolly/colly/v2"
)

func ScrapeCOINPhilibertnet() (map[string]any, []map[string]any, error) {
	URLs := []string{
		// "https://www.philibertnet.com/en/12004-coin",
		// "https://www.philibertnet.com/en/gmt/136213-a-gest-of-robin-hood-817054012725.html",
		// "https://www.philibertnet.com/en/gmt/130362-vijayanagara-the-deccan-empires-of-medieval-india-1290-1398-817054010721.html",
		// "https://www.philibertnet.com/en/gmt/112640-fire-in-the-lake-fall-of-saigon-expansion-817054012206.html",
		// "https://www.philibertnet.com/en/gmt/123470-people-power-817054012473.html?search_query=people+power&results=272",
		// "https://www.philibertnet.com/en/10546-boardgames-with-miniatures",
		// "https://www.philibertnet.com/en/50-jeux-de-societe",
		// "https://www.philibertnet.com/en/10780-kid-games",
		// "https://www.philibertnet.com/en/882-wargames",
		"https://www.philibertnet.com/en/16388-vente-flash",
	}

	rs := []map[string]any{}
	var info map[string]any

	for _, item := range URLs {
		temp_info, temp, err := scrapeCOINPhilibertnet(item)
		if err != nil {
			return nil, nil, err
		}

		info = temp_info
		rs = append(rs, temp...)
	}

	info["scraped"] = len(rs)

	return info, rs, nil
}

func scrapeCOINPhilibertnet(url string) (map[string]any, []map[string]any, error) {
	store_id := int64(41)
	rs := []map[string]any{}

	collector := colly.NewCollector(
		colly.AllowedDomains("www.philibertnet.com"),
		user_agent,
		colly.CacheDir(CacheDir),
	)

	collector.OnHTML(".ajax_block_product", func(e *colly.HTMLElement) {
		raw_price := e.ChildText(".price")
		old_price := e.ChildText(".price-discount")
		if old_price == "" {
			old_price = raw_price
		}

		var stock int
		if e.ChildAttr(".ajax_add_to_cart_button", "disabled") == "disabled" {
			stock = 2
		} else {
			stock = 0
		}

		item := map[string]any{
			"name":           e.ChildText(".s_title_block"),
			"store_id":       store_id,
			"store_thumb":    e.ChildAttr(".product_img_link img", "src"),
			"stock":          stock,
			"price":          getPrice(raw_price),
			"original_price": getPrice(old_price),
			"url":            e.ChildAttr(".s_title_block a", "href"),
			// "tag":            "COIN",
		}

		rs = append(rs, item)
	})

	collector.OnHTML("#center_column", func(e *colly.HTMLElement) {
		raw_price := e.ChildText("#our_price_display")

		var stock int
		if e.ChildText("#availability_value") == "Out of stock" {
			stock = 2
		} else {
			stock = 0
		}

		item := map[string]any{
			"name":        e.ChildText("#product_name"),
			"store_id":    store_id,
			"store_thumb": e.ChildAttr("#bigpic", "src"),
			"stock":       stock,
			"price":       getPrice(raw_price),
			"url":         e.Request.URL.String(),
			"tag":         "Philibert",
			// "tag":         "COIN",
		}

		rs = append(rs, item)
	})

	collector.OnHTML("#pagination_next > a", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))

		if Debug {
			log.Println("Visiting: " + link)
		}

		collector.Visit(link)
	})

	collector.Visit(url)
	collector.Wait()

	return map[string]interface{}{
		"name":    "Philibert",
		"id":      store_id,
		"scraped": len(rs),
	}, rs, nil
}
