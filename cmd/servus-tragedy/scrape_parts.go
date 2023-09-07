package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/DictumMortuum/servus-extapi/pkg/util"
	"github.com/gocolly/colly/v2"
	"github.com/urfave/cli/v2"
)

var (
	lib_url = "https://www.greek-language.gr/digitalResources/ancient_greek/library"
	re_url  = regexp.MustCompile(`browse.html\?text_id=(\d+)&page=(\d+)`)
)

type TragedyPart struct {
	Author  string `json:"author,omitempty"`
	Title   string `json:"title,omitempty"`
	Section string `json:"section,omitempty"`
	Part    string `json:"part,omitempty"`
	TextId  int    `json:"text_id,omitempty"`
	Page    int    `json:"page,omitempty"`
	Text    string `json:"text,omitempty"`
}

func (p TragedyPart) Browse() string {
	return fmt.Sprintf("%s/browse.html?text_id=%d&page=%d", lib_url, p.TextId, p.Page)
}

func (p *TragedyPart) GetText() string {
	collector := colly.NewCollector(
		colly.CacheDir("/tmp"),
	)

	collector.OnHTML(".left-part td p", func(e *colly.HTMLElement) {
		p.Text, _ = e.DOM.Html()
	})

	collector.Visit(p.Browse())
	collector.Wait()

	return p.Text
}

func scrapeParts(c *cli.Context) ([]TragedyPart, error) {
	author_id := c.Int("author_id")
	rs := []TragedyPart{}

	collector := colly.NewCollector(
		colly.CacheDir("/tmp"),
	)

	collector.OnHTML(".library .row", func(e *colly.HTMLElement) {
		_title := strings.Split(e.ChildText(".span4 .well h2"), " - ")
		author := _title[0]
		title := _title[1]

		e.ForEach(".span8 .well .info li", func(idx int, tragedy *colly.HTMLElement) {
			tragedy.ForEach(".unstyled li", func(idx2 int, section *colly.HTMLElement) {
				_section := tragedy.ChildText("h4")
				part := section.ChildText("a")
				url := section.ChildAttr("a", "href")
				text_id := -1
				page := -1

				refs := re_url.FindAllStringSubmatch(url, -1)
				if len(refs) > 0 {
					match := refs[0]
					text_id = util.Atoi(match[1])
					page = util.Atoi(match[2])
				}

				rs = append(rs, TragedyPart{
					Author:  author,
					Title:   title,
					Section: _section,
					Part:    part,
					TextId:  text_id,
					Page:    page,
				})
			})
		})
	})

	collector.Visit(fmt.Sprintf(lib_url+"/index.html?start=0&author_id=%d", author_id))
	collector.Visit(fmt.Sprintf(lib_url+"/index.html?start=5&author_id=%d", author_id))
	collector.Visit(fmt.Sprintf(lib_url+"/index.html?start=10&author_id=%d", author_id))
	collector.Wait()

	return rs, nil
}
