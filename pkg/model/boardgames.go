package model

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/DictumMortuum/servus/pkg/models"
	"gorm.io/gorm"
)

type Boardgame struct {
	Id             int64                `gorm:"primaryKey" json:"id"`
	Uuid           string               `json:"uuid"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	Name           string               `gorm:"index" json:"name"`
	Year           int64                `json:"year"`
	MinPlayers     int                  `json:"minplayers"`
	MaxPlayers     int                  `json:"maxplayers"`
	Square200      string               `json:"square200"`
	BggDataNotNull bool                 `gorm:"-" json:"bgg_data_not_null"`
	Rank           models.JsonNullInt64 `gorm:"index" json:"rank"`
	RankNotNull    bool                 `gorm:"-" json:"rank_not_null"`
	Cooperative    bool                 `gorm:"cooperative" json:"cooperative"`
	Solitaire      bool                 `gorm:"solitaire" json:"solitaire"`
	Prices         []BoardgamePrice     `json:"prices"`
	Data           models.Json          `db:"data" json:"data"`
}

func (Boardgame) TableName() string {
	return "tboardgames"
}

func (Boardgame) DefaultFilter(db *gorm.DB) *gorm.DB {
	return db
}

// func pricesPreload(c *gin.Context, col string) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		raw := c.Query("filter")

// 		var payload map[string]any
// 		err := json.Unmarshal([]byte(raw), &payload)
// 		if err != nil {
// 			return db
// 		} else {
// 			for key, val := range payload {
// 				switch val.(type) {
// 				case map[string]any:
// 					fmt.Println(val)
// 					if key == col {
// 						for nested_key, nested_val := range val.(map[string]any) {
// 							db = argToGorm(db, nested_key, nested_val)
// 						}
// 					}
// 				default:
// 					fmt.Println(val)
// 				}
// 			}

// 			return db
// 		}
// 	}
// }

// var boardgames []Boardgame
// rs := db.Debug().Scopes(Filter(c), Paginate(c), Sort(c)).Preload("Prices", pricesPreload(c, "prices")).Find(&boardgames)
// if rs.Error != nil {
// 	return nil, rs.Error
// }

func (Boardgame) List(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (any, error) {
	var data []Boardgame
	rs := db.Scopes(scopes...).Find(&data)
	return data, rs.Error
}

func (Boardgame) Get(db *gorm.DB, id int64) (any, error) {
	var data Boardgame
	rs := db.First(&data, id)
	return data, rs.Error
}

func (obj Boardgame) Update(db *gorm.DB, id int64, body []byte) (any, error) {
	model := Boardgame{
		Id: id,
	}

	payload := Boardgame{
		Id: id,
	}
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	fmt.Println(payload)

	rs := db.Model(&model).Save(payload)
	if rs.Error != nil {
		return nil, err
	}

	return obj.Get(db, id)
}

// func (Boardgame) Create(db *gorm.DB, body []byte) (any, error) {
// 	var payload Boardgame
// 	err := json.Unmarshal(body, &payload)
// 	if err != nil {
// 		return nil, err
// 	}

// 	rs := db.Create(&payload)
// 	if rs.Error != nil {
// 		return nil, err
// 	}

// 	return payload, nil
// }

func (Boardgame) Create(db *gorm.DB, body []byte) (any, error) {
	var payload map[string]any
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	var id int64
	if val, ok := payload["objectid"]; ok {
		id = int64(val.(float64))
	}

	rs := db.Exec("insert into tboardgames (id, bgg_data) values (?, ?) on duplicate key update bgg_data = ?", id, string(body), string(body))
	if rs.Error != nil {
		return nil, err
	}

	rs = db.Exec(`
		update tboardgames set
			rank = CONVERT(REPLACE(json_extract(bgg_data, '$.rankinfo[0].rank'), '"', ''), INT),
			year = CONVERT(REPLACE(json_extract(bgg_data, '$.yearpublished'), '"', ''), INT),
			min_players = CONVERT(REPLACE(json_extract(bgg_data, '$.minplayers'), '"', ''), INT),
			max_players = CONVERT(REPLACE(json_extract(bgg_data, '$.maxplayers'), '"', ''), INT),
			square200 = REPLACE(json_extract(bgg_data, '$.images.square200'), '"', ''),
			name = REPLACE(json_extract(bgg_data, '$.name'), '"', ''),
			updated_at = NOW()
		where id = ?
	`, id)
	if rs.Error != nil {
		return nil, err
	}

	return payload, nil
}

func (obj Boardgame) Delete(db *gorm.DB, id int64) (any, error) {
	data, err := obj.Get(db, id)
	if err != nil {
		return nil, err
	}

	rs := db.Delete(&Boardgame{}, id)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return data, nil
}

type column struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Tiebreak string  `json:"tiebreak,omitempty"`
	Factor   float64 `json:"factor,omitempty"`
	Exclude  bool    `json:"exclude,omitempty"`
}

type settings struct {
	Columns     []column `json:"columns,omitempty"`
	Cooperative bool     `json:"cooperative,omitempty"`
	Autowin     []string `json:"autowin,omitempty"`
}

func (obj Boardgame) Score(stat *Stat) (float64, error) {
	if stat == nil {
		return 0, nil
	}

	var stat_data map[string]any
	err := json.Unmarshal(stat.Data, &stat_data)
	if err != nil {
		return 0, err
	}

	if val, ok := stat_data["score"].(float64); ok {
		return val, nil
	}

	if val, ok := stat_data["winner"].(bool); ok {
		if val {
			return 1, nil
		} else {
			return 0, nil
		}
	}

	var s settings
	err = obj.Data.Unmarshal(&s)
	if err != nil {
		return 0, err
	}

	log.Println(s)
	score := 0.0
	base := 0.01

	for _, col := range s.Columns {
		if col.Tiebreak != "" {
			factor := col.Factor
			if factor == 0 {
				factor = base
			}

			if val, ok := stat_data[col.Name].(float64); ok {
				if col.Tiebreak == "asc" {
					score += val * base
				} else if col.Tiebreak == "desc" {
					score -= val * base
				}
			}
			base *= 0.1
		}

		if col.Exclude {
			continue
		}

		switch col.Type {
		case "int":
			if val, ok := stat_data[col.Name].(float64); ok {
				score += val
			}
		case "negint":
			if val, ok := stat_data[col.Name].(float64); ok {
				score -= val
			}
		default:
			if val, ok := stat_data[col.Name].(float64); ok {
				score += val
			}
		}
	}

	return score, nil
}
