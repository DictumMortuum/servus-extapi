package model

import (
	"encoding/json"
	"github.com/DictumMortuum/servus/pkg/models"
	"gorm.io/gorm"
)

type BGStatsGame struct {
	Id               int64                `json:"id"`
	RefId            int64                `json:"ref_id"`
	Uuid             string               `gorm:"uniqueIndex" json:"uuid"`
	Name             string               `json:"name"`
	ModificationDate string               `json:"modification_date"`
	Cooperative      bool                 `json:"cooperative"`
	HighestWins      bool                 `json:"highest_wins"`
	NoPoints         bool                 `json:"no_points"`
	UsesTeams        bool                 `json:"uses_teams"`
	BoardgameId      models.JsonNullInt64 `json:"boardgame_id"`
	Boardgame        Boardgame            `json:"boardgame"`
}

func (BGStatsGame) TableName() string {
	return "tbgstatsgames"
}

func (BGStatsGame) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (BGStatsGame) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []BGStatsGame
	rs := db.Scopes(scopes...).Preload("Boardgame").Find(&data)
	return data, rs.Error
}

func (BGStatsGame) Get(db *gorm.DB, id int64) (any, error) {
	var data BGStatsGame
	rs := db.Preload("Boardgame").First(&data, id)
	return data, rs.Error
}

func (obj BGStatsGame) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := BGStatsGame{
		Id: id,
	}

	var payload BGStatsGame
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

func (BGStatsGame) Create(db *gorm.DB, body []byte) (any, error) {
	var payload BGStatsGame
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

func (obj BGStatsGame) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&BGStatsGame{}, id)
	if err != nil {
		return nil, rs.Error
	}

	return data, nil
}
