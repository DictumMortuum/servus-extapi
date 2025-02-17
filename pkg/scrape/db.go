package scrape

import (
	"github.com/jmoiron/sqlx"
)

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

func UpdatePages(DB *sqlx.DB, id int64, count int) error {
	q := `update tscrapeurl set last_pages = ? where id = ?`
	_, err := DB.Exec(q, count, id)
	if err != nil {
		return err
	}

	return nil
}
