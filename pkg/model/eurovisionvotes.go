package model

import (
	"encoding/json"
	"errors"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EurovisionVote struct {
	Id       int64          `gorm:"primaryKey" json:"id"`
	UserId   string         `json:"user_id"`
	Email    string         `json:"email"`
	Included bool           `json:"included"`
	Votes    datatypes.JSON `gorm:"serializer:json" json:"votes"`
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

	// https://stackoverflow.com/questions/56653423/gorm-doesnt-update-boolean-field-to-false
	rs := db.Model(&model).Updates(map[string]any{
		"UserId":   payload.UserId,
		"Email":    payload.Email,
		"Included": payload.Included,
	})
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

func GetEurovisionVotes(req *Map, res *Map) error {
	DB, err := req.GetGorm()
	if err != nil {
		return err
	}

	var data []EurovisionVote
	rs := DB.Find(&data, "included = true")
	if rs.Error != nil {
		return rs.Error
	}

	type temp struct {
		Boardgame Boardgame `json:"boardgame"`
	}

	type temp_res struct {
		Flag  string `json:"flag"`
		Name  string `json:"name"`
		Votes int    `json:"votes"`
	}

	votes := []int{12, 10, 8, 7, 6, 5, 4, 3, 2, 1}
	result := map[string]temp_res{}

	for _, vote := range data {
		var raw []temp
		err = json.Unmarshal(vote.Votes, &raw)
		if err != nil {
			return err
		}

		for i, item := range raw {
			if i < len(votes) && i >= 0 {
				tmp := result[item.Boardgame.Name]
				tmp.Votes += votes[i]
				tmp.Name = item.Boardgame.Name
				tmp.Flag = item.Boardgame.Square200
				result[item.Boardgame.Name] = tmp
			}
		}
	}

	final := []temp_res{}
	for _, item := range result {
		final = append(final, item)
	}

	res.Set("data", final)
	return nil
}

// {"flag": "224517", "name": " Brass: Birmingham", "votes": 60}
