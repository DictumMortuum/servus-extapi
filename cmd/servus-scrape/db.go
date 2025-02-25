package main

import (
	"database/sql"

	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/jmoiron/sqlx"
)

var (
	nstmt *sqlx.NamedStmt
)

func Exists(DB *sqlx.DB, payload map[string]any) (int64, error) {
	if nstmt == nil {
		q := `
			select
				id
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
			return -1, err
		}

		nstmt = tx
	}

	var rs int64
	err := nstmt.Get(&rs, payload)
	if err == sql.ErrNoRows {
		return -1, nil
	}
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	}

	return rs, nil
}

type Price struct {
	Id      int64  `json:"id,omitempty"`
	StoreId int64  `json:"store_id,omitempty"`
	Name    string `json:"name,omitempty"`
}

func Get(DB *sqlx.DB, payload map[string]any) (*Price, error) {
	id, _ := payload["store_id"]
	name, _ := payload["name"]

	var rs Price
	err := DB.QueryRowx(`
		select
			id,
			store_id,
			name
		from
			tprices
		where
			name = ? and
			store_id = ?
	`, name, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func UpdateCounts(DB *sqlx.DB, store_id int64) error {
	q := `
		update
			tboardgamestores
		set
			count = (select count(*) from tprices where store_id = ?),
			latest_count = (select count(*) from tprices where store_id = ? and updated > NOW() - interval 7 day)
		where
			id = ?
	`

	_, err := DB.Exec(q, store_id, store_id, store_id)
	if err != nil {
		return err
	}

	return nil
}

func Insert(DB *sqlx.DB, payload map[string]any) (int64, error) {
	if val, ok := payload["price"]; ok {
		if val.(float64) == 0 {
			return -1, nil
		}
	}

	if val, ok := payload["store_thumb"]; ok {
		if val.(string) == "" {
			payload["store_thumb"] = "https://placehold.co/200x200"
		}
	}

	if _, ok := payload["tag"]; !ok {
		payload["tag"] = ""
	}

	exists, err := Exists(DB, payload)
	if err != nil {
		return -1, err
	}

	if exists != -1 {
		err := Fresh(DB, exists)
		if err != nil {
			return -1, err
		}

		return -1, nil
	}

	// p, _ := Get(DB, payload)
	// log.Println(p)

	q := `
		insert into tprices (
			name,
			store_id,
			store_thumb,
			price,
			original_price,
			stock,
			url,
			deleted,
			boardgame_id,
			created,
			updated,
			tag
		) values (
			:name,
			:store_id,
			:store_thumb,
			:price,
			:original_price,
			:stock,
			:url,
			0,
			NULL,
			NOW(),
			NOW(),
			:tag
		) on duplicate key update
		 updated = NOW(),
		 store_thumb = :store_thumb,
		 price = :price,
		 original_price = :original_price,
		 stock = :stock,
		 deleted = 0
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

func Stale(DB *sqlx.DB, id int64) error {
	q := `update tprices set deleted = 1 where store_id = ?`
	_, err := DB.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}

func Fresh(DB *sqlx.DB, id int64) error {
	q := `update tprices set deleted = 0 where id = ?`
	_, err := DB.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}

func Cleanup(DB *sqlx.DB, id int64) error {
	q := `delete from tprices where deleted = 1 and store_id = ?`
	_, err := DB.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}

func Delete(DB *sqlx.DB, id int64) error {
	q := `delete from tprices where store_id = ?`
	_, err := DB.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}

func GetScrapes(DB *sqlx.DB) ([]model.Scrape, error) {
	rs := []model.Scrape{}

	sql := `
		select
			sc.*,
			st.name as store_name
		from
			tscrape sc,
			tboardgamestores st
		where
			sc.store_id = st.id
	`

	err := DB.Select(&rs, sql)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func GetURLs(DB *sqlx.DB, id int64) ([]model.ScrapeUrl, error) {
	rs := []model.ScrapeUrl{}

	sql := `
		select
			*
		from
			tscrapeurl
		where
			scrape_id = ?
	`

	err := DB.Select(&rs, sql, id)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
