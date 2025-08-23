package model

import (
	"encoding/json"
	"log"
	"sort"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Play struct {
	Id          int64          `gorm:"primaryKey" json:"id"`
	Uuid        string         `json:"uuid"`
	BoardgameId int64          `json:"boardgame_id"`
	Boardgame   Boardgame      `json:"boardgame"`
	Date        time.Time      `json:"date"`
	PlayData    datatypes.JSON `gorm:"serializer:json" json:"play_data"`
	Stats       []Stat         `json:"stats"`
	LocationId  int64          `json:"location_id"`
	Location    Location       `json:"location"`
}

type PlayData struct {
	Draws       []bool    `json:"draws"`
	Players     []int64   `json:"players"`
	Winners     []int64   `json:"winners"`
	Solo        bool      `json:"solo"`
	Cooperative bool      `json:"cooperative"`
	Teams       [][]int64 `json:"teams"`
}

func (user *Play) BeforeCreate(tx *gorm.DB) (err error) {
	user.Uuid = uuid.NewString()
	return
}

func (Play) TableName() string {
	return "tboardgameplays"
}

func (Play) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

func (Play) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []Play
	rs := db.Scopes(scopes...).Preload("Boardgame").Preload("Stats").Preload("Location").Preload("Stats.Player").Find(&data)
	return data, rs.Error
}

func (Play) Get(db *gorm.DB, id int64) (any, error) {
	var data Play
	rs := db.Preload("Boardgame").Preload("Stats.Player").Preload("Location").First(&data, id)
	return data, rs.Error
}

func (obj Play) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	var payload Play
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	var play_data map[string]any
	err = json.Unmarshal(payload.PlayData, &play_data)
	if err != nil {
		return nil, err
	}

	var play_data2 PlayData
	err = json.Unmarshal(payload.PlayData, &play_data2)
	if err != nil {
		return nil, err
	}

	log.Println(play_data2)

	stats := payload.Stats
	sort.Slice(stats, func(i, j int) bool {
		scorei, err := payload.Boardgame.Score(&stats[i])
		if err != nil {
			log.Fatal(err)
		}

		scorej, err := payload.Boardgame.Score(&stats[j])
		if err != nil {
			log.Fatal(err)
		}

		return scorei > scorej
	})

	players := []int64{}
	for _, item := range stats {
		players = append(players, item.PlayerId)
	}

	data := map[string]any{
		"players": players,
	}

	if payload.Boardgame.Solitaire && len(payload.Stats) == 1 {
		data["solo"] = true
	}

	if payload.Boardgame.Cooperative {
		data["cooperative"] = true
	}

	if payload.Boardgame.Cooperative {
		winners := []int64{}

		for _, stat := range payload.Stats {
			var stat_data map[string]any
			err := json.Unmarshal(stat.Data, &stat_data)
			if err != nil {
				return nil, err
			}

			if val, ok := stat_data["won"]; ok {
				if val.(bool) {
					winners = append(winners, stat.PlayerId)
				}
			}
		}

		data["winners"] = winners
	} else {
		draws := []bool{}
		// if there are teams present
		if val, ok := play_data["teams"]; ok {
			data["teams"] = val

			team_scores := []float64{}
			log.Println(data["teams"])
			for _, team := range data["teams"].([]any) {
				team_score := 0.0
				for _, id := range team.([]any) {
					for _, stat := range payload.Stats {
						if int64(id.(float64)) == stat.PlayerId {
							score, err := payload.Boardgame.Score(&stat)
							if err != nil {
								log.Fatal(err)
							}
							team_score += score
						}
					}
				}
				team_scores = append(team_scores, team_score)
			}

			for i := 0; i < len(team_scores)-1; i++ {
				draws = append(draws, team_scores[i] == team_scores[i+1])
			}
		} else {
			for i := 0; i < len(stats)-1; i++ {
				scorei, err := payload.Boardgame.Score(&stats[i])
				if err != nil {
					log.Fatal(err)
				}

				scorej, err := payload.Boardgame.Score(&stats[i+1])
				if err != nil {
					log.Fatal(err)
				}

				draws = append(draws, scorei == scorej)
			}

			winners := []int64{payload.Stats[0].PlayerId}
			for i, draw := range draws {
				if draw {
					winners = append(winners, payload.Stats[i+1].PlayerId)
				} else {
					break
				}
			}
			data["winners"] = winners
		}
		data["draws"] = draws
	}

	// for _, item := range stats {
	// 	log.Println(item)
	// }

	payload.Stats = stats

	raw_json, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	payload.PlayData = datatypes.JSON(raw_json)

	rs := db.Debug().Omit("Boardgame").Omit("Location").Omit("Stats.Player").Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Session(&gorm.Session{
		FullSaveAssociations: true,
	}).Create(&payload)
	if rs.Error != nil {
		return nil, err
	}

	return obj.Get(db, id)
}

func (obj Play) Create(db *gorm.DB, body []byte) (any, error) {
	var payload Play
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	rs := db.Debug().Omit("Boardgame").Omit("Location").Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Session(&gorm.Session{
		FullSaveAssociations: true,
	}).Create(&payload)
	if rs.Error != nil {
		return nil, err
	}

	// log.Println(obj.Get(db, payload.Id))

	// log.Println(payload)

	return payload, nil
}

func (obj Play) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&Play{}, id)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return data, nil
}
