package queue

import (
	"log"
	"net/http"
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/gin-gonic/gin"
)

func GetStats(connection rmq.Connection, layout string, refresh string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		queues, err := connection.GetOpenQueues()
		if err != nil {
			log.Fatal(err)
		}

		stats, err := connection.CollectStats(queues)
		if err != nil {
			log.Fatal(err)
		}

		ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(stats.GetHtml(layout, refresh)))
	}
}

func Cleaner(connection rmq.Connection) {
	cleaner := rmq.NewCleaner(connection)

	for range time.Tick(time.Minute) {
		cleaner.Clean()
	}
}
