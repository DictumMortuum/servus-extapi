package main

import (
	"encoding/json"
	"fmt"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
)

type player struct {
	Id               int64  `json:"id"`
	Uuid             string `json:"uuid"`
	Name             string `json:"name"`
	IsAnonymous      bool   `json:"isAnonymous"`
	ModificationDate string `json:"modificationDate"`
	BggUsername      string `json:"bggUsername"`
}

type location struct {
	Id               int64  `json:"id"`
	Uuid             string `json:"uuid"`
	Name             string `json:"name"`
	ModificationDate string `json:"modificationDate"`
}

type game struct {
	Id               int64  `json:"id"`
	RefId            int64  `json:"ref_id"`
	Uuid             string `json:"uuid"`
	Name             string `json:"name"`
	ModificationDate string `json:"modificationDate"`
	Cooperative      bool   `json:"cooperative"`
	HighestWins      bool   `json:"highestWins"`
	NoPoints         bool   `json:"noPoints"`
	UsesTeams        bool   `json:"usesTeams"`
	BoardgameId      int64  `json:"bggId"`
}

type play struct {
	Id               int64  `json:"id"`
	Uuid             string `json:"uuid"`
	ModificationDate string `json:"modificationDate"`
	PlayDate         string `json:"playDate"`
	UsesTeams        bool   `json:"usesTeams"`
	Ignored          bool   `json:"ignored"`
	ManualWinner     bool   `json:"manualWinner"`
	Rounds           uint   `json:"rounds"`
	LocationId       *int64 `json:"locationRefId"`
	GameId           *int64 `json:"gameRefId"`
	Stats            []stat `json:"playerScores"`
}

type stat struct {
	Score       *string `json:"score"`
	Winner      bool    `json:"winner"`
	NewPlayer   bool    `json:"newPlayer"`
	StartPlayer bool    `json:"startPlayer"`
	PlayerId    int64   `json:"playerRefId"`
	Rank        int     `json:"rank"`
	SeatOrder   uint    `json:"seatOrder"`
}

type BGStatsUpload struct {
	Players   []player   `json:"players"`
	Locations []location `json:"locations"`
	Games     []game     `json:"games"`
	Plays     []play     `json:"plays"`
}

func CreateBGStats(c *gin.Context, db *gorm.DB) (interface{}, error) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}

	var payload BGStatsUpload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	players := []model.BGStatsPlayer{}
	for _, player := range payload.Players {
		players = append(players, model.BGStatsPlayer{
			RefId:            player.Id,
			Uuid:             player.Uuid,
			Name:             player.Name,
			ModificationDate: player.ModificationDate,
			BggUsername:      player.BggUsername,
			PlayerId: models.JsonNullInt64{
				Int64: 0,
				Valid: false,
			},
		})
	}

	rs := db.Debug().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{}),
	}).Create(&players)
	if rs.Error != nil {
		return nil, err
	}

	players_map := map[int64]int64{}
	for _, player := range players {
		players_map[player.RefId] = player.Id
	}

	locations := []model.BGStatsLocation{}
	for _, location := range payload.Locations {
		locations = append(locations, model.BGStatsLocation{
			RefId:            location.Id,
			Uuid:             location.Uuid,
			Name:             location.Name,
			ModificationDate: location.ModificationDate,
			LocationId: models.JsonNullInt64{
				Int64: 0,
				Valid: false,
			},
		})
	}

	rs = db.Debug().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{}),
	}).Create(&locations)
	if rs.Error != nil {
		return nil, err
	}

	locations_map := map[int64]int64{}
	for _, location := range locations {
		locations_map[location.RefId] = location.Id
	}

	games := []model.BGStatsGame{}
	for _, game := range payload.Games {
		games = append(games, model.BGStatsGame{
			RefId:            game.Id,
			Uuid:             game.Uuid,
			Name:             game.Name,
			ModificationDate: game.ModificationDate,
			Cooperative:      game.Cooperative,
			HighestWins:      game.HighestWins,
			NoPoints:         game.NoPoints,
			UsesTeams:        game.UsesTeams,
			BoardgameId: models.JsonNullInt64{
				Int64: game.BoardgameId,
				Valid: true,
			},
		})
	}

	games_map := map[int64]int64{}
	for _, game := range games {
		rs = db.Debug().Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "uuid"}},
			DoUpdates: clause.Assignments(map[string]interface{}{}),
		}).Create(&game)
		if rs.Error != nil {
			return nil, err
		}
		games_map[game.RefId] = game.Id
	}

	for _, play := range payload.Plays {
		if play.GameId == nil {
			fmt.Println(play)
			continue
		}

		location := models.JsonNullInt64{
			Int64: 0,
			Valid: false,
		}

		game := models.JsonNullInt64{
			Int64: 0,
			Valid: false,
		}

		if play.LocationId != nil {
			location.Valid = true
			location.Int64 = locations_map[*play.LocationId]
		}

		if play.GameId != nil {
			game.Valid = true
			game.Int64 = games_map[*play.GameId]
		}

		payload := model.BGStatsPlay{
			Uuid:             play.Uuid,
			ModificationDate: play.ModificationDate,
			PlayDate:         play.PlayDate,
			GameId:           game,
			LocationId:       location,
			Rounds:           play.Rounds,
			ManualWinner:     play.ManualWinner,
			Ignored:          play.Ignored,
			UsesTeams:        play.UsesTeams,
		}

		rs = db.Debug().Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "uuid"}},
			DoUpdates: clause.Assignments(map[string]interface{}{}),
		}).Create(&payload)
		if rs.Error != nil {
			return nil, err
		}

		for _, stat := range play.Stats {
			score := models.JsonNullString{
				String: "",
				Valid:  false,
			}

			if stat.Score != nil {
				score.Valid = true
				score.String = *stat.Score
			}

			rs = db.Debug().Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "play_id"}, {Name: "player_id"}},
				DoNothing: true,
			}).Create(&model.BGStat{
				PlayId:      payload.Id,
				Score:       score,
				Winner:      stat.Winner,
				NewPlayer:   stat.NewPlayer,
				StartPlayer: stat.StartPlayer,
				PlayerId:    players_map[stat.PlayerId],
				Rank:        stat.Rank,
				SeatOrder:   stat.SeatOrder,
			})
			if rs.Error != nil {
				return nil, err
			}
		}
	}

	return payload, nil
}
