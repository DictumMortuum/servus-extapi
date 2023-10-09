package main

import (
	"log"
	"os"
	"time"

	"github.com/go-rod/rod"
)

func scrapeVideo(url string) string {
	browser := rod.New().MustConnect().Trace(false).Timeout(120 * time.Second)
	defer browser.MustClose()

	video := browser.MustPage(url)
	defer video.Close()
	video.MustWaitStable()
	rs := video.MustElement(`video source`).MustProperty("src")

	err := os.RemoveAll("/tmp/rod")
	if err != nil {
		log.Fatal(err)
	}

	return rs.String()
}
