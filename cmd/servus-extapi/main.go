package main

import (
	"log"

	"github.com/DictumMortuum/servus-extapi/pkg/bgg"
	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/middleware"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Databases map[string]string `config:"databases"`
}

var (
	Cfg Config
)

func Version(c *gin.Context) {
	rs := map[string]any{
		"version": "v0.0.8",
	}
	c.AbortWithStatusJSON(200, rs)
}

func Route(router *gin.RouterGroup, endpoint string, obj Routable) {
	group := router.Group("/" + endpoint)
	group.Use(func(c *gin.Context) {
		c.Set("apimodel", obj)
		c.Next()
	})
	{
		group.GET("", OpenDB, CountMany, GetMany(obj.List), CloseDB)
		group.GET("/:id", OpenDB, Id, GetOne(obj.Get), CloseDB)
		group.PUT("/:id", OpenDB, Id, Body, UpdateOne(obj.Update), CloseDB)
		group.POST("", OpenDB, Body, CreateOne(obj.Create), CloseDB)
		group.DELETE("/:id", OpenDB, Id, DeleteOne(obj.Delete), CloseDB)
	}
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
	Route(g, "players", model.Player{})
	Route(g, "plays", model.Play{})
	Route(g, "stats", model.Stat{})
	Route(g, "prices", model.Price{})
	Route(g, "locations", model.Location{})
	Route(g, "stores", model.Store{})
	Route(g, "boardgames", model.Boardgame{})
	Route(g, "bgstatsplayers", model.BGStatsPlayer{})
	Route(g, "bgstatslocations", model.BGStatsLocation{})
	Route(g, "bgstatsgames", model.BGStatsGame{})
	Route(g, "bgstats", model.BGStat{})
	Route(g, "bgstatsplays", model.BGStatsPlay{})
	Route(g, "ignoredprices", model.IgnoredPrice{})
	Route(g, "ignorednames", model.IgnoredName{})
	Route(g, "cachedprices", model.CachedPrice{})

	g.POST("/bgstatsupload", G(CreateBGStats))

	cachedPrices := model.CachedPrice{}
	g.GET("/cachedprices/search/:id", OpenDB, Id, LoadOne(cachedPrices.Get), G(bgg.SearchCachedPriceOnBgg), CloseDB)
	g.POST("/cachedprices/create/:id", OpenDB, Id, G(cachedPrices.CreatePrice), CloseDB)

	prices := model.Price{}
	g.GET("/prices/search/:id", OpenDB, Id, LoadOne(prices.Get), G(bgg.SearchCachedPriceOnBgg), CloseDB)

	log.Fatal(r.Run(":10000"))
}
