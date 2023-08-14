package main

import (
	"fmt"
	"sort"

	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/DictumMortuum/servus/pkg/models"
)

type Player struct {
	Id      int64   `gorm:"primaryKey" json:"id"`
	Name    string  `json:"name"`
	Surname string  `json:"surname"`
	Email   *string `json:"email"`
	Avatar  *string `json:"avatar"`
	Hidden  bool    `json:"hidden"`
}

func GetPlayers(req *model.Map, res *model.Map) error {
	db, err := req.GetDB()
	if err != nil {
		return err
	}

	var rs []Network
	err = db.Select(&rs, fmt.Sprintf(`
		select
			pl.id,
			CONCAT(pl.name, " ", pl.surname) name,
			pl.avatar url,
			count(*) count
		from
			tboardgameplayers pl,
			tboardgamestats s,
			tboardgameplays p
		where
			pl.id = s.player_id and
			s.play_id = p.id and
			pl.hidden = 0
			%s
		group by
			1
		having
			count(*) > 5
		order by
			4
	`, YearConstraint(req, "and")))
	if err != nil {
		return err
	}

	sort.Slice(rs, func(i int, j int) bool {
		return rs[i].Count >= rs[j].Count
	})

	res.Set("players", rs)
	return nil
}

func GetPlayerDetail(req *model.Map, res *model.Map) error {
	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	db, err := req.GetDB()
	if err != nil {
		return err
	}

	var rs Player
	err = db.QueryRowx(`
		select
			*
		from
			tboardgameplayers
		where
			id = ?
	`, id).StructScan(&rs)
	if err != nil {
		return err
	}

	res.Set("player", rs)
	return nil
}

type Play struct {
	Id             int64                 `json:"id,omitempty"`
	Winners        models.JsonArray      `json:"winners,omitempty"`
	Cooperative    models.JsonNullString `json:"cooperative,omitempty"`
	CooperativeWin models.JsonNullString `json:"cooperative_win,omitempty"`
}

func GetPlayerPlays(req *model.Map, res *model.Map) error {
	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	db, err := req.GetDB()
	if err != nil {
		return err
	}

	var rs []Play
	err = db.Select(&rs, fmt.Sprintf(`
		select
			p.id,
			json_extract(p.play_data, '$.winners') winners,
			json_extract(p.play_data, '$.cooperative') cooperative,
			json_extract(s.data, '$.won') cooperative_win
		from
			tboardgamestats s,
			tboardgameplays p
		where
			s.player_id = ? and
			p.id = s.play_id %s
	`, YearConstraint(req, "and")), id)
	if err != nil {
		return err
	}

	won := 0
	count := 0
	cooperative_count := 0
	cooperative_won := 0

	for _, play := range rs {
		var winners []int64

		err := play.Winners.Unmarshal(&winners)
		if err != nil {
			return err
		}

		if play.Cooperative.Valid && play.Cooperative.String == "true" {
			cooperative_count++

			for _, winner := range winners {
				if winner == id {
					cooperative_won++
				}
			}
		} else {
			count++

			for _, winner := range winners {
				if winner == id {
					won++
				}
			}
		}
	}

	res.Set("cooperative", cooperative_count)
	res.Set("cooperative_won", cooperative_won)
	res.Set("plays_count", count)
	res.Set("plays_won", won)

	if cooperative_count > 0 {
		res.Set("cooperative_per", float64(cooperative_won)/float64(cooperative_count))
	} else {
		res.Set("cooperative_per", 0)
	}

	if count > 0 {
		res.Set("per", float64(won)/float64(count))
	} else {
		res.Set("per", 0)
	}

	return nil
}
