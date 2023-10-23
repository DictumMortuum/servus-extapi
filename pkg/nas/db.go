package nas

import (
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

func Insert(DB *sqlx.DB, url string) (int64, error) {
	q := `
		insert into tyoutubedl (
			url
		) values (
			:url
		)
	`
	row, err := DB.NamedExec(q, map[string]any{
		"url": url,
	})
	if err != nil {
		return -1, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}
