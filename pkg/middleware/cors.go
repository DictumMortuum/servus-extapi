package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() func(*gin.Context) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"range", "content-type"}
	config.ExposeHeaders = []string{"X-Total-Count", "Content-Range", "Content-Description", "Content-Disposition", "Filename"}
	return cors.New(config)
}
