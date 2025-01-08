package main

import (
	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
)

func getMappings() ([]model.Device, error) {
	DB, err := db.DatabaseX()
	if err != nil {
		return nil, err
	}
	defer DB.Close()

	rs := []model.Device{}
	q := `select id, mac, alias from tdevices`
	err = DB.Select(&rs, q)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
