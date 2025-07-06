package model

import (
	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func OmitMultiple(resources []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, resource := range resources {
			db = db.Omit(resource)
		}

		return db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Session(&gorm.Session{
			FullSaveAssociations: true,
		})
	}
}

func GetSqlx(db *gorm.DB) (*sqlx.DB, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlxDB := sqlx.NewDb(sqlDB, "mysql")
	return sqlxDB, nil
}
