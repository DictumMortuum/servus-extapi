package main

import (
	"log"

	"github.com/DictumMortuum/servus-extapi/pkg/adapter"
	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/middleware"
	"github.com/DictumMortuum/servus-extapi/pkg/queries"
	"github.com/gin-gonic/gin"
)

func Version(c *gin.Context) {
	rs := map[string]any{
		"version": "v0.0.5",
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
		middleware.Url,
		middleware.BindYear,
		adapter.A(GetPlayedGames),
		middleware.Result,
	)

	g.GET(
		"/:id",
		middleware.Id,
		middleware.Url,
		middleware.BindYear,
		middleware.BindN,
		adapter.A(GetBoardgamePlays),
		adapter.A(GetLatestBoardgames),
		adapter.A(GetBoardgameDetail),
		adapter.A(queries.GetPlayers),
		middleware.Result,
	)

	g.GET(
		"/info/:id",
		middleware.Id,
		adapter.A(GetBoardgameInfo),
		middleware.Result,
	)

	r.Run(":10002")
}
