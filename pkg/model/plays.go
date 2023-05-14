package model

import (
	"encoding/json"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Play struct {
	Id          int64          `gorm:"primaryKey" json:"id"`
	BoardgameId int64          `json:"boardgame_id"`
	Boardgame   Boardgame      `json:"boardgame"`
	Date        time.Time      `json:"date"`
	PlayData    datatypes.JSON `gorm:"serializer:json" json:"play_data"`
	Stats       []Stat         `json:"stats"`
	LocationId  int64          `json:"location_id"`
	Location    Location       `json:"location"`
}

func (Play) TableName() string {
	return "tboardgameplays"
}

func (Play) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (Play) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []Play
	rs := db.Scopes(scopes...).Preload("Boardgame").Preload("Stats").Preload("Location").Preload("Stats.Player").Find(&data)
	return data, rs.Error
}

func (Play) Get(db *gorm.DB, id int64) (any, error) {
	var data Play
	rs := db.Preload("Boardgame").Preload("Stats.Player").Preload("Location").First(&data, id)
	return data, rs.Error
}

func (obj Play) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	var payload Play
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	rs := db.Debug().Omit("Boardgame").Omit("Location").Clauses(clause.OnConflict{UpdateAll: true}).Session(&gorm.Session{FullSaveAssociations: true}).Create(&payload)
	if rs.Error != nil {
		return nil, err
	}

	return obj.Get(db, id)
}

func (Play) Create(db *gorm.DB, body []byte) (any, error) {
	var payload Play
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	rs := db.Debug().Omit("Boardgame").Omit("Location").Clauses(clause.OnConflict{UpdateAll: true}).Session(&gorm.Session{FullSaveAssociations: true}).Create(&payload)
	if rs.Error != nil {
		return nil, err
	}

	return payload, nil
}

func (obj Play) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&Play{}, id)
	if err != nil {
		return nil, rs.Error
	}

	return data, nil
}
