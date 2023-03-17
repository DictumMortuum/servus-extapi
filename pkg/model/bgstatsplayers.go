package model

import (
	"encoding/json"
	"github.com/DictumMortuum/servus/pkg/models"
	"gorm.io/gorm"
)

type BGStatsPlayer struct {
	Id               int64                `json:"id"`
	RefId            int64                `json:"ref_id"`
	Uuid             string               `gorm:"uniqueIndex" json:"uuid"`
	Name             string               `json:"name"`
	IsAnonymous      bool                 `json:"is_anonymous"`
	ModificationDate string               `json:"modification_date"`
	BggUsername      string               `json:"bgg_username"`
	PlayerId         models.JsonNullInt64 `json:"player_id"`
	Player           Player               `json:"player"`
}

func (BGStatsPlayer) TableName() string {
	return "tbgstatsplayers"
}

func (BGStatsPlayer) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []BGStatsPlayer
	rs := db.Scopes(scopes...).Preload("Player").Find(&data)
	return data, rs.Error
}

func (BGStatsPlayer) Get(db *gorm.DB, id int64) (any, error) {
	var data BGStatsPlayer
	rs := db.Preload("Player").First(&data, id)
	return data, rs.Error
}

func (obj BGStatsPlayer) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := BGStatsPlayer{
		Id: id,
	}

	var payload BGStatsPlayer
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

func (BGStatsPlayer) Create(db *gorm.DB, body []byte) (any, error) {
	var payload BGStatsPlayer
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

func (obj BGStatsPlayer) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&BGStatsPlayer{}, id)
	if err != nil {
		return nil, rs.Error
	}

	return data, nil
}

// rs := db.Debug().Clauses(clause.OnConflict{
// 	Columns:   []clause.Column{{Name: "uuid"}},
// 	DoUpdates: clause.Assignments(map[string]interface{}{}),
// }).Create(&payload)
// if rs.Error != nil {
// 	return nil, err
// }
