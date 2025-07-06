package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Table struct {
	Id           int64              `gorm:"primaryKey" json:"id,omitempty"`
	BoardgameId  int64              `json:"boardgame_id,omitempty"`
	Boardgame    Boardgame          `json:"boardgame,omitempty"`
	CreatorId    string             `json:"creator_id,omitempty"`
	Creator      string             `json:"creator,omitempty"`
	Location     string             `json:"location,omitempty"`
	Seats        int                `json:"seats"`
	Teaching     bool               `json:"teaching"`
	Date         time.Time          `json:"date,omitempty"`
	Participants []TableParticipant `json:"participants"`
}

func (Table) TableName() string {
	return "ttables"
}

func (Table) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (Table) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []Table
	rs := db.Scopes(scopes...).Preload("Boardgame").Preload("Participants").Find(&data)
	return data, rs.Error
}

func (Table) Get(db *gorm.DB, id int64) (any, error) {
	var data Table
	rs := db.Preload("Boardgame").Preload("Participants").First(&data, id)
	return data, rs.Error
}

func (obj Table) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := Table{
		Id: id,
	}

	var payload Table
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

func (Table) Create(db *gorm.DB, body []byte) (any, error) {
	var payload Table
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	rs := db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&payload)
	if rs.Error != nil {
		return nil, err
	}

	return payload, nil
}

func (obj Table) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Where("table_id = ?", id).Delete(&TableParticipant{})
	if rs.Error != nil {
		return nil, rs.Error
	}

	rs = db.Delete(&Table{}, id)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return data, nil
}
