package model

import (
	"encoding/json"
	"time"

	"github.com/DictumMortuum/servus/pkg/models"
	"gorm.io/gorm"
)

type BoardgamePrice struct {
	Id          int64                `gorm:"primaryKey" json:"id"`
	BoardgameId models.JsonNullInt64 `gorm:"foreignkey" json:"boardgame_id"`
	CrDate      time.Time            `json:"date"`
	StoreId     int64                `gorm:"foreignkey" json:"store_id"`
	StoreThumb  string               `json:"store_thumb"`
	Name        string               `gorm:"index"  json:"name"`
	Price       float64              `json:"BoardgamePrice"`
	Stock       int                  `json:"stock"`
	Url         string               `json:"url"`
	Mapped      bool                 `json:"mapped"`
	Ignored     bool                 `json:"ignored"`
	Batch       int                  `json:"batch"`
}

func (BoardgamePrice) TableName() string {
	return "tboardgameprices"
}

func (BoardgamePrice) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (BoardgamePrice) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []BoardgamePrice
	rs := db.Scopes(scopes...).Find(&data)
	return data, rs.Error
}

func (BoardgamePrice) Get(db *gorm.DB, id int64) (any, error) {
	var data BoardgamePrice
	rs := db.First(&data, id)
	return data, rs.Error
}

func (obj BoardgamePrice) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := BoardgamePrice{
		Id: id,
	}

	var payload BoardgamePrice
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

func (BoardgamePrice) Create(db *gorm.DB, body []byte) (any, error) {
	var payload BoardgamePrice
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

func (obj BoardgamePrice) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&BoardgamePrice{}, id)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return data, nil
}
