package main

import (
	"log"
	"net/http"
	"os"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Version(c *gin.Context) {
	rs := map[string]any{
		"version": "v0.0.7",
	}
	c.AbortWithStatusJSON(200, rs)
}

func readinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func livenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
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

	err = os.MkdirAll(config.Cfg.Deco.Folder, 0777)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/version", Version)
	r.GET("/metrics", gin.WrapF(Metrics()))
	r.GET("/readiness", gin.WrapF(readinessHandler()))
	r.GET("/liveness", gin.WrapF(livenessHandler()))
	log.Fatal(r.Run(config.Cfg.FileExporterPort))
}
