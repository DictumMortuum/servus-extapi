package model

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Wishlist struct {
	Id         int64  `gorm:"primaryKey" json:"id"`
	UserId     string `json:"user_id"`
	Url        string `json:"url"`
	Reserved   bool   `json:"reserved"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Screenshot string `json:"screenshot" gorm:"-"`
}

func (Wishlist) TableName() string {
	return "twishlist"
}

func (Wishlist) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (Wishlist) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []Wishlist
	rs := db.Scopes(scopes...).Find(&data)
	return data, rs.Error
}

func (Wishlist) Get(db *gorm.DB, id int64) (any, error) {
	var data Wishlist
	rs := db.First(&data, id)
	return data, rs.Error
}

func (obj Wishlist) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := Wishlist{
		Id: id,
	}

	var payload Wishlist
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

func (Wishlist) Create(db *gorm.DB, body []byte) (any, error) {
	var payload Wishlist
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

	raw, err := http.Get(payload.Screenshot)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(fmt.Sprintf("/data/cache/wish-%d.jpg", payload.Id))
	if err != nil {
		return nil, err
	}

	defer file.Close()

	_, err = io.Copy(file, raw.Body)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (obj Wishlist) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&Wishlist{}, id)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return data, nil
}
