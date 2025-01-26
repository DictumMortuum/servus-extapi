package model

import (
	"encoding/json"

	"github.com/DictumMortuum/servus/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Scrape struct {
	Id                int64                 `gorm:"primaryKey" json:"id"`
	StoreId           int64                 `gorm:"foreignkey" json:"store_id"`
	StoreName         string                `gorm:"-" json:"store_name"`
	SelItem           string                `json:"sel_item"`
	SelName           string                `json:"sel_name"`
	SelItemThumb      string                `json:"sel_item_thumb"`
	SelItemInstock    models.JsonNullString `json:"sel_item_instock"`
	SelItemPreorder   models.JsonNullString `json:"sel_item_preorder"`
	SelItemOutofstock models.JsonNullString `json:"sel_item_outofstock"`
	SelPrice          string                `json:"sel_price"`
	SelAltPrice       models.JsonNullString `json:"sel_alt_price"`
	SelOriginalPrice  string                `json:"sel_original_price"`
	SelUrl            string                `json:"sel_url"`
	SelNext           string                `json:"sel_next"`
	Tag               string                `json:"tag"`
	AllowedDomain     string                `json:"allowed_domain"`
	AbsoluteNextUrl   bool                  `json:"absolute_next_url"`
	URLs              []ScrapeUrl           `json:"urls"`
}

func (Scrape) TableName() string {
	return "tscrape"
}

func (Scrape) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (c Scrape) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []Scrape
	rs := db.Scopes(scopes...).Preload("URLs").Scopes(c.DefaultFilter).Find(&data)
	return data, rs.Error
}

func (Scrape) Get(db *gorm.DB, id int64) (any, error) {
	var data Scrape
	rs := db.Preload("URLs").First(&data, id)
	return data, rs.Error
}

func (obj Scrape) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	var payload Scrape
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	// delete(payload, "store_name")

	// rs := db.Model(&model).Save(payload)
	// if rs.Error != nil {
	// 	return nil, err
	// }

	err = db.Debug().Model(&payload).Association("URLs").Unscoped().Replace(payload.URLs)
	if err != nil {
		return err, nil
	}

	rs := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Session(&gorm.Session{
		FullSaveAssociations: true,
	}).Create(&payload)
	if rs.Error != nil {
		return nil, err
	}

	return obj.Get(db, id)
}

func (Scrape) Create(db *gorm.DB, body []byte) (any, error) {
	var payload Scrape
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	rs := db.Debug().Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Session(&gorm.Session{
		FullSaveAssociations: true,
	}).Create(&payload)
	if rs.Error != nil {
		return nil, err
	}

	return payload, nil
}

func (obj Scrape) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&Scrape{}, id)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return data, nil
}
