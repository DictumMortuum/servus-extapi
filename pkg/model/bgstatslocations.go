package model

import (
	"encoding/json"
	"github.com/DictumMortuum/servus/pkg/models"
	"gorm.io/gorm"
)

type BGStatsLocation struct {
	Id               int64                `json:"id"`
	RefId            int64                `json:"ref_id"`
	Uuid             string               `gorm:"uniqueIndex" json:"uuid"`
	Name             string               `json:"name"`
	ModificationDate string               `json:"modification_date"`
	LocationId       models.JsonNullInt64 `json:"location_id"`
	Location         Location             `json:"location"`
}

func (BGStatsLocation) TableName() string {
	return "tbgstatslocations"
}

func (BGStatsLocation) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []BGStatsLocation
	rs := db.Scopes(scopes...).Preload("Location").Find(&data)
	return data, rs.Error
}

func (BGStatsLocation) Get(db *gorm.DB, id int64) (any, error) {
	var data BGStatsLocation
	rs := db.Preload("Location").First(&data, id)
	return data, rs.Error
}

func (obj BGStatsLocation) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := BGStatsLocation{
		Id: id,
	}

	var payload BGStatsLocation
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

func (BGStatsLocation) Create(db *gorm.DB, body []byte) (any, error) {
	var payload BGStatsLocation
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

func (obj BGStatsLocation) Delete(db *gorm.DB, id int64) (any, error) {
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
