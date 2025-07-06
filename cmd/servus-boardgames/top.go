package main

import (
	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
)

// type topRequest struct {
// 	MinPlayers    int     `json:"min_players"`
// 	MaxPlayers    int     `json:"max_players"`
// 	IsCooperative bool    `json:"is_cooperative"`
// 	IsSoloable    bool    `json:"is_soloable"`
// 	MinRating     float64 `json:"min_rating"`
// 	MaxRating     float64 `json:"max_rating"`
// }

func GetPopularGames(req *model.Map, res *model.Map) error {
	DB, err := req.GetDB()
	if err != nil {
		return err
	}

	RDB, err := req.GetRedis()
	if err != nil {
		return err
	}

	rs := []Boardgame{}
	err = db.CachedSelect(DB, RDB, "GetPopularGames", &rs, `
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
			json_extract(g.bgg_data, '$.polls.userplayers.best[0].min') best_min_players,
			json_extract(g.bgg_data, '$.polls.userplayers.best[0].max') best_max_players,
			g.min_players,
			g.max_players
		from
			tboardgames g
		where
			g.rank <= 100 and g.rank > 0
		order by
			rank
	`)
	if err != nil {
		return err
	}

	res.Set("options", rs)

	return nil
}
