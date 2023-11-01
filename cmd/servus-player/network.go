package main

import (
	"fmt"
	"sort"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
)

type Network struct {
	Id     int64  `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Count  int    `json:"count,omitempty"`
	Url    string `json:"url,omitempty"`
	Hidden bool   `json:"hidden"`
}

func GetNetwork(req *model.Map, res *model.Map) error {
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

	rs := []Network{}
	err = db.CachedSelect(DB, RDB, "GetNetwork"+url, &rs, fmt.Sprintf(`
		select
			pl.id,
			CONCAT(pl.name, " ", pl.surname) name,
			pl.avatar url,
			pl.hidden,
			count(*) count
		from
			tboardgamestats s,
			tboardgameplays p,
			tboardgameplayers pl
		where
			s.play_id in (select play_id from tboardgamestats where player_id = ?) and
			s.player_id = pl.id and
			s.play_id = p.id
			%s
		group by
			1
		order by
		  4
	`, db.YearConstraint(req, "and")), id)
	if err != nil {
		return err
	}

	sort.Slice(rs, func(i int, j int) bool {
		return rs[i].Count >= rs[j].Count
	})

	res.Set("network_length", len(rs))

	network := []Network{}
	for _, player := range rs {
		if !player.Hidden {
			network = append(network, player)
		}
	}

	res.Set("network", network)

	return nil
}
