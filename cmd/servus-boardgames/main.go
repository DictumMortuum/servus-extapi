package main

import (
	"log"

	"github.com/DictumMortuum/servus-extapi/pkg/adapter"
	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/middleware"
	"github.com/DictumMortuum/servus-extapi/pkg/queries"
	"github.com/DictumMortuum/servus-extapi/pkg/queue"
	"github.com/adjust/rmq/v5"
	"github.com/gin-gonic/gin"
)

func Version(c *gin.Context) {
	rs := map[string]any{
		"version": "v0.0.18",
	}
	c.AbortWithStatusJSON(200, rs)
}

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	connection, err := rmq.OpenConnection("handler", "tcp", config.Cfg.Databases["redis"], 2, nil)
	if err != nil {
		log.Fatal(err)
	}

	go queue.Cleaner(connection)
	go Consumer(connection)

	r := gin.Default()
	r.Use(middleware.Cors())
	g := r.Group("/boardgames")
	g.GET("/version", Version)
	g.GET("/queue", queue.GetStats(connection, "", ""))

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
		adapter.A(GetBoardgameScores),
		adapter.A(queries.GetPlayers),
		middleware.Result,
	)

	g.GET(
		"/info/:id",
		middleware.Id,
		adapter.A(GetBoardgameInfo),
		middleware.Result,
	)

	g.GET(
		"/options/:num",
		middleware.Num,
		adapter.A(GetPopularGamesForNum),
		middleware.Result,
	)

	g.POST(
		"/top",
		middleware.Body,
		adapter.A(GetPopularGames),
		middleware.Result,
	)

	g.POST(
		"/scrape/url/:id",
		middleware.Id,
		middleware.Body,
		adapter.A(ScrapeUrl),
		middleware.Result,
	)

	g.POST(
		"/scrape/:id",
		middleware.Id,
		middleware.Body,
		adapter.A(ScrapeStore),
		middleware.Result,
	)

	r.Run(":10002")
}
