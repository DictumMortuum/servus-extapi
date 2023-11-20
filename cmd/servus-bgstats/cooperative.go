package main

import (
	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"
)

func updateCooperativeBoardgame(DB *sqlx.DB, id int64) error {
	q := `
	update
		tboardgames
	set
		cooperative = true
	where
		id = ?
	`

	_, err := DB.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}

func updateSolitaireBoardgame(DB *sqlx.DB, id int64) error {
	q := `
	update
		tboardgames
	set
		solitaire = true
	where
		id = ?
	`

	_, err := DB.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}

func cooperative(c *cli.Context) error {
	DB, err := db.DatabaseX()
	if err != nil {
		return err
	}
	defer DB.Close()

	type boardgame struct {
		Id        int64            `json:"id,omitempty"`
		Mechanics models.JsonArray `json:"mechanics,omitempty"`
	}

	q := `
	select
		g.id,
		json_extract(g.bgg_data, '$.links.boardgamemechanic') mechanics
	from
		tboardgames g
	`

	rs := []boardgame{}
	err = DB.Select(&rs, q)
	if err != nil {
		return err
	}

	for _, item := range rs {
		for _, mechanic := range item.Mechanics {
			switch v := mechanic.(type) {
			case map[string]any:
				if v["name"].(string) == "Cooperative Game" {
					err := updateCooperativeBoardgame(DB, item.Id)
					if err != nil {
						return err
					}
				} else if v["name"].(string) == "Solo / Solitaire Game" {
					err := updateSolitaireBoardgame(DB, item.Id)
					if err != nil {
						return err
					}
				}
			default:
				continue
			}
		}

	}

	return nil
}
