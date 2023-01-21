package main

import (
	"context"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	"github.com/itsjamie/gin-cors"
	"log"
	"time"
)

type Config struct {
	Databases map[string]string `config:"databases"`
}

var (
	Cfg Config
)

func SetConfig() gin.HandlerFunc {
	return cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type, Bearer, range, apikey",
		ExposedHeaders:  "x-total-count, Content-Range",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	})
}

func Version(c *gin.Context) {
	rs := map[string]any{
		"version": "v0.0.1",
	}
	c.AbortWithStatusJSON(200, rs)
}

func Route(router *gin.Engine, endpoint string, obj Routable) {
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
	loader := confita.NewLoader(
		file.NewBackend("/etc/conf.d/servusrc.yml"),
	)

	err := loader.Load(context.Background(), &Cfg)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.Use(SetConfig())
	// router.Use(CORSMiddleware())
	router.GET("/version", Version)
	Route(router, "players", model.Player{})
	Route(router, "plays", model.Play{})
	Route(router, "stats", model.Stat{})
	Route(router, "prices", model.Price{})
	Route(router, "locations", model.Location{})
	Route(router, "stores", model.Store{})
	Route(router, "boardgames", model.Boardgame{})
	Route(router, "bgstatsplayers", model.BGStatsPlayer{})
	Route(router, "bgstatslocations", model.BGStatsLocation{})
	Route(router, "bgstatsgames", model.BGStatsGame{})
	Route(router, "bgstats", model.BGStat{})
	Route(router, "bgstatsplays", model.BGStatsPlay{})

	router.POST("/bgstatsupload", G(CreateBGStats))

	log.Fatal(router.Run(":10000"))
}
