package nas

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type Download struct {
	Url       string `json:"url,omitempty"`
	Processed bool   `json:"processed,omitempty"`
}

func Exists(DB *sqlx.DB, url string) (*Download, error) {
	var rs Download
	err := DB.QueryRowx(`
		select
			url,
			processed
		from
			tyoutubedl
		where
			url = ?
	`, url).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func Insert(DB *sqlx.DB, series, url string) (int64, error) {
	q := `
		insert into tyoutubedl (
			series,
			url,
			processed
		) values (
			:series,
			:url,
			0
		) on duplicate key update id = id
	`
	row, err := DB.NamedExec(q, map[string]any{
		"series": series,
		"url":    url,
	})
	if err != nil {
		return -1, err
	}

	log.Println(series, url)

	id, err := row.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}
