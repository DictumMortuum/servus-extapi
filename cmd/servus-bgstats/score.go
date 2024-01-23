package main

import (
	"log"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"
)

func updateScore(DB *sqlx.DB, id int64, score float64) error {
	q := `
		update
			tboardgamestats
		set
			data = JSON_SET(
				IFNULL(data, JSON_OBJECT()),'$.score', :score
			)
		where
			id = :id
	`
	_, err := DB.NamedExec(q, map[string]any{
		"id":    id,
		"score": score,
	})
	if err != nil {
		return err
	}

	return nil
}

func score(c *cli.Context) error {
	DB, err := db.DatabaseX()
	if err != nil {
		return err
	}
	defer DB.Close()

	type stats struct {
		Id   int64       `json:"id,omitempty"`
		Data models.Json `json:"data,omitempty"`
	}

	// q := `
	// select
	// 	s.id,
	// 	s.data
	// from
	// 	tboardgamestats s
	// where
	// 	json_extract(s.data, '$.score') is null
	// `

	q := `
		select
			s.id,
			s.data
		from
			tboardgamestats s
		where
			json_extract(s.data, '$.score') = 1 and json_extract(s.data, '$.winner') = "false"
	`

	rs := []stats{}
	err = DB.Select(&rs, q)
	if err != nil {
		return err
	}

	for _, item := range rs {
		score := 0.0
		flag := true
		// for _, val := range item.Data {
		// 	switch val := val.(type) {
		// 	case int:
		// 	case float64:
		// 		{
		// 			flag = true
		// 			score += val
		// 		}
		// 	case bool:
		// 		{
		// 			flag = true
		// 			if val {
		// 				score = 1
		// 			} else {
		// 				score = 0
		// 			}
		// 		}
		// 	default:
		// 		{
		// 			// log.Println(val, "not int")
		// 		}
		// 	}
		// }

		if flag {
			log.Println(item.Data, score)
			err := updateScore(DB, item.Id, 0)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
