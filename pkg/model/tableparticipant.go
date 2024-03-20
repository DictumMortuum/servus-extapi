package model

import (
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TableParticipant struct {
	Id      int64  `gorm:"primaryKey" json:"id,omitempty"`
	TableId int64  `json:"table_id,omitempty"`
	UserId  string `json:"user_id,omitempty"`
	Name    string `json:"name,omitempty"`
}

func (TableParticipant) TableName() string {
	return "ttableparticipants"
}

func (TableParticipant) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (TableParticipant) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []TableParticipant
	rs := db.Scopes(scopes...).Find(&data)
	return data, rs.Error
}

func (TableParticipant) Get(db *gorm.DB, id int64) (any, error) {
	var data TableParticipant
	rs := db.First(&data, id)
	return data, rs.Error
}

func (obj TableParticipant) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := TableParticipant{
		Id: id,
	}

	var payload TableParticipant
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

func (TableParticipant) Create(db *gorm.DB, body []byte) (any, error) {
	var payload TableParticipant
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

func (obj TableParticipant) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&TableParticipant{}, id)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return data, nil
}
