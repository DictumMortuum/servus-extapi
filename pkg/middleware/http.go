package middleware

import (
	"io"
	"net/http"

	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/gin-gonic/gin"
)

func Id(c *gin.Context) {
	id := c.Param("id")

	m, err := model.ToMap(c, "req")
	if err != nil {
		c.Error(err)
		return
	}

	m.Set("id", id)
	m.Set("url", c.Request.URL.String())
}

func Body(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
		return
	}

	m, err := model.ToMap(c, "req")
	if err != nil {
		c.Error(err)
		return
	}

	m.Set("body", body)
}

func Result(c *gin.Context) {
	req, err := model.ToMap(c, "req")
	if err != nil {
		c.Error(err)
	}

	res, err := model.ToMap(c, "res")
	if err != nil {
		c.Error(err)
	}

	err = req.Close()
	if err != nil {
		c.Error(err)
	}

	if len(c.Errors) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": c.Errors,
		})
	} else {
		c.JSON(http.StatusOK, res.Internal)
	}
}

func BindYear(c *gin.Context) {
	m, err := model.ToMap(c, "req")
	if err != nil {
		c.Error(err)
		return
	}

	type Args struct {
		Year     string `form:"year"`
		YearFlag bool   `form:"year_flag"`
	}

	var payload Args
	c.ShouldBind(&payload)

	m.Set("year", payload.Year)
	m.Set("year_flag", payload.YearFlag)
}

func BindN(c *gin.Context) {
	m, err := model.ToMap(c, "req")
	if err != nil {
		c.Error(err)
		return
	}

	type Args struct {
		Count int64 `form:"count"`
	}

	var payload Args
	c.ShouldBind(&payload)

	m.Set("n", payload.Count)
}
