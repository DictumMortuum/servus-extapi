package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gin-gonic/gin"
)

type Boardgame struct {
	Id         int64                 `json:"id,omitempty"`
	Name       string                `json:"name,omitempty"`
	Rank       models.JsonNullInt64  `json:"rank,omitempty"`
	Count      int                   `json:"count,omitempty"`
	Square200  models.JsonNullString `json:"url,omitempty"`
	Mechanics  models.JsonArray      `json:"mechanics,omitempty"`
	Designers  models.JsonArray      `json:"designers,omitempty"`
	Categories models.JsonArray      `json:"categories,omitempty"`
	Subdomains models.JsonArray      `json:"subdomains,omitempty"`
	Families   models.JsonArray      `json:"families,omitempty"`
	Weight     float64               `json:"weight,omitempty"`
	Average    string                `json:"average,omitempty"`
}

func BindYear(c *gin.Context) {
	m, err := model.ToMap(c, "req")
	if err != nil {
		c.Error(err)
		return
	}

	type Args struct {
		Year     string `form:"year"`
		YearFlag bool   `form:"year_flag"`
	}

	var payload Args
	c.ShouldBind(&payload)

	m.Set("year", payload.Year)
	m.Set("year_flag", payload.YearFlag)
	m.Set("n", 12)
}

func YearConstraint(req *model.Map, start string) string {
	q := ""

	yearFlag, err := req.GetBool("year_flag")
	if err == nil && yearFlag {
		year, err := req.GetInt64("year")
		if err == nil && year != 0 {
			q = fmt.Sprintf("%s date >= '%d-01-01' and date < '%d-01-01'", start, year, year+1)
		}
	}

	return q
}

func GetPlayerGames(req *model.Map, res *model.Map) error {
	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	db, err := req.GetDB()
	if err != nil {
		return err
	}

	rs := []Boardgame{}
	err = db.Select(&rs, fmt.Sprintf(`
		select
			g.id,
			g.name,
			g.rank,
			g.square200,
			json_extract(g.bgg_data, '$.links.boardgamemechanic') mechanics,
			json_extract(g.bgg_data, '$.links.boardgamedesigner') designers,
			json_extract(g.bgg_data, '$.links.boardgamecategory') categories,
			json_extract(g.bgg_data, '$.links.boardgamefamily') families,
			json_extract(g.bgg_data, '$.links.boardgamesubdomain') subdomains,
			json_extract(g.bgg_data, '$.polls.boardgameweight.averageweight') weight,
			json_extract(g.bgg_data, '$.stats.average') average,
			count(*) count
		from
			tboardgames g,
			tboardgamestats s,
			tboardgameplays p
		where
			g.id = s.boardgame_id and
			s.play_id = p.id and
			s.player_id = ? %s
		group by
			1
	`, YearConstraint(req, "and")), id)
	if err != nil {
		return err
	}

	sort.Slice(rs, func(i int, j int) bool {
		return rs[i].Count >= rs[j].Count
	})

	matches := 0
	avg_weight := 0.0
	average := 0.0
	for _, stat := range rs {
		matches += stat.Count
		avg_weight += stat.Weight * float64(stat.Count)

		tmp := strings.Trim(stat.Average, "\"\\")
		s, err := strconv.ParseFloat(tmp, 64)
		if err == nil {
			average += s * float64(stat.Count)
		} else {
			fmt.Println(s, err)
		}
	}

	req.Set("stats", rs)
	res.Set("games", len(rs))
	res.Set("matches", matches)

	if matches > 0 {
		res.Set("weight", avg_weight/float64(matches))
		res.Set("average", average/float64(matches))
	} else {
		res.Set("weight", 0.0)
		res.Set("average", 0.0)
	}

	n, err := req.GetInt64("n")
	if err != nil {
		return err
	}

	if len(rs) > int(n) {
		res.Set("played", rs[0:n])
	} else {
		res.Set("played", rs)
	}

	return nil
}
