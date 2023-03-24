package model

import (
	"encoding/json"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type CachedPrice struct {
	Id         int64   `gorm:"primaryKey" json:"id"`
	StoreId    int64   `gorm:"foreignkey" json:"store_id"`
	StoreThumb string  `json:"store_thumb"`
	Name       string  `gorm:"index"  json:"name"`
	Price      float64 `json:"price"`
	Stock      int     `json:"stock"`
	Url        string  `json:"url"`
}

func (CachedPrice) TableName() string {
	return "tboardgamepricescached"
}

func (CachedPrice) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []CachedPrice
	rs := db.Scopes(scopes...).Find(&data)
	return data, rs.Error
}

func (CachedPrice) Get(db *gorm.DB, id int64) (any, error) {
	var data CachedPrice
	rs := db.First(&data, id)
	return data, rs.Error
}

func (obj CachedPrice) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := CachedPrice{
		Id: id,
	}

	var payload CachedPrice
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

func (CachedPrice) Create(db *gorm.DB, body []byte) (any, error) {
	var payload CachedPrice
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

func (obj CachedPrice) CreatePrice(c *gin.Context, db *gorm.DB) (any, error) {
	id := c.GetInt64("apiid")

	var data CachedPrice
	db.First(&data, id)

	payload := Price{
		CrDate:     time.Now(),
		StoreId:    data.StoreId,
		StoreThumb: data.StoreThumb,
		Name:       data.Name,
		Price:      data.Price,
		Stock:      data.Stock,
		Url:        data.Url,
		Mapped:     false,
		Ignored:    false,
		Batch:      0,
		BoardgameId: models.JsonNullInt64{
			Int64: 0,
			Valid: false,
		},
	}

	rs := db.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{}),
	}).Create(&payload)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return payload, nil
}

func (obj CachedPrice) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&CachedPrice{}, id)
	if err != nil {
		return nil, rs.Error
	}

	return data, nil
}
