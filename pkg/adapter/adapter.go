package adapter

import (
	"net/http"
	"time"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func A(f func(*model.Map, *model.Map) error) func(*gin.Context) {
	return func(c *gin.Context) {
		if len(c.Errors) > 0 {
			return
		}

		req, err := model.ToMap(c, "req")
		if err != nil {
			c.Error(err)
			return
		}

		res, err := model.ToMap(c, "res")
		if err != nil {
			c.Error(err)
			return
		}

		err = f(req, res)
		if err != nil {
			c.Error(err)
			return
		}
	}
}

func C(store *persistence.InMemoryStore, t time.Duration, f func(*model.Map, *model.Map) error) func(*gin.Context) {
	return cache.CachePage(store, t, A(f))
}

func G(f func(*gin.Context, *gorm.DB) (interface{}, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		DB, m, err := db.Gorm()
		if err != nil {
			c.JSON(http.StatusFailedDependency, gin.H{"error": err.Error()})
			return
		}
		defer m.Close()

		data, err := f(c, DB)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, data)
	}
}
