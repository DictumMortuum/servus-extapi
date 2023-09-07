package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/DictumMortuum/servus-extapi/pkg/util"
	"github.com/gocolly/colly/v2"
)

type TragedyPage struct {
	Roles []string
}

func scrapePages(url string) (*TragedyPage, error) {
	roles := []string{}

	collector := colly.NewCollector(
		colly.CacheDir("/tmp"),
	)

	collector.OnHTML(".left-part span.name", func(e *colly.HTMLElement) {
		roles = append(roles, strings.TrimSpace(e.Text))
	})

	collector.OnHTML(".right-part i", func(e *colly.HTMLElement) {
		log.Println("bla", e.Text)
	})

	collector.OnHTML(".left-part td p", func(e *colly.HTMLElement) {
		log.Println(e.DOM.Html())
	})

	collector.Visit(url)
	collector.Wait()

	fmt.Println(util.Unique(roles))

	return &TragedyPage{
		Roles: util.Unique(roles),
	}, nil
}
