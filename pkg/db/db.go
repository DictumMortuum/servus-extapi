package db

import (
	"database/sql"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/util"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

func Gorm() (*gorm.DB, *sql.DB, error) {
	sqlDB, err := Database()
	if err != nil {
		return nil, nil, err
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:              sqlDB,
		DefaultStringSize: 512,
	}), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, nil, err
	}

	return db, sqlDB, nil
}
