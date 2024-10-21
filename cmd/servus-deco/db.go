package main

import (
	"github.com/DictumMortuum/servus-extapi/pkg/db"
)

type Mapping struct {
	Id    int64  `json:"id"`
	Mac   string `json:"mac"`
	Alias string `json:"alias"`
}

func getMappings() ([]Mapping, error) {
	DB, err := db.DatabaseX()
	if err != nil {
		return nil, err
	}
	defer DB.Close()

	rs := []Mapping{}
	q := `select id, mac, alias from tdevices`
	err = DB.Select(&rs, q)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
