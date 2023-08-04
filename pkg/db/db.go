package db

import (
	"database/sql"
	"github.com/DictumMortuum/servus-extapi/pkg/config"
)

func Database(dsn string) (*sql.DB, error) {
	sqlDB, err := sql.Open("mysql", config.Cfg.Databases["mysql"])
	if err != nil {
		return nil, err
	}

	return sqlDB, nil
}
