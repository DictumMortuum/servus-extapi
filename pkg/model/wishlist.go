package model

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/DictumMortuum/servus-extapi/pkg/screenshot"
	"github.com/gin-gonic/gin"
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

	// https://stackoverflow.com/questions/56653423/gorm-doesnt-update-boolean-field-to-false
	rs := db.Model(&model).Updates(map[string]any{
		"UserId":   payload.UserId,
		"Email":    payload.Email,
		"Reserved": payload.Reserved,
		"Name":     payload.Name,
		"Url":      payload.Url,
	})
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

	// raw, err := http.Get(payload.Screenshot)
	// if err != nil {
	// 	return nil, err
	// }

	// file, err := os.Create(fmt.Sprintf("/data/cache/wish-%d.jpg", payload.Id))
	// if err != nil {
	// 	return nil, err
	// }
	// defer file.Close()

	// _, err = io.Copy(file, raw.Body)
	// if err != nil {
	// 	return nil, err
	// }

	err = screenshot.Do(payload.Url, payload.Id)
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

func GetWishlistScreenshot(c *gin.Context) {
	req, err := ToMap(c, "req")
	if err != nil {
		c.Error(err)
		return
	}

	id, err := req.GetInt64("id")
	if err != nil {
		c.Error(err)
		return
	}

	data, err := screenshot.Get(id)
	if err != nil {
		c.Error(err)
		return
	}

	info, err := data.Stat()
	if err != nil {
		c.Error(err)
		return
	}

	c.DataFromReader(http.StatusOK, info.Size, "application/octet-stream", data, map[string]string{})
}

func UpdateWishlistScreenshot(c *gin.Context) {
	req, err := ToMap(c, "req")
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err,
		})
		return
	}

	id, err := req.GetInt64("id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err,
		})
		return
	}

	DB, err := req.GetGorm()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err,
		})
		return
	}

	var data Wishlist
	rs := DB.First(&data, "id = ? ", id)
	if errors.Is(rs.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, map[string]any{
			"status": "not found",
		})
		return
	}
	if rs.Error != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"error": rs.Error,
		})
		return
	}

	err = screenshot.Do(data.Url, data.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"error": rs.Error,
		})
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "updated",
	})
}
