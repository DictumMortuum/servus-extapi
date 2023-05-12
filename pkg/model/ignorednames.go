package model

import (
	"encoding/json"
	"gorm.io/gorm"
)

type IgnoredName struct {
	Id   int64  `gorm:"primaryKey" json:"id"`
	Name string `gorm:"index"  json:"name"`
}

func (IgnoredName) TableName() string {
	return "tboardgamenamesignored"
}

func (IgnoredName) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []IgnoredName
	rs := db.Scopes(scopes...).Find(&data)
	return data, rs.Error
}

func (IgnoredName) Get(db *gorm.DB, id int64) (any, error) {
	var data IgnoredName
	rs := db.First(&data, id)
	return data, rs.Error
}

func (obj IgnoredName) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := IgnoredName{
		Id: id,
	}

	var payload IgnoredName
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

func (IgnoredName) Create(db *gorm.DB, body []byte) (any, error) {
	var payload IgnoredName
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

func (obj IgnoredName) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&IgnoredName{}, id)
	if err != nil {
		return nil, rs.Error
	}

	return data, nil
}
