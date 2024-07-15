package main

import (
	"log"
	"net/http"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Version(c *gin.Context) {
	rs := map[string]any{
		"version": "v0.0.5",
	}
	c.AbortWithStatusJSON(200, rs)
}

func Metrics() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := readFiles()
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}

		promhttp.Handler().ServeHTTP(writer, request)
	}
}

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	g := r.Group("/exporter/file")
	g.GET("/version", Version)
	g.GET("/metrics", gin.WrapF(Metrics()))
	log.Fatal(r.Run(config.Cfg.FileExporterPort))
}
