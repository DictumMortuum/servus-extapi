package main

import (
	"database/sql"

	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/jmoiron/sqlx"
)

func Exists(DB *sqlx.DB, mac string) (int64, error) {
	q := `
		select
			id
		from
			tdevices
		where
			mac = ?
	`

	var rs int64
	err := DB.Get(&rs, q, mac)
	if err == sql.ErrNoRows {
		return -1, nil
	}
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	}

	return rs, nil
}

func Insert(DB *sqlx.DB, payload Val) (int64, error) {
	exists, err := Exists(DB, payload.Metric.MAC)
	if err != nil {
		return -1, err
	}

	if exists != -1 {
		return exists, nil
	}

	q := `
		insert into tdevices (
			mac,
			alias,
			online
		) values (
			:mac,
			:nickname,
			1
		)
	`
	row, err := DB.NamedExec(q, payload.Metric)
	if err != nil {
		return -1, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func Status(DB *sqlx.DB, mac string, online int) error {
	_, err := DB.Exec(`update tdevices set online = ? where mac = ?`, online, mac)
	if err != nil {
		return err
	}

	return nil
}

func Reset(DB *sqlx.DB) error {
	_, err := DB.Exec(`update tdevices set online = 0`)
	if err != nil {
		return err
	}

	return nil
}

func Select(DB *sqlx.DB) ([]model.Device, error) {
	var rs []model.Device
	err := DB.Select(&rs, `select * from tdevices`)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
