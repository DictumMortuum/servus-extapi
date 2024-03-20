package model

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Book struct {
	Id          int64  `gorm:"primaryKey" json:"id,omitempty"`
	Category    string `json:"category,omitempty"`
	Name        string `json:"name,omitempty"`
	ISBN        string `json:"isbn,omitempty"`
	Publisher   string `json:"publisher,omitempty"`
	Translation string `json:"translation,omitempty"`
	Writer      string `json:"writer,omitempty"`
}

func (Book) TableName() string {
	return "tbooks"
}

func (Book) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (c Book) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []Book
	rs := db.Scopes(scopes...).Scopes(c.DefaultFilter).Find(&data)
	return data, rs.Error
}

func (Book) Get(db *gorm.DB, id int64) (any, error) {
	var data Book
	rs := db.First(&data, id)
	return data, rs.Error
}

func (obj Book) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := Book{
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

func (Book) Create(db *gorm.DB, body []byte) (any, error) {
	var payload Book
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

func (obj Book) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&IgnoredName{}, id)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return data, nil
}
