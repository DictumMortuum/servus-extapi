package model

import (
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EurovisionParticipation struct {
	Id          int64     `gorm:"primaryKey" json:"id,omitempty"`
	UserId      string    `json:"user_id,omitempty"`
	Email       string    `json:"email"`
	BoardgameId int64     `json:"boardgame_id"`
	Boardgame   Boardgame `json:"boardgame"`
}

func (EurovisionParticipation) TableName() string {
	return "teurovisionparticipation"
}

func (EurovisionParticipation) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (EurovisionParticipation) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []EurovisionParticipation
	rs := db.Scopes(scopes...).Preload("Boardgame").Find(&data)
	return data, rs.Error
}

func (EurovisionParticipation) Get(db *gorm.DB, id int64) (any, error) {
	var data EurovisionParticipation
	rs := db.Preload("Boardgame").First(&data, id)
	return data, rs.Error
}

func (obj EurovisionParticipation) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := EurovisionParticipation{
		Id: id,
	}

	var payload EurovisionParticipation
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

func (EurovisionParticipation) Create(db *gorm.DB, body []byte) (any, error) {
	var payload EurovisionParticipation
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	rs := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&payload)
	if rs.Error != nil {
		return nil, err
	}

	return payload, nil
}

func (obj EurovisionParticipation) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&EurovisionParticipation{}, id)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return data, nil
}
