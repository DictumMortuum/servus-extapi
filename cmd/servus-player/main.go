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
		"version": "v0.0.13",
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
	g := r.Group("/player")
	g.GET("/version", Version)

	g.GET(
		"/all",
		middleware.Url,
		middleware.BindYear,
		adapter.A(GetPlayers),
		middleware.Result,
	)

	g.GET(
		"/:id",
		middleware.Id,
		middleware.Url,
		middleware.BindYear,
		middleware.BindN,
		adapter.A(GetPlayerDetail),
		adapter.A(GetPlayerGames),
		adapter.A(GetPlayerPlays),
		adapter.A(GetPlayers),
		adapter.A(GetLatestGames),
		adapter.A(GetPlayerScores),
		adapter.A(ProcessMechanics),
		adapter.A(ProcessDesigners),
		adapter.A(ProcessCategories),
		adapter.A(ProcessFamilies),
		adapter.A(ProcessSubdomains),
		adapter.A(ProcessLocations),
		adapter.A(GetNetwork),
		middleware.Result,
	)

	r.Run(":10001")
}
