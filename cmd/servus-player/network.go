package main

import (
	"fmt"
	"sort"

	"github.com/DictumMortuum/servus-extapi/pkg/model"
)

type Network struct {
	Id    int64  `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Count int    `json:"count,omitempty"`
	Url   string `json:"url,omitempty"`
}

func GetNetwork(req *model.Map, res *model.Map) error {
	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	q := ""
	year, err := req.GetInt64("year")
	if err == nil && year != 0 {
		q = fmt.Sprintf("and date >= '%d-01-01' and date < '%d-01-01'", year, year+1)
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
	`, q), id, id)
	if err != nil {
		return err
	}

	sort.Slice(rs, func(i int, j int) bool {
		return rs[i].Count >= rs[j].Count
	})

	res.Set("network", rs)
	return nil
}
