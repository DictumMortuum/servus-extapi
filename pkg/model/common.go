package model

import (
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
