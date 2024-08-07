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
}

func Force(c *gin.Context) {
	f := c.Query("force")

	m, err := model.ToMap(c, "req")
	if err != nil {
		c.Error(err)
		return
	}

	m.Set("force", f == "true")
}

func Url(c *gin.Context) {
	m, err := model.ToMap(c, "req")
	if err != nil {
		c.Error(err)
		return
	}

	m.Set("url", c.Request.URL.String())
}

func Num(c *gin.Context) {
	id := c.Param("num")

	m, err := model.ToMap(c, "req")
	if err != nil {
		c.Error(err)
		return
	}

	m.Set("num", id)
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

	for key, val := range res.Headers {
		c.Header(key, val)
	}

	if len(c.Errors) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": c.Errors,
		})
	} else {
		c.JSON(http.StatusOK, res.Internal)
	}
}

func ResultRa(c *gin.Context) {
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

	for key, val := range res.Headers {
		c.Header(key, val)
	}

	if len(c.Errors) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": c.Errors,
		})
	} else {
		if val, ok := res.Internal["data"]; ok {
			c.JSON(http.StatusOK, val)
		} else {
			c.JSON(http.StatusOK, []any{})
		}
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
