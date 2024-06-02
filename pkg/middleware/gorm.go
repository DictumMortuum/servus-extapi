package middleware

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/mrz1836/go-sanitize"
	"gorm.io/gorm"
)

func Paginate(c *gin.Context) {
	m, err := model.ToMap(c, "req")
	if err != nil {
		c.Error(err)
		return
	}

	raw_range := c.Query("range")
	var payload []int
	var offset int
	var limit int
	err = json.Unmarshal([]byte(raw_range), &payload)
	if err != nil {
		offset = 0
		limit = 50
		m.Set("range", fmt.Sprintf("%d-%d", 0, 50))
	} else {
		offset = payload[0]
		limit = payload[1] - payload[0] + 1
		m.Set("range", fmt.Sprintf("%d-%d", payload[0], payload[1]))
	}

	m.Paginate = func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit)
	}
}

func argToGorm(db *gorm.DB, key string, val any, whole bool) *gorm.DB {
	if key == "whole" {
		return db
	}

	if strings.Contains(key, "@not") {
		key = strings.Split(key, "@")[0]
		return db.Not(key, val)
	} else if strings.Contains(key, "@like") {
		key = strings.Split(key, "@")[0]
		term := val.(string)

		if term != "%%" {
			return db.Where(key+" LIKE ?", val)
		}
	} else if strings.Contains(key, "@gt") {
		key = strings.Split(key, "@")[0]
		return db.Where(key+" >= ?", val)
	} else if strings.Contains(key, "@autolike") {
		key = strings.Split(key, "@")[0]
		temp := strings.TrimSpace(sanitize.AlphaNumeric(val.(string), true))
		terms := strings.Split(temp, " ")

		for _, term := range terms {
			if whole {
				db = db.Where(key + " COLLATE utf8mb4_unicode_ci REGEXP '[[:<:]]" + term + "[[:>:]]'")
			} else {
				db = db.Where(key+" COLLATE utf8mb4_unicode_ci LIKE ?", "%"+term+"%")
			}
		}

		return db
	} else if strings.Contains(key, "@rightlike") {
		key = strings.Split(key, "@")[0]
		temp := strings.TrimSpace(sanitize.AlphaNumeric(val.(string), true))
		terms := strings.Split(temp, " ")

		for _, term := range terms {
			db = db.Where(key+" LIKE ?", term+"%")
		}

		return db
	} else {
		return db.Where(key, val)
	}

	return db
}

func Filter(c *gin.Context) {
	m, err := model.ToMap(c, "req")
	if err != nil {
		c.Error(err)
		return
	}

	m.Filter = filter(c)
}

func filter(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		raw := c.Query("filter")
		var payload map[string]any
		err := json.Unmarshal([]byte(raw), &payload)
		if err != nil {
			return db
		} else {
			whole := false
			if val, ok := payload["whole"]; ok {
				if val.(float64) == 1 {
					whole = true
				}
			}

			for key, val := range payload {
				switch val.(type) {
				case map[string]any:
					continue
				default:
					db = argToGorm(db, key, val, whole)
				}
			}

			return db
		}
	}
}

func Sort(c *gin.Context) {
	m, err := model.ToMap(c, "req")
	if err != nil {
		c.Error(err)
		return
	}

	m.Sort = sort(c)
}

func sort(c *gin.Context) func(db *gorm.DB) *gorm.DB {
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
