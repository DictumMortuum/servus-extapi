package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/DictumMortuum/servus/pkg/models"
)

type Boardgame struct {
	Id         int64                  `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Rank       models.JsonNullInt64   `json:"rank,omitempty"`
	Count      int                    `json:"count,omitempty"`
	Square200  models.JsonNullString  `json:"url,omitempty"`
	Mechanics  models.JsonArray       `json:"mechanics,omitempty"`
	Designers  models.JsonArray       `json:"designers,omitempty"`
	Categories models.JsonArray       `json:"categories,omitempty"`
	Subdomains models.JsonArray       `json:"subdomains,omitempty"`
	Families   models.JsonArray       `json:"families,omitempty"`
	Weight     models.JsonNullFloat64 `json:"weight,omitempty"`
	Average    models.JsonNullString  `json:"average,omitempty"`
}

func GetPlayerGames(req *model.Map, res *model.Map) error {
	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	url, err := req.GetString("url")
	if err != nil {
		return err
	}

	DB, err := req.GetDB()
	if err != nil {
		return err
	}

	RDB, err := req.GetRedis()
	if err != nil {
		return err
	}

	rs := []Boardgame{}
	err = db.CachedSelect(DB, RDB, "GetPlayerGames"+url, &rs, fmt.Sprintf(`
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
	`, db.YearConstraint(req, "and")), id)
	if err != nil {
		return err
	}

	sort.Slice(rs, func(i int, j int) bool {
		return rs[i].Count >= rs[j].Count
	})

	matches := 0
	avg_weight_matches := 0
	avg_weight := 0.0
	average_matches := 0
	average := 0.0
	for _, stat := range rs {
		matches += stat.Count

		if stat.Weight.Valid {
			avg_weight_matches += stat.Count
			avg_weight += stat.Weight.Float64 * float64(stat.Count)
		}

		if stat.Average.Valid {
			tmp := strings.Trim(stat.Average.String, "\"\\")
			s, err := strconv.ParseFloat(tmp, 64)
			if err == nil {
				average_matches += stat.Count
				average += s * float64(stat.Count)
			} else {
				fmt.Println(s, err)
			}
		}
	}

	req.Set("stats", rs)
	res.Set("games", len(rs))
	res.Set("matches", matches)

	if avg_weight_matches > 0 {
		res.Set("weight", avg_weight/float64(avg_weight_matches))
	} else {
		res.Set("weight", 0.0)
	}

	if average_matches > 0 {
		res.Set("average", average/float64(average_matches))
	} else {
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

type LatestBoardgame struct {
	Id        int64                 `json:"id,omitempty"`
	Name      string                `json:"name,omitempty"`
	Square200 models.JsonNullString `json:"url,omitempty"`
	Winners   models.JsonArray      `json:"winners,omitempty"`
	Players   models.JsonArray      `json:"players,omitempty"`
	Date      time.Time             `json:"date,omitempty"`
}

func GetLatestGames(req *model.Map, res *model.Map) error {
	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	url, err := req.GetString("url")
	if err != nil {
		return err
	}

	n, err := req.GetInt64("n")
	if err != nil {
		return err
	}

	DB, err := req.GetDB()
	if err != nil {
		return err
	}

	RDB, err := req.GetRedis()
	if err != nil {
		return err
	}

	rs := []LatestBoardgame{}
	err = db.CachedSelect(DB, RDB, "GetLatestGames"+url, &rs, fmt.Sprintf(`
		select
			g.id,
			g.name,
			g.square200,
			json_extract(p.play_data, '$.winners') winners,
			json_extract(p.play_data, '$.players') players,
			p.date date
		from
			tboardgames g,
			tboardgamestats s,
			tboardgameplays p
		where
			g.id = s.boardgame_id and
			s.play_id = p.id and
			s.player_id = ? %s
		order by
			p.date desc, p.id desc
		limit %d
	`, db.YearConstraint(req, "and"), n), id)
	if err != nil {
		return err
	}

	res.Set("latest", rs)

	return nil
}
