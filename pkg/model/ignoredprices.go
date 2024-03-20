package model

import (
	"encoding/json"

	"gorm.io/gorm"
)

type IgnoredPrice struct {
	Id      int64  `gorm:"primaryKey" json:"id"`
	StoreId int64  `gorm:"foreignkey" json:"store_id"`
	Name    string `gorm:"index"  json:"name"`
}

func (IgnoredPrice) TableName() string {
	return "tboardgamepricesignored"
}

func (IgnoredPrice) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (IgnoredPrice) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []IgnoredPrice
	rs := db.Scopes(scopes...).Find(&data)
	return data, rs.Error
}

func (IgnoredPrice) Get(db *gorm.DB, id int64) (any, error) {
	var data IgnoredPrice
	rs := db.First(&data, id)
	return data, rs.Error
}

func (obj IgnoredPrice) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := IgnoredPrice{
		Id: id,
	}

	var payload IgnoredPrice
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

func (IgnoredPrice) Create(db *gorm.DB, body []byte) (any, error) {
	var payload IgnoredPrice
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

func (obj IgnoredPrice) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&IgnoredPrice{}, id)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return data, nil
}
