package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/middleware"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDB(c *gin.Context) *gorm.DB {
	val, ok := c.Get("apidb")

	if ok && val != nil {
		db, _ := val.(*gorm.DB)
		return db
	}

	return nil
}

func GetModel(c *gin.Context) model.Routable {
	val, ok := c.Get("apimodel")

	if ok && val != nil {
		return val.(model.Routable)
	}

	return nil
}

func GetBody(c *gin.Context) []byte {
	val, ok := c.Get("apibody")

	if ok && val != nil {
		rs, _ := val.([]byte)
		return rs
	}

	return nil
}

func GetOne(f func(*gorm.DB, int64) (any, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.GetInt64("apiid")
		db := GetDB(c)

		data, err := f(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusFailedDependency, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, data)
		c.Next()
	}
}

func LoadOne(f func(*gorm.DB, int64) (any, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.GetInt64("apiid")
		db := GetDB(c)

		data, err := f(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusFailedDependency, gin.H{"error": err.Error()})
		}

		raw, err := json.Marshal(data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusFailedDependency, gin.H{"error": err.Error()})
		}

		var m map[string]any
		err = json.Unmarshal(raw, &m)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusFailedDependency, gin.H{"error": err.Error()})
		}

		c.Set("data", data)
		c.Set("mapped_data", m)
		c.Next()
	}
}

func CountMany(c *gin.Context) {
	db := GetDB(c)
	model := GetModel(c)

	var count int64
	rs := db.Model(model).Scopes(middleware.Filter(c), model.DefaultFilter).Count(&count)
	if rs.Error != nil {
		c.AbortWithStatusJSON(http.StatusFailedDependency, gin.H{"error": rs.Error.Error()})
	}
	c.Header("Content-Range", fmt.Sprintf("%d", count))

	c.Next()
}

func GetMany(f func(*gorm.DB, ...func(*gorm.DB) *gorm.DB) (any, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		db := GetDB(c)

		data, err := f(db, middleware.Filter(c), middleware.Paginate(c), middleware.Sort(c))
		// data, err := f(db, Filter(c), Paginate(c))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusFailedDependency, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, data)

		c.Next()
	}
}

func UpdateOne(f func(*gorm.DB, int64, []byte) (any, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.GetInt64("apiid")
		db := GetDB(c)
		body := GetBody(c)

		data, err := f(db, id, body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusFailedDependency, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, data)

		c.Next()
	}
}

func CreateOne(f func(*gorm.DB, []byte) (any, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		db := GetDB(c)
		body := GetBody(c)

		data, err := f(db, body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusFailedDependency, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, data)

		c.Next()
	}
}

func DeleteOne(f func(*gorm.DB, int64) (any, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.GetInt64("apiid")
		db := GetDB(c)

		data, err := f(db, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusFailedDependency, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, data)
		c.Next()
	}
}

func OpenDB(c *gin.Context) {
	db, conn, err := db.Gorm()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusFailedDependency, gin.H{"error": err.Error()})
	}
	c.Set("apidb", db)
	c.Set("apiconn", conn)
	c.Next()
}

func CloseDB(c *gin.Context) {
	raw, ok := c.Get("apiconn")
	if ok {
		conn, _ := raw.(*sql.DB)
		conn.Close()
		c.AbortWithStatus(http.StatusOK)
	} else {
		c.AbortWithStatusJSON(http.StatusFailedDependency, gin.H{"error": "db conn not found"})
	}

	// c.Next()
}

func Id(c *gin.Context) {
	raw_id := c.Param("id")

	id, err := strconv.ParseInt(raw_id, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
	}
	c.Set("apiid", id)

	c.Next()
}

func Body(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
	}
	c.Set("apibody", body)

	c.Next()
}
