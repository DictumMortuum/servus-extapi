package main

import (
	"encoding/json"
	"strings"

	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mrz1836/go-sanitize"
)

func insertTrace(DB *sqlx.DB, name string) error {
	q := `
		insert into ttraces (
			cr_date,
			type,
			raw
		) values (
			NOW(),
			:type,
			:raw
		) on duplicate key update cr_date = NOW()
	`
	_, err := DB.NamedExec(q, map[string]any{
		"raw":  name,
		"type": "search",
	})
	if err != nil {
		return err
	}

	return nil
}

func searchFilter(c *gin.Context) {
	m, err := model.ToMap(c, "req")
	if err != nil {
		c.Error(err)
		return
	}

	DB, err := m.GetDB()
	if err != nil {
		c.Error(err)
		return
	}

	raw := c.Query("filter")
	var payload map[string]any
	err = json.Unmarshal([]byte(raw), &payload)
	if err != nil {
		c.Error(err)
		return
	}

	for key, val := range payload {
		if key == "name@autolike" {
			sanitized := sanitize.AlphaNumeric(val.(string), true)
			sanitized = strings.TrimSpace(sanitized)

			if len(sanitized) < 4 {
				continue
			}

			err := insertTrace(DB, sanitized)
			if err != nil {
				c.Error(err)
				return
			}
		}
	}
}
