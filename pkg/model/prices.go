package model

import (
	"encoding/json"
	"github.com/DictumMortuum/servus/pkg/models"
	"gorm.io/gorm"
	"time"
)

type Price struct {
	Id          int64                `gorm:"primaryKey" json:"id"`
	BoardgameId models.JsonNullInt64 `gorm:"foreignkey" json:"boardgame_id"`
	CrDate      time.Time            `json:"date"`
	StoreId     int64                `gorm:"foreignkey" json:"store_id"`
	StoreThumb  string               `json:"store_thumb"`
	Name        string               `gorm:"index"  json:"name"`
	Price       float64              `json:"price"`
	Stock       int                  `json:"stock"`
	Url         string               `json:"url"`
	Mapped      bool                 `json:"mapped"`
	Ignored     bool                 `json:"ignored"`
	Batch       int                  `json:"batch"`
}

func (Price) TableName() string {
	return "tboardgameprices"
}

func (Price) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []Price
	rs := db.Scopes(scopes...).Find(&data)
	return data, rs.Error
}

func (Price) Get(db *gorm.DB, id int64) (any, error) {
	var data Price
	rs := db.First(&data, id)
	return data, rs.Error
}

func (obj Price) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := Price{
		Id: id,
	}

	var payload Price
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

func (Price) Create(db *gorm.DB, body []byte) (any, error) {
	var payload Price
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

func (obj Price) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&Price{}, id)
	if err != nil {
		return nil, rs.Error
	}

	return data, nil
}
