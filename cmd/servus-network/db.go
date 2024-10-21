package main

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
)

type Player struct {
	Id     int64                 `gorm:"primaryKey" json:"id"`
	Name   string                `json:"name"`
	Email  models.JsonNullString `json:"email"`
	Hidden bool                  `json:"hidden"`
	Count  int                   `json:"count"`
}

func NetworkPlayersInYear(DB *sqlx.DB) ([]Player, error) {
	rs := []Player{}

	sql := `
		select
			pl.id,
			CONCAT(pl.name, " ", pl.surname) name,
			pl.email email,
			pl.hidden hidden,
			count(*) count
		from
			tboardgameplayers pl,
			tboardgamestats s,
			tboardgameplays p
		where
			pl.id = s.player_id and
			s.play_id = p.id and
			pl.hidden = 0 and
			p.date >= "2024-01-01"
		group by
			1, 2, 3, 4
		order by
			5
	`

	err := DB.Select(&rs, sql)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
