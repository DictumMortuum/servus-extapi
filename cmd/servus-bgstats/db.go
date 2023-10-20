package main

import (
	"time"

	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
)

type BGPlay struct {
	Id          int64
	BoardgameId int64
	LocationId  models.JsonNullInt64
	PlayId      models.JsonNullInt64
	PlayDate    string
	BGStats     []BGStat
}

func getBGPlays(DB *sqlx.DB) ([]BGPlay, error) {
	rs := []BGPlay{}
	q := `
		select
			p.id,
			g.boardgame_id,
			l.location_id,
			p.play_id,
			p.play_date
		from
			tbgstatsplays p,
			tbgstatsgames g,
			tbgstatslocations l
		where
			p.game_id = g.id and
			p.location_id = l.id
	`
	err := DB.Select(&rs, q)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func getMappedBGPlays(DB *sqlx.DB) ([]BGPlay, error) {
	rs := []BGPlay{}
	q := `
		select
			p.id,
			g.boardgame_id,
			l.location_id,
			p.play_id,
			p.play_date
		from
			tbgstatsplays p,
			tbgstatsgames g,
			tbgstatslocations l
		where
			p.game_id = g.id and
			p.location_id = l.id and
			p.play_id is not null
	`
	err := DB.Select(&rs, q)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func getUnmappedBGPlays(DB *sqlx.DB) ([]BGPlay, error) {
	rs := []BGPlay{}
	q := `
		select
			p.id,
			g.boardgame_id,
			l.location_id,
			p.play_id,
			p.play_date
		from
			tbgstatsplays p,
			tbgstatsgames g,
			tbgstatslocations l
		where
			p.game_id = g.id and
			p.location_id = l.id and
			p.play_id is null and
			p.play_date > '2021-01-1 00:00:00'
	`
	err := DB.Select(&rs, q)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func getSimilarBGPlays(DB *sqlx.DB, play BGPlay) ([]BGPlay, error) {
	rs := []BGPlay{}
	q := `
		select
			p.id,
			g.boardgame_id,
			l.location_id,
			p.play_id,
			p.play_date
		from
			tbgstatsplays p,
			tbgstatsgames g,
			tbgstatslocations l
		where
			p.game_id = g.id and
			g.boardgame_id = ? and
			p.location_id = l.id and
			l.location_id = ? and
			p.play_date = ?
	`
	err := DB.Select(&rs, q, play.BoardgameId, play.LocationId, play.PlayDate)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

type BGStat struct {
	PlayerId int64
	Score    models.JsonNullInt64
	Winner   bool
}

func getBGStats(DB *sqlx.DB, play BGPlay) ([]BGStat, error) {
	rs := []BGStat{}
	q := `
		select
			pl.player_id,
			cast(s.score as decimal(10, 0)) score,
			s.winner
		from
			tbgstatsplays p,
			tbgstats s,
			tbgstatsplayers pl
		where
			s.play_id = p.id and
			s.player_id = pl.id and
			p.id = ?
	`
	err := DB.Select(&rs, q, play.Id)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

type Play struct {
	Id          int64
	BoardgameId int64
	LocationId  models.JsonNullInt64
	Players     models.JsonArray
	Date        time.Time
	Stats       []Stat
}

func findMappedPlay(DB *sqlx.DB, play BGPlay) ([]Play, error) {
	rs := []Play{}

	q := `
		select
			p.id,
			p.boardgame_id,
			json_extract(p.play_data, '$.players') players,
			p.date
		from
			tboardgameplays p,
			tboardgamestats s
		where
			p.boardgame_id = ? and
			s.play_id = p.id and
			p.date between ? - interval 2 day and ? + interval 2 day
		group by 1
		having count(*) = ?
	`
	err := DB.Select(&rs, q, play.BoardgameId, play.PlayDate, play.PlayDate, len(play.BGStats))
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func getPlay(DB *sqlx.DB, play BGPlay) (*Play, error) {
	var rs Play

	q := `
		select
			p.id,
			p.boardgame_id,
			p.location_id,
			json_extract(p.play_data, '$.players') players,
			p.date
		from
			tboardgameplays p
		where
			p.id = ?
	`
	err := DB.Get(&rs, q, play.PlayId)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

type Stat struct {
	Id       int64
	PlayerId int64
	Score    models.JsonNullInt64
	Winner   bool
}

func getStats(DB *sqlx.DB, play Play) ([]Stat, error) {
	rs := []Stat{}
	q := `
		select
			s.id,
			s.player_id,
			json_extract(s.data, '$.score') score,
			IFNULL(json_extract(s.data, '$.winner'), false) winner
		from
			tboardgamestats s
		where
			s.play_id = ?
	`
	err := DB.Select(&rs, q, play.Id)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func updateLink(DB *sqlx.DB, bgplay BGPlay, play Play) error {
	q := `
		update
			tbgstatsplays
		set
			play_id = :play_id
		where
			id = :id
	`
	_, err := DB.NamedExec(q, map[string]any{
		"id":      bgplay.Id,
		"play_id": play.Id,
	})
	if err != nil {
		return err
	}

	return nil
}

func updateLocation(DB *sqlx.DB, bgplay BGPlay, play Play) error {
	q := `
		update
			tboardgameplays
		set
			location_id = :location_id
		where
			id = :id
	`
	_, err := DB.NamedExec(q, map[string]any{
		"id":          play.Id,
		"location_id": bgplay.LocationId,
	})
	if err != nil {
		return err
	}

	return nil
}

func insertPlay(DB *sqlx.DB, bgplay BGPlay) (int64, error) {
	q := `
		insert into tboardgameplays (
			boardgame_id,
			date,
			location_id
		) values (
			:boardgame_id,
			:date,
			:location_id
		)
	`
	row, err := DB.NamedExec(q, map[string]any{
		"date":         bgplay.PlayDate,
		"boardgame_id": bgplay.BoardgameId,
		"location_id":  bgplay.LocationId,
	})
	if err != nil {
		return -1, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

//update tboardgameplays set play_data = JSON_SET(IFNULL(play_data, JSON_OBJECT()),'$.winner', 96) where id = 1035

func insertStats(DB *sqlx.DB, play_id int64, boardgame_id int64, bgstat BGStat) error {
	q := `
		insert into tboardgamestats (
			play_id,
			boardgame_id,
			player_id
		) values (
			:play_id,
			:boardgame_id,
			:player_id
		)
	`
	row, err := DB.NamedExec(q, map[string]any{
		"boardgame_id": boardgame_id,
		"play_id":      play_id,
		"player_id":    bgstat.PlayerId,
	})
	if err != nil {
		return err
	}

	// {"score":null,"won":true,"new_player":false,"seat_order":0,"start_player":false,"winner":true}

	id, err := row.LastInsertId()
	if err != nil {
		return err
	}

	_, err = DB.NamedExec(`
		update tboardgamestats set data = JSON_SET(IFNULL(data, JSON_OBJECT()),'$.score', :score) where id = :id
	`, map[string]any{
		"id":    id,
		"score": bgstat.Score,
	})
	if err != nil {
		return err
	}

	return nil
}

func updateWinner(DB *sqlx.DB, bgstat BGStat, stat Stat) error {
	q := `
		update
			tboardgamestats
		set
			data = JSON_SET(
				IFNULL(data, JSON_OBJECT()),'$.winner', if(:winner mod 2 = 0, FALSE, TRUE)
			)
		where
			id = :id
	`
	_, err := DB.NamedExec(q, map[string]any{
		"id":     stat.Id,
		"winner": bgstat.Winner,
	})
	if err != nil {
		return err
	}

	return nil
}
