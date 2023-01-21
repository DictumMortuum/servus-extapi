package model

import (
	"encoding/json"
	"github.com/DictumMortuum/servus/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BGStatsPlay struct {
	Id               int64                `json:"id"`
	Uuid             string               `json:"uuid" gorm:"uniqueIndex"`
	ModificationDate string               `json:"modification_date"`
	PlayDate         string               `json:"play_date"`
	UsesTeams        bool                 `json:"uses_teams"`
	Ignored          bool                 `json:"ignored"`
	ManualWinner     bool                 `json:"manual_winner"`
	Rounds           uint                 `json:"rounds"`
	Stats            []BGStat             `json:"stats" gorm:"foreignKey:PlayId"`
	LocationId       models.JsonNullInt64 `json:"location_id"`
	Location         BGStatsLocation      `json:"location"`
	GameId           models.JsonNullInt64 `json:"game_id"`
	Game             BGStatsGame          `json:"game"`
	PlayId           models.JsonNullInt64 `json:"play_id"`
	Play             Play                 `json:"play"`
}

func (BGStatsPlay) TableName() string {
	return "tbgstatsplays"
}

func (BGStatsPlay) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []BGStatsPlay
	rs := db.Scopes(scopes...).Preload("Location").Preload("Stats").Preload("Stats.Player").Preload("Game").Find(&data)
	return data, rs.Error
}

func (BGStatsPlay) Get(db *gorm.DB, id int64) (any, error) {
	var data BGStatsPlay
	rs := db.Preload("Location").Preload("Stats").Preload("Stats.Player").First(&data, id)
	return data, rs.Error
}

func (obj BGStatsPlay) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	// model := BGStatsPlay{
	// 	Id: id,
	// }

	var payload BGStatsPlay
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	payload.Id = id

	rs := db.Debug().Omit("Location").Omit("Stats.Player").Omit("Play").Clauses(clause.OnConflict{UpdateAll: true}).Session(&gorm.Session{FullSaveAssociations: true}).Updates(&payload)
	if rs.Error != nil {
		return nil, err
	}

	return obj.Get(db, id)
}

func (BGStatsPlay) Create(db *gorm.DB, body []byte) (any, error) {
	var payload BGStatsPlay
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

func (obj BGStatsPlay) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&Player{}, id)
	if err != nil {
		return nil, rs.Error
	}

	return data, nil
}
