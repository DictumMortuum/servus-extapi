package main

import (
	"fmt"
	"time"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/DictumMortuum/servus/pkg/models"
)

type DistinctBoardgame struct {
	Id        int64                 `json:"id,omitempty"`
	Name      string                `json:"name,omitempty"`
	Square200 models.JsonNullString `json:"url,omitempty"`
	Date      time.Time             `json:"date,omitempty"`
}

func GetDistinctGames(req *model.Map, res *model.Map) error {
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

	rs := []DistinctBoardgame{}
	err = db.CachedSelect(DB, RDB, "GetDistinctGames"+url, &rs, fmt.Sprintf(`
		select
			g.id,
			g.name,
			g.square200,
			max(p.date) date
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
		order by
			date desc
	`, db.YearConstraint(req, "and")), id)
	if err != nil {
		return err
	}

	res.Set("distinct", rs)

	return nil
}
