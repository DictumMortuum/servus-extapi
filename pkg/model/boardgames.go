package model

import (
	"encoding/json"
	"github.com/DictumMortuum/servus/pkg/models"
	"gorm.io/gorm"
)

type Boardgame struct {
	Id             int64                `gorm:"primaryKey" json:"id"`
	Name           string               `gorm:"index" json:"name"`
	Year           int64                `json:"year"`
	MinPlayers     int                  `json:"minplayers"`
	MaxPlayers     int                  `json:"maxplayers"`
	Square200      string               `json:"square200"`
	BggDataNotNull bool                 `gorm:"-" json:"bgg_data_not_null"`
	Rank           models.JsonNullInt64 `gorm:"index" json:"rank"`
	RankNotNull    bool                 `gorm:"-" json:"rank_not_null"`
	Prices         []Price              `json:"prices"`
}

func (Boardgame) TableName() string {
	return "tboardgames"
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

	var payload Boardgame
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

	rs := db.Debug().Exec("insert into tboardgames (id, bgg_data) values (?, ?) on duplicate key update bgg_data = bgg_data", id, string(body))
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
			name = REPLACE(json_extract(bgg_data, '$.name'), '"', '')
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

	rs := db.Delete(&Player{}, id)
	if err != nil {
		return nil, rs.Error
	}

	return data, nil
}
