package main

import (
	"github.com/jmoiron/sqlx"
)

func Insert(DB *sqlx.DB, payload map[string]any) (int64, error) {
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
		) on duplicate key update updated = NOW()
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
