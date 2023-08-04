package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Database() (*gorm.DB, *sql.DB, error) {
	dsn := config.Cfg.Databases["mariadb"]

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, nil, err
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:              sqlDB,
		DefaultStringSize: 512,
	}), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, nil, err
	}

	db.AutoMigrate(&model.Boardgame{}, &model.Play{}, &model.Player{}, &model.Stat{}, &model.Price{}, &model.Store{}, &model.Location{}, &model.BGStatsPlayer{}, &model.BGStatsLocation{}, &model.BGStatsGame{}, &model.BGStatsPlay{}, &model.BGStat{}, &model.IgnoredPrice{}, &model.CachedPrice{}, &model.IgnoredName{})
	return db, sqlDB, nil
}

func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		raw_range := c.Query("range")

		var payload []int
		var offset int
		var limit int
		err := json.Unmarshal([]byte(raw_range), &payload)
		if err != nil {
			offset = 0
			limit = 50
		} else {
			offset = payload[0]
			limit = payload[1] - payload[0] + 1
		}

		return db.Offset(offset).Limit(limit)
	}
}

func argToGorm(db *gorm.DB, key string, val any) *gorm.DB {
	if strings.Contains(key, "@not") {
		key = strings.Split(key, "@")[0]
		return db.Not(key, val)
	} else if strings.Contains(key, "@like") {
		key = strings.Split(key, "@")[0]
		term := val.(string)

		if term != "%%" {
			return db.Where(key+" LIKE ?", val)
		}
	} else {
		return db.Where(key, val)
	}

	return db
}

func EmbeddedFilter(c *gin.Context, col string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		raw := c.Query("filter")

		var payload map[string]any
		err := json.Unmarshal([]byte(raw), &payload)
		if err != nil {
			return db
		} else {
			for key, val := range payload {
				switch val.(type) {
				case map[string]any:
					fmt.Println(val)
					if key == col {
						for nested_key, nested_val := range val.(map[string]any) {
							db = argToGorm(db, nested_key, nested_val)
						}
					}
				default:
					fmt.Println(val)
				}
			}

			return db
		}
	}
}

func Filter(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		raw := c.Query("filter")

		var payload map[string]any
		err := json.Unmarshal([]byte(raw), &payload)
		if err != nil {
			return db
		} else {
			for key, val := range payload {
				switch val.(type) {
				case map[string]any:
					continue
				default:
					db = argToGorm(db, key, val)
				}
			}

			return db
		}
	}
}

func Sort(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		raw := c.Query("sort")

		var payload []string
		err := json.Unmarshal([]byte(raw), &payload)
		if err != nil {
			return db
		} else {
			return db.Order(strings.Join(payload, " "))
		}
	}
}
