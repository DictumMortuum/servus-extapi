package main

import (
	"fmt"
	"sort"

	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/DictumMortuum/servus-extapi/pkg/util"
)

type Mechanic struct {
	Id    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Count int    `json:"count,omitempty"`
}

func (m *Mechanic) New(payload map[string]any) *Mechanic {
	for key, val := range payload {
		if key == "objectid" {
			m.Id = val.(string)
		} else if key == "name" {
			m.Name = val.(string)
		}
	}

	return m
}

func ProcessDesigners(req *model.Map, res *model.Map) error {
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
		for _, item := range stat.Designers {
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
		res.Set("designers", rs[0:n])
	} else {
		res.Set("designers", rs)
	}

	return nil
}
