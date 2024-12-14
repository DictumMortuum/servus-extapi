package model

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Configuration struct {
	Id     int64  `gorm:"primaryKey" json:"id"`
	Config string `json:"config"`
	Value  bool   `json:"value"`
}

func (Configuration) TableName() string {
	return "tconfig"
}

func (Configuration) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (c Configuration) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []Configuration
	rs := db.Scopes(scopes...).Scopes(c.DefaultFilter).Find(&data)
	return data, rs.Error
}

func (Configuration) Get(db *gorm.DB, id int64) (any, error) {
	var data Configuration
	rs := db.First(&data, id)
	return data, rs.Error
}

func (obj Configuration) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := Configuration{
		Id: id,
	}

	var payload Configuration
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	// https://stackoverflow.com/questions/56653423/gorm-doesnt-update-boolean-field-to-false
	rs := db.Model(&model).Updates(map[string]any{
		"config": payload.Config,
		"value":  payload.Value,
	})
	if rs.Error != nil {
		return nil, err
	}

	return obj.Get(db, id)
}

func (Configuration) Create(db *gorm.DB, body []byte) (any, error) {
	var payload Configuration
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

func (obj Configuration) Delete(db *gorm.DB, id int64) (any, error) {
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
