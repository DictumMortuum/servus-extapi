package adapter

import (
	"time"

	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
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
