package model

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Device struct {
	Id    int64  `json:"id"`
	Mac   string `json:"mac"`
	Alias string `json:"alias"`
}

func (Device) TableName() string {
	return "tdevices"
}

func (Device) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (Device) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []Device
	rs := db.Scopes(scopes...).Find(&data)
	return data, rs.Error
}

func (Device) Get(db *gorm.DB, id int64) (any, error) {
	var data Device
	rs := db.First(&data, id)
	return data, rs.Error
}

func (obj Device) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := Device{
		Id: id,
	}

	var payload Device
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

func (Device) Create(db *gorm.DB, body []byte) (any, error) {
	var payload Device
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

func (obj Device) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&Device{}, id)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return data, nil
}
