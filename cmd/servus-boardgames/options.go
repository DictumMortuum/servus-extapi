package main

import (
	"fmt"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
)

func GetPopularGamesForNum(req *model.Map, res *model.Map) error {
	DB, err := req.GetDB()
	if err != nil {
		return err
	}

	num, err := req.GetInt64("num")
	if err != nil {
		return err
	}

	RDB, err := req.GetRedis()
	if err != nil {
		return err
	}

	rs := []Boardgame{}
	err = db.CachedSelect(DB, RDB, fmt.Sprintf("GetPopularGamesForNum%d", num), &rs, fmt.Sprintf(`
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
			g.min_players,
			g.max_players,
			max(p.date) last_played,
			count(*) count
		from
			tboardgames g,
			tboardgameplays p
		where
			g.id = p.boardgame_id and
			g.min_players <= %d and
			g.max_players >= %d %s
		group by
			1
		order by
			rank, last_played, weight desc, count desc
	`, num, num, db.YearConstraint(req, "and")))
	if err != nil {
		return err
	}

	res.Set("options", rs)

	return nil
}
