package main

import (
	"fmt"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/DictumMortuum/servus/pkg/models"
)

func GetBoardgameDetail(req *model.Map, res *model.Map) error {
	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	DB, err := req.GetDB()
	if err != nil {
		return err
	}

	var rs Boardgame
	err = DB.QueryRowx(`
		select
			g.id,
			g.name,
			g.rank,
			g.square200,
			json_extract(g.bgg_data, '$.polls.boardgameweight.averageweight') weight,
			json_extract(g.bgg_data, '$.stats.average') average
		from
			tboardgames g
		where
			g.id = ?
	`, id).StructScan(&rs)
	if err != nil {
		return err
	}

	res.Set("boardgame", rs)
	return nil
}

type Play struct {
	Id          int64                 `json:"id,omitempty"`
	BoardgameId int64                 `json:"boardgame_id,omitempty"`
	Name        string                `json:"name,omitempty"`
	Square200   models.JsonNullString `json:"url,omitempty"`
	Winners     models.JsonArray      `json:"winners,omitempty"`
	Players     models.JsonArray      `json:"players,omitempty"`
	Cooperative models.JsonNullString `json:"cooperative,omitempty"`
}

func GetBoardgamePlays(req *model.Map, res *model.Map) error {
	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	DB, err := req.GetDB()
	if err != nil {
		return err
	}

	var rs []Play
	err = DB.Select(&rs, fmt.Sprintf(`
		select
			p.id,
			p.boardgame_id,
			g.name,
			g.square200,
			json_extract(p.play_data, '$.winners') winners,
			json_extract(p.play_data, '$.players') players,
			json_extract(p.play_data, '$.cooperative') cooperative
		from
			tboardgameplays p,
			tboardgames g
		where
			p.boardgame_id = ? and
			g.id = p.boardgame_id %s
	`, db.YearConstraint(req, "and")), id)
	if err != nil {
		return err
	}

	player_counts := map[int]int{}
	player_wins := map[int64]int{}

	for _, play := range rs {
		var winners []int64
		var players []int64

		err := play.Winners.Unmarshal(&winners)
		if err != nil {
			return err
		}

		err = play.Players.Unmarshal(&players)
		if err != nil {
			return err
		}

		n := len(players)
		if _, ok := player_counts[n]; !ok {
			player_counts[n] = 1
		}
		player_counts[n]++

		for _, winner := range winners {
			if _, ok := player_wins[winner]; !ok {
				player_wins[winner] = 0
			}

			player_wins[winner]++
		}
	}

	res.Set("player_counts", player_counts)
	res.Set("player_wins", player_wins)
	res.Set("count", len(rs))

	return nil
}
