package db

import (
	"database/sql"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/util"
	"github.com/jmoiron/sqlx"
)

func Database() (*sql.DB, error) {
	sqlDB, err := sql.Open("mysql", config.Cfg.Databases["mariadb"])
	if err != nil {
		return nil, err
	}

	return sqlDB, nil
}

func DatabaseX() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", config.Cfg.Databases["mariadb"])
	if err != nil {
		return nil, err
	}

	db.MapperFunc(util.ToSnake)
	return db, nil
}
