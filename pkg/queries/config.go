package queries

import (
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/jmoiron/sqlx"
)

func GetConfig(DB *sqlx.DB, c string) (*model.Configuration, error) {
	var cfg model.Configuration
	err := DB.Get(&cfg, `select * from tconfig where config = ?`, c)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
