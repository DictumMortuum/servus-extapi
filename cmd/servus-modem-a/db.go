package main

import (
	"encoding/json"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/jmoiron/sqlx"
)

func saveStats(s *model.Modem, id string) error {
	payload, err := json.Marshal(s)
	if err != nil {
		return err
	}

	db, err := sqlx.Connect("mysql", config.Cfg.Databases["mariadb"])
	if err != nil {
		return err
	}
	defer db.Close()

	q := `update tkeyval set json = :json, date = NOW() where id = :id`
	_, err = db.NamedExec(q, map[string]any{
		"id":   id,
		"json": string(payload),
	})
	if err != nil {
		return err
	}

	return nil
}
