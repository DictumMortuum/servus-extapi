package model

import (
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type FinderUser struct {
	Id          int64          `gorm:"primaryKey" json:"id"`
	BggUsername string         `json:"bgg_username"`
	Collection  datatypes.JSON `gorm:"serializer:json" json:"collection"`
}

func (FinderUser) TableName() string {
	return "tfinderusers"
}

func (FinderUser) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (FinderUser) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []FinderUser
	rs := db.Scopes(scopes...).Find(&data)

	var retval []FinderUser
	for _, item := range data {
		item.Collection = nil
		retval = append(retval, item)
	}

	return retval, rs.Error
}

func (FinderUser) Get(db *gorm.DB, id int64) (any, error) {
	var data FinderUser
	rs := db.First(&data, id)
	return data, rs.Error
}

func (obj FinderUser) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := FinderUser{
		Id: id,
	}

	var payload FinderUser
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

func (FinderUser) Create(db *gorm.DB, body []byte) (any, error) {
	var payload FinderUser
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

func (obj FinderUser) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&FinderUser{}, id)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return data, nil
}
