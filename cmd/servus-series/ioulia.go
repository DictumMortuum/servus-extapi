package main

import (
	"log"
	"os"
	"time"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/nas"
	"github.com/go-rod/rod"
	"github.com/urfave/cli/v2"
)

func ioulia(ctx *cli.Context) error {
	episodes := scrapeIouliaEpisodes("https://www.alphatv.gr/show/to-proxenio-tis-ioulias/?vtype=player&vid=57348&showId=1530&year=2023")

	if len(episodes) == 0 {
		return nil
	}

	DB, err := db.DatabaseX()
	if err != nil {
		return err
	}
	defer DB.Close()

	for _, episode := range episodes {
		// rs, _ := nas.Exists(DB, episode)
		// if rs == nil {
		_, err := nas.Insert(DB, "ioulia", episode)
		if err != nil {
			return err
		}

		// payload := map[string]any{
		// 	"url":   episode,
		// 	"path":  "/volume1/plex/greek series/Ioulia/",
		// 	"owner": "dimitris@dictummortuum.com",
		// 	"group": "dimitris@dictummortuum.com",
		// }

		// err = nas.YoutubeDL(payload)
		// if err != nil {
		// 	return err
		// }

		// log.Println(payload)
		// }
	}

	return nil
}

func scrapeIouliaEpisodes(url string) []string {
	urls := []string{}

	browser := rod.New().MustConnect().Trace(false).Timeout(120 * time.Second)
	defer browser.MustClose()

	video := browser.MustPage(url)
	defer video.Close()
	video.MustWaitStable()
	rs := video.MustElement(`#currentvideourl`).MustAttribute("data-plugin-player")
	log.Println(*rs)

	err := os.RemoveAll("/tmp/rod")
	if err != nil {
		log.Fatal(err)
	}

	return urls
}
