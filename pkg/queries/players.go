package queries

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

type Player struct {
	Id      int64   `gorm:"primaryKey" json:"id"`
	Name    string  `json:"name"`
	Surname string  `json:"surname"`
	Email   *string `json:"email"`
	Avatar  *string `json:"avatar"`
	Hidden  bool    `json:"hidden"`
}

func GetPlayers(req *model.Map, res *model.Map) error {
	DB, err := req.GetDB()
	if err != nil {
		return err
	}

	RDB, err := req.GetRedis()
	if err != nil {
		return err
	}

	url, err := req.GetString("url")
	if err != nil {
		return err
	}

	rs := []Network{}
	err = db.CachedSelect(DB, RDB, "GetPlayers"+url, &rs, fmt.Sprintf(`
		select
			pl.id,
			CONCAT(pl.name, " ", pl.surname) name,
			pl.avatar url,
			pl.hidden hidden,
			count(*) count
		from
			tboardgameplayers pl,
			tboardgamestats s,
			tboardgameplays p
		where
			pl.id = s.player_id and
			s.play_id = p.id
			%s
		group by
			1
		order by
			4
	`, db.YearConstraint(req, "and")))
	if err != nil {
		return err
	}

	sort.Slice(rs, func(i int, j int) bool {
		return rs[i].Count >= rs[j].Count
	})

	res.Set("players", rs)
	return nil
}
