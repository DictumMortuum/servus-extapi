package bgg

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/DictumMortuum/servus-extapi/pkg/util"
	"github.com/gocolly/colly/v2"
)

type Rating struct {
	Id     int
	Name   string
	Rating float64
}

func GetRatingsFromBgg(name string) ([]Rating, error) {
	rs := map[string]Rating{}
	retval := []Rating{}

	collector := colly.NewCollector(
		colly.AllowedDomains("boardgamegeek.com"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
		colly.CacheDir("/tmp/scrape"),
	)

	collector.OnHTML("#collectionitems tr", func(e *colly.HTMLElement) {
		name := e.ChildText(".primary")
		href := e.ChildAttr(".primary", "href")
		re := regexp.MustCompile("[0-9]+")
		id := re.FindAllString(href, -1)
		rating := util.Atof(e.ChildText(".rating"))

		if len(id) == 1 {
			rs[name] = Rating{
				Name:   name,
				Rating: rating,
				Id:     util.Atoi(id[0]),
			}
		}
	})

	collector.OnHTML(".geekpages", func(e *colly.HTMLElement) {
		raw := e.ChildText("a:last-child")
		n := util.Atoi(raw)

		for i := 2; i <= n; i++ {
			url := fmt.Sprintf("https://boardgamegeek.com/collection/user/%s?pageID=%d&rated=1", strings.Replace(name, " ", "%20", -1), i)
			collector.Visit(url)
		}
	})

	url := fmt.Sprintf("https://boardgamegeek.com/collection/user/%s?rated=1", strings.Replace(name, " ", "%20", -1))
	collector.Visit(url)
	collector.Wait()

	for _, val := range rs {
		retval = append(retval, val)
	}

	return retval, nil
}
