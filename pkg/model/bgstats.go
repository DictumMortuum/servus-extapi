package model

import (
	"encoding/json"
	"github.com/DictumMortuum/servus/pkg/models"
	"gorm.io/gorm"
)

type BGStat struct {
	Id          int64                 `json:"id"`
	PlayId      int64                 `json:"play_id" gorm:"foreignkey,uniqueIndex"`
	Score       models.JsonNullString `json:"score"`
	Winner      bool                  `json:"winner"`
	NewPlayer   bool                  `json:"new_player"`
	StartPlayer bool                  `json:"start_player"`
	PlayerId    int64                 `json:"player_id" gorm:"foreignkey,uniqueIndex"`
	Player      BGStatsPlayer         `json:"player"`
	Rank        int                   `json:"rank"`
	SeatOrder   uint                  `json:"seat_order"`
}

func (BGStat) TableName() string {
	return "tbgstats"
}

func (BGStat) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []BGStat
	rs := db.Scopes(scopes...).Preload("Player").Find(&data)
	return data, rs.Error
}

func (BGStat) Get(db *gorm.DB, id int64) (any, error) {
	var data BGStat
	rs := db.Preload("Player").First(&data, id)
	return data, rs.Error
}

func (obj BGStat) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := BGStat{
		Id: id,
	}

	var payload BGStat
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

func (BGStat) Create(db *gorm.DB, body []byte) (any, error) {
	var payload BGStat
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

func (obj BGStat) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&Player{}, id)
	if err != nil {
		return nil, rs.Error
	}

	return data, nil
}
