package main

import (
	"log"

	"github.com/DictumMortuum/servus-extapi/pkg/adapter"
	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func Version(c *gin.Context) {
	rs := map[string]any{
		"version": "v0.0.1",
	}
	c.AbortWithStatusJSON(200, rs)
}

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.Use(middleware.Cors())
	g := r.Group("/boardgames")
	g.GET("/version", Version)

	g.GET(
		"/all",
		middleware.BindYear,
		adapter.A(GetPlayedGames),
		middleware.Result,
	)

	g.GET(
		"/:id",
		middleware.Id,
		middleware.BindYear,
		middleware.Result,
	)

	r.Run(":10002")
}
