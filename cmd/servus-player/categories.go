package main

import (
	"fmt"
	"sort"

	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/DictumMortuum/servus-extapi/pkg/util"
)

func ProcessCategories(req *model.Map, res *model.Map) error {
	val, ok := req.Get("stats")
	if !ok {
		return fmt.Errorf("could not find stats")
	}

	stats, err := util.ToArray(val, Boardgame{})
	if err != nil {
		return err
	}

	col := map[string]Mechanic{}
	for _, stat := range stats {
		for _, item := range stat.Categories {
			n := Mechanic{Count: stat.Count}
			n.New(item.(map[string]any))

			if entry, ok := col[n.Id]; ok {
				entry.Count += stat.Count
				col[n.Id] = entry
			} else {
				col[n.Id] = n
			}
		}
	}

	rs := []Mechanic{}
	for _, val := range col {
		rs = append(rs, val)
	}

	sort.Slice(rs, func(i int, j int) bool {
		return rs[i].Count >= rs[j].Count
	})

	n, err := req.GetInt64("n")
	if err != nil {
		return err
	}

	if len(rs) > int(n) {
		res.Set("categories", rs[0:n])
	} else {
		res.Set("categories", rs)
	}

	return nil
}
