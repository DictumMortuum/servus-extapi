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
	Id      int64            `json:"id,omitempty"`
	Winners models.JsonArray `json:"winners,omitempty"`
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
			json_extract(p.play_data, '$.winners') winners
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
	for _, play := range rs {
		var winners []int64

		err := play.Winners.Unmarshal(&winners)
		if err != nil {
			return err
		}

		for _, winner := range winners {
			if winner == id {
				won++
			}
		}
	}

	res.Set("plays_count", len(rs))
	res.Set("plays_won", won)
	return nil
}
