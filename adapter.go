package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type Routable interface {
	List(*gorm.DB, ...func(*gorm.DB) *gorm.DB) (any, error)
	Get(*gorm.DB, int64) (any, error)
	Update(*gorm.DB, int64, []byte) (any, error)
	Create(*gorm.DB, []byte) (any, error)
	Delete(*gorm.DB, int64) (any, error)
}

func G(f func(*gin.Context, *gorm.DB) (interface{}, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		db, m, err := Database()
		if err != nil {
			c.JSON(http.StatusFailedDependency, gin.H{"error": err.Error()})
			return
		}
		defer m.Close()

		data, err := f(c, db)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, data)
	}
}
