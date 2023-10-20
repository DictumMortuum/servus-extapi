package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/DictumMortuum/servus/pkg/models"
)

type Score struct {
	PlayerId    int64                 `json:"player_id,omitempty"`
	BoardgameId int64                 `json:"boardgame_id,omitempty"`
	Name        string                `json:"name,omitempty"`
	Square200   models.JsonNullString `json:"url,omitempty"`
	Score       models.JsonNullInt64  `json:"score,omitempty"`
	Winner      models.JsonNullString `json:"winner,omitempty"`
	Date        time.Time             `json:"date,omitempty"`
	Cooperative models.JsonNullString `json:"cooperative,omitempty"`
}

func GetPlayerScores(req *model.Map, res *model.Map) error {
	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	DB, err := req.GetDB()
	if err != nil {
		return err
	}

	var rs []Score
	err = DB.Select(&rs, fmt.Sprintf(`
		select
			s.player_id,
			p.boardgame_id,
			g.square200,
			g.name,
			p.date,
			json_extract(s.data, '$.winner') winner,
			max(cast(json_extract(s.data, '$.score') as decimal(10, 0))) score,
			json_extract(p.play_data, '$.cooperative') cooperative
		from
			tboardgamestats s,
			tboardgameplays p,
			tboardgames g
		where
			g.id = s.boardgame_id and
			p.id = s.play_id %s
		group by 2, 1
	`, db.YearConstraint(req, "and")))
	if err != nil {
		return err
	}

	retval := map[int64]Score{}

	for _, item := range rs {
		if item.Cooperative.Valid {
			continue
		}

		if val, ok := retval[item.BoardgameId]; ok {
			if val.Score.Int64 < item.Score.Int64 {
				retval[item.BoardgameId] = item
			}
		} else {
			retval[item.BoardgameId] = item
		}
	}

	temp := []Score{}

	for _, val := range retval {
		if val.PlayerId == id {
			temp = append(temp, val)
		}
	}

	sort.Slice(temp, func(i, j int) bool {
		return temp[i].Date.After(temp[j].Date)
	})

	if len(temp) > 6 {
		res.Set("scores", temp[0:6])
	} else {
		res.Set("scores", temp)
	}

	return nil
}
