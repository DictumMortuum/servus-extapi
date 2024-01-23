package main

import (
	"fmt"
	"time"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/DictumMortuum/servus/pkg/models"
)

type Stat struct {
	Id          int64            `json:"id,omitempty"`
	BoardgameId int64            `json:"boardgame_id,omitempty"`
	PlayerId    int64            `json:"player_id,omitempty"`
	Winners     models.JsonArray `json:"winner,omitempty"`
	Score       float64          `json:"score,omitempty"`
	Date        time.Time        `json:"date,omitempty"`
}

func GetBoardgameScores(req *model.Map, res *model.Map) error {
	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	DB, err := req.GetDB()
	if err != nil {
		return err
	}

	var rs []Stat
	err = DB.Select(&rs, fmt.Sprintf(`
		select
			s.id,
			p.boardgame_id,
			s.player_id,
			json_extract(s.data, '$.score') score,
			json_extract(p.play_data, '$.winners') winners,
			p.date date
		from
			tboardgameplays p,
			tboardgamestats s
		where
			p.boardgame_id = ? and
			p.boardgame_id = s.boardgame_id and
			p.id = s.play_id %s
	`, db.YearConstraint(req, "and")), id)
	if err != nil {
		return err
	}

	average_score := 0.0
	average_winning_score := 0.0
	wins := 0
	max_score := 0.0
	var max_player_id int64
	var max_time time.Time

	for _, item := range rs {
		average_score += item.Score

		if item.Score > max_score {
			max_score = item.Score
			max_player_id = item.PlayerId
			max_time = item.Date
		}

		for _, winner := range item.Winners {
			if int(winner.(float64)) == int(item.PlayerId) {
				average_winning_score += item.Score
				wins++
			}
		}
	}

	if len(rs) > 0 {
		res.Set("average_score", average_score/float64(len(rs)))
		res.Set("average_winning_score", average_winning_score/float64(wins))
		res.Set("max_score", max_score)
		res.Set("max_player_id", max_player_id)
		res.Set("max_time", max_time)
	} else {
		res.Set("average_score", 0)
		res.Set("average_winning_score", 0)
		res.Set("max_score", 0)
		res.Set("max_player_id", -1)
		res.Set("max_time", time.Now())
	}

	return nil
}
