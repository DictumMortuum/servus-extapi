package main

import (
	"fmt"
	"sort"

	"github.com/DictumMortuum/servus-extapi/pkg/model"
)

func ProcessLocations(req *model.Map, res *model.Map) error {
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

	rs := []Mechanic{}
	err = db.Select(&rs, fmt.Sprintf(`
		select
			l.id,
			l.name,
			count(*) count
		from
			tboardgamestats s,
			tboardgameplays p,
			tlocation l
		where
			s.play_id = p.id and
			p.location_id = l.id and
			s.player_id = ? %s
		group by
			1
	`, q), id)
	if err != nil {
		return err
	}

	sort.Slice(rs, func(i int, j int) bool {
		return rs[i].Count >= rs[j].Count
	})

	res.Set("locations", rs)
	return nil
}
