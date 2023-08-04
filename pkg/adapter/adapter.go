package adapter

import (
	"github.com/DictumMortuum/servus-extapi/pkg/model"
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
