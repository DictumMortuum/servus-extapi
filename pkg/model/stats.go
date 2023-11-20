package model

import (
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Stat struct {
	Id          int64          `gorm:"primaryKey" json:"id"`
	BoardgameId int64          `json:"boardgame_id"`
	PlayId      int64          `json:"play_id" gorm:"foreignkey"`
	PlayerId    int64          `json:"player_id"`
	Player      Player         `json:"player"`
	Data        datatypes.JSON `gorm:"serializer:json" json:"data"`
}

func (Stat) TableName() string {
	return "tboardgamestats"
}

func (Stat) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (Stat) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []Stat
	rs := db.Scopes(scopes...).Preload("Player").Find(&data)
	return data, rs.Error
}

func (Stat) Get(db *gorm.DB, id int64) (any, error) {
	var data Stat
	rs := db.Preload("Player").First(&data, id)
	return data, rs.Error
}

func (obj Stat) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := Stat{
		Id: id,
	}

	var payload Stat
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	rs := db.Model(&model).Updates(payload)
	if rs.Error != nil {
		return nil, err
	}

	return obj.Get(db, id)
}

func (Stat) Create(db *gorm.DB, body []byte) (any, error) {
	var payload Stat
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	rs := db.Create(&payload)
	if rs.Error != nil {
		return nil, err
	}

	return payload, nil
}

func (obj Stat) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&Stat{}, id)
	if err != nil {
		return nil, rs.Error
	}

	return data, nil
}
