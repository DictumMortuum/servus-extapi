package main

import (
	"log"

	"github.com/DictumMortuum/servus-extapi/pkg/adapter"
	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/middleware"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/gin-gonic/gin"
)

func Version(c *gin.Context) {
	rs := map[string]any{
		"version": "v0.0.11",
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
	g := r.Group("/rest")
	g.GET("/version", Version)
	adapter.Route(g, "players", model.Player{})
	adapter.Route(g, "plays", model.Play{})
	adapter.Route(g, "stats", model.Stat{})
	adapter.Route(g, "boardgameprices", model.BoardgamePrice{})
	adapter.Route(g, "prices", model.Price{})
	adapter.Route(g, "locations", model.Location{})
	adapter.Route(g, "stores", model.Store{})
	adapter.Route(g, "boardgames", model.Boardgame{})
	adapter.Route(g, "bgstatsplayers", model.BGStatsPlayer{})
	adapter.Route(g, "bgstatslocations", model.BGStatsLocation{})
	adapter.Route(g, "bgstatsgames", model.BGStatsGame{})
	adapter.Route(g, "bgstats", model.BGStat{})
	adapter.Route(g, "bgstatsplays", model.BGStatsPlay{})
	adapter.Route(g, "ignoredprices", model.IgnoredPrice{})
	adapter.Route(g, "ignorednames", model.IgnoredName{})
	adapter.Route(g, "cachedprices", model.CachedPrice{})

	g.POST("/bgstatsupload", adapter.G(CreateBGStats))

	// cachedPrices := model.CachedPrice{}
	// g.GET("/cachedprices/search/:id", OpenDB, Id, LoadOne(cachedPrices.Get), adapter.G(bgg.SearchCachedPriceOnBgg), CloseDB)
	// g.POST("/cachedprices/create/:id", OpenDB, Id, adapter.G(cachedPrices.CreatePrice), CloseDB)

	// prices := model.BoardgamePrice{}
	// g.GET("/prices/search/:id", OpenDB, Id, LoadOne(prices.Get), adapter.G(bgg.SearchCachedPriceOnBgg), CloseDB)

	log.Fatal(r.Run(":10000"))
}
