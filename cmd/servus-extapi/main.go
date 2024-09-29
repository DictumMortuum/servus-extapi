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
		"version": "v0.0.30",
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
	adapter.RaRoute(g, "players", model.Player{})
	adapter.RaRoute(g, "plays", model.Play{})
	adapter.RaRoute(g, "stats", model.Stat{})
	adapter.RaRoute(g, "boardgameprices", model.BoardgamePrice{})
	adapter.RaRoute(g, "prices", model.Price{})
	adapter.RaRoute(g, "locations", model.Location{})
	adapter.RaRoute(g, "stores", model.Store{})
	adapter.RaRoute(g, "boardgames", model.Boardgame{})
	adapter.RaRoute(g, "bgstatsplayers", model.BGStatsPlayer{})
	adapter.RaRoute(g, "bgstatslocations", model.BGStatsLocation{})
	adapter.RaRoute(g, "bgstatsgames", model.BGStatsGame{})
	adapter.RaRoute(g, "bgstats", model.BGStat{})
	adapter.RaRoute(g, "bgstatsplays", model.BGStatsPlay{})
	adapter.RaRoute(g, "ignoredprices", model.IgnoredPrice{})
	adapter.RaRoute(g, "ignorednames", model.IgnoredName{})
	adapter.RaRoute(g, "cachedprices", model.CachedPrice{})
	adapter.RaRoute(g, "youtubedl", model.YoutubeDL{})
	adapter.RaRoute(g, "tables", model.Table{})
	adapter.RaRoute(g, "tableparticipants", model.TableParticipant{})
	adapter.RaRoute(g, "books", model.Book{})
	adapter.RaRoute(g, "eurovisionparticipations", model.EurovisionParticipation{})
	adapter.RaRoute(g, "eurovisionvotes", model.EurovisionVote{})
	adapter.RaRoute(g, "finderusers", model.FinderUser{})
	adapter.RaRoute(g, "wishlist", model.Wishlist{})

	// jwt := middleware.Jwt("http://sol.dictummortuum.com:3567/.well-known/jwks.json")

	g.POST("/bgstatsupload", adapter.G(CreateBGStats))
	g.GET("/eurovisionvotes/user/:id", middleware.Id, adapter.A(model.GetEurovisionVoteByUserId), middleware.ResultRa)
	g.GET("/eurovisionvotes/all", adapter.A(model.GetEurovisionVotes), middleware.ResultRa)
	g.GET("/eurovisionparticipations/user/:id", middleware.Id, adapter.A(model.GetEurovisionParticipationsByUserId), middleware.ResultRa)
	g.GET("/players/email/:id", middleware.Id, adapter.A(model.GetPlayerByEmail), middleware.ResultRa)

	// cachedPrices := model.CachedPrice{}
	// g.GET("/cachedprices/search/:id", OpenDB, Id, LoadOne(cachedPrices.Get), adapter.G(bgg.SearchCachedPriceOnBgg), CloseDB)
	// g.POST("/cachedprices/create/:id", OpenDB, Id, adapter.G(cachedPrices.CreatePrice), CloseDB)

	// prices := model.BoardgamePrice{}
	// g.GET("/prices/search/:id", OpenDB, Id, LoadOne(prices.Get), adapter.G(bgg.SearchCachedPriceOnBgg), CloseDB)

	log.Fatal(r.Run(":10000"))
}
