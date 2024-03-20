package model

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Player struct {
	Id      int64   `gorm:"primaryKey" json:"id"`
	Name    string  `json:"name"`
	Surname string  `json:"surname"`
	Email   *string `json:"email"`
	Hidden  bool    `json:"hidden"`
	// BGStatsPlayers []BGStatsPlayer `json:"bg_stats_players"`
}

func (Player) TableName() string {
	return "tboardgameplayers"
}

func (Player) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (Player) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []Player
	// rs := db.Scopes(scopes...).Preload("BGStatsPlayers").Find(&data)
	rs := db.Scopes(scopes...).Find(&data)
	return data, rs.Error
}

func (Player) Get(db *gorm.DB, id int64) (any, error) {
	var data Player
	// rs := db.Preload("BGStatsPlayers").First(&data, id)
	rs := db.First(&data, id)
	return data, rs.Error
}

func (obj Player) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := Player{
		Id: id,
	}

	var payload Player
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

func (Player) Create(db *gorm.DB, body []byte) (any, error) {
	var payload Player
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

func (obj Player) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&Player{}, id)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return data, nil
}
