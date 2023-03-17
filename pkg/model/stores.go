package model

import (
	"encoding/json"
	"gorm.io/gorm"
)

type Store struct {
	Id   int64  `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

func (Store) TableName() string {
	return "tboardgamestores"
}

func (Store) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []Store
	rs := db.Scopes(scopes...).Find(&data)
	return data, rs.Error
}

func (Store) Get(db *gorm.DB, id int64) (any, error) {
	var data Store
	rs := db.First(&data, id)
	return data, rs.Error
}

func (obj Store) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := Store{
		Id: id,
	}

	var payload Store
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

func (Store) Create(db *gorm.DB, body []byte) (any, error) {
	var payload Store
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

func (obj Store) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&Store{}, id)
	if err != nil {
		return nil, rs.Error
	}

	return data, nil
}
