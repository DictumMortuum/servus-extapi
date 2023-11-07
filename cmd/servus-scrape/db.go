package main

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var (
	nstmt *sqlx.NamedStmt
)

func Exists(DB *sqlx.DB, payload map[string]any) (bool, error) {
	if nstmt == nil {
		q := `
			select
				1
			from
				tprices
			where
				name = :name and
				store_id = :store_id and
				store_thumb = :store_thumb and
				price = :price and
				stock = :stock
		`

		tx, err := DB.PrepareNamed(q)
		if err != nil {
			return false, err
		}

		nstmt = tx
	}

	var rs bool
	err := nstmt.Get(&rs, payload)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return rs, nil
}

func Insert(DB *sqlx.DB, payload map[string]any) (int64, error) {
	exists, err := Exists(DB, payload)
	if err != nil {
		return -1, err
	}
	if exists {
		return -1, nil
	}

	q := `
		insert into tprices (
			name,
			store_id,
			store_thumb,
			price,
			stock,
			url,
			deleted,
			boardgame_id,
			created,
			updated
		) values (
			:name,
			:store_id,
			:store_thumb,
			:price,
			:stock,
			:url,
			0,
			NULL,
			NOW(),
			NOW()
		) on duplicate key update updated = NOW(), store_thumb = :store_thumb, price = :price, stock = :stock
	`
	row, err := DB.NamedExec(q, payload)
	if err != nil {
		return -1, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}
