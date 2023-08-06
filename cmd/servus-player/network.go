package main

import (
	"fmt"
	"sort"

	"github.com/DictumMortuum/servus-extapi/pkg/model"
)

type Network struct {
	Id     int64  `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Count  int    `json:"count,omitempty"`
	Url    string `json:"url,omitempty"`
	Hidden bool   `json:"hidden,omitempty"`
}

func GetNetwork(req *model.Map, res *model.Map) error {
	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	db, err := req.GetDB()
	if err != nil {
		return err
	}

	rs := []Network{}
	err = db.Select(&rs, fmt.Sprintf(`
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
			play_id in (select play_id from tboardgamestats where player_id = ?) and
			s.player_id = pl.id and
			s.play_id = p.id and
			s.player_id <> ?
			%s
		group by
			1
		order by
		  4
	`, YearConstraint(req, "and")), id, id)
	if err != nil {
		return err
	}

	sort.Slice(rs, func(i int, j int) bool {
		return rs[i].Count >= rs[j].Count
	})

	n, err := req.GetInt64("n")
	if err != nil {
		return err
	}

	res.Set("network_length", len(rs))

	network := []Network{}
	for _, player := range rs {
		if !player.Hidden {
			network = append(network, player)
		}
	}

	if len(network) > int(n) {
		res.Set("network", network[0:n])
	} else {
		res.Set("network", network)
	}

	return nil
}
