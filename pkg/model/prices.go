package model

import (
	"encoding/json"
	"time"

	"github.com/DictumMortuum/servus/pkg/models"
	"gorm.io/gorm"
)

type Price struct {
	Id          int64                `gorm:"primaryKey" json:"id"`
	Created     time.Time            `json:"created"`
	Updated     time.Time            `json:"updated"`
	Name        string               `gorm:"index" json:"name"`
	StoreId     int64                `gorm:"foreignkey" json:"store_id"`
	StoreThumb  string               `json:"store_thumb"`
	Price       float64              `json:"price"`
	Stock       int                  `json:"stock"`
	Url         string               `json:"url"`
	Deleted     bool                 `json:"deleted"`
	BoardgameId models.JsonNullInt64 `gorm:"foreignkey" json:"boardgame_id"`
}

func (Price) TableName() string {
	return "tprices"
}

func (Price) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db.Where("deleted = 0")
}

func (c Price) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []Price
	rs := db.Scopes(scopes...).Scopes(c.DefaultFilter).Find(&data)
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

	var payload map[string]any
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	rs := db.Model(&model).Save(payload)
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
	var data Price
	rs := db.First(&data, id)
	if rs.Error != nil {
		return nil, rs.Error
	}
	data.Deleted = true

	rs = db.Model(&data).Updates(data)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return data, nil
}
