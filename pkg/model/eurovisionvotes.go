package model

import (
	"encoding/json"
	"errors"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EurovisionVote struct {
	Id     int64          `gorm:"primaryKey" json:"id"`
	UserId string         `json:"user_id"`
	Email  string         `json:"email"`
	Votes  datatypes.JSON `gorm:"serializer:json" json:"votes"`
}

func (EurovisionVote) TableName() string {
	return "teurovisionvotes"
}

func (EurovisionVote) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (EurovisionVote) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []EurovisionVote
	rs := db.Scopes(scopes...).Find(&data)
	return data, rs.Error
}

func (EurovisionVote) Get(db *gorm.DB, id int64) (any, error) {
	var data EurovisionVote
	rs := db.First(&data, id)
	return data, rs.Error
}

func GetEurovisionVoteByUserId(req *Map, res *Map) error {
	DB, err := req.GetGorm()
	if err != nil {
		return err
	}

	id, err := req.GetString("id")
	if err != nil {
		return err
	}

	var data EurovisionVote
	rs := DB.First(&data, "user_id = ? ", id)
	if errors.Is(rs.Error, gorm.ErrRecordNotFound) {
		res.Set("data", nil)
		return nil
	}
	if rs.Error != nil {
		return rs.Error
	}

	res.Set("data", data)
	return nil
}

func (obj EurovisionVote) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := EurovisionVote{
		Id: id,
	}

	var payload EurovisionVote
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

func (EurovisionVote) Create(db *gorm.DB, body []byte) (any, error) {
	var payload EurovisionVote
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	rs := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&payload)
	if rs.Error != nil {
		return nil, err
	}

	return payload, nil
}

func (obj EurovisionVote) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&EurovisionVote{}, id)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return data, nil
}
