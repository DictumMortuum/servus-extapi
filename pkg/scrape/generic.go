package scrape

import (
	"errors"
	"log"
	"strings"

	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
)

func countStock(col []map[string]any, status int) int {
	count := 0

	for _, item := range col {
		if val, ok := item["stock"]; ok {
			if val.(int) == status {
				count++
			}
		}
	}

	return count
}

func extractText(e *colly.HTMLElement, s string) string {
	tmp := strings.Split(s, ",")

	if len(tmp) == 1 {
		if s[0] == '.' || s[0] == '#' {
			return e.ChildText(s)
		} else {
			return e.Attr(s)
		}
	} else if len(tmp) == 2 {
		return e.ChildAttr(tmp[0], tmp[1])
	}

	return ""
}

func extractCmp(e *colly.HTMLElement, s string) bool {
	cmp := strings.Split(s, "@")

	if len(cmp) != 2 {
		log.Println("comparison not defined", s)
		return false
	}

	tmp := strings.Split(cmp[1], ",")

	switch cmp[0] {
	case "childHasClass":
		{
			return childHasClass(e, tmp[0], tmp[1])
		}
	case "hasClass":
		{
			return hasClass(e, tmp[0])
		}
	default:
		{
			if len(tmp) == 2 {
				return e.ChildText(tmp[0]) == tmp[1]
			} else if len(tmp) == 3 {
				return e.ChildAttr(tmp[0], tmp[1]) == tmp[2]
			}
		}
	}

	return false
}

func extractURL(e *colly.HTMLElement, s string, isAbsolute bool) string {
	if isAbsolute {
		return e.Request.AbsoluteURL(extractText(e, s))
	} else {
		return extractText(e, s)
	}
}

type GenericScrapeRequest struct {
	ScrapeUrl model.ScrapeUrl
	Cache     bool
	ListOnly  bool
}

func GenericScrape(scraper model.Scrape, DB *sqlx.DB, req GenericScrapeRequest) (map[string]any, []map[string]any, error) {
	store_id := scraper.StoreId
	rs := []map[string]any{}
	pages := []string{req.ScrapeUrl.Url}

	collector := colly.NewCollector(
		colly.AllowedDomains(scraper.AllowedDomain),
		user_agent,
	)

	if req.Cache {
		collector.CacheDir = CacheDir
	}

	collector.OnHTML(scraper.SelItem, func(e *colly.HTMLElement) {
		raw_price := extractText(e, scraper.SelPrice)
		if scraper.SelAltPrice.Valid {
			if raw_price == "" {
				raw_price = extractText(e, scraper.SelAltPrice.String)
			}
		}

		old_price := extractText(e, scraper.SelOriginalPrice)
		if old_price == "" {
			old_price = raw_price
		}

		stock := 0

		if scraper.SelItemInstock.Valid {
			if extractCmp(e, scraper.SelItemInstock.String) {
				stock = 0
			}
		}

		if scraper.SelItemPreorder.Valid {
			if extractCmp(e, scraper.SelItemPreorder.String) {
				stock = 1
			}
		}

		if scraper.SelItemOutofstock.Valid {
			if extractCmp(e, scraper.SelItemOutofstock.String) {
				stock = 2
			}
		}

		item := map[string]any{
			"name":           extractText(e, scraper.SelName),
			"store_id":       store_id,
			"store_thumb":    extractText(e, scraper.SelItemThumb),
			"stock":          stock,
			"price":          getPrice(raw_price),
			"original_price": getPrice(old_price),
			"url":            extractURL(e, scraper.SelUrl, scraper.AbsoluteNextUrl),
		}

		rs = append(rs, item)
	})

	collector.OnHTML(scraper.SelNext, func(e *colly.HTMLElement) {
		link := extractURL(e, "href", scraper.AbsoluteNextUrl)

		if Debug {
			log.Println("Visiting: " + link)
		}

		pages = append(pages, link)
		collector.Visit(link)
	})

	collector.Visit(req.ScrapeUrl.Url)
	collector.Wait()

	uniqueRs := unique(rs)

	inserted := 0
	if !req.ListOnly {
		for _, item := range rs {
			id, err := Insert(DB, item)
			if err != nil {
				return nil, nil, errors.New("could not insert item")
			}

			if id != -1 {
				inserted++
			}
		}
	}

	return map[string]interface{}{
		"id":            req.ScrapeUrl.Id,
		"name":          scraper.StoreName,
		"store_id":      store_id,
		"scraped":       len(uniqueRs),
		"instock":       countStock(uniqueRs, 0),
		"preorder":      countStock(uniqueRs, 1),
		"outofstock":    countStock(uniqueRs, 2),
		"pages_visited": len(pages),
		"inserted":      inserted,
	}, uniqueRs, nil
}
