package main

import (
	"fmt"
	"sort"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
)

func ProcessLocations(req *model.Map, res *model.Map) error {
	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	DB, err := req.GetDB()
	if err != nil {
		return err
	}

	rs := []Mechanic{}
	err = DB.Select(&rs, fmt.Sprintf(`
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
	`, db.YearConstraint(req, "and")), id)
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

	if len(rs) > int(n) {
		res.Set("locations", rs[0:n])
	} else {
		res.Set("locations", rs)
	}

	return nil
}
