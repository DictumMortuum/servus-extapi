package model

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Location struct {
	Id     int64  `gorm:"primaryKey" json:"id"`
	Name   string `json:"name"`
	Hidden bool   `json:"hidden"`
}

func (Location) TableName() string {
	return "tlocation"
}

func (Location) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (Location) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []Location
	rs := db.Scopes(scopes...).Find(&data)
	return data, rs.Error
}

func (Location) Get(db *gorm.DB, id int64) (any, error) {
	var data Location
	rs := db.First(&data, id)
	return data, rs.Error
}

func (obj Location) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := Location{
		Id: id,
	}

	var payload Location
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

func (Location) Create(db *gorm.DB, body []byte) (any, error) {
	var payload Location
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

func (obj Location) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&Location{}, id)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return data, nil
}
