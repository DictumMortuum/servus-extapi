package scrape

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"mvdan.cc/xurls/v2"
)

var (
	Debug      = false
	CacheDir   = "/tmp/scrape"
	price      = regexp.MustCompile("([0-9]+[,.][0-9]+)")
	pages      = regexp.MustCompile("([0-9]+) Σελ")
	user_agent = colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	Scrapers   = map[string]any{}
	IDs        = map[string]int64{
		"avalongames":     25,
		"boardsofmadness": 16,
		"crystallotus":    24,
		"efantasy":        8,
		"epitrapezio":     15,
		"fantasyshop":     28,
		"gameexplorers":   22,
		"gamerules":       4,
		"gamesuniverse":   20,
		"genx":            27,
		"hobbytheory":     23,
		"kaissaeu":        6,
		"kaissagames":     9,
		"meepleonboard":   10,
		"meepleplanet":    7,
		"mysterybay":      3,
		"ozon":            17,
		"politeia":        12,
		"rollnplay":       26,
		// "vgames":          5,
		"xrysoftero": 21,
		"innkeeper":  30,
		// "kaissapagkrati":  31,
		"fantasygate": 2,
		"gamescom":    18,
		"nolabelx":    32,
		// "efantasycrete":   33,
		"dragonseye":     34,
		"playce":         35,
		"rollntrade":     36,
		"mythicvault":    37,
		"kaissaioannina": 38,
		"kaissachania":   39,
		// "gametheory":     40,
		"philibert": 41,
		"fanen":     42,
		"gamershq":  43,
		"hexasim":   44,
		// "udogrebe":        45,
		"myfriendsgames": 46,
		"milan":          47,
	}
)

// func removeAccents(s string) string {
// 	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
// 	output, _, e := transform.String(t, s)
// 	if e != nil {
// 		panic(e)
// 	}
// 	return output
// }

// func unique(intSlice []int64) []int64 {
// 	keys := make(map[int64]bool)
// 	list := []int64{}
// 	for _, entry := range intSlice {
// 		if _, value := keys[entry]; !value {
// 			keys[entry] = true
// 			list = append(list, entry)
// 		}
// 	}
// 	return list
// }

func hasClass(e *colly.HTMLElement, c string) bool {
	raw := e.Attr("class")
	classes := strings.Split(raw, " ")

	for _, class := range classes {
		if class == c {
			return true
		}
	}

	return false
}

func childHasClass(e *colly.HTMLElement, child string, c string) bool {
	raw := e.ChildAttr(child, "class")
	classes := strings.Split(raw, " ")

	for _, class := range classes {
		if class == c {
			return true
		}
	}

	return false
}

func getPrice(raw string) float64 {
	raw = strings.ReplaceAll(raw, ",", ".")
	match := price.FindStringSubmatch(raw)

	if len(match) > 0 {
		price, _ := strconv.ParseFloat(match[1], 64)
		return price
	} else {
		return 0.0
	}
}

func getPages(raw string) int {
	match := pages.FindStringSubmatch(raw)

	if len(match) > 0 {
		page, _ := strconv.ParseInt(match[1], 10, 64)
		return int(page)
	} else {
		return 0
	}
}

func getURL(raw string) []string {
	xurl := xurls.Strict()
	return xurl.FindAllString(raw, -1)
}

func unique(col []map[string]any) []map[string]any {
	temp := map[string]map[string]any{}

	for _, item := range col {
		if val, ok := item["name"]; ok {
			if name, ok := val.(string); ok {
				name = strings.TrimSpace(name)

				if name == "" {
					continue
				}

				temp[name] = item
			}
		}
	}

	rs := []map[string]any{}
	for _, val := range temp {
		rs = append(rs, val)
	}

	sort.Slice(rs, func(i int, j int) bool {
		return rs[i]["name"].(string) > rs[j]["name"].(string)
	})

	return rs
}
