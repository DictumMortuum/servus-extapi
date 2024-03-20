package model

import (
	"encoding/json"

	"gorm.io/gorm"
)

type YoutubeDL struct {
	Id        int64  `gorm:"primaryKey" json:"id"`
	Series    string `json:"series"`
	Url       string `json:"url"`
	Processed bool   `json:"processed"`
}

func (YoutubeDL) TableName() string {
	return "tyoutubedl"
}

func (YoutubeDL) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (YoutubeDL) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []YoutubeDL
	rs := db.Scopes(scopes...).Find(&data)
	return data, rs.Error
}

func (YoutubeDL) Get(db *gorm.DB, id int64) (any, error) {
	var data YoutubeDL
	rs := db.First(&data, id)
	return data, rs.Error
}

func (obj YoutubeDL) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := YoutubeDL{
		Id: id,
	}

	var payload YoutubeDL
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

func (YoutubeDL) Create(db *gorm.DB, body []byte) (any, error) {
	var payload YoutubeDL
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

func (obj YoutubeDL) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&YoutubeDL{}, id)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return data, nil
}
