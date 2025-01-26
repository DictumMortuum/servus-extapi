package model

import (
	"database/sql"
	"encoding/json"

	"gorm.io/gorm"
)

type ScrapeUrl struct {
	Id             int64         `gorm:"primaryKey" json:"id"`
	ScrapeId       int64         `gorm:"foreignkey" json:"scrape_id"`
	Url            string        `json:"url"`
	LastScraped    sql.NullInt32 `json:"last_scraped"`
	LastInstock    sql.NullInt32 `json:"last_instock"`
	LastPreorder   sql.NullInt32 `json:"last_preorder"`
	LastOutofstock sql.NullInt32 `json:"last_outofstock"`
	LastPages      sql.NullInt32 `json:"last_pages"`
}

func (ScrapeUrl) TableName() string {
	return "tscrapeurl"
}

func (ScrapeUrl) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (c ScrapeUrl) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []ScrapeUrl
	rs := db.Scopes(scopes...).Scopes(c.DefaultFilter).Find(&data)
	return data, rs.Error
}

func (ScrapeUrl) Get(db *gorm.DB, id int64) (any, error) {
	var data ScrapeUrl
	rs := db.First(&data, id)
	return data, rs.Error
}

func (obj ScrapeUrl) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := ScrapeUrl{
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

func (ScrapeUrl) Create(db *gorm.DB, body []byte) (any, error) {
	var payload ScrapeUrl
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

func (obj ScrapeUrl) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&ScrapeUrl{}, id)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return data, nil
}
